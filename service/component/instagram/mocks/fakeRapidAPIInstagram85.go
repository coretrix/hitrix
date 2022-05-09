package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/coretrix/hitrix/service/component/instagram"
)

type FakeRapidAPIInstagram85 struct {
	mock.Mock
}

func (f *FakeRapidAPIInstagram85) GetName() string {
	return f.Called().Get(0).(string)
}

func (f *FakeRapidAPIInstagram85) GetAccount(_ string) (*instagram.Account, error) {
	args := f.Called()
	return args.Get(0).(*instagram.Account), args.Error(1)
}

func (f *FakeRapidAPIInstagram85) GetFeed(_ int64, _ string) ([]*instagram.Post, string, error) {
	args := f.Called()
	return args.Get(0).([]*instagram.Post), args.String(1), args.Error(2)
}
