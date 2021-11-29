// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	beeorm "github.com/latolukasz/beeorm"
	mock "github.com/stretchr/testify/mock"

	otp "github.com/coretrix/hitrix/service/component/otp"
)

// OTPService is an autogenerated mock type for the IOTP type
type OTPService struct {
	mock.Mock
}

// Call provides a mock function with given fields: ormService, phone, customMessage
func (_m *OTPService) Call(ormService *beeorm.Engine, phone *otp.Phone, customMessage string) (string, error) {
	ret := _m.Called(ormService, phone, customMessage)

	var r0 string
	if rf, ok := ret.Get(0).(func(*beeorm.Engine, *otp.Phone, string) string); ok {
		r0 = rf(ormService, phone, customMessage)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*beeorm.Engine, *otp.Phone, string) error); ok {
		r1 = rf(ormService, phone, customMessage)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendSMS provides a mock function with given fields: ormService, phone
func (_m *OTPService) SendSMS(ormService *beeorm.Engine, phone *otp.Phone) (string, error) {
	ret := _m.Called(ormService, phone)

	var r0 string
	if rf, ok := ret.Get(0).(func(*beeorm.Engine, *otp.Phone) string); ok {
		r0 = rf(ormService, phone)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*beeorm.Engine, *otp.Phone) error); ok {
		r1 = rf(ormService, phone)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VerifyOTP provides a mock function with given fields: ormService, phone, code
func (_m *OTPService) VerifyOTP(ormService *beeorm.Engine, phone *otp.Phone, code string) (bool, bool, error) {
	ret := _m.Called(ormService, phone, code)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*beeorm.Engine, *otp.Phone, string) bool); ok {
		r0 = rf(ormService, phone, code)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(*beeorm.Engine, *otp.Phone, string) bool); ok {
		r1 = rf(ormService, phone, code)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*beeorm.Engine, *otp.Phone, string) error); ok {
		r2 = rf(ormService, phone, code)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}