package mocks

import (
	"github.com/stretchr/testify/mock"

	dynamiclink "github.com/coretrix/hitrix/service/component/dynamic_link"
)

type FakeDynamicLinksGenerator struct {
	mock.Mock
}

func (f *FakeDynamicLinksGenerator) GenerateDynamicLink(s string) (*dynamiclink.GenerateResponse, error) {
	args := f.Called(s)

	return args.Get(0).(*dynamiclink.GenerateResponse), args.Error(1)
}
