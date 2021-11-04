package featureflag

import (
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/latolukasz/beeorm"
)

type ServiceFeatureFlagInterface interface {
	IsActive(ormService *beeorm.Engine, name string) bool
	FailIfIsNotActive(ormService *beeorm.Engine, name string) error
	Activate(ormService *beeorm.Engine, name string) error
	DeActivate(ormService *beeorm.Engine, name string) error
	Create(ormService *beeorm.Engine, clockService clock.IClock, name string, isActive bool) error
	Delete(ormService *beeorm.Engine, name string) error
	GetAll(ormService *beeorm.Engine, pager *beeorm.Pager) []*entity.FeatureFlagEntity
}
