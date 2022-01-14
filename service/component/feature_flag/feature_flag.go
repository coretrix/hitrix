package featureflag

import (
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/clock"
)

type ServiceFeatureFlagInterface interface {
	IsActive(ormService *beeorm.Engine, name string) bool
	FailIfIsNotActive(ormService *beeorm.Engine, name string) error
	Enable(ormService *beeorm.Engine, name string) error
	Disable(ormService *beeorm.Engine, name string) error
	GetScriptsSingleInstance(ormService *beeorm.Engine) []app.IScript
	GetScriptsMultiInstance(ormService *beeorm.Engine) []app.IScript
	Register(featureFlags ...IFeatureFlag)
	Sync(ormService *beeorm.Engine, clockService clock.IClock)
}
