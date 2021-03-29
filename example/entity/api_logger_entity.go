package entity

import (
	"time"

	"github.com/latolukasz/orm"
)

const (
	APILogStatusNew       = "new"
	APILogStatusCompleted = "completed"
	APILogStatusFailed    = "failed"

	APILogTypeApple = "Apple"
)

type apiLogType struct {
	orm.EnumModel
	Apple string
}

var APILogTypeAll = &apiLogType{
	Apple: APILogTypeApple,
}

type apiLogStatus struct {
	orm.EnumModel
	New       string
	Completed string
	Failed    string
}

var APILogStatusAll = &apiLogStatus{
	New:       APILogStatusNew,
	Completed: APILogStatusCompleted,
	Failed:    APILogStatusFailed,
}

type APILogEntity struct {
	orm.ORM   `orm:"table=api_log;redisCache"`
	ID        uint64
	Type      string `orm:"enum=entity.APILogTypeAll;required"`
	Status    string `orm:"enum=entity.APILogStatusAll;required"`
	Request   interface{}
	Response  interface{}
	Message   string
	CreatedAt time.Time `orm:"time=true"`
}

func (e *APILogEntity) SetID(value uint64) {
	e.ID = value
}

func (e *APILogEntity) SetType(value string) {
	e.Type = value
}

func (e *APILogEntity) SetStatus(value string) {
	e.Status = value
}

func (e *APILogEntity) SetRequest(value interface{}) {
	e.Request = value
}

func (e *APILogEntity) SetResponse(value interface{}) {
	e.Response = value
}

func (e *APILogEntity) SetMessage(value string) {
	e.Message = value
}

func (e *APILogEntity) SetCreatedAt(value time.Time) {
	e.CreatedAt = value
}
