package queue

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/latolukasz/beeorm/v2"

	"github.com/coretrix/hitrix"
	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/service"
)

const (
	obtainLockRetryDuration = time.Second
)

type ConsumerOneByModulo interface {
	GetMaxModulo() int
	Consume(ormService *datalayer.ORM, event beeorm.Event) error
	GetQueueName(moduloID int) string
	GetGroupName(moduloID int, suffix *string) string
}

type ConsumerManyByModulo interface {
	GetMaxModulo() int
	Consume(ormService *datalayer.ORM, events []beeorm.Event) error
	GetQueueName(moduloID int) string
	GetGroupName(moduloID int, suffix *string) string
}

type ConsumerOne interface {
	Consume(ormService *datalayer.ORM, event beeorm.Event) error
	GetQueueName() string
	GetGroupName(suffix *string) string
}

type ConsumerMany interface {
	Consume(ormService *datalayer.ORM, events []beeorm.Event) error
	GetQueueName() string
	GetGroupName(suffix *string) string
}

type ConsumerRunner struct {
	ctx context.Context
}

func NewConsumerRunner(ctx context.Context) *ConsumerRunner {
	return &ConsumerRunner{ctx: ctx}
}

func (r *ConsumerRunner) RunConsumerMany(consumer ConsumerMany, groupNameSuffix *string, prefetchCount int) {
	queueName := consumer.GetQueueName()

	log.Printf("RunConsumerMany initialized (%s)", queueName)

	ormService := service.DI().OrmEngine().Clone()
	eventsConsumer := ormService.GetEventBroker().Consumer(consumer.GetGroupName(groupNameSuffix))

	service.DI().App().Add(1)

	defer service.DI().App().Done()

	for {
		// eventsConsumer.Consume should block and not return anything
		// if it returns true => this consumer is exited with no errors, but still not consuming
		// if it returns false => this consumer is exited with error "could not obtain lock", so we should retry
		if exitedWithNoErrors := eventsConsumer.Consume(r.ctx, prefetchCount, func(events []beeorm.Event) {
			log.Printf("We have %d new dirty events in %s", len(events), queueName)

			if err := consumer.Consume(ormService, events); err != nil {
				panic(err)
			}

			log.Printf("We consumed %d dirty events in %s", len(events), queueName)
		}); !exitedWithNoErrors {
			log.Printf("RunConsumerMany failed to start (%s) - retrying in %.1f seconds", queueName, obtainLockRetryDuration.Seconds())
			time.Sleep(obtainLockRetryDuration)

			continue
		}

		log.Println("eventsConsumer.Consume returned true")
		log.Printf("RunConsumerMany exited (%s)", queueName)

		break
	}
}

func (r *ConsumerRunner) RunConsumerOne(consumer ConsumerOne, groupNameSuffix *string, prefetchCount int) {
	queueName := consumer.GetQueueName()

	log.Printf("RunConsumerOne initialized (%s)", queueName)

	ormService := service.DI().OrmEngine().Clone()
	eventsConsumer := ormService.GetEventBroker().Consumer(consumer.GetGroupName(groupNameSuffix))

	service.DI().App().Add(1)

	defer service.DI().App().Done()

	for {
		// eventsConsumer.Consume should block and not return anything
		// if it returns true => this consumer is exited with no errors, but still not consuming
		// if it returns false => this consumer is exited with error "could not obtain lock", so we should retry
		if exitedWithNoErrors := eventsConsumer.Consume(r.ctx, prefetchCount, func(events []beeorm.Event) {
			log.Printf("We have %d new dirty events in %s", len(events), queueName)

			for _, event := range events {
				if err := consumer.Consume(ormService, event); err != nil {
					panic(err)
				}
				event.Ack()
			}

			log.Printf("We consumed %d dirty events in %s", len(events), queueName)
		}); !exitedWithNoErrors {
			log.Printf("RunConsumerOne failed to start (%s) - retrying in %.1f seconds", queueName, obtainLockRetryDuration.Seconds())
			time.Sleep(obtainLockRetryDuration)

			continue
		}

		log.Println("eventsConsumer.Consume returned true")

		break
	}
	log.Printf("RunConsumerOne exited (%s)", queueName)
}

