package entity

import (
	"time"

	"github.com/latolukasz/beeorm"
)

type ReverseGeocodingEntity struct {
	beeorm.ORM  `orm:"table=reverse_geocoding;redisCache"`
	ID          uint64
	Lat         float64 `orm:"required;unique=Lat_Lng_Language:1"`
	Lng         float64 `orm:"required;unique=Lat_Lng_Language:2"`
	Address     string
	Language    string `orm:"required;unique=Lat_Lng_Language:3"`
	Provider    string
	RawResponse interface{}
	ExpiresAt   time.Time `orm:"time=true;index=ExpiresAt"`
	CreatedAt   time.Time `orm:"time=true"`

	CachedQueryLatLngLanguage *beeorm.CachedQuery `queryOne:":Lat = ? AND :Lng = ? AND :Language = ?"`
}
