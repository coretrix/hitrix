package entity

import (
	"time"

	"github.com/latolukasz/beeorm/v2"
)

type FeatureFlagEntity struct {
	beeorm.ORM `orm:"table=feature_flags;redisCache;redisSearch=search_pool"`
	ID         uint64
	Name       string     `orm:"length=100;required;unique=Name;searchable"`
	Registered bool       `orm:"searchable"`
	Enabled    bool       `orm:"searchable"`
	UpdatedAt  *time.Time `orm:"time=true"`
	CreatedAt  time.Time  `orm:"time=true"`
}
