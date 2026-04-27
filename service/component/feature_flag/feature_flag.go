package featureflag

import (
	"github.com/coretrix/trixorm"

	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/clock"
)

type ServiceFeatureFlagInterface interface {
	IsActive(ormService *trixorm.Engine, name string) bool
	FailIfIsNotActive(ormService *trixorm.Engine, name string) error
	Enable(ormService *trixorm.Engine, name string) error
	Disable(ormService *trixorm.Engine, name string) error
	GetScriptsSingleInstance(ormService *trixorm.Engine) []app.IScript
	GetScriptsMultiInstance(ormService *trixorm.Engine) []app.IScript
	Register(featureFlags ...IFeatureFlag)
	Sync(ormService *trixorm.Engine, clockService clock.IClock)
}
