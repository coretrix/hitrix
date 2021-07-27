package queue

import (
	"context"
	"log"

	"github.com/coretrix/hitrix/service"

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

var consumerErrorHandler = func(err error, event beeorm.Event) {
	defer func() {
		if newError := recover(); newError != nil {
			log.Printf("Error: %v\nNew Error: %v", err, newError)
		}
	}()

	errorLoggerService, has := service.DI().ErrorLogger()
	if has {
		errorLoggerService.LogError(err)
	}
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
	eventsConsumer.SetErrorHandler(consumerErrorHandler)

	eventsConsumer.Consume(prefetchCount, func(events []beeorm.Event) {
		if err := consumer.Consume(events); err != nil {
			panic(err)
		}
	})
}

func (r *ConsumerRunner) RunConsumerOne(consumer ConsumerOne, groupNameSuffix *string, prefetchCount int) {
	eventsConsumer := r.ormService.GetEventBroker().Consumer(consumer.GetGroupName(groupNameSuffix))
	eventsConsumer.SetErrorHandler(consumerErrorHandler)

	eventsConsumer.Consume(prefetchCount, func(events []beeorm.Event) {
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

		go func() {
			eventsConsumer := r.ormService.GetEventBroker().Consumer(consumer.GetGroupName(currentModulo, groupNameSuffix))
			eventsConsumer.SetErrorHandler(consumerErrorHandler)
			eventsConsumer.Consume(prefetchCount, func(events []beeorm.Event) {
				for _, event := range events {
					if err := consumer.Consume(event); err != nil {
						panic(err)
					}
					event.Ack()
				}
			})
		}()
	}
}

func (r *ConsumerRunner) RunConsumerManyByModulo(consumer ConsumerManyByModulo, groupNameSuffix *string, prefetchCount int) {
	for moduloID := 1; moduloID <= consumer.GetMaxModulo(); moduloID++ {
		currentModulo := moduloID
		go func() {
			eventsConsumer := r.ormService.GetEventBroker().Consumer(consumer.GetGroupName(currentModulo, groupNameSuffix))
			eventsConsumer.SetErrorHandler(consumerErrorHandler)
			eventsConsumer.Consume(prefetchCount, func(events []beeorm.Event) {
				if err := consumer.Consume(events); err != nil {
					panic(err)
				}
			})
		}()
	}
}
