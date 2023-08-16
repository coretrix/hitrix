package otp

import (
	//nolint //G501: Blocklisted import crypto/md5: weak cryptographic primitive
	"crypto/md5"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/dongri/phonenumber"

	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/pkg/helper"
	"github.com/coretrix/hitrix/pkg/queue/streams"
)

type IOTP interface {
	SendSMS(ormService *datalayer.DataLayer, phone *Phone) (string, error)
	VerifyOTP(ormService *datalayer.DataLayer, phone *Phone, code string) (bool, bool, error)
	Call(ormService *datalayer.DataLayer, phone *Phone, customMessage string) (string, error)
	GetGatewayRegistry() map[string]IOTPSMSGateway
}

type OTP struct {
	GatewayPriority         []IOTPSMSGateway
	GatewayName             map[string]IOTPSMSGateway
	GatewayPhonePrefixRegex map[*regexp.Regexp]IOTPSMSGateway
	RetryOTP                bool
}

type Phone struct {
	Number  string
	ISO3166 phonenumber.ISO3166
}

func NewOTP(retryOTP bool, gateways ...IOTPSMSGateway) *OTP {
	otp := &OTP{
		GatewayPriority:         make([]IOTPSMSGateway, 0),
		GatewayName:             map[string]IOTPSMSGateway{},
		GatewayPhonePrefixRegex: map[*regexp.Regexp]IOTPSMSGateway{},
		RetryOTP:                retryOTP,
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

func (o *OTP) SendSMS(ormService *datalayer.DataLayer, phone *Phone) (string, error) {
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
		} else if o.RetryOTP {
			ormService.GetEventBroker().Publish(streams.StreamMsgRetryOTP, &RetryDTO{
				Code:               code,
				Phone:              phone,
				OTPTrackerEntityID: otpTrackerEntity.ID,
				Gateway:            gateway.GetName(),
			}, nil)
		}
	}

	return code, err
}

func (o *OTP) Call(ormService *datalayer.DataLayer, phone *Phone, customMessage string) (string, error) {
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

func (o *OTP) VerifyOTP(ormService *datalayer.DataLayer, phone *Phone, code string) (bool, bool, error) {
	otpTrackerEntity, err := o.getOTPTrackerEntity(ormService, phone)

	if err != nil {
		return false, false, err
	}

	gateway := o.GatewayName[otpTrackerEntity.GatewayName]

	var otpRequestValid bool
	var otpCodeValid bool

	otpTrackerEntity.GatewayVerifyRequest, otpTrackerEntity.GatewayVerifyResponse, otpRequestValid, otpCodeValid, err =
		gateway.VerifyOTP(phone, code, otpTrackerEntity.Code)

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

func (o *OTP) GetGatewayRegistry() map[string]IOTPSMSGateway {
	return o.GatewayName
}

func (o *OTP) getOTPTrackerEntity(ormService *datalayer.DataLayer, phone *Phone) (*entity.OTPTrackerEntity, error) {
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

type RetryDTO struct {
	Code               string
	Phone              *Phone
	OTPTrackerEntityID uint64
	Gateway            string
}
