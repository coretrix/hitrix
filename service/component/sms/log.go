package sms

import (
	"time"

	"github.com/latolukasz/beeorm"
)

type LogEntity interface {
	beeorm.Entity
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
	ormService *beeorm.Engine
	logEntity  LogEntity
}

func (db *DBLog) Do() {
	db.ormService.Flush(db.logEntity)
}

type Logger interface {
	Do()
}

func NewSmsLog(ormService *beeorm.Engine, entity LogEntity) Logger {
	return &DBLog{ormService: ormService, logEntity: entity}
}
