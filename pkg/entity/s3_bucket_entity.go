package entity

import (
	"github.com/latolukasz/orm"
)

type S3BucketCounterEntity struct {
	orm.ORM `orm:"table=s3_buckets_counters"`
	ID      uint64
	Counter uint64 `orm:"required"`
}
