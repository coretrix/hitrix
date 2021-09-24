package otp

import (
	"crypto/md5"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/pkg/helper"
	"github.com/dongri/phonenumber"
	"github.com/latolukasz/beeorm"
)

type IOTP interface {
	SendSMS(ormService *beeorm.Engine, phone *Phone) (string, error)
	VerifySMS(ormService *beeorm.Engine, phone *Phone, code string) (bool, bool, error)
}

type OTP struct {
	GatewayPriority []IOTPSMSGateway
	GatewayName     map[string]IOTPSMSGateway
}

type Phone struct {
	Number  string
	ISO3166 phonenumber.ISO3166
}

func NewOTP(gateways ...IOTPSMSGateway) *OTP {
	otp := &OTP{
		GatewayPriority: make([]IOTPSMSGateway, len(gateways)),
		GatewayName:     map[string]IOTPSMSGateway{},
	}

	for i, gateway := range gateways {
		_, has := otp.GatewayName[gateway.GetName()]

		if has {
			panic("OTPProviders duplicated for name: " + gateway.GetName())
		}

		otp.GatewayPriority[i] = gateway
		otp.GatewayName[gateway.GetName()] = gateway
	}

	return otp
}

func (o *OTP) SendSMS(ormService *beeorm.Engine, phone *Phone) (string, error) {
	var code string
	var err error

	for priority, gateway := range o.GatewayPriority {
		code = gateway.GetCode()

		otpTrackerEntity := &entity.OTPTrackerEntity{
			Type:              entity.OTPTrackerTypeSMS,
			To:                phone.Number,
			Code:              code,
			GatewayName:       gateway.GetName(),
			GatewayPriority:   uint8(priority),
			GatewaySendStatus: entity.OTPTrackerGatewaySendStatusNew,
			SentAt:            time.Now(), //TODO ClockService
		}

		otpTrackerEntity.GatewaySendRequest, otpTrackerEntity.GatewaySendResponse, err = gateway.SendOTP(phone, code)

		if err != nil {
			otpTrackerEntity.GatewaySendStatus = entity.OTPTrackerGatewaySendStatusGatewayError
		} else {
			otpTrackerEntity.GatewaySendStatus = entity.OTPTrackerGatewaySendStatusSent
		}

		ormService.Flush(otpTrackerEntity)

		if err == nil {
			ormService.GetRedis().Set(o.getRedisKey(phone), otpTrackerEntity.ID, helper.Hour)
			break
		}
	}

	return code, err
}

func (o *OTP) VerifySMS(ormService *beeorm.Engine, phone *Phone, code string) (bool, bool, error) {
	otpTrackerEntity, err := o.getOTPTrackerEntity(ormService, phone)

	if err != nil {
		return false, false, err
	}

	gateway := o.GatewayName[otpTrackerEntity.GatewayName]

	var otpRequestValid bool
	var otpCodeValid bool

	otpTrackerEntity.GatewayVerifyRequest, otpTrackerEntity.GatewayVerifyResponse, otpRequestValid, otpCodeValid, err = gateway.VerifyOTP(phone, code)

	//TODO add error to tracker
	if err != nil {
		otpTrackerEntity.GatewayVerifyStatus = entity.OTPTrackerGatewayVerifyStatusGatewayError
	} else if !otpRequestValid {
		otpTrackerEntity.GatewayVerifyStatus = entity.OTPTrackerGatewayVerifyStatusExpired
	} else if !otpCodeValid {
		otpTrackerEntity.GatewayVerifyStatus = entity.OTPTrackerGatewayVerifyStatusInvalidCode
	} else {
		otpTrackerEntity.GatewayVerifyStatus = entity.OTPTrackerGatewayVerifyStatusSuccess
	}

	ormService.Flush(otpTrackerEntity)

	return otpRequestValid, otpCodeValid, err
}

func (o *OTP) getOTPTrackerEntity(ormService *beeorm.Engine, phone *Phone) (*entity.OTPTrackerEntity, error) {
	otpTrackerEntityIDString, has := ormService.GetRedis().Get(o.getRedisKey(phone))

	if !has {
		return nil, errors.New("OTP: redis key expired")
	}

	otpTrackerEntityID, err := strconv.ParseUint(otpTrackerEntityIDString, 10, 64)

	if err != nil {
		return nil, errors.New("OTP: " + err.Error())
	}

	otpTrackerEntity := &entity.OTPTrackerEntity{}

	found := ormService.LoadByID(otpTrackerEntityID, otpTrackerEntity)

	if !found {
		return nil, errors.New("OTP tracker not found")
	}

	return otpTrackerEntity, nil
}

func (o *OTP) getRedisKey(phone *Phone) string {
	// #nosec
	return fmt.Sprintf("%x", md5.Sum([]byte(phone.Number)))
}
