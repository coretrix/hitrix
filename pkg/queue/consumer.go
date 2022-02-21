package queue

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/jpillora/backoff"
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix"
	"github.com/coretrix/hitrix/service"
)

const (
	maxBackoffAttempts = 200
	minBackoffDuration = 200 * time.Millisecond
	maxBackoffDuration = 30 * time.Second
	backoffFactor      = 2
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

	b := &backoff.Backoff{
		Min:    minBackoffDuration,
		Max:    maxBackoffDuration,
		Factor: backoffFactor,
	}

	totalAttempts := 0
	for {
		totalAttempts++

		// eventsConsumer.Consume should block and not return anything
		// if it returns true => this consumer is exited with no errors, but still not consuming
		// if it returns false => this consumer is exited with error "could not obtain lock", so we should retry using exponential backoff
		if started := eventsConsumer.Consume(r.ctx, prefetchCount, func(events []beeorm.Event) {
			log.Printf("We have %d new dirty events in %s", len(events), queueName)

			if err := consumer.Consume(ormService, events); err != nil {
				panic(err)
			}

			log.Printf("We consumed %d dirty events in %s", len(events), queueName)
		}); !started {
			if totalAttempts > maxBackoffAttempts {
				service.DI().ErrorLogger().LogError(
					fmt.Sprintf("RunConsumerMany failed to start consumer after %d attempts (%s)", maxBackoffAttempts, queueName),
				)
				b.Reset()
				break
			}

			log.Printf("RunConsumerMany failed to start (%s) - retrying...", queueName)
			b.Duration()
			continue
		} else {
			b.Reset()
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

	b := &backoff.Backoff{
		Min:    minBackoffDuration,
		Max:    maxBackoffDuration,
		Factor: backoffFactor,
	}

	totalAttempts := 0
	for {
		totalAttempts++

		// eventsConsumer.Consume should block and not return anything
		// if it returns true => this consumer is exited with no errors, but still not consuming
		// if it returns false => this consumer is exited with error "could not obtain lock", so we should retry using exponential backoff
		if started := eventsConsumer.Consume(r.ctx, prefetchCount, func(events []beeorm.Event) {
			log.Printf("We have %d new dirty events in %s", len(events), queueName)

			for _, event := range events {
				if err := consumer.Consume(ormService, event); err != nil {
					panic(err)
				}
				event.Ack()
			}

			log.Printf("We consumed %d dirty events in %s", len(events), queueName)
		}); !started {
			if totalAttempts > maxBackoffAttempts {
				service.DI().ErrorLogger().LogError(
					fmt.Sprintf("RunConsumerOne failed to start consumer after %d attempts (%s)", maxBackoffAttempts, queueName),
				)
				b.Reset()
				break
			}

			log.Printf("RunConsumerOne failed to start (%s) - retrying...", queueName)
			b.Duration()
			continue
		} else {
			b.Reset()
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
			outerWG.Done()

			b := &backoff.Backoff{
				Min:    minBackoffDuration,
				Max:    maxBackoffDuration,
				Factor: backoffFactor,
			}

			totalAttempts := 0
			for {
				totalAttempts++

				// eventsConsumer.Consume should block and not return anything
				// if it returns true => this consumer is exited with no errors, but still not consuming
				// if it returns false => this consumer is exited with error "could not obtain lock", so we should retry using exponential backoff
				if started := eventsConsumer.Consume(r.ctx, prefetchCount, func(events []beeorm.Event) {
					log.Printf("We have %d new dirty events in %s", len(events), consumerGroupName)

					for _, event := range events {
						if err := consumer.Consume(ormService, event); err != nil {
							panic(err)
						}
						event.Ack()
					}

					log.Printf("We consumed %d dirty events in %s", len(events), consumerGroupName)
				}); !started {
					if totalAttempts > maxBackoffAttempts {
						service.DI().ErrorLogger().LogError(
							fmt.Sprintf("RunConsumerOneByModulo failed to start consumer after %d attempts (%s)", maxBackoffAttempts, queueName),
						)
						b.Reset()
						break
					}

					log.Printf("RunConsumerOneByModulo failed to start for goroutine %d (%s) - retrying...", currentModulo, queueName)
					b.Duration()
					continue
				} else {
					b.Reset()
					break
				}
			}

			log.Printf("RunConsumerOneByModulo exited for goroutine %d (%s)", currentModulo, queueName)
		})

		innerWG.Wait()
	}

	outerWG.Wait()
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
			outerWG.Done()

			b := &backoff.Backoff{
				Min:    minBackoffDuration,
				Max:    maxBackoffDuration,
				Factor: backoffFactor,
			}

			totalAttempts := 0
			for {
				totalAttempts++

				// eventsConsumer.Consume should block and not return anything
				// if it returns true => this consumer is exited with no errors, but still not consuming
				// if it returns false => this consumer is exited with error "could not obtain lock", so we should retry using exponential backoff
				if started := eventsConsumer.Consume(r.ctx, prefetchCount, func(events []beeorm.Event) {
					log.Printf("We have %d new dirty events in %s", len(events), consumerGroupName)

					if err := consumer.Consume(ormService, events); err != nil {
						panic(err)
					}

					log.Printf("We consumed %d dirty events in %s", len(events), consumerGroupName)
				}); !started {
					if totalAttempts > maxBackoffAttempts {
						service.DI().ErrorLogger().LogError(
							fmt.Sprintf("RunConsumerManyByModulo failed to start consumer after %d attempts (%s)", maxBackoffAttempts, queueName),
						)
						b.Reset()
						break
					}

					log.Printf("RunConsumerManyByModulo failed to start for goroutine %d (%s) - retrying...", currentModulo, queueName)
					b.Duration()
					continue
				} else {
					b.Reset()
					break
				}
			}

			log.Printf("RunConsumerManyByModulo exited for goroutine %d (%s)", currentModulo, queueName)
		})

		innerWG.Wait()
	}

	outerWG.Wait()
}
