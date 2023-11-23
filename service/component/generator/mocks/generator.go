package mocks

import (
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

func (f *FakeGenerator) RandomPasswordGenerator(l int) string {
	return f.Called(l).String(0)
}
func (f *FakeGenerator) RandomPINCodeGenerator(l int) string {
	return f.Called(l).String(0)
}
