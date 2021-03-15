package entity

import (
	"github.com/latolukasz/orm"
)

type OSSBucketCounterEntity struct {
	orm.ORM `orm:"table=oss_buckets_counters"`
	ID      uint64
	Counter uint64 `orm:"required"`
}
