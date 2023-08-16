package mocks

import (
	mock "github.com/stretchr/testify/mock"

	"github.com/coretrix/hitrix/datalayer"
	otp "github.com/coretrix/hitrix/service/component/otp"
)

type OTPService struct {
	mock.Mock
}

func (o *OTPService) Call(ormService *datalayer.ORM, phone *otp.Phone, customMessage string) (string, error) {
	ret := o.Called(ormService, phone, customMessage)

	var r0 string
	if rf, ok := ret.Get(0).(func(*datalayer.ORM, *otp.Phone, string) string); ok {
		r0 = rf(ormService, phone, customMessage)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*datalayer.ORM, *otp.Phone, string) error); ok {
		r1 = rf(ormService, phone, customMessage)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (o *OTPService) SendSMS(ormService *datalayer.ORM, phone *otp.Phone) (string, error) {
	ret := o.Called(ormService, phone)

	var r0 string
	if rf, ok := ret.Get(0).(func(*datalayer.ORM, *otp.Phone) string); ok {
		r0 = rf(ormService, phone)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*datalayer.ORM, *otp.Phone) error); ok {
		r1 = rf(ormService, phone)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (o *OTPService) VerifyOTP(ormService *datalayer.ORM, phone *otp.Phone, code string) (bool, bool, error) {
	ret := o.Called(ormService, phone, code)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*datalayer.ORM, *otp.Phone, string) bool); ok {
		r0 = rf(ormService, phone, code)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(*datalayer.ORM, *otp.Phone, string) bool); ok {
		r1 = rf(ormService, phone, code)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*datalayer.ORM, *otp.Phone, string) error); ok {
		r2 = rf(ormService, phone, code)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

func (o *OTPService) GetGatewayRegistry() map[string]otp.IOTPSMSGateway {
	return o.Called().Get(0).(map[string]otp.IOTPSMSGateway)
}
