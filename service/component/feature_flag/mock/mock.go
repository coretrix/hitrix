package mock

import (
	"github.com/coretrix/trixorm"
	"github.com/stretchr/testify/mock"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/clock"
	featureflag "github.com/coretrix/hitrix/service/component/feature_flag"
)

type FakeServiceFeatureFlag struct {
	mock.Mock
}

func (s *FakeServiceFeatureFlag) IsActive(_ *trixorm.Engine, name string) bool {
	called := s.Called(name)

	return called.Bool(0)
}

func (s *FakeServiceFeatureFlag) FailIfIsNotActive(_ *trixorm.Engine, name string) error {
	called := s.Called(name)

	return called.Error(0)
}

func (s *FakeServiceFeatureFlag) Enable(_ *trixorm.Engine, name string) error {
	called := s.Called(name)

	return called.Error(0)
}

func (s *FakeServiceFeatureFlag) Disable(_ *trixorm.Engine, name string) error {
	called := s.Called(name)

	return called.Error(0)
}

func (s *FakeServiceFeatureFlag) GetAll(_ *trixorm.Engine, pager *trixorm.Pager) []*entity.FeatureFlagEntity {
	called := s.Called(pager)

	return called.Get(0).([]*entity.FeatureFlagEntity)
}

func (s *FakeServiceFeatureFlag) GetScriptsSingleInstance(ormService *trixorm.Engine) []app.IScript {
	called := s.Called(ormService)

	return called.Get(0).([]app.IScript)
}

func (s *FakeServiceFeatureFlag) GetScriptsMultiInstance(ormService *trixorm.Engine) []app.IScript {
	called := s.Called(ormService)

	return called.Get(0).([]app.IScript)
}

func (s *FakeServiceFeatureFlag) Register(featureFlags ...featureflag.IFeatureFlag) {
	s.Called(featureFlags)
}

func (s *FakeServiceFeatureFlag) Sync(ormService *trixorm.Engine, clockService clock.IClock) {
	s.Called(ormService, clockService)
}
