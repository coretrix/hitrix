package delayedqueue

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/xorcare/pointer"

	"github.com/coretrix/hitrix/pkg/dto/delayedqueue"
	"github.com/coretrix/hitrix/service"
)

func Get(ctx context.Context) *delayedqueue.List {
	appService := service.DI().App()
	ormService := service.DI().OrmEngineForContext(ctx)
	now := service.DI().Clock().Now()

	result := &delayedqueue.List{Rows: make([]delayedqueue.Row, len(appService.RedisDelayedQueues))}
	for i, queue := range appService.RedisDelayedQueues {
		result.Rows[i].Queue = queue
		result.Rows[i].Total = ormService.GetRedis(appService.RedisPools.Persistent).ZCount(queue, "-inf", "+inf")
		values := ormService.GetRedis(appService.RedisPools.Persistent).ZRangeArgsWithScores(redis.ZRangeArgs{
			Key:     queue,
			Start:   -1,
			Stop:    now.AddDate(10, 0, 0).UnixMilli(),
			Offset:  1,
			Count:   1,
			ByScore: true})

		if len(values) > 0 {
			result.Rows[i].LatestItem = pointer.Uint64(uint64(values[0].Score))
		}
	}

	return result
}
