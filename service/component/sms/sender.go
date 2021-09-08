package sms

import (
	"fmt"

	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"

	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/latolukasz/beeorm"
)

type ISender interface {
	SendOTPSMS(ormService *beeorm.Engine, errorLoggerService errorlogger.ErrorLogger, otp *OTP) error
	SendOTPCallout(ormService *beeorm.Engine, errorLoggerService errorlogger.ErrorLogger, otp *OTP) error
	SendMessage(ormService *beeorm.Engine, errorLoggerService errorlogger.ErrorLogger, message *Message) error
	SendVerificationSMS(ormService *beeorm.Engine, errorLoggerService errorlogger.ErrorLogger, otp *OTP) error
	SendVerificationCallout(ormService *beeorm.Engine, errorLoggerService errorlogger.ErrorLogger, otp *OTP) error
	VerifyCode(ormService *beeorm.Engine, errorLoggerService errorlogger.ErrorLogger, otp *OTP) error
}

type Sender struct {
	Clock          clock.Clock
	GatewayFactory map[string]Gateway
	Logger         LogEntity
}

func (s *Sender) SendOTPSMS(ormService *beeorm.Engine, errorLoggerService errorlogger.ErrorLogger, otp *OTP) error {
	primaryProvider := otp.Provider.Primary
	secondaryProvider := otp.Provider.Secondary

	primaryGateway, ok := s.GatewayFactory[primaryProvider]
	if !ok {
		return fmt.Errorf("primary provider not supported")
	}

	smsTrackerEntity := s.Logger
	smsTrackerEntity.SetTo(otp.Number)
	smsTrackerEntity.SetType(entity.SMSTrackerTypeSMS)
	smsTrackerEntity.SetText(otp.OTP)
	smsTrackerEntity.SetFromPrimaryGateway(primaryProvider)
	smsTrackerEntity.SetSentAt(s.Clock.Now())

	trySecondaryProvider := false
	status, err := primaryGateway.SendOTPSMS(otp)
	if err != nil {
		trySecondaryProvider = true
		smsTrackerEntity.SetPrimaryGatewayError(err.Error())
		errorLoggerService.LogError(err)
	}

	if trySecondaryProvider {
		secondaryGateway, ok := s.GatewayFactory[secondaryProvider]
		if !ok {
			return fmt.Errorf("secondary provider not supported")
		}

		smsTrackerEntity.SetFromSecondaryGateway(secondaryProvider)

		status, err = secondaryGateway.SendOTPSMS(otp)
		if err != nil {
			smsTrackerEntity.SetSecondaryGatewayError(err.Error())
			errorLoggerService.LogError(err)
		}
	}

	smsTrackerEntity.SetStatus(status)
	logger := NewSmsLog(ormService, smsTrackerEntity)
	logger.Do()

	if status != success {
		return fmt.Errorf("sending sms failed")
	}

	return nil
}

func (s *Sender) SendOTPCallout(ormService *beeorm.Engine, errorLoggerService errorlogger.ErrorLogger, otp *OTP) error {
	primaryProvider := otp.Provider.Primary
	secondaryProvider := otp.Provider.Secondary

	primaryGateway, ok := s.GatewayFactory[primaryProvider]
	if !ok {
		return fmt.Errorf("primary provider not supported")
	}

	smsTrackerEntity := s.Logger
	smsTrackerEntity.SetTo(otp.Number)
	smsTrackerEntity.SetType(entity.SMSTrackerTypeCallout)
	smsTrackerEntity.SetText(otp.OTP)
	smsTrackerEntity.SetFromPrimaryGateway(primaryProvider)
	smsTrackerEntity.SetSentAt(s.Clock.Now())

	trySecondaryProvider := false
	status, err := primaryGateway.SendOTPCallout(otp)
	if err != nil {
		trySecondaryProvider = true
		smsTrackerEntity.SetPrimaryGatewayError(err.Error())
		errorLoggerService.LogError(err)
	}

	if trySecondaryProvider {
		secondaryGateway, ok := s.GatewayFactory[secondaryProvider]
		if !ok {
			return fmt.Errorf("secondary provider not supported")
		}

		smsTrackerEntity.SetFromSecondaryGateway(secondaryProvider)

		status, err = secondaryGateway.SendOTPCallout(otp)
		if err != nil {
			smsTrackerEntity.SetSecondaryGatewayError(err.Error())
			errorLoggerService.LogError(err)
		}
	}

	smsTrackerEntity.SetStatus(status)
	logger := NewSmsLog(ormService, smsTrackerEntity)
	logger.Do()

	if status != success {
		return fmt.Errorf("sending sms failed")
	}

	return nil
}

