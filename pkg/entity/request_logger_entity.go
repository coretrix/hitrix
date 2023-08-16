package entity

import (
	"time"

	"github.com/latolukasz/beeorm/v2"
)

type RequestLoggerEntity struct {
	beeorm.ORM      `orm:"table=request_logger"`
	ID              uint64
	URL             string `orm:"length=500;index=URL"`
	UserID          uint64 `orm:"index=UserID"`
	AppName         string `orm:"required;index=AppName"`
	Request         []byte `orm:"mediumblob"`
	Response        []byte `orm:"mediumblob"`
	Log             []byte `orm:"mediumblob"`
	Status          int
	RequestDuration int64
	CreatedAt       time.Time `orm:"time=true;index=CreatedAt"`
}
