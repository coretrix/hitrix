package mock

import (
	"github.com/latolukasz/beeorm"
	"github.com/stretchr/testify/mock"

	"github.com/coretrix/hitrix/pkg/entity"
)

type FakeServiceSetting struct {
	mock.Mock
}

func (s *FakeServiceSetting) Get(_ *beeorm.Engine, key string) (*entity.SettingsEntity, bool) {
	called := s.Called(key)

	return called.Get(0).(*entity.SettingsEntity), called.Bool(1)
}
func (s *FakeServiceSetting) GetString(_ *beeorm.Engine, key string) (string, bool) {
	called := s.Called(key)

	return called.String(0), called.Bool(1)
}
func (s *FakeServiceSetting) GetInt(_ *beeorm.Engine, key string) (int, bool) {
	called := s.Called(key)

	return called.Int(0), called.Bool(1)
}
func (s *FakeServiceSetting) GetUint(_ *beeorm.Engine, key string) (uint, bool) {
	called := s.Called(key)

	return called.Get(0).(uint), called.Bool(1)
}
func (s *FakeServiceSetting) GetInt64(_ *beeorm.Engine, key string) (int64, bool) {
	called := s.Called(key)

	return called.Get(0).(int64), called.Bool(1)
}
func (s *FakeServiceSetting) GetUint64(_ *beeorm.Engine, key string) (uint64, bool) {
	called := s.Called(key)

	return called.Get(0).(uint64), called.Bool(1)
}
func (s *FakeServiceSetting) GetFloat64(_ *beeorm.Engine, key string) (float64, bool) {
	called := s.Called(key)

	return called.Get(0).(float64), called.Bool(1)
}
func (s *FakeServiceSetting) GetBool(_ *beeorm.Engine, key string) (bool, bool) {
	called := s.Called(key)

	return called.Bool(0), called.Bool(1)
}
