package entity

import (
	"time"

	"github.com/latolukasz/beeorm/v2"
)

type GeocodingCacheEntity struct {
	beeorm.ORM  `orm:"table=geocoding_cache;redisCache"`
	ID          uint64
	Lat         float64 `orm:"decimal=8,5"`
	Lng         float64 `orm:"decimal=8,5"`
	Address     string  `orm:"required"`
	AddressHash string  `orm:"length=32;required;unique=AddressHash_Language:1"`
	Language    string  `orm:"required;enum=entity.LanguageValueAll;unique=AddressHash_Language:2"`
	Provider    string
	RawResponse interface{}
	ExpiresAt   time.Time `orm:"time=true;index=ExpiresAt"`
	CreatedAt   time.Time `orm:"time=true"`

	CachedQueryAddressHashLanguage *beeorm.CachedQuery `queryOne:":AddressHash = ? AND :Language = ?"`
}
