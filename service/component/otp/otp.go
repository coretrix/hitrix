package otp

import (
	//nolint //G501: Blocklisted import crypto/md5: weak cryptographic primitive
	"crypto/md5"
	"errors"
	"fmt"
	"math"
	mail2 "net/mail"
	"regexp"
	"strconv"
	"strings"

	"github.com/dongri/phonenumber"
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/pkg/helper"
	"github.com/coretrix/hitrix/pkg/queue/streams"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/generator"
	"github.com/coretrix/hitrix/service/component/mail"
)

type IOTP interface {
	Send(ormService *beeorm.Engine, send Send) (string, error)
	Verify(ormService *beeorm.Engine, verify Verify) (bool, bool, error)
	GetSMSGatewayRegistry() map[string]IOTPSMSGateway
}

type Phone struct {
	Number  string
	ISO3166 phonenumber.ISO3166
}

type Email struct {
	Email string
}

type Send struct {
	Phone            *Phone
	SMSCustomMessage string
	Email            *Email
	EmailConfig      *EmailConfig
}

type EmailConfig struct {
	From         string
	Subject      string
	TemplateName string
}

type Verify struct {
	Phone *Phone
	Email *Email
	Code  string
}

type Config struct {
	ClockService     clock.IClock
	SMSConfig        SMSConfig
	MailConfig       MailConfig
	CodeLength       int
	GeneratorService generator.IGenerator
}

type SMSConfig struct {
	SMSGateways []IOTPSMSGateway
	RetryOTP    bool
}

type MailConfig struct {
	Sender mail.ISender
}

type OTP struct {
	ClockService               clock.IClock
	GeneratorService           generator.IGenerator
	SMSGatewayPriority         []IOTPSMSGateway
	SMSGatewayName             map[string]IOTPSMSGateway
	SMSGatewayPhonePrefixRegex map[*regexp.Regexp]IOTPSMSGateway
	SMSRetryOTP                bool
	MailSender                 mail.ISender
	CodeLength                 int
}

func NewOTP(config Config) *OTP {
	otp := &OTP{
		SMSGatewayPriority:         make([]IOTPSMSGateway, 0),
		SMSGatewayName:             map[string]IOTPSMSGateway{},
		SMSGatewayPhonePrefixRegex: map[*regexp.Regexp]IOTPSMSGateway{},
		SMSRetryOTP:                config.SMSConfig.RetryOTP,
		MailSender:                 config.MailConfig.Sender,
		GeneratorService:           config.GeneratorService,
		ClockService:               config.ClockService,
	}

	for _, gateway := range config.SMSConfig.SMSGateways {
		_, has := otp.SMSGatewayName[gateway.GetName()]

		if has {
			panic("OTPProviders duplicated for name: " + gateway.GetName())
		}

		otp.SMSGatewayName[gateway.GetName()] = gateway

		phonePrefixes := gateway.GetPhonePrefixes()
		if phonePrefixes == nil {
			otp.SMSGatewayPriority = append(otp.SMSGatewayPriority, gateway)
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

			otp.SMSGatewayPhonePrefixRegex[compiledRegex] = gateway
		}
	}

	otp.CodeLength = 5

	if config.CodeLength > 0 {
		otp.CodeLength = config.CodeLength
	}

	return otp
}

func (o *OTP) Send(ormService *beeorm.Engine, send Send) (string, error) {
	if send.Phone != nil {
		return o.sendSMS(ormService, send)
	} else if send.Email != nil {
		return o.sendEmail(ormService, send)
	}

	panic("no send method selected")
}

func (o *OTP) Verify(ormService *beeorm.Engine, verify Verify) (bool, bool, error) {
	if verify.Phone != nil {
		return o.verifySMS(ormService, verify)
	} else if verify.Email != nil {
		return o.verifyEmail(ormService, verify)
	}

	panic("no verify method selected")
}

func (o *OTP) GetSMSGatewayRegistry() map[string]IOTPSMSGateway {
	return o.SMSGatewayName
}

