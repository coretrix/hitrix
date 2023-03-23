package entity

import (
	"time"

	"github.com/latolukasz/beeorm"
)

type RequestLoggerEntity struct {
	beeorm.ORM      `orm:"table=request_logger;redisCache;redisSearch=search_pool"`
	ID              uint64 `orm:"sortable"`
	URL             string `orm:"searchable"`
	UserID          uint64 `orm:"searchable"`
	AppName         string `orm:"required;searchable"`
	Request         []byte `orm:"mediumblob"`
	RequestText     string `orm:"length=max"`
	Response        []byte `orm:"mediumblob"`
	ResponseText    string `orm:"length=max"`
	Log             []byte `orm:"mediumblob"`
	Status          int
	RequestDuration int64     `orm:"sortable"`
	CreatedAt       time.Time `orm:"time=true;searchable"`
}
