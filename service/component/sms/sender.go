package sms

import (
	"fmt"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/clock"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
)

type ISender interface {
	SendMessage(ormService *beeorm.Engine, message *Message) error
}

type Sender struct {
	ClockService       clock.IClock
	ErrorLoggerService errorlogger.ErrorLogger
	PrimaryProvider    IProvider
	SecondaryProvider  IProvider
	SandboxMode        bool
	TrackerEnabled     bool
}

func (s *Sender) SendMessage(ormService *beeorm.Engine, message *Message) error {
	var primaryProvider IProvider
	var secondaryProvider IProvider

	if message.Provider != nil {
		primaryProvider = message.Provider.Primary
		secondaryProvider = message.Provider.Secondary
	} else {
		primaryProvider = s.PrimaryProvider
		secondaryProvider = s.SecondaryProvider
	}

	if primaryProvider == nil {
		return fmt.Errorf("primary provider not defined")
	}

	smsTrackerEntity := entity.SmsTrackerEntity{}
	smsTrackerEntity.To = message.Number
	smsTrackerEntity.Type = entity.SMSTrackerTypeSMS
	smsTrackerEntity.Text = message.Text
	smsTrackerEntity.FromPrimaryGateway = primaryProvider.GetName()
	smsTrackerEntity.SentAt = s.ClockService.Now()

	trySecondaryProvider := false
	var status string
	var err error

	if !s.SandboxMode {
		status, err = primaryProvider.SendSMSMessage(message)
		if err != nil {
			trySecondaryProvider = true

			smsTrackerEntity.PrimaryGatewayError = err.Error()
			s.ErrorLoggerService.LogError(err)
		}
	} else {
		status = success
	}

	if trySecondaryProvider && secondaryProvider != nil {
		smsTrackerEntity.FromSecondaryGateway = secondaryProvider.GetName()

		status, err = secondaryProvider.SendSMSMessage(message)
		if err != nil {
			smsTrackerEntity.SecondaryGatewayError = err.Error()
			s.ErrorLoggerService.LogError(err)
		}
	}

	smsTrackerEntity.Status = status

	if s.TrackerEnabled {
		ormService.Flush(&smsTrackerEntity)
	}

	if status != success {
		return fmt.Errorf("sending sms failed")
	}

	return nil
}
