package apilogger

import (
	"reflect"
	"time"

	"github.com/coretrix/hitrix/datalayer"
)

type mysqlDBLog struct {
	logEntity  ILogEntity
	currentLog ILogEntity
}

func NewMysqlAPILogger(entity ILogEntity) IAPILogger {
	return &mysqlDBLog{logEntity: entity}
}

func (l *mysqlDBLog) LogStart(ormService *datalayer.DataLayer, logType string, request interface{}) {
	var logEntity ILogEntity

	if l.logEntity.GetID() == 0 {
		logEntity = l.logEntity
	} else {
		logEntity = reflect.New(reflect.ValueOf(l.logEntity).Elem().Type()).Interface().(ILogEntity)
	}

	logEntity.SetType(logType)
	logEntity.SetRequest(request)
	logEntity.SetStatus("new")
	logEntity.SetCreatedAt(time.Now())

	ormService.Flush(logEntity)

	l.currentLog = logEntity
}

func (l *mysqlDBLog) LogError(ormService *datalayer.DataLayer, message string, response interface{}) {
	if l.currentLog == nil {
		panic("log is not created")
	}

	currentLog := l.currentLog
	currentLog.SetMessage(message)
	currentLog.SetResponse(response)
	currentLog.SetStatus("failed")

	ormService.Flush(currentLog)
}

func (l *mysqlDBLog) LogSuccess(ormService *datalayer.DataLayer, response interface{}) {
	if l.currentLog == nil {
		panic("log is not created")
	}

	currentLog := l.currentLog

	currentLog.SetStatus("completed")
	currentLog.SetResponse(response)

	ormService.Flush(currentLog)
}
