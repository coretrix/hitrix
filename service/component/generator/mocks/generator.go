package mocks

import (
	"github.com/AmirSoleimani/VoucherCodeGenerator/vcgen"
	"github.com/stretchr/testify/mock"
)

type FakeGenerator struct {
	mock.Mock
}

func (f *FakeGenerator) GenerateRandomRangeNumber(min, max int64) int64 {
	return int64(f.Called(min, max).Int(0))
}

func (f *FakeGenerator) GenerateSha256Hash(input string) string {
	return f.Called(input).String(0)
}

func (f *FakeGenerator) GenerateRandomCode(generator *vcgen.Generator) string {
	return f.Called(generator).String(0)
}
