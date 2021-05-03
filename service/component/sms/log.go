package sms

import (
	"time"

	"github.com/latolukasz/orm"
)

type LogEntity interface {
	orm.Entity
	SetStatus(string)
	SetTo(string)
	SetText(string)
	SetFromPrimaryGateway(string)
	SetFromSecondaryGateway(string)
	SetPrimaryGatewayError(string)
	SetSecondaryGatewayError(string)
	SetType(string)
	SetSentAt(time time.Time)
}

type DBLog struct {
	ormService *orm.Engine
	logEntity  LogEntity
}

func (db *DBLog) Do() {
	db.ormService.Flush(db.logEntity)
}

type Logger interface {
	Do()
}

func NewSmsLog(ormService *orm.Engine, entity LogEntity) Logger {
	return &DBLog{ormService: ormService, logEntity: entity}
}
