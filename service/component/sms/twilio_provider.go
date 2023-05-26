package sms

import (
	"errors"

	"github.com/kevinburke/twilio-go"

	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
)

const Twilio = "twilio"

type TwilioProvider struct {
	SID        string
	Token      string
	FromNumber string
}

func NewTwilioProvider(configService config.IConfig, _ clock.IClock) (IProvider, error) {
	sid, ok := configService.String("sms.twilio.sid")
	if !ok {
		return nil, errors.New("missing sms.twilio.sid")
	}

	token, ok := configService.String("sms.twilio.token")
	if !ok {
		return nil, errors.New("missing sms.twilio.token")
	}

	fromNumberTwilio, ok := configService.String("sms.twilio.from_number")
	if !ok {
		return nil, errors.New("missing sms.twilio.from_number")
	}

	return &TwilioProvider{
		SID:        sid,
		Token:      token,
		FromNumber: fromNumberTwilio,
	}, nil
}

func (g *TwilioProvider) GetName() string {
	return Twilio
}

func (g *TwilioProvider) SendSMSMessage(message *Message) (string, error) {
	api := twilio.NewClient(g.SID, g.Token, nil)

	_, err := api.Messages.SendMessage(g.FromNumber, message.Number, message.Text, nil)
	if err != nil {
		return err.Error(), err
	}

	return success, nil
}
