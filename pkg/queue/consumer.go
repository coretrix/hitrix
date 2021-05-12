package queue

import (
	"context"
	"log"

	"github.com/coretrix/hitrix/service"

	"github.com/latolukasz/orm"
)

type ConsumerOneByModulo interface {
	GetMaxModulo() int
	Consume(event orm.Event) error
	GetQueueName(moduloID int) string
	GetGroupName(moduloID int, suffix *string) string
}

type ConsumerManyByModulo interface {
	GetMaxModulo() int
	Consume(events []orm.Event) error
	GetQueueName(moduloID int) string
	GetGroupName(moduloID int, suffix *string) string
}

type ConsumerOne interface {
	Consume(event orm.Event) error
	GetQueueName() string
	GetGroupName(suffix *string) string
}

type ConsumerMany interface {
	Consume(events []orm.Event) error
	GetQueueName() string
	GetGroupName(suffix *string) string
}

var consumerErrorHandler = func(err interface{}, event orm.Event) error {
	defer func() {
		if newError := recover(); newError != nil {
			log.Printf("Error: %v\nNew Error: %v", err, newError)
		}
	}()

	errorLoggerService, has := service.DI().ErrorLogger()
	if has {
		errorLoggerService.LogError(err)
	}
	return nil
}

type ConsumerRunner struct {
	ctx        context.Context
	ormService *orm.Engine
}

func NewConsumerRunner(ctx context.Context, ormService *orm.Engine) *ConsumerRunner {
	return &ConsumerRunner{ctx: ctx, ormService: ormService}
}

func (r *ConsumerRunner) RunConsumerMany(consumer ConsumerMany, groupNameSuffix *string, prefetchCount int, blocking bool) {
	eventsConsumer := r.ormService.GetEventBroker().Consumer(consumer.GetQueueName(), consumer.GetGroupName(groupNameSuffix))
	eventsConsumer.SetLimit(10)
	eventsConsumer.SetErrorHandler(consumerErrorHandler)

	eventsConsumer.Consume(r.ctx, prefetchCount, blocking, func(events []orm.Event) {
		if err := consumer.Consume(events); err != nil {
			panic(err)
		}
	})
}

func (r *ConsumerRunner) RunConsumerOne(consumer ConsumerOne, groupNameSuffix *string, prefetchCount int, blocking bool) {
	eventsConsumer := r.ormService.GetEventBroker().Consumer(consumer.GetQueueName(), consumer.GetGroupName(groupNameSuffix))
	eventsConsumer.SetLimit(10)
	eventsConsumer.SetErrorHandler(consumerErrorHandler)

	eventsConsumer.Consume(r.ctx, prefetchCount, blocking, func(events []orm.Event) {
		for _, event := range events {
			if err := consumer.Consume(event); err != nil {
				panic(err)
			}
			event.Ack()
		}
	})
}

func (r *ConsumerRunner) RunConsumerOneByModulo(consumer ConsumerOneByModulo, groupNameSuffix *string, prefetchCount int, blocking bool) {
	for moduloID := 1; moduloID <= consumer.GetMaxModulo(); moduloID++ {
		currentModulo := moduloID

		go func() {
			eventsConsumer := r.ormService.GetEventBroker().Consumer(consumer.GetQueueName(currentModulo), consumer.GetGroupName(currentModulo, groupNameSuffix))
			eventsConsumer.SetLimit(10)
			eventsConsumer.SetErrorHandler(consumerErrorHandler)
			eventsConsumer.Consume(r.ctx, prefetchCount, blocking, func(events []orm.Event) {
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

func (r *ConsumerRunner) RunConsumerManyByModulo(consumer ConsumerManyByModulo, groupNameSuffix *string, prefetchCount int, blocking bool) {
	for moduloID := 1; moduloID <= consumer.GetMaxModulo(); moduloID++ {
		currentModulo := moduloID
		go func() {
			eventsConsumer := r.ormService.GetEventBroker().Consumer(consumer.GetQueueName(currentModulo), consumer.GetGroupName(currentModulo, groupNameSuffix))
			eventsConsumer.SetLimit(10)
			eventsConsumer.SetErrorHandler(consumerErrorHandler)
			eventsConsumer.Consume(r.ctx, prefetchCount, blocking, func(events []orm.Event) {
				if err := consumer.Consume(events); err != nil {
					panic(err)
				}
			})
		}()
	}
}
