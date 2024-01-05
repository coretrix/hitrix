package sms

import (
	"fmt"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
)

type ISender interface {
	SendMessage(ormService *beeorm.Engine, message *Message) error
}

type Sender struct {
	ConfigService      config.IConfig
	ClockService       clock.IClock
	ErrorLoggerService errorlogger.ErrorLogger
	PrimaryProvider    IProvider
	SecondaryProvider  IProvider
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
		return fmt.Errorf("primary provider not supported")
	}

	smsTrackerEntity := &entity.SmsTrackerEntity{}
	smsTrackerEntity.SetTo(message.Number)
	smsTrackerEntity.SetType(entity.SMSTrackerTypeSMS)
	smsTrackerEntity.SetText(message.Text)
	smsTrackerEntity.SetFromPrimaryProvider(primaryProvider.GetName())
	smsTrackerEntity.SetSentAt(s.ClockService.Now())

	trySecondaryProvider := false

	sandBoxMode, _ := s.ConfigService.Bool("sms.sandbox_mode")

	var status string
	var err error

	if !sandBoxMode {
		status, err = primaryProvider.SendSMSMessage(message)
		if err != nil {
			trySecondaryProvider = true

			smsTrackerEntity.SetPrimaryProviderError(err.Error())
			s.ErrorLoggerService.LogError(err)
		}
	} else {
		status = success
	}

	if trySecondaryProvider && secondaryProvider != nil {
		smsTrackerEntity.SetFromSecondaryProvider(secondaryProvider.GetName())

		status, err = secondaryProvider.SendSMSMessage(message)
		if err != nil {
			smsTrackerEntity.SetSecondaryProviderError(err.Error())
			s.ErrorLoggerService.LogError(err)
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
