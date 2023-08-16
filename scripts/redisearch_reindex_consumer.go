package scripts

import (
	"context"

	"github.com/coretrix/hitrix/pkg/queue"
	"github.com/coretrix/hitrix/pkg/queue/consumers"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
)

type ReindexConsumerScript struct {
}

func (script *ReindexConsumerScript) Run(ctx context.Context, _ app.IExit) {
	queue.NewConsumerRunner(ctx).RunConsumerOne(consumers.NewReindexConsumer(service.DI().OrmEngine().RedisSearch), nil, 1)
}

func (script *ReindexConsumerScript) Infinity() bool {
	return true
}

func (script *ReindexConsumerScript) Unique() bool {
	return true
}

func (script *ReindexConsumerScript) Description() string {
	return "redisearch reindex consumer"
}
