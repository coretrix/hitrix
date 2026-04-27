package apilogger

import (
	"time"

	"github.com/coretrix/trixorm"
)

type IAPILogger interface {
	LogStart(ormService *trixorm.Engine, logType string, request interface{})
	LogError(ormService *trixorm.Engine, message string, response interface{})
	LogSuccess(ormService *trixorm.Engine, response interface{})
}

type ILogEntity interface {
	trixorm.Entity
	SetID(value uint64)
	SetType(value string)
	SetStatus(value string)
	SetRequest(value interface{})
	SetResponse(value interface{})
	SetMessage(value string)
	SetCreatedAt(value time.Time)
}