func (r *ConsumerRunner) RunConsumerOneByModulo(consumer ConsumerOneByModulo, groupNameSuffix *string, prefetchCount int) {
	maxModulo := consumer.GetMaxModulo()

	baseQueueName := ""

	queueNameParts := strings.Split(consumer.GetQueueName(maxModulo), "_")
	if len(queueNameParts) > 0 {
		baseQueueName = queueNameParts[0]
	}

	log.Printf("RunConsumerOneByModulo initialized (%s)", baseQueueName)

	for moduloID := 1; moduloID <= maxModulo; moduloID++ {
		currentModulo := moduloID

		hitrix.GoroutineWithRestart(func() {
			queueName := consumer.GetQueueName(currentModulo)
			consumerGroupName := consumer.GetGroupName(currentModulo, groupNameSuffix)

			log.Printf("RunConsumerOneByModulo started goroutine %d (%s)", currentModulo, queueName)

			ormService := service.DI().OrmEngine().Clone()
			eventsConsumer := ormService.GetEventBroker().Consumer(consumerGroupName)
			service.DI().App().Add(1)
			defer service.DI().App().Done()

			for {
				// eventsConsumer.Consume should block and not return anything
				// if it returns true => this consumer is exited with no errors, but still not consuming
				// if it returns false => this consumer is exited with error "could not obtain lock", so we should retry
				if exitedWithNoErrors := eventsConsumer.Consume(r.ctx, prefetchCount, func(events []beeorm.Event) {
					log.Printf("We have %d new dirty events in %s", len(events), consumerGroupName)

					for _, event := range events {
						if err := consumer.Consume(ormService, event); err != nil {
							panic(err)
						}
						event.Ack()
					}

					log.Printf("We consumed %d dirty events in %s", len(events), consumerGroupName)
				}); !exitedWithNoErrors {
					log.Printf(
						"RunConsumerOneByModulo failed to start for goroutine %d (%s) - retrying in %.1f seconds",
						currentModulo,
						queueName,
						obtainLockRetryDuration.Seconds())
					time.Sleep(obtainLockRetryDuration)

					continue
				}
				log.Printf("eventsConsumer.Consume returned true for goroutine %d (%s)", currentModulo, queueName)
				log.Printf("RunConsumerOneByModulo exited (%s)", baseQueueName)

				break
			}

			log.Printf("RunConsumerOneByModulo goroutine %d (%s) exited", currentModulo, queueName)
		})

		time.Sleep(time.Second)
	}
}

func (r *ConsumerRunner) RunConsumerManyByModulo(consumer ConsumerManyByModulo, groupNameSuffix *string, prefetchCount int) {
	maxModulo := consumer.GetMaxModulo()

	baseQueueName := ""
	queueNameParts := strings.Split(consumer.GetQueueName(maxModulo), "_")

	if len(queueNameParts) > 0 {
		baseQueueName = queueNameParts[0]
	}

	log.Printf("RunConsumerManyByModulo initialized (%s)", baseQueueName)

	for moduloID := 1; moduloID <= maxModulo; moduloID++ {
		currentModulo := moduloID

		hitrix.GoroutineWithRestart(func() {
			queueName := consumer.GetQueueName(currentModulo)
			consumerGroupName := consumer.GetGroupName(currentModulo, groupNameSuffix)

			log.Printf("RunConsumerManyByModulo started goroutine %d (%s)", currentModulo, queueName)

			ormService := service.DI().OrmEngine().Clone()
			eventsConsumer := ormService.GetEventBroker().Consumer(consumerGroupName)
			service.DI().App().Add(1)
			defer service.DI().App().Done()

			for {
				// eventsConsumer.Consume should block and not return anything
				// if it returns true => this consumer is exited with no errors, but still not consuming
				// if it returns false => this consumer is exited with error "could not obtain lock", so we should retry
				if exitedWithNoErrors := eventsConsumer.Consume(r.ctx, prefetchCount, func(events []beeorm.Event) {
					log.Printf("We have %d new dirty events in %s", len(events), consumerGroupName)

					if err := consumer.Consume(ormService, events); err != nil {
						panic(err)
					}

					log.Printf("We consumed %d dirty events in %s", len(events), consumerGroupName)
				}); !exitedWithNoErrors {
					log.Printf(
						"RunConsumerManyByModulo failed to start for goroutine %d (%s) - retrying in %.1f seconds",
						currentModulo,
						queueName,
						obtainLockRetryDuration.Seconds())
					time.Sleep(obtainLockRetryDuration)

					continue
				}

				log.Printf("eventsConsumer.Consume returned true for goroutine %d (%s)", currentModulo, queueName)
				log.Printf("RunConsumerManyByModulo exited (%s)", baseQueueName)

				break
			}

			log.Printf("RunConsumerManyByModulo goroutine %d (%s) exited", currentModulo, queueName)
		})

		time.Sleep(time.Second)
	}
}

