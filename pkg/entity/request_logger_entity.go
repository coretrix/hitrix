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
	Content         []byte `orm:"mediumblob"`
	ResponseContent []byte `orm:"mediumblob"`
	Text            string `orm:"length=max"`
	ResponseText    string `orm:"length=max"`
	Log             []byte `orm:"mediumblob"`
	Status          int
	CreatedAt       time.Time `orm:"time=true;searchable"`
}
