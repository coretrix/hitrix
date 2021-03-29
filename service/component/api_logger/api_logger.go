package apilogger

import (
	"reflect"
	"time"

	"github.com/latolukasz/orm"
)

type LogEntity interface {
	orm.Entity
	SetID(value uint64)
	SetType(value string)
	SetStatus(value string)
	SetRequest(value interface{})
	SetResponse(value interface{})
	SetMessage(value string)
	SetCreatedAt(value time.Time)
}

type APILogger interface {
	LogStart(logType string, request interface{})
	LogError(message string, response interface{})
	LogSuccess(response interface{})
}

type DBLog struct {
	ormService *orm.Engine
	logEntity  LogEntity
	currentLog LogEntity
}

func NewAPILog(ormService *orm.Engine, entity LogEntity) APILogger {
	return &DBLog{ormService: ormService, logEntity: entity}
}

func (l *DBLog) LogStart(logType string, request interface{}) {
	var logEntity LogEntity

	if l.logEntity.GetID() == 0 {
		logEntity = l.logEntity
	} else {
		logEntity = reflect.New(reflect.ValueOf(l.logEntity).Elem().Type()).Interface().(LogEntity)
	}

	logEntity.SetType(logType)
	logEntity.SetRequest(request)
	logEntity.SetStatus("new")
	logEntity.SetCreatedAt(time.Now())

	l.ormService.Flush(logEntity)

	l.currentLog = logEntity
}

func (l *DBLog) LogError(message string, response interface{}) {
	if l.currentLog == nil {
		panic("log is not created")
	}

	currentLog := l.currentLog
	currentLog.SetMessage(message)
	currentLog.SetResponse(response)
	currentLog.SetStatus("failed")

	l.ormService.Flush(currentLog)
}

func (l *DBLog) LogSuccess(response interface{}) {
	if l.currentLog == nil {
		panic("log is not created")
	}

	currentLog := l.currentLog

	currentLog.SetStatus("completed")
	currentLog.SetResponse(response)

	l.ormService.Flush(currentLog)
}
