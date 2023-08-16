package datalayer

import (
	redisearch "github.com/coretrix/beeorm-redisearch-plugin"
	"github.com/latolukasz/beeorm/v2"
)

type DataLayer struct {
	beeorm.Engine
	*redisearch.RedisSearch
}

func (d *DataLayer) Clone() *DataLayer {
	return &DataLayer{
		Engine:      d.Engine.Clone(),
		RedisSearch: d.RedisSearch,
	}
}
