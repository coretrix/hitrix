package redissearch

import (
	"fmt"
	"strconv"

	"github.com/coretrix/trixorm"
)

// EntityIndex returns only ORM-managed indexes with a known numeric entity ID.
// Custom indexers need their own ID lookup contract and must not use this action.
func EntityIndex(engine *trixorm.Engine, pool, name string) *trixorm.RedisSearchIndex {
	registry := engine.GetRegistry()
	if _, found := registry.GetEntities()[name]; !found {
		return nil
	}

	schema := registry.GetTableSchema(name)
	search, found := schema.GetRedisSearch(engine)
	if !found || search.GetPoolConfig().GetCode() != pool {
		return nil
	}

	for _, index := range registry.GetRedisSearchIndices()[pool] {
		if index.Name == name && index.Indexer != nil && len(index.Prefixes) == 1 {
			return index
		}
	}

	return nil
}

// ReindexEntity uses the ORM's database indexer and field mappings, but queues
// Redis writes only for the requested document. It never drops an index or
// updates MySQL. The ORM indexer reads at most one normal indexing page.
func ReindexEntity(engine *trixorm.Engine, pool, name string, id uint64) (bool, error) {
	if id == 0 {
		return false, fmt.Errorf("entity ID must be positive")
	}

	index := EntityIndex(engine, pool, name)
	if index == nil {
		return false, fmt.Errorf("single-record reindexing is not supported for this index")
	}

	indexExists := false
	for _, indexName := range engine.GetRedisSearch(pool).ListIndices() {
		if indexName == name {
			indexExists = true
			break
		}
	}
	if !indexExists {
		return false, fmt.Errorf("index does not exist")
	}

	pusher := &entityReindexPusher{
		newPusher: func() trixorm.RedisSearchIndexPusher {
			return engine.NewRedisSearchIndexPusher(pool)
		},
		targetKey: index.Prefixes[0] + strconv.FormatUint(id, 10),
	}
	index.Indexer(engine, id-1, pusher)
	if pusher.selectedPusher == nil {
		return false, nil
	}

	pusher.Flush()
	// A direct SQL repair may also have left the entity's own cache stale.
	engine.ClearCacheByIDs(engine.GetRegistry().GetTableSchema(name).NewEntity(), id)

	return true, nil
}

// Embedding preserves the ORM's private field-mapping method. A fresh pusher
// per document isolates fields as well as pipeline commands for skipped rows.
type entityReindexPusher struct {
	trixorm.RedisSearchIndexPusher
	newPusher      func() trixorm.RedisSearchIndexPusher
	targetKey      string
	currentKey     string
	selectedPusher trixorm.RedisSearchIndexPusher
}

func (pusher *entityReindexPusher) NewDocument(key string) {
	pusher.RedisSearchIndexPusher = pusher.newPusher()
	pusher.RedisSearchIndexPusher.NewDocument(key)
	pusher.currentKey = key
}

func (pusher *entityReindexPusher) PushDocument() {
	if pusher.currentKey != pusher.targetKey {
		return
	}

	pusher.RedisSearchIndexPusher.PushDocument()
	pusher.selectedPusher = pusher.RedisSearchIndexPusher
}

func (pusher *entityReindexPusher) DeleteDocuments(_ ...string) {
	// The ORM entity indexer never deletes documents. Do not allow an indexer
	// to queue unrelated deletions through this restricted adapter.
}

func (pusher *entityReindexPusher) Flush() {
	if pusher.selectedPusher != nil {
		pusher.selectedPusher.Flush()
	}
}