type ScalableConsumerRunner struct {
	ctx       context.Context
	redisPool string
}

func NewScalableConsumerRunner(ctx context.Context, redisPool string) *ScalableConsumerRunner {
	return &ScalableConsumerRunner{ctx: ctx, redisPool: redisPool}
}

func (r *ScalableConsumerRunner) RunScalableConsumerMany(consumer ConsumerMany, groupNameSuffix *string, prefetchCount int) {
	ormService := service.DI().OrmEngine().Clone()
	redis := ormService.GetRedis(r.redisPool)

	queueName := consumer.GetQueueName()
	consumerGroupName := consumer.GetGroupName(groupNameSuffix)

	currentIndex := addConsumerGroup(redis, consumerGroupName)

	log.Printf("RunScalableConsumerMany index (%d) initialized (%s)", currentIndex, queueName)

	eventsConsumer := ormService.GetEventBroker().Consumer(consumerGroupName)

	service.DI().App().Add(1)

	defer service.DI().App().Done()

	for {
		// eventsConsumer.ConsumeMany should block and not return anything
		// if it returns true => this consumer is exited with no errors, but still not consuming
		// if it returns false => this consumer is exited with error "could not obtain lock", so we should retry
		if exitedWithNoErrors := eventsConsumer.ConsumeMany(r.ctx, currentIndex, prefetchCount, func(events []beeorm.Event) {
			log.Printf("We have %d new dirty events in %s", len(events), queueName)

			if err := consumer.Consume(ormService, events); err != nil {
				removeConsumerGroup(eventsConsumer, redis, consumerGroupName, currentIndex)
				panic(err)
			}

			log.Printf("We consumed %d dirty events in %s", len(events), queueName)
		}); !exitedWithNoErrors {
			log.Printf("RunScalableConsumerMany failed to start (%s) - retrying in %.1f seconds", queueName, obtainLockRetryDuration.Seconds())
			time.Sleep(obtainLockRetryDuration)

			continue
		}

		log.Println("eventsConsumer.ConsumeMany returned true")

		break
	}

	removeConsumerGroup(eventsConsumer, redis, consumerGroupName, currentIndex)
}

