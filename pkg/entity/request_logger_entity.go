package entity

import (
	"time"

	"github.com/latolukasz/beeorm"
)

type RequestLoggerEntity struct {
	beeorm.ORM      `orm:"table=request_logger;redisCache;redisSearch=search_pool"`
	ID              uint64 `orm:"sortable"`
	URL             string `orm:"searchable;length=500"`
	UserID          uint64 `orm:"searchable"`
	AppName         string `orm:"required;searchable"`
	Request         []byte `orm:"mediumblob"`
	Response        []byte `orm:"mediumblob"`
	Log             []byte `orm:"mediumblob"`
	Status          int
	RequestDuration int64     `orm:"sortable"`
	CreatedAt       time.Time `orm:"time=true;searchable"`
}
