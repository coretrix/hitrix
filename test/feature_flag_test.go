package main

import (
	"testing"

	"github.com/coretrix/hitrix/service/component/app"

	featureflag "github.com/coretrix/hitrix/service/component/feature_flag"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/registry"
	"github.com/stretchr/testify/assert"
)

func initFeatures(flagInterface featureflag.ServiceFeatureFlagInterface) {
	flagInterface.Register(
		&ProductCollection{},
	)
}

const productCollectionFeature = "product_collection"

type ProductCollection struct {
}

func (f *ProductCollection) GetName() string {
	return productCollectionFeature
}

func (f *ProductCollection) ScriptsSingleInstance() []app.IScript {
	return nil
}

func (f *ProductCollection) ScriptsMultiInstance() []app.IScript {
	return nil
}

func TestFeatureFlag(t *testing.T) {
	createContextMyApp(t, "my-app", nil,
		[]*service.DefinitionGlobal{
			registry.ServiceProviderErrorLogger(),
			registry.ServiceProviderClock(),
			registry.ServiceProviderFeatureFlag(initFeatures),
		},
		nil,
	)

	featureFlagService := service.DI().FeatureFlag()

	ormService := service.DI().OrmEngine()
	clockService := service.DI().Clock()
	assert.Nil(t, featureFlagService.Create(ormService, clockService, productCollectionFeature, true))

	assert.True(t, featureFlagService.IsActive(ormService, productCollectionFeature))
	assert.Nil(t, featureFlagService.FailIfIsNotActive(ormService, productCollectionFeature))

	assert.Nil(t, featureFlagService.DeActivate(ormService, productCollectionFeature))
	assert.False(t, featureFlagService.IsActive(ormService, productCollectionFeature))

	assert.Nil(t, featureFlagService.Activate(ormService, productCollectionFeature))

	assert.Len(t, featureFlagService.GetAll(ormService, beeorm.NewPager(1, 10)), 1)

	assert.Nil(t, featureFlagService.Delete(ormService, productCollectionFeature))
}
