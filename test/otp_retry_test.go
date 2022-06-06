package main

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/pkg/queue/consumers"
	"github.com/coretrix/hitrix/service"
	mockClockComponent "github.com/coretrix/hitrix/service/component/clock/mocks"
	"github.com/coretrix/hitrix/service/component/otp"
	"github.com/coretrix/hitrix/service/component/otp/mocks"
	mockClockRegistry "github.com/coretrix/hitrix/service/registry/mocks"
)

func TestOTPRetry(t *testing.T) {
	clock := &mockClockComponent.FakeSysClock{}
	clock.On("Now").Return(time.Unix(1, 0))

	createContextMyApp(t, "my-app", nil,
		[]*service.DefinitionGlobal{
			mockClockRegistry.ServiceProviderMockClock(clock),
		},
		nil)

	ormService := service.DI().OrmEngine()
	clockService := service.DI().Clock()

	code := "code"
	phoneNumber := "0123456789"
	phone := &otp.Phone{Number: phoneNumber}

	otpTrackerEntity := &entity.OTPTrackerEntity{
		Type:                entity.OTPTrackerTypeSMS,
		To:                  phoneNumber,
		Code:                code,
		GatewayName:         "gateway",
		GatewayPriority:     0,
		GatewaySendStatus:   entity.OTPTrackerGatewaySendStatusGatewayError,
		GatewaySendRequest:  "request",
		GatewaySendResponse: "response",
		RetryCount:          0,
		MaxRetriesReached:   false,
		SentAt:              clockService.Now(),
	}

	ormService.Flush(otpTrackerEntity)

	gateway := &mocks.FakeGateway{}
	gateway.On("SendOTP", phone, code).Return("request1", "response1", nil)

	dto := &otp.RetryDTO{
		Code:               code,
		Phone:              phone,
		OTPTrackerEntityID: otpTrackerEntity.ID,
		Gateway:            "g",
	}

	registry := map[string]otp.IOTPSMSGateway{
		"g": gateway,
	}

	consumers.RetryOTP(ormService, registry, dto, otpTrackerEntity, 10)

	otpTrackerEntity = &entity.OTPTrackerEntity{}
	ormService.LoadByID(1, otpTrackerEntity)

	assert.Equal(t, entity.OTPTrackerGatewaySendStatusSent, otpTrackerEntity.GatewaySendStatus)
	assert.Equal(t, "request1", otpTrackerEntity.GatewaySendRequest)
	assert.Equal(t, "response1", otpTrackerEntity.GatewaySendResponse)
	assert.Equal(t, 1, otpTrackerEntity.RetryCount)
	assert.Equal(t, false, otpTrackerEntity.MaxRetriesReached)

	clock.AssertExpectations(t)
	gateway.AssertExpectations(t)
}

