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
	Clock              clock.IClock
	PrimaryProvider    IProvider
	SecondaryProvider  IProvider
	ErrorLoggerService errorlogger.ErrorLogger
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
	smsTrackerEntity.SetSentAt(s.Clock.Now())

	trySecondaryProvider := false

	status, err := primaryProvider.SendSMSMessage(message)
	if err != nil {
		trySecondaryProvider = true

		smsTrackerEntity.SetPrimaryProviderError(err.Error())
		s.ErrorLoggerService.LogError(err)
	}

	if trySecondaryProvider {
		if secondaryProvider == nil {
			return fmt.Errorf("secondary provider not supported")
		}

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
