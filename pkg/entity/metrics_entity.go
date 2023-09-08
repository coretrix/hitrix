package entity

import (
	"github.com/latolukasz/beeorm"
	"time"
)

type MetricsEntity struct {
	beeorm.ORM `orm:"table=metrics"`
	ID         uint64
	Metrics    string    `orm:"mediumblob"`
	CreatedAt  time.Time `orm:"time=true;"`
}
