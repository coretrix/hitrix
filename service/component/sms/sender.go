package sms

import (
	"fmt"
	"strings"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
)

type ISender interface {
	SendMessage(ormService *beeorm.Engine, message *Message, params map[string]interface{}) error
}

type Sender struct {
	ConfigService      config.IConfig
	ClockService       clock.IClock
	ErrorLoggerService errorlogger.ErrorLogger
	PrimaryProvider    IProvider
	SecondaryProvider  IProvider
}

func (s *Sender) SendMessage(ormService *beeorm.Engine, message *Message, params map[string]interface{}) error {
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

	text := message.Text

	for paramName, value := range params {
		text = strings.Replace(text, fmt.Sprintf("[[%s]]", paramName), fmt.Sprintf("%v", value), -1)
	}

	smsTrackerEntity := &entity.SmsTrackerEntity{}
	smsTrackerEntity.SetTo(message.Number)
	smsTrackerEntity.SetType(entity.SMSTrackerTypeSMS)
	smsTrackerEntity.SetText(text)
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
