package consumers

import (
	"fmt"
	"log"

	redisearch "github.com/coretrix/beeorm-redisearch-plugin"
	"github.com/latolukasz/beeorm/v2"
	"github.com/latolukasz/beeorm/v2/plugins/crud_stream"

	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/example/redis"
	redissearch "github.com/coretrix/hitrix/example/redis_search"
	"github.com/coretrix/hitrix/pkg/queue/streams"
)

type DirtyConsumer struct {
}

func NewDirtyConsumer() *DirtyConsumer {
	return &DirtyConsumer{}
}

func (c *DirtyConsumer) GetQueueName() string {
	return crud_stream.ChannelName
}

func (c *DirtyConsumer) GetGroupName(suffix *string) string {
	return streams.GetGroupName(c.GetQueueName(), suffix)
}

func (c *DirtyConsumer) Consume(ormService *datalayer.DataLayer, events []beeorm.Event) error {
	entityNameIDsMapping := map[string][]uint64{}

	for _, dirtyEvent := range events {
		crudEvent := &crud_stream.CrudEvent{}
		dirtyEvent.Unserialize(crudEvent)

		ids, ok := entityNameIDsMapping[crudEvent.EntityName]
		if !ok {
			ids = make([]uint64, 0)
		}

		ids = append(ids, crudEvent.ID)
		entityNameIDsMapping[crudEvent.EntityName] = ids
	}

	for entityName, ids := range entityNameIDsMapping {
		processor, ok := dirtyEntityProcessorFactory[entityName]
		if !ok {
			panic(fmt.Errorf("dirty processor for entity %s not registered", entityName))
		}

		if err := processor(ormService, ids); err != nil {
			panic(err)
		}
	}

	return nil
}

var dirtyEntityProcessorFactory = map[string]func(*datalayer.DataLayer, []uint64) error{
	"entity.DevPanelUserEntity": processDevPanelUserEntity,
}

//nolint // info
func processDevPanelUserEntity(ormService *datalayer.DataLayer, ids []uint64) error {
	log.Printf("indexing %d dev panel users", len(ids))

	devPanelUserEntities := make([]*entity.DevPanelUserEntity, 0)
	ormService.LoadByIDs(ids, &devPanelUserEntities)

	pusher := redisearch.NewRedisSearchIndexPusher(ormService.Engine, redis.SearchPool)

	redissearch.SetDevPanelUserIndexFields(pusher, devPanelUserEntities)

	pusher.Flush()

	return nil
}