func (s *Sender) SendMessage(ormService *beeorm.Engine, errorLoggerService errorlogger.ErrorLogger, message *Message) error {
	primaryProvider := message.Provider.Primary
	secondaryProvider := message.Provider.Secondary

	primaryGateway, ok := s.GatewayFactory[primaryProvider]
	if !ok {
		return fmt.Errorf("primary provider not supported")
	}

	smsTrackerEntity := s.Logger
	smsTrackerEntity.SetTo(message.Number)
	smsTrackerEntity.SetType(entity.SMSTrackerTypeCallout)
	smsTrackerEntity.SetText(message.Text)
	smsTrackerEntity.SetFromPrimaryGateway(primaryProvider)
	smsTrackerEntity.SetSentAt(s.Clock.Now())

	trySecondaryProvider := false
	status, err := primaryGateway.SendSMSMessage(message)
	if err != nil {
		trySecondaryProvider = true
		smsTrackerEntity.SetPrimaryGatewayError(err.Error())
		errorLoggerService.LogError(err)
	}

	if trySecondaryProvider {
		secondaryGateway, ok := s.GatewayFactory[secondaryProvider]
		if !ok {
			return fmt.Errorf("secondary provider not supported")
		}

		smsTrackerEntity.SetFromSecondaryGateway(secondaryProvider)

		status, err = secondaryGateway.SendSMSMessage(message)
		if err != nil {
			smsTrackerEntity.SetSecondaryGatewayError(err.Error())
			errorLoggerService.LogError(err)
		}
	}

	smsTrackerEntity.SetStatus(status)
	logger := NewSmsLog(ormService, smsTrackerEntity)
	logger.Do()

	if status != success {
		return fmt.Errorf("sending sms failed")
	}

	return nil
}

func (s *Sender) SendVerificationSMS(ormService *beeorm.Engine, errorLoggerService errorlogger.ErrorLogger, otp *OTP) error {
	primaryProvider := otp.Provider.Primary
	secondaryProvider := otp.Provider.Secondary

	primaryGateway, ok := s.GatewayFactory[primaryProvider]
	if !ok {
		return fmt.Errorf("primary provider not supported")
	}

	smsTrackerEntity := s.Logger
	smsTrackerEntity.SetTo(otp.Number)
	smsTrackerEntity.SetType(entity.SMSTrackerTypeSMS)
	smsTrackerEntity.SetText(otp.OTP)
	smsTrackerEntity.SetFromPrimaryGateway(primaryProvider)
	smsTrackerEntity.SetSentAt(s.Clock.Now())

	trySecondaryProvider := false
	status, err := primaryGateway.SendVerificationSMS(otp)
	if err != nil {
		trySecondaryProvider = true
		smsTrackerEntity.SetPrimaryGatewayError(err.Error())
		errorLoggerService.LogError(err)
	}

	if trySecondaryProvider {
		secondaryGateway, ok := s.GatewayFactory[secondaryProvider]
		if !ok {
			return fmt.Errorf("secondary provider not supported")
		}

		smsTrackerEntity.SetFromSecondaryGateway(secondaryProvider)

		status, err = secondaryGateway.SendVerificationSMS(otp)
		if err != nil {
			smsTrackerEntity.SetSecondaryGatewayError(err.Error())
			errorLoggerService.LogError(err)
		}
	}

	smsTrackerEntity.SetStatus(status)
	logger := NewSmsLog(ormService, smsTrackerEntity)
	logger.Do()

	if status != success {
		return fmt.Errorf("sending sms failed")
	}

	return nil
}

