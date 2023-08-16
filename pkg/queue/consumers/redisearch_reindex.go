package consumers

import (
	redisearch "github.com/coretrix/beeorm-redisearch-plugin"
	"github.com/latolukasz/beeorm/v2"

	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/pkg/queue/streams"
)

type ReindexConsumer struct {
	redisearch *redisearch.RedisSearch
}

func NewReindexConsumer(redisearch *redisearch.RedisSearch) *ReindexConsumer {
	return &ReindexConsumer{redisearch: redisearch}
}

func (c *ReindexConsumer) GetQueueName() string {
	return redisearch.RedisSearchIndexerChannel
}

func (c *ReindexConsumer) GetGroupName(suffix *string) string {
	return streams.GetGroupName(c.GetQueueName(), suffix)
}

func (c *ReindexConsumer) Consume(_ *datalayer.DataLayer, event beeorm.Event) error {
	indexerEvent := &redisearch.IndexerEventRedisearch{}

	event.Unserialize(indexerEvent)

	c.redisearch.HandleRedisIndexerEvent(indexerEvent.Index)

	return nil
}
