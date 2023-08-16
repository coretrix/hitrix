package datalayer

import (
	redisearch "github.com/coretrix/beeorm-redisearch-plugin"
	"github.com/latolukasz/beeorm/v2"
)

type ORM struct {
	beeorm.Engine
	*redisearch.RedisSearchEngine
}

func (d *ORM) Clone() *ORM {
	return &ORM{
		Engine:            d.Engine.Clone(),
		RedisSearchEngine: d.RedisSearchEngine,
	}
}
