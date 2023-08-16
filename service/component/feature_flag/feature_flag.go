package featureflag

import (
	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/clock"
)

type ServiceFeatureFlagInterface interface {
	IsActive(ormService *datalayer.ORM, name string) bool
	FailIfIsNotActive(ormService *datalayer.ORM, name string) error
	Enable(ormService *datalayer.ORM, name string) error
	Disable(ormService *datalayer.ORM, name string) error
	GetScriptsSingleInstance(ormService *datalayer.ORM) []app.IScript
	GetScriptsMultiInstance(ormService *datalayer.ORM) []app.IScript
	Register(featureFlags ...IFeatureFlag)
	Sync(ormService *datalayer.ORM, clockService clock.IClock)
}
