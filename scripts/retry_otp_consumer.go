package scripts

import (
	"context"

	"github.com/coretrix/hitrix/pkg/queue"
	"github.com/coretrix/hitrix/pkg/queue/consumers"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
)

type RetryOTPConsumer struct {
}

func (script *RetryOTPConsumer) Run(ctx context.Context, _ app.IExit) {
	ormService := service.DI().OrmEngine()
	configService := service.DI().Config()

	maxRetries, ok := configService.Int("sms.max_retries")
	if !ok {
		panic("missing sms.max_retries")
	}

	queue.NewConsumerRunner(ctx).RunConsumerOne(consumers.NewOTPRetryConsumer(ormService, maxRetries), nil, 1)
}

func (script *RetryOTPConsumer) Infinity() bool {
	return true
}

func (script *RetryOTPConsumer) Unique() bool {
	return true
}

func (script *RetryOTPConsumer) Description() string {
	return "retry otp consumer"
}
