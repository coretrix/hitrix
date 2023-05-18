package sms

import (
	"errors"
	"fmt"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"

	"github.com/kavenegar/kavenegar-go"
)

const Kavenegar = "kavenegar"

type KavenegarProvider struct {
	APIKey string
	Sender string
}

func NewKavenegarProvider(configService config.IConfig, _ clock.IClock) (IProvider, error) {
	// register kavenegar
	apiKey, ok := configService.String("sms.kavenegar.api_key")
	if !ok {
		return nil, errors.New("missing sms.kavenegar.api_key")
	}
	sender, ok := configService.String("sms.kavenegar.sender")
	if !ok {
		return nil, errors.New("missing sms.kavenegar.sender")
	}

	return &KavenegarProvider{
		APIKey: apiKey,
		Sender: sender,
	}, nil
}

func (g *KavenegarProvider) GetName() string {
	return Kavenegar
}

func (g *KavenegarProvider) SendSMSMessage(message *Message) (string, error) {
	api := kavenegar.New(g.APIKey)

	res, err := api.Message.Send(g.Sender, []string{message.Number}, message.Text, nil)
	if err != nil {
		return err.Error(), err
	}

	if len(res) < 1 {
		e := fmt.Errorf("there was a problem sending sms")

		return e.Error(), e
	}

	return res[0].StatusText, nil
}
