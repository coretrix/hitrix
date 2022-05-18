package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/coretrix/hitrix/service/component/instagram"
)

type FakeInstagram struct {
	mock.Mock
}

func (f *FakeInstagram) GetRandomProvider() instagram.IProvider {
	return f.Called().Get(0).(instagram.IProvider)
}

func (f *FakeInstagram) GetProvider(_ string) instagram.IProvider {
	return f.Called().Get(0).(instagram.IProvider)
}
