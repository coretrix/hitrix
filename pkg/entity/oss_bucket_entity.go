package entity

import (
	"github.com/coretrix/trixorm"
)

type OSSBucketCounterEntity struct {
	trixorm.ORM `orm:"table=oss_buckets_counters"`
	ID          uint64
	Counter     uint64 `orm:"required"`
}
