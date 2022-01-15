package queue

import (
	"context"
	"log"
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
			if err := consumer.Consume(ormService, events); err != nil {
				panic(err)
			}
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
		for _, event := range events {
			if err := consumer.Consume(ormService, event); err != nil {
				panic(err)
			}
			event.Ack()
		}
	})
}

func (r *ConsumerRunner) RunConsumerOneByModulo(consumer ConsumerOneByModulo, groupNameSuffix *string, prefetchCount int) {
	for moduloID := 1; moduloID <= consumer.GetMaxModulo(); moduloID++ {
		currentModulo := moduloID

		hitrix.GoroutineWithRestart(func() {
			ormService := service.DI().OrmEngine().Clone()

			eventsConsumer := ormService.GetEventBroker().Consumer(consumer.GetGroupName(currentModulo, groupNameSuffix))
			eventsConsumer.Consume(r.ctx, prefetchCount, func(events []beeorm.Event) {
				for _, event := range events {
					if err := consumer.Consume(ormService, event); err != nil {
						panic(err)
					}
					event.Ack()
				}
			})
		})
	}
}

func (r *ConsumerRunner) RunConsumerManyByModulo(consumer ConsumerManyByModulo, groupNameSuffix *string, prefetchCount int) {
	for moduloID := 1; moduloID <= consumer.GetMaxModulo(); moduloID++ {
		currentModulo := moduloID

		hitrix.GoroutineWithRestart(func() {
			ormService := service.DI().OrmEngine().Clone()
			eventsConsumer := ormService.GetEventBroker().Consumer(consumer.GetGroupName(currentModulo, groupNameSuffix))
			eventsConsumer.Consume(r.ctx, prefetchCount, func(events []beeorm.Event) {
				if err := consumer.Consume(ormService, events); err != nil {
					panic(err)
				}
			})
		})
	}
}
