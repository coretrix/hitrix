package mock

import (
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/clock"
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

func (s *FakeServiceFeatureFlag) Activate(_ *beeorm.Engine, name string) error {
	called := s.Called(name)
	return called.Error(0)
}

func (s *FakeServiceFeatureFlag) DeActivate(_ *beeorm.Engine, name string) error {
	called := s.Called(name)
	return called.Error(0)
}

func (s *FakeServiceFeatureFlag) Create(_ *beeorm.Engine, _ clock.IClock, name string, isActive bool) error {
	called := s.Called(name, isActive)
	return called.Error(0)
}

func (s *FakeServiceFeatureFlag) Delete(_ *beeorm.Engine, name string) error {
	called := s.Called(name)
	return called.Error(0)
}

func (s *FakeServiceFeatureFlag) GetAll(_ *beeorm.Engine, pager *beeorm.Pager) []*entity.FeatureFlagEntity {
	called := s.Called(pager)
	return called.Get(0).([]*entity.FeatureFlagEntity)
}
