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
	ormService := service.DI().OrmEngine().Clone()
	eventsConsumer := ormService.GetEventBroker().Consumer(consumer.GetGroupName(groupNameSuffix))

	var attempts uint8

	for {
		attempts++
		log.Printf("Starting %s", consumer.GetQueueName())

		started := eventsConsumer.Consume(r.ctx, prefetchCount, func(events []beeorm.Event) {
			log.Printf("We have %d new dirty events in %s", len(events), consumer.GetQueueName())

			if err := consumer.Consume(ormService, events); err != nil {
				panic(err)
			}

			log.Printf("We consumed %d dirty events in %s", len(events), consumer.GetQueueName())
		})

		if started {
			break
		}

		if attempts == 9 {
			service.DI().ErrorLogger().LogError("failed to start stream consumer: " + consumer.GetQueueName())
			break
		}
		log.Printf("failed to start %s. Sleeping for 90 seconds", consumer.GetQueueName())
		time.Sleep(time.Second * 92)
	}
}

func (r *ConsumerRunner) RunConsumerOne(consumer ConsumerOne, groupNameSuffix *string, prefetchCount int) {
	ormService := service.DI().OrmEngine().Clone()

	eventsConsumer := ormService.GetEventBroker().Consumer(consumer.GetGroupName(groupNameSuffix))

	eventsConsumer.Consume(r.ctx, prefetchCount, func(events []beeorm.Event) {
		log.Printf("We have %d new dirty events in %s", len(events), consumer.GetQueueName())

		for _, event := range events {
			if err := consumer.Consume(ormService, event); err != nil {
				panic(err)
			}
			event.Ack()
		}

		log.Printf("We consumed %d dirty events in %s", len(events), consumer.GetQueueName())
	})
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
			queueName := consumer.GetQueueName(moduloID)
			consumerGroupName := consumer.GetGroupName(moduloID, groupNameSuffix)

			log.Printf("RunConsumerOneByModulo started goroutine %d (%s)", moduloID, queueName)

			ormService := service.DI().OrmEngine().Clone()
			eventsConsumer := ormService.GetEventBroker().Consumer(consumerGroupName)

			log.Printf("RunConsumerOneByModulo is ready to consume events (%s)", queueName)

			innerWG.Done()
			outerWG.Done()

			eventsConsumer.Consume(r.ctx, prefetchCount, func(events []beeorm.Event) {
				log.Printf("We have %d new dirty events in %s", len(events), consumerGroupName)

				for _, event := range events {
					if err := consumer.Consume(ormService, event); err != nil {
						panic(err)
					}
					event.Ack()
				}

				log.Printf("We consumed %d dirty events in %s", len(events), consumerGroupName)
			})

			log.Printf("RunConsumerOneByModulo exited unexpectedly for goroutine %d (%s)", moduloID, queueName)
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
			queueName := consumer.GetQueueName(moduloID)
			consumerGroupName := consumer.GetGroupName(moduloID, groupNameSuffix)

			log.Printf("RunConsumerManyByModulo started goroutine %d (%s)", moduloID, queueName)

			ormService := service.DI().OrmEngine().Clone()
			eventsConsumer := ormService.GetEventBroker().Consumer(consumerGroupName)

			log.Printf("RunConsumerManyByModulo is ready to consume events (%s)", queueName)

			innerWG.Done()
			outerWG.Done()

			eventsConsumer.Consume(r.ctx, prefetchCount, func(events []beeorm.Event) {
				log.Printf("We have %d new dirty events in %s", len(events), consumerGroupName)

				if err := consumer.Consume(ormService, events); err != nil {
					panic(err)
				}

				log.Printf("We consumed %d dirty events in %s", len(events), consumerGroupName)
			})

			log.Printf("RunConsumerManyByModulo exited unexpectedly for goroutine %d (%s)", moduloID, queueName)
		})

		innerWG.Wait()
	}

	outerWG.Wait()

	log.Printf("RunConsumerManyByModulo exited (%s)", baseQueueName)
}
