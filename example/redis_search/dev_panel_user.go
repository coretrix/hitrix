package redissearch

import (
	"strconv"

	redisearch "github.com/coretrix/beeorm-redisearch-plugin"
	"github.com/latolukasz/beeorm/v2"

	"github.com/coretrix/hitrix/example/entity"
)

const DevPanelUserIndex = "custom_index_dev_panel_user"

func GetDevPanelUserIndex(redisSearchPool string) *redisearch.RedisSearchIndex {
	index := &redisearch.RedisSearchIndex{}
	index.Name = DevPanelUserIndex
	index.RedisPool = redisSearchPool
	index.Prefixes = []string{DevPanelUserIndex + ":"}

	// document fields
	index.AddNumericField("ID", true, false)
	index.AddTextField("Email", 1, true, false, false)
	index.AddTextField("Password", 1, true, false, false)

	// force reindex func
	index.Indexer = devPanelUserIndexer

	return index
}

func SetDevPanelUserIndexFields(pusher redisearch.RedisSearchIndexPusher, devPanelUserEntities []*entity.DevPanelUserEntity) {
	deletedIDs := make([]string, 0)

	for _, devPanelUserEntity := range devPanelUserEntities {
		id := DevPanelUserIndex + ":" + strconv.FormatUint(devPanelUserEntity.ID, 10)

		if devPanelUserEntity.FakeDelete {
			deletedIDs = append(deletedIDs, id)

			continue
		}

		pusher.NewDocument(id)
		pusher.SetUint("ID", devPanelUserEntity.ID)
		pusher.SetString("Email", devPanelUserEntity.Email)
		pusher.SetString("Password", devPanelUserEntity.Password)
		pusher.PushDocument()
	}

	if len(deletedIDs) != 0 {
		pusher.DeleteDocuments(deletedIDs...)
	}
}

func devPanelUserIndexer(engine beeorm.Engine, lastID uint64, pusher redisearch.RedisSearchIndexPusher) (newID uint64, hasMore bool) {
	where := beeorm.NewWhere("ID > ? ORDER BY ID ASC", lastID)

	devPanelUserEntities := make([]*entity.DevPanelUserEntity, 0)
	engine.Search(where, beeorm.NewPager(1, 1000), &devPanelUserEntities)

	if len(devPanelUserEntities) == 0 {
		return lastID, false
	}

	SetDevPanelUserIndexFields(pusher, devPanelUserEntities)
	pusher.Flush()

	lastID = devPanelUserEntities[len(devPanelUserEntities)-1].ID

	return lastID, !(len(devPanelUserEntities) < 1000)
}
