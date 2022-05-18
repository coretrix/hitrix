package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/coretrix/hitrix/service/component/instagram"
)

type FakeRapidAPIInstagram struct {
	mock.Mock
}

func (f *FakeRapidAPIInstagram) GetName() string {
	return f.Called().Get(0).(string)
}

func (f *FakeRapidAPIInstagram) GetAccount(_ string) (*instagram.Account, error) {
	args := f.Called()
	return args.Get(0).(*instagram.Account), args.Error(1)
}

func (f *FakeRapidAPIInstagram) GetFeed(_ int64, _ string) ([]*instagram.Post, string, error) {
	args := f.Called()
	return args.Get(0).([]*instagram.Post), args.String(1), args.Error(2)
}