func TestOTPWithMultipleRetry(t *testing.T) {
	clock := &mockClockComponent.FakeSysClock{}
	clock.On("Now").Return(time.Unix(1, 0))

	createContextMyApp(t, "my-app", nil,
		[]*service.DefinitionGlobal{
			mockClockRegistry.ServiceProviderMockClock(clock),
		},
		nil)

	ormService := service.DI().OrmEngine()
	clockService := service.DI().Clock()

	code := "code"
	phoneNumber := "0123456789"
	phone := &otp.Phone{Number: phoneNumber}

	otpTrackerEntity := &entity.OTPTrackerEntity{
		Type:                entity.OTPTrackerTypeSMS,
		To:                  phoneNumber,
		Code:                code,
		GatewayName:         "gateway",
		GatewayPriority:     0,
		GatewaySendStatus:   entity.OTPTrackerGatewaySendStatusGatewayError,
		GatewaySendRequest:  "request",
		GatewaySendResponse: "response",
		RetryCount:          0,
		MaxRetriesReached:   false,
		SentAt:              clockService.Now(),
	}

	ormService.Flush(otpTrackerEntity)

	gateway := &mocks.FakeGateway{}
	gateway.On("SendOTP", phone, code).Return("request1", "response1", errors.New("error")).Once()
	gateway.On("SendOTP", phone, code).Return("request2", "response2", errors.New("error")).Once()
	gateway.On("SendOTP", phone, code).Return("request3", "response3", errors.New("error")).Once()
	gateway.On("SendOTP", phone, code).Return("request4", "response4", nil).Once()

	dto := &otp.RetryDTO{
		Code:               code,
		Phone:              phone,
		OTPTrackerEntityID: otpTrackerEntity.ID,
		Gateway:            "g",
	}

	registry := map[string]otp.IOTPSMSGateway{
		"g": gateway,
	}

	consumers.RetryOTP(ormService, registry, dto, otpTrackerEntity, 10)

	otpTrackerEntity = &entity.OTPTrackerEntity{}
	ormService.LoadByID(1, otpTrackerEntity)

	assert.Equal(t, entity.OTPTrackerGatewaySendStatusSent, otpTrackerEntity.GatewaySendStatus)
	assert.Equal(t, "request4", otpTrackerEntity.GatewaySendRequest)
	assert.Equal(t, "response4", otpTrackerEntity.GatewaySendResponse)
	assert.Equal(t, 4, otpTrackerEntity.RetryCount)
	assert.Equal(t, false, otpTrackerEntity.MaxRetriesReached)

	clock.AssertExpectations(t)
	gateway.AssertExpectations(t)
}

func TestOTPRetryWithMaxReached(t *testing.T) {
	clock := &mockClockComponent.FakeSysClock{}
	clock.On("Now").Return(time.Unix(1, 0))

	createContextMyApp(t, "my-app", nil,
		[]*service.DefinitionGlobal{
			mockClockRegistry.ServiceProviderMockClock(clock),
		},
		nil)

	ormService := service.DI().OrmEngine()
	clockService := service.DI().Clock()

	code := "code"
	phoneNumber := "0123456789"
	phone := &otp.Phone{Number: phoneNumber}

	otpTrackerEntity := &entity.OTPTrackerEntity{
		Type:                entity.OTPTrackerTypeSMS,
		To:                  phoneNumber,
		Code:                code,
		GatewayName:         "gateway",
		GatewayPriority:     0,
		GatewaySendStatus:   entity.OTPTrackerGatewaySendStatusGatewayError,
		GatewaySendRequest:  "request",
		GatewaySendResponse: "response",
		RetryCount:          0,
		MaxRetriesReached:   false,
		SentAt:              clockService.Now(),
	}

	ormService.Flush(otpTrackerEntity)

	gateway := &mocks.FakeGateway{}
	gateway.On("SendOTP", phone, code).Return("request1", "response1", errors.New("error")).Once()
	gateway.On("SendOTP", phone, code).Return("request2", "response2", errors.New("error")).Once()
	gateway.On("SendOTP", phone, code).Return("request3", "response3", errors.New("error")).Once()

	dto := &otp.RetryDTO{
		Code:               code,
		Phone:              phone,
		OTPTrackerEntityID: otpTrackerEntity.ID,
		Gateway:            "g",
	}

	registry := map[string]otp.IOTPSMSGateway{
		"g": gateway,
	}

	consumers.RetryOTP(ormService, registry, dto, otpTrackerEntity, 3)

	otpTrackerEntity = &entity.OTPTrackerEntity{}
	ormService.LoadByID(1, otpTrackerEntity)

	assert.Equal(t, entity.OTPTrackerGatewaySendStatusGatewayError, otpTrackerEntity.GatewaySendStatus)
	assert.Equal(t, "request3", otpTrackerEntity.GatewaySendRequest)
	assert.Equal(t, "response3", otpTrackerEntity.GatewaySendResponse)
	assert.Equal(t, 3, otpTrackerEntity.RetryCount)
	assert.Equal(t, true, otpTrackerEntity.MaxRetriesReached)

	clock.AssertExpectations(t)
	gateway.AssertExpectations(t)
}
