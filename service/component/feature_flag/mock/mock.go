package mock

import (
	"github.com/latolukasz/beeorm/v2"
	"github.com/stretchr/testify/mock"

	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/clock"
	featureflag "github.com/coretrix/hitrix/service/component/feature_flag"
)

type FakeServiceFeatureFlag struct {
	mock.Mock
}

func (s *FakeServiceFeatureFlag) IsActive(_ *datalayer.DataLayer, name string) bool {
	called := s.Called(name)

	return called.Bool(0)
}

func (s *FakeServiceFeatureFlag) FailIfIsNotActive(_ *datalayer.DataLayer, name string) error {
	called := s.Called(name)

	return called.Error(0)
}

func (s *FakeServiceFeatureFlag) Enable(_ *datalayer.DataLayer, name string) error {
	called := s.Called(name)

	return called.Error(0)
}

func (s *FakeServiceFeatureFlag) Disable(_ *datalayer.DataLayer, name string) error {
	called := s.Called(name)

	return called.Error(0)
}

func (s *FakeServiceFeatureFlag) GetAll(_ *datalayer.DataLayer, pager *beeorm.Pager) []*entity.FeatureFlagEntity {
	called := s.Called(pager)

	return called.Get(0).([]*entity.FeatureFlagEntity)
}

func (s *FakeServiceFeatureFlag) GetScriptsSingleInstance(ormService *datalayer.DataLayer) []app.IScript {
	called := s.Called(ormService)

	return called.Get(0).([]app.IScript)
}

func (s *FakeServiceFeatureFlag) GetScriptsMultiInstance(ormService *datalayer.DataLayer) []app.IScript {
	called := s.Called(ormService)

	return called.Get(0).([]app.IScript)
}

func (s *FakeServiceFeatureFlag) Register(featureFlags ...featureflag.IFeatureFlag) {
	s.Called(featureFlags)
}

func (s *FakeServiceFeatureFlag) Sync(ormService *datalayer.DataLayer, clockService clock.IClock) {
	s.Called(ormService, clockService)
}
