package entity

import (
	"time"

	"github.com/coretrix/trixorm"
)

type FeatureFlagEntity struct {
	trixorm.ORM `orm:"table=feature_flags;localCache"`
	ID          uint64
	Name        string     `orm:"length=100;required;unique=Name"`
	Registered  bool       `orm:"index=Registered_Enabled:1"`
	Enabled     bool       `orm:"index=Registered_Enabled:2"`
	UpdatedAt   *time.Time `orm:"time=true"`
	CreatedAt   time.Time  `orm:"time=true"`

	CachedQueryAll               *trixorm.CachedQuery `query:"1 ORDER BY ID"`
	CachedQueryName              *trixorm.CachedQuery `queryOne:":Name = ?"`
	CachedQueryRegisteredEnabled *trixorm.CachedQuery `query:":Registered = ? AND :Enabled = ?"`
}
