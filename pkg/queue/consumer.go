package queue

import (
	"context"

	"github.com/coretrix/hitrix"
	"github.com/latolukasz/beeorm"
)

type ConsumerOneByModulo interface {
	GetMaxModulo() int
	Consume(event beeorm.Event) error
	GetQueueName(moduloID int) string
	GetGroupName(moduloID int, suffix *string) string
}

type ConsumerManyByModulo interface {
	GetMaxModulo() int
	Consume(events []beeorm.Event) error
	GetQueueName(moduloID int) string
	GetGroupName(moduloID int, suffix *string) string
}

type ConsumerOne interface {
	Consume(event beeorm.Event) error
	GetQueueName() string
	GetGroupName(suffix *string) string
}

type ConsumerMany interface {
	Consume(events []beeorm.Event) error
	GetQueueName() string
	GetGroupName(suffix *string) string
}

type ConsumerRunner struct {
	ctx        context.Context
	ormService *beeorm.Engine
}

func NewConsumerRunner(ctx context.Context, ormService *beeorm.Engine) *ConsumerRunner {
	return &ConsumerRunner{ctx: ctx, ormService: ormService}
}

func (r *ConsumerRunner) RunConsumerMany(consumer ConsumerMany, groupNameSuffix *string, prefetchCount int) {
	eventsConsumer := r.ormService.GetEventBroker().Consumer(consumer.GetGroupName(groupNameSuffix))

	eventsConsumer.Consume(r.ctx, prefetchCount, func(events []beeorm.Event) {
		if err := consumer.Consume(events); err != nil {
			panic(err)
		}
	})
}

func (r *ConsumerRunner) RunConsumerOne(consumer ConsumerOne, groupNameSuffix *string, prefetchCount int) {
	eventsConsumer := r.ormService.GetEventBroker().Consumer(consumer.GetGroupName(groupNameSuffix))

	eventsConsumer.Consume(r.ctx, prefetchCount, func(events []beeorm.Event) {
		for _, event := range events {
			if err := consumer.Consume(event); err != nil {
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
			eventsConsumer := r.ormService.GetEventBroker().Consumer(consumer.GetGroupName(currentModulo, groupNameSuffix))
			eventsConsumer.Consume(r.ctx, prefetchCount, func(events []beeorm.Event) {
				for _, event := range events {
					if err := consumer.Consume(event); err != nil {
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
			eventsConsumer := r.ormService.GetEventBroker().Consumer(consumer.GetGroupName(currentModulo, groupNameSuffix))
			eventsConsumer.Consume(r.ctx, prefetchCount, func(events []beeorm.Event) {
				if err := consumer.Consume(events); err != nil {
					panic(err)
				}
			})
		})
	}
}
