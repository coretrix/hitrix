package apilogger

import (
	"time"

	"github.com/latolukasz/beeorm/v2"

	"github.com/coretrix/hitrix/datalayer"
)

type IAPILogger interface {
	LogStart(ormService *datalayer.DataLayer, logType string, request interface{})
	LogError(ormService *datalayer.DataLayer, message string, response interface{})
	LogSuccess(ormService *datalayer.DataLayer, response interface{})
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
