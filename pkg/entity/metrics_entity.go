package entity

import (
	"time"

	"github.com/coretrix/trixorm"
)

type MetricsEntity struct {
	trixorm.ORM `orm:"table=metrics"`
	ID          uint64
	AppName     string
	Metrics     string    `orm:"mediumblob"`
	CreatedAt   time.Time `orm:"time=true;"`
}
