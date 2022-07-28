package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/coretrix/hitrix/service/component/otp"
)

type FakeGateway struct {
	mock.Mock
}

func (f *FakeGateway) GetName() string {
	return f.Called().String(0)
}

func (f *FakeGateway) GetCode() string {
	return f.Called().String(0)
}

func (f *FakeGateway) GetPhonePrefixes() []string {
	return f.Called().Get(0).([]string)
}

func (f *FakeGateway) SendOTP(phone *otp.Phone, code string) (string, string, error) {
	args := f.Called(phone, code)

	return args.String(0), args.String(1), args.Error(2)
}

func (f *FakeGateway) Call(phone *otp.Phone, code string, customMessage string) (string, string, error) {
	args := f.Called(phone, code, customMessage)

	return args.String(0), args.String(1), args.Error(2)
}

func (f *FakeGateway) VerifyOTP(phone *otp.Phone, code, generatedCode string) (string, string, bool, bool, error) {
	args := f.Called(phone, code, generatedCode)

	return args.String(0), args.String(1), args.Bool(2), args.Bool(3), args.Error(4)
}
