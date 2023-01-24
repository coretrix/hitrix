package entity

import (
	"time"

	"github.com/latolukasz/beeorm"
)

type GeocodingEntity struct {
	beeorm.ORM `orm:"table=geocoding;redisCache;redisSearch=search_pool"`
	ID         uint64    `orm:"sortable"`
	Lat        float64   `orm:"searchable"`
	Lng        float64   `orm:"searchable"`
	Address    string    `orm:"searchable"`
	CreatedAt  time.Time `orm:"time=true"`
}
