package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/pkg/entity"
)

type FakeServiceSetting struct {
	mock.Mock
}

func (s *FakeServiceSetting) Get(_ *datalayer.ORM, key string) (*entity.SettingsEntity, bool) {
	called := s.Called(key)

	return called.Get(0).(*entity.SettingsEntity), called.Bool(1)
}
func (s *FakeServiceSetting) GetString(_ *datalayer.ORM, key string) (string, bool) {
	called := s.Called(key)

	return called.String(0), called.Bool(1)
}
func (s *FakeServiceSetting) GetInt(_ *datalayer.ORM, key string) (int, bool) {
	called := s.Called(key)

	return called.Int(0), called.Bool(1)
}
func (s *FakeServiceSetting) GetUint(_ *datalayer.ORM, key string) (uint, bool) {
	called := s.Called(key)

	return called.Get(0).(uint), called.Bool(1)
}
func (s *FakeServiceSetting) GetInt64(_ *datalayer.ORM, key string) (int64, bool) {
	called := s.Called(key)

	return called.Get(0).(int64), called.Bool(1)
}
func (s *FakeServiceSetting) GetUint64(_ *datalayer.ORM, key string) (uint64, bool) {
	called := s.Called(key)

	return called.Get(0).(uint64), called.Bool(1)
}
func (s *FakeServiceSetting) GetFloat64(_ *datalayer.ORM, key string) (float64, bool) {
	called := s.Called(key)

	return called.Get(0).(float64), called.Bool(1)
}
func (s *FakeServiceSetting) GetBool(_ *datalayer.ORM, key string) (bool, bool) {
	called := s.Called(key)

	return called.Bool(0), called.Bool(1)
}