func (r *ScalableConsumerRunner) RunScalableConsumerOne(consumer ConsumerOne, groupNameSuffix *string, prefetchCount int) {
	ormService := service.DI().OrmEngine().Clone()
	redis := ormService.GetRedis(r.redisPool)

	queueName := consumer.GetQueueName()
	consumerGroupName := consumer.GetGroupName(groupNameSuffix)

	currentIndex := addConsumerGroup(redis, consumerGroupName)

	log.Printf("RunScalableConsumerOne index (%d) initialized (%s)", currentIndex, queueName)

	eventsConsumer := ormService.GetEventBroker().Consumer(consumerGroupName)

	service.DI().App().Add(1)

	defer service.DI().App().Done()

	for {
		// eventsConsumer.ConsumeMany should block and not return anything
		// if it returns true => this consumer is exited with no errors, but still not consuming
		// if it returns false => this consumer is exited with error "could not obtain lock", so we should retry
		if exitedWithNoErrors := eventsConsumer.ConsumeMany(r.ctx, currentIndex, prefetchCount, func(events []beeorm.Event) {
			log.Printf("We have %d new dirty events in %s", len(events), queueName)

			for _, event := range events {
				if err := consumer.Consume(ormService, event); err != nil {
					removeConsumerGroup(eventsConsumer, redis, consumerGroupName, currentIndex)
					panic(err)
				}
				event.Ack()
			}

			log.Printf("We consumed %d dirty events in %s", len(events), queueName)
		}); !exitedWithNoErrors {
			log.Printf("RunScalableConsumerOne failed to start (%s) - retrying in %.1f seconds", queueName, obtainLockRetryDuration.Seconds())
			time.Sleep(obtainLockRetryDuration)

			continue
		}

		log.Println("eventsConsumer.ConsumeMany returned true")

		break
	}

	removeConsumerGroup(eventsConsumer, redis, consumerGroupName, currentIndex)

	log.Printf("RunScalableConsumerOne exited (%s)", queueName)
}

const consumerGroupsKey = "consumer_groups"

type indexer struct {
	LatestIndex           int
	ActiveConsumerIndexes map[int]*struct{}
}

func addConsumerGroup(redis beeorm.RedisCache, consumerGroupName string) int {
	indexerValue, err := getConsumerGroupIndexer(redis, consumerGroupName)
	if err != nil {
		panic(err)
	}

	if indexerValue == nil {
		indexerValue = &indexer{}
	}

	if indexerValue.ActiveConsumerIndexes == nil {
		indexerValue.ActiveConsumerIndexes = map[int]*struct{}{}
	}

	indexerValue.LatestIndex++
	indexerValue.ActiveConsumerIndexes[indexerValue.LatestIndex] = &struct{}{}

	err = setConsumerGroupIndexer(redis, consumerGroupName, indexerValue)
	if err != nil {
		panic(err)
	}

	return indexerValue.LatestIndex
}

func removeConsumerGroup(consumer beeorm.EventsConsumer, redis beeorm.RedisCache, consumerGroupName string, indexToRemove int) {
	indexerValue, err := getConsumerGroupIndexer(redis, consumerGroupName)
	if err != nil {
		panic(err)
	}

	delete(indexerValue.ActiveConsumerIndexes, indexToRemove)

	err = setConsumerGroupIndexer(redis, consumerGroupName, indexerValue)
	if err != nil {
		panic(err)
	}

	// transfer pending items from stopped consumer to another if available as per:
	// https://beeorm.io/guide/event_broker.html#consumers-scaling
	if len(indexerValue.ActiveConsumerIndexes) != 0 {
		indexToTransferClaimedItems := 0
		for index := range indexerValue.ActiveConsumerIndexes {
			indexToTransferClaimedItems = index

			break
		}

		if indexToTransferClaimedItems != 0 {
			log.Printf("claiming from %d to %d", indexToRemove, indexToTransferClaimedItems)
			consumer.Claim(indexToRemove, indexToTransferClaimedItems)
		}
	}
}

func setConsumerGroupIndexer(redis beeorm.RedisCache, consumerGroupName string, indexer *indexer) error {
	marshaled, err := json.Marshal(indexer)
	if err != nil {
		return err
	}

	redis.HSet(consumerGroupsKey, consumerGroupName, marshaled)

	return err
}

func getConsumerGroupIndexer(redis beeorm.RedisCache, consumerGroupName string) (*indexer, error) {
	marshaled, has := redis.HGet(consumerGroupsKey, consumerGroupName)
	if !has {
		return nil, nil
	}

	indexer := &indexer{}
	if err := json.Unmarshal([]byte(marshaled), indexer); err != nil {
		return nil, err
	}

	return indexer, nil
}
