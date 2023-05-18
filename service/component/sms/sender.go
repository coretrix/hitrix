package sms

import (
	"fmt"
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/clock"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
)

type ISender interface {
	SendMessage(ormService *beeorm.Engine, errorLoggerService errorlogger.ErrorLogger, message *Message) error
}

type Sender struct {
	Clock             clock.IClock
	ProviderContainer map[string]IProvider
	Logger            LogEntity
}

func (s *Sender) SendMessage(ormService *beeorm.Engine, errorLoggerService errorlogger.ErrorLogger, message *Message) error {
	primaryProvider := message.Provider.Primary
	secondaryProvider := message.Provider.Secondary

	primaryGateway, ok := s.ProviderContainer[primaryProvider]
	if !ok {
		return fmt.Errorf("primary provider not supported")
	}

	smsTrackerEntity := s.Logger
	smsTrackerEntity.SetTo(message.Number)
	smsTrackerEntity.SetType(entity.SMSTrackerTypeCallout)
	smsTrackerEntity.SetText(message.Text)
	smsTrackerEntity.SetFromPrimaryProvider(primaryProvider)
	smsTrackerEntity.SetSentAt(s.Clock.Now())

	trySecondaryProvider := false

	status, err := primaryGateway.SendSMSMessage(message)
	if err != nil {
		trySecondaryProvider = true

		smsTrackerEntity.SetPrimaryProviderError(err.Error())
		errorLoggerService.LogError(err)
	}

	if trySecondaryProvider {
		secondaryGateway, ok := s.ProviderContainer[secondaryProvider]
		if !ok {
			return fmt.Errorf("secondary provider not supported")
		}

		smsTrackerEntity.SetFromSecondaryProvider(secondaryProvider)

		status, err = secondaryGateway.SendSMSMessage(message)
		if err != nil {
			smsTrackerEntity.SetSecondaryProviderError(err.Error())
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
