package featureflag

import (
	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/clock"
)

type ServiceFeatureFlagInterface interface {
	IsActive(ormService *datalayer.DataLayer, name string) bool
	FailIfIsNotActive(ormService *datalayer.DataLayer, name string) error
	Enable(ormService *datalayer.DataLayer, name string) error
	Disable(ormService *datalayer.DataLayer, name string) error
	GetScriptsSingleInstance(ormService *datalayer.DataLayer) []app.IScript
	GetScriptsMultiInstance(ormService *datalayer.DataLayer) []app.IScript
	Register(featureFlags ...IFeatureFlag)
	Sync(ormService *datalayer.DataLayer, clockService clock.IClock)
}