func (o *OTP) sendSMS(ormService *beeorm.Engine, send Send) (string, error) {
	//validate phone
	phone := send.Phone
	var code string
	var err error

	gatewayPriority := make([]IOTPSMSGateway, 0)

	for regex, gateway := range o.SMSGatewayPhonePrefixRegex {
		if regex.MatchString(phone.Number) {
			gatewayPriority = append(gatewayPriority, gateway)
		}
	}

	if len(gatewayPriority) == 0 {
		gatewayPriority = o.SMSGatewayPriority
	}

	for priority, gateway := range gatewayPriority {
		code = gateway.GetCode()
		if send.SMSCustomMessage != "" {
			code = strings.Replace(send.SMSCustomMessage, "_CODE_", code, -1)
		}

		otpTrackerEntity := &entity.OTPTrackerEntity{
			Type:              entity.OTPTrackerTypeSMS,
			To:                phone.Number,
			Code:              code,
			GatewayName:       gateway.GetName(),
			GatewayPriority:   uint8(priority),
			GatewaySendStatus: entity.OTPTrackerGatewaySendStatusNew,
			SentAt:            o.ClockService.Now(),
		}

		otpTrackerEntity.GatewaySendRequest, otpTrackerEntity.GatewaySendResponse, err = gateway.SendOTP(phone, code)

		if err != nil {
			otpTrackerEntity.GatewaySendStatus = entity.OTPTrackerGatewaySendStatusGatewayError
		} else {
			otpTrackerEntity.GatewaySendStatus = entity.OTPTrackerGatewaySendStatusSent
		}

		ormService.Flush(otpTrackerEntity)

		if err == nil {
			ormService.GetRedis().Set(o.getRedisKey(phone.Number), otpTrackerEntity.ID, helper.Hour)

			break
		} else if o.SMSRetryOTP {
			ormService.GetEventBroker().Publish(streams.StreamMsgRetryOTP, &RetryDTO{
				Code:               code,
				Phone:              phone,
				OTPTrackerEntityID: otpTrackerEntity.ID,
				Gateway:            gateway.GetName(),
			})
		}
	}

	return code, err
}

func (o *OTP) sendEmail(ormService *beeorm.Engine, send Send) (string, error) {
	if o.MailSender == nil {
		panic("mail sender not defined")
	}

	email := send.Email.Email

	_, err := mail2.ParseAddress(email)
	if err != nil {
		return "", errors.New("mail address not valid")
	}

	code := o.getCode()

	otpTrackerEntity := &entity.OTPTrackerEntity{
		Type:              entity.OTPTrackerTypeSMS,
		To:                email,
		Code:              code,
		GatewaySendStatus: entity.OTPTrackerGatewaySendStatusNew,
		SentAt:            o.ClockService.Now(),
	}

	err = o.MailSender.SendTemplate(ormService, &mail.Message{
		From:         send.EmailConfig.From,
		To:           email,
		Subject:      send.EmailConfig.Subject,
		TemplateName: send.EmailConfig.TemplateName,
		TemplateData: map[string]interface{}{"code": code},
	})

	if err != nil {
		otpTrackerEntity.GatewaySendStatus = entity.OTPTrackerGatewaySendStatusGatewayError
	} else {
		otpTrackerEntity.GatewaySendStatus = entity.OTPTrackerGatewaySendStatusSent
	}

	ormService.Flush(otpTrackerEntity)

	ormService.GetRedis().Set(o.getRedisKey(email), otpTrackerEntity.ID, helper.Hour)

	return code, err
}

func (o *OTP) verifySMS(ormService *beeorm.Engine, verify Verify) (bool, bool, error) {
	otpTrackerEntity, err := o.getOTPTrackerEntity(ormService, verify.Phone.Number)
	if err != nil {
		return false, false, err
	}

	gateway := o.SMSGatewayName[otpTrackerEntity.GatewayName]

	var otpRequestValid bool
	var otpCodeValid bool

	otpTrackerEntity.GatewayVerifyRequest, otpTrackerEntity.GatewayVerifyResponse, otpRequestValid, otpCodeValid, err =
		gateway.VerifyOTP(verify.Phone, verify.Code, otpTrackerEntity.Code)

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

func (o *OTP) verifyEmail(ormService *beeorm.Engine, verify Verify) (bool, bool, error) {
	otpTrackerEntity, err := o.getOTPTrackerEntity(ormService, verify.Email.Email)
	if err != nil {
		return false, false, err
	}

	if otpTrackerEntity.Code == verify.Code {
		otpTrackerEntity.GatewayVerifyStatus = entity.OTPTrackerGatewayVerifyStatusSuccess
	} else {
		otpTrackerEntity.GatewayVerifyStatus = entity.OTPTrackerGatewayVerifyStatusInvalidCode
	}

	ormService.Flush(otpTrackerEntity)

	return true, otpTrackerEntity.Code == verify.Code, nil
}

func (o *OTP) getOTPTrackerEntity(ormService *beeorm.Engine, verifyKey string) (*entity.OTPTrackerEntity, error) {
	otpTrackerEntityIDString, has := ormService.GetRedis().Get(o.getRedisKey(verifyKey))

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

func (o *OTP) getRedisKey(verifyKey string) string {
	// #nosec
	return fmt.Sprintf("otp_%x", md5.Sum([]byte(verifyKey)))
}

func (o *OTP) getCode() string {
	return strconv.FormatInt(
		o.GeneratorService.GenerateRandomRangeNumber(
			int64(math.Pow(10, float64(o.CodeLength-1))),
			int64(math.Pow(10, float64(o.CodeLength)))-1,
		),
		10)
}

type RetryDTO struct {
	Code               string
	Phone              *Phone
	OTPTrackerEntityID uint64
	Gateway            string
}
