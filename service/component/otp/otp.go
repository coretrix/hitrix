package otp

import (
	"crypto/md5" // nolint
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/dongri/phonenumber"
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/pkg/helper"
)

type IOTP interface {
	SendSMS(ormService *beeorm.Engine, phone *Phone) (string, error)
	VerifyOTP(ormService *beeorm.Engine, phone *Phone, code string) (bool, bool, error)
	Call(ormService *beeorm.Engine, phone *Phone, customMessage string) (string, error)
}

type OTP struct {
	GatewayPriority         []IOTPSMSGateway
	GatewayName             map[string]IOTPSMSGateway
	GatewayPhonePrefixRegex map[*regexp.Regexp]IOTPSMSGateway
}

type Phone struct {
	Number  string
	ISO3166 phonenumber.ISO3166
}

func NewOTP(gateways ...IOTPSMSGateway) *OTP {
	otp := &OTP{
		GatewayPriority:         make([]IOTPSMSGateway, 0),
		GatewayName:             map[string]IOTPSMSGateway{},
		GatewayPhonePrefixRegex: map[*regexp.Regexp]IOTPSMSGateway{},
	}

	for _, gateway := range gateways {
		_, has := otp.GatewayName[gateway.GetName()]

		if has {
			panic("OTPProviders duplicated for name: " + gateway.GetName())
		}

		otp.GatewayName[gateway.GetName()] = gateway

		phonePrefixes := gateway.GetPhonePrefixes()
		if phonePrefixes == nil {
			otp.GatewayPriority = append(otp.GatewayPriority, gateway)
		} else {
			regex := "^"
			for i, phonePrefix := range phonePrefixes {
				regex += "(\\" + phonePrefix + ")"
				if i < len(phonePrefixes)-1 {
					regex += "|"
				}
			}

			compiledRegex, err := regexp.Compile(regex)
			if err != nil {
				panic(err)
			}

			otp.GatewayPhonePrefixRegex[compiledRegex] = gateway
		}
	}

	return otp
}

func (o *OTP) SendSMS(ormService *beeorm.Engine, phone *Phone) (string, error) {
	var code string
	var err error

	gatewayPriority := make([]IOTPSMSGateway, 0)

	for regex, gateway := range o.GatewayPhonePrefixRegex {
		if regex.MatchString(phone.Number) {
			gatewayPriority = append(gatewayPriority, gateway)
		}
	}

	if len(gatewayPriority) == 0 {
		gatewayPriority = o.GatewayPriority
	}

	for priority, gateway := range gatewayPriority {
		code = gateway.GetCode()

		otpTrackerEntity := &entity.OTPTrackerEntity{
			Type:              entity.OTPTrackerTypeSMS,
			To:                phone.Number,
			Code:              code,
			GatewayName:       gateway.GetName(),
			GatewayPriority:   uint8(priority),
			GatewaySendStatus: entity.OTPTrackerGatewaySendStatusNew,
			SentAt:            time.Now(), // TODO ClockService
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

func (o *OTP) Call(ormService *beeorm.Engine, phone *Phone, customMessage string) (string, error) {
	var code string
	var err error

	for priority, gateway := range o.GatewayPriority {
		code = gateway.GetCode()

		otpTrackerEntity := &entity.OTPTrackerEntity{
			Type:              entity.OTPTrackerTypeCallout,
			To:                phone.Number,
			Code:              code,
			GatewayName:       gateway.GetName(),
			GatewayPriority:   uint8(priority),
			GatewaySendStatus: entity.OTPTrackerGatewaySendStatusNew,
			SentAt:            time.Now(), // TODO ClockService
		}

		otpTrackerEntity.GatewaySendRequest, otpTrackerEntity.GatewaySendResponse, err = gateway.Call(phone, code, customMessage)

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

func (o *OTP) VerifyOTP(ormService *beeorm.Engine, phone *Phone, code string) (bool, bool, error) {
	otpTrackerEntity, err := o.getOTPTrackerEntity(ormService, phone)

	if err != nil {
		return false, false, err
	}

	gateway := o.GatewayName[otpTrackerEntity.GatewayName]

	var otpRequestValid bool
	var otpCodeValid bool

	otpTrackerEntity.GatewayVerifyRequest, otpTrackerEntity.GatewayVerifyResponse, otpRequestValid, otpCodeValid, err = gateway.VerifyOTP(phone, code, otpTrackerEntity.Code)

	// TODO add error to tracker
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
