package main

import (
	"testing"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/registry"
	"github.com/stretchr/testify/assert"
)

func TestFeatureFlag(t *testing.T) {
	createContextMyApp(t, "my-app", nil,
		[]*service.DefinitionGlobal{
			registry.ServiceProviderClock(),
			registry.ServiceProviderFeatureFlag(),
		},
		nil,
	)

	featureName := "bundle"
	featureFlagService := service.DI().FeatureFlag()

	ormService := service.DI().OrmEngine()
	clockService := service.DI().Clock()
	assert.Nil(t, featureFlagService.Create(ormService, clockService, featureName, true))

	assert.True(t, featureFlagService.IsActive(ormService, featureName))
	assert.Nil(t, featureFlagService.FailIfIsNotActive(ormService, featureName))

	assert.Nil(t, featureFlagService.DeActivate(ormService, featureName))
	assert.False(t, featureFlagService.IsActive(ormService, featureName))

	assert.Nil(t, featureFlagService.Activate(ormService, featureName))

	assert.Len(t, featureFlagService.GetAll(ormService, beeorm.NewPager(1, 10)), 1)

	assert.Nil(t, featureFlagService.Delete(ormService, featureName))
}
