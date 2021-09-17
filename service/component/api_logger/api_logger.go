package apilogger

import (
	"time"

	"github.com/latolukasz/beeorm"
)

type IAPILogger interface {
	LogStart(ormService *beeorm.Engine, logType string, request interface{})
	LogError(ormService *beeorm.Engine, message string, response interface{})
	LogSuccess(ormService *beeorm.Engine, response interface{})
}

type ILogEntity interface {
	beeorm.Entity
	SetID(value uint64)
	SetType(value string)
	SetStatus(value string)
	SetRequest(value interface{})
	SetResponse(value interface{})
	SetMessage(value string)
	SetCreatedAt(value time.Time)
}
