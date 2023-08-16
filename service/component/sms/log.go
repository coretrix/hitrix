package sms

import (
	"time"

	"github.com/latolukasz/beeorm/v2"

	"github.com/coretrix/hitrix/datalayer"
)

type LogEntity interface {
	beeorm.Entity
	SetStatus(string)
	SetTo(string)
	SetText(string)
	SetFromPrimaryProvider(string)
	SetFromSecondaryProvider(string)
	SetPrimaryProviderError(string)
	SetSecondaryProviderError(string)
	SetType(string)
	SetSentAt(time time.Time)
}

type DBLog struct {
	ormService *datalayer.ORM
	logEntity  LogEntity
}

func (db *DBLog) Do() {
	db.ormService.Flush(db.logEntity)
}

type Logger interface {
	Do()
}

func NewSmsLog(ormService *datalayer.ORM, entity LogEntity) Logger {
	return &DBLog{ormService: ormService, logEntity: entity}
}
