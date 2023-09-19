package delayedqueues

import (
	"context"
	"github.com/coretrix/hitrix/pkg/dto/metrics"
	"github.com/coretrix/hitrix/service"
)

func Get(ctx context.Context) *delayedqueue.List {
	appService := service.DI().App()
	ormService := service.DI().OrmEngineForContext(ctx)

	result := &delayedqueue.List{Rows: make([]delayedqueue.Row, len(appService.RedisDelayedQueues))}
	for i, queue := range appService.RedisDelayedQueues {
		result.Rows[i].Queue = queue
		result.Rows[i].Total = ormService.GetRedis(appService.RedisPools.Persistent).ZCount(queue, "-inf", "+inf")
	}

	return result
}
