package queue

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix"
	"github.com/coretrix/hitrix/service"
)

const (
	obtainLockRetryDuration = time.Second
)

type ConsumerOneByModulo interface {
	GetMaxModulo() int
	Consume(ormService *beeorm.Engine, event beeorm.Event) error
	GetQueueName(moduloID int) string
	GetGroupName(moduloID int, suffix *string) string
}

type ConsumerManyByModulo interface {
	GetMaxModulo() int
	Consume(ormService *beeorm.Engine, events []beeorm.Event) error
	GetQueueName(moduloID int) string
	GetGroupName(moduloID int, suffix *string) string
}

type ConsumerOne interface {
	Consume(ormService *beeorm.Engine, event beeorm.Event) error
	GetQueueName() string
	GetGroupName(suffix *string) string
}

type ConsumerMany interface {
	Consume(ormService *beeorm.Engine, events []beeorm.Event) error
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
		} else {
			log.Println("eventsConsumer.Consume returned true")
			break
		}
	}

	log.Printf("RunConsumerMany exited (%s)", queueName)
}

func (r *ConsumerRunner) RunConsumerOne(consumer ConsumerOne, groupNameSuffix *string, prefetchCount int) {
	queueName := consumer.GetQueueName()

	log.Printf("RunConsumerOne initialized (%s)", queueName)

	ormService := service.DI().OrmEngine().Clone()
	eventsConsumer := ormService.GetEventBroker().Consumer(consumer.GetGroupName(groupNameSuffix))

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
		} else {
			log.Println("eventsConsumer.Consume returned true")
			break
		}
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

	outerWG := sync.WaitGroup{}
	outerWG.Add(maxModulo)

	for moduloID := 1; moduloID <= maxModulo; moduloID++ {
		innerWG := sync.WaitGroup{}
		innerWG.Add(1)

		hitrix.GoroutineWithRestart(func() {
			currentModulo := moduloID

			queueName := consumer.GetQueueName(currentModulo)
			consumerGroupName := consumer.GetGroupName(currentModulo, groupNameSuffix)

			log.Printf("RunConsumerOneByModulo started goroutine %d (%s)", currentModulo, queueName)

			ormService := service.DI().OrmEngine().Clone()
			eventsConsumer := ormService.GetEventBroker().Consumer(consumerGroupName)

			innerWG.Done()

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
					log.Printf("RunConsumerOneByModulo failed to start for goroutine %d (%s) - retrying in %.1f seconds", currentModulo, queueName, obtainLockRetryDuration.Seconds())
					time.Sleep(obtainLockRetryDuration)
					continue
				} else {
					log.Printf("eventsConsumer.Consume returned true for goroutine %d (%s)", currentModulo, queueName)
					outerWG.Done()
					break
				}
			}

			log.Printf("RunConsumerOneByModulo goroutine %d (%s) exited", currentModulo, queueName)
		})

		innerWG.Wait()
	}

	outerWG.Wait()
	log.Printf("RunConsumerOneByModulo exited (%s)", baseQueueName)
}

func (r *ConsumerRunner) RunConsumerManyByModulo(consumer ConsumerManyByModulo, groupNameSuffix *string, prefetchCount int) {
	maxModulo := consumer.GetMaxModulo()

	baseQueueName := ""
	queueNameParts := strings.Split(consumer.GetQueueName(maxModulo), "_")
	if len(queueNameParts) > 0 {
		baseQueueName = queueNameParts[0]
	}

	log.Printf("RunConsumerManyByModulo initialized (%s)", baseQueueName)

	outerWG := sync.WaitGroup{}
	outerWG.Add(maxModulo)

	for moduloID := 1; moduloID <= maxModulo; moduloID++ {
		innerWG := sync.WaitGroup{}
		innerWG.Add(1)

		hitrix.GoroutineWithRestart(func() {
			currentModulo := moduloID

			queueName := consumer.GetQueueName(currentModulo)
			consumerGroupName := consumer.GetGroupName(currentModulo, groupNameSuffix)

			log.Printf("RunConsumerManyByModulo started goroutine %d (%s)", currentModulo, queueName)

			ormService := service.DI().OrmEngine().Clone()
			eventsConsumer := ormService.GetEventBroker().Consumer(consumerGroupName)

			innerWG.Done()

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
					log.Printf("RunConsumerManyByModulo failed to start for goroutine %d (%s) - retrying in %.1f seconds", currentModulo, queueName, obtainLockRetryDuration.Seconds())
					time.Sleep(obtainLockRetryDuration)
					continue
				} else {
					log.Printf("eventsConsumer.Consume returned true for goroutine %d (%s)", currentModulo, queueName)
					outerWG.Done()
					break
				}
			}

			log.Printf("RunConsumerManyByModulo goroutine %d (%s) exited", currentModulo, queueName)
		})

		innerWG.Wait()
	}

	outerWG.Wait()
	log.Printf("RunConsumerManyByModulo exited (%s)", baseQueueName)
}
