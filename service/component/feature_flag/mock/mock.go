package mock

import (
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/clock"
	featureflag "github.com/coretrix/hitrix/service/component/feature_flag"
	"github.com/latolukasz/beeorm"
	"github.com/stretchr/testify/mock"
)

type FakeServiceFeatureFlag struct {
	mock.Mock
}

func (s *FakeServiceFeatureFlag) IsActive(_ *beeorm.Engine, name string) bool {
	called := s.Called(name)
	return called.Bool(1)
}

func (s *FakeServiceFeatureFlag) FailIfIsNotActive(_ *beeorm.Engine, name string) error {
	called := s.Called(name)
	return called.Error(0)
}

func (s *FakeServiceFeatureFlag) Enable(_ *beeorm.Engine, name string) error {
	called := s.Called(name)
	return called.Error(0)
}

func (s *FakeServiceFeatureFlag) Disable(_ *beeorm.Engine, name string) error {
	called := s.Called(name)
	return called.Error(0)
}

func (s *FakeServiceFeatureFlag) GetAll(_ *beeorm.Engine, pager *beeorm.Pager) []*entity.FeatureFlagEntity {
	called := s.Called(pager)
	return called.Get(0).([]*entity.FeatureFlagEntity)
}

func (s *FakeServiceFeatureFlag) GetScriptsSingleInstance(ormService *beeorm.Engine) []app.IScript {
	called := s.Called(ormService)
	return called.Get(0).([]app.IScript)
}

func (s *FakeServiceFeatureFlag) GetScriptsMultiInstance(ormService *beeorm.Engine) []app.IScript {
	called := s.Called(ormService)
	return called.Get(0).([]app.IScript)
}

func (s *FakeServiceFeatureFlag) Register(featureFlags ...featureflag.IFeatureFlag) {
	s.Called(featureFlags)
}

func (s *FakeServiceFeatureFlag) Sync(ormService *beeorm.Engine, clockService clock.IClock) {
	s.Called(ormService, clockService)
}
