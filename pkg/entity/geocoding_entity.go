package entity

import (
	"time"

	"github.com/latolukasz/beeorm"
)

type GeocodingEntity struct {
	beeorm.ORM  `orm:"table=geocoding;redisCache"`
	ID          uint64
	Lat         float64
	Lng         float64
	Address     string `orm:"required;unique=Address_Language:1"`
	Language    string `orm:"required;enum=entity.LanguageValueAll;unique=Address_Language:2"`
	Provider    string
	RawResponse interface{}
	ExpiresAt   time.Time `orm:"time=true;index=ExpiresAt"`
	CreatedAt   time.Time `orm:"time=true"`

	CachedQueryAddressLanguage *beeorm.CachedQuery `queryOne:":Address = ? AND :Language = ?"`
}
