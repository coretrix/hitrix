package entity

import (
	"time"

	"github.com/coretrix/trixorm"
)

type GeocodingReverseCacheEntity struct {
	trixorm.ORM              `orm:"table=geocoding_reverse_cache;localCache;redisCache"`
	ID                       uint64
	Lat                      float64 `orm:"decimal=8,5;required;unique=Lat_Lng_Language:1"`
	Lng                      float64 `orm:"decimal=8,5;required;unique=Lat_Lng_Language:2"`
	AdministrativeAreaLevel1 string  `orm:"required"`
	CityName                 string  `orm:"required"`
	Address                  string
	Language                 string `orm:"required;enum=entity.LanguageValueAll;unique=Lat_Lng_Language:3"`
	Provider                 string
	RawResponse              interface{}
	ExpiresAt                time.Time `orm:"time=true;index=ExpiresAt"`
	CreatedAt                time.Time `orm:"time=true"`

	CachedQueryLatLngLanguage *trixorm.CachedQuery `queryOne:":Lat = ? AND :Lng = ? AND :Language = ?"`
}