func (s *Sender) SendVerificationCallout(ormService *beeorm.Engine, errorLoggerService errorlogger.ErrorLogger, otp *OTP) error {
	primaryProvider := otp.Provider.Primary
	secondaryProvider := otp.Provider.Secondary

	primaryGateway, ok := s.GatewayFactory[primaryProvider]
	if !ok {
		return fmt.Errorf("primary provider not supported")
	}

	smsTrackerEntity := s.Logger
	smsTrackerEntity.SetTo(otp.Number)
	smsTrackerEntity.SetType(entity.SMSTrackerTypeCallout)
	smsTrackerEntity.SetText(otp.OTP)
	smsTrackerEntity.SetFromPrimaryGateway(primaryProvider)
	smsTrackerEntity.SetSentAt(s.Clock.Now())

	trySecondaryProvider := false
	status, err := primaryGateway.SendVerificationCallout(otp)
	if err != nil {
		trySecondaryProvider = true
		smsTrackerEntity.SetPrimaryGatewayError(err.Error())
		errorLoggerService.LogError(err)
	}

	if trySecondaryProvider {
		secondaryGateway, ok := s.GatewayFactory[secondaryProvider]
		if !ok {
			return fmt.Errorf("secondary provider not supported")
		}

		smsTrackerEntity.SetFromSecondaryGateway(secondaryProvider)

		status, err = secondaryGateway.SendVerificationCallout(otp)
		if err != nil {
			smsTrackerEntity.SetSecondaryGatewayError(err.Error())
			errorLoggerService.LogError(err)
		}
	}

	smsTrackerEntity.SetStatus(status)
	logger := NewSmsLog(ormService, smsTrackerEntity)
	logger.Do()

	if status != success {
		return fmt.Errorf("sending sms failed")
	}

	return nil
}

func (s *Sender) VerifyCode(ormService *beeorm.Engine, errorLoggerService errorlogger.ErrorLogger, otp *OTP) error {
	primaryProvider := otp.Provider.Primary
	secondaryProvider := otp.Provider.Secondary

	primaryGateway, ok := s.GatewayFactory[primaryProvider]
	if !ok {
		return fmt.Errorf("primary provider not supported")
	}

	//TODO: create new log type for check the code
	smsTrackerEntity := s.Logger
	smsTrackerEntity.SetTo(otp.Number)
	smsTrackerEntity.SetType(entity.SMSTrackerTypeSMS)
	smsTrackerEntity.SetText(fmt.Sprintf(otp.Template, otp.OTP))
	smsTrackerEntity.SetFromPrimaryGateway(primaryProvider)
	smsTrackerEntity.SetSentAt(s.Clock.Now())

	trySecondaryProvider := false
	status, err := primaryGateway.VerifyCode(otp)
	if err != nil {
		trySecondaryProvider = true
		smsTrackerEntity.SetPrimaryGatewayError(err.Error())
		errorLoggerService.LogError(err)
	}

	if trySecondaryProvider {
		secondaryGateway, ok := s.GatewayFactory[secondaryProvider]
		if !ok {
			return fmt.Errorf("secondary provider not supported")
		}

		smsTrackerEntity.SetFromSecondaryGateway(secondaryProvider)

		status, err = secondaryGateway.VerifyCode(otp)
		if err != nil {
			smsTrackerEntity.SetSecondaryGatewayError(err.Error())
			errorLoggerService.LogError(err)
		}
	}

	smsTrackerEntity.SetStatus(status)
	logger := NewSmsLog(ormService, smsTrackerEntity)
	logger.Do()

	if status != success {
		return fmt.Errorf("sending sms failed")
	}

	return nil
}
