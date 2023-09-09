package entity

import (
	"time"

	"github.com/latolukasz/beeorm"
)

type MetricsEntity struct {
	beeorm.ORM `orm:"table=metrics"`
	ID         uint64
	AppName    string
	Metrics    string    `orm:"mediumblob"`
	CreatedAt  time.Time `orm:"time=true;"`
}
