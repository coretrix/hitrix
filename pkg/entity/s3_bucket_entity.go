package entity

import (
	"github.com/latolukasz/beeorm"
)

type S3BucketCounterEntity struct {
	beeorm.ORM `orm:"table=s3_buckets_counters"`
	ID         uint64
	Counter    uint64 `orm:"required"`
}
