package sms

import (
	"fmt"

	"github.com/coretrix/hitrix/example/entity"

	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/latolukasz/beeorm"
)

type ISender interface {
	SendOTPSMS(ormService *beeorm.Engine, otp *OTP) error
	SendOTPCallout(ormService *beeorm.Engine, otp *OTP) error
	SendMessage(ormService *beeorm.Engine, message *Message) error
}

type Sender struct {
	Clock          clock.Clock
	GatewayFactory map[string]Gateway
	Logger         LogEntity
}

func (s *Sender) SendOTPSMS(ormService *beeorm.Engine, otp *OTP) error {
	primaryProvider := otp.Provider.Primary
	secondaryProvider := otp.Provider.Secondary

	primaryGateway, ok := s.GatewayFactory[primaryProvider]
	if !ok {
		return fmt.Errorf("primary provider not supported")
	}

	otp.OTP = fmt.Sprintf(otp.Template, otp.OTP)

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
		}
	}

	smsTrackerEntity.SetStatus(status)
	logger := NewSmsLog(ormService, smsTrackerEntity)
	logger.Do()
	return nil
}

func (s *Sender) SendOTPCallout(ormService *beeorm.Engine, otp *OTP) error {
	primaryProvider := otp.Provider.Primary
	secondaryProvider := otp.Provider.Secondary

	primaryGateway, ok := s.GatewayFactory[primaryProvider]
	if !ok {
		return fmt.Errorf("primary provider not supported")
	}

	otp.OTP = fmt.Sprintf(otp.Template, otp.OTP)

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
		}
	}

	smsTrackerEntity.SetStatus(status)
	logger := NewSmsLog(ormService, smsTrackerEntity)
	logger.Do()
	return nil
}

func (s *Sender) SendMessage(ormService *beeorm.Engine, message *Message) error {
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
		}
	}

	smsTrackerEntity.SetStatus(status)
	logger := NewSmsLog(ormService, smsTrackerEntity)
	logger.Do()
	return nil
}
