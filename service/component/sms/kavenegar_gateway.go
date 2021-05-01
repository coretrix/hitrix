package sms

import (
	"fmt"

	"github.com/kavenegar/kavenegar-go"
)

type KavenegarGateway struct {
	APIKey string
	Sender string
}

func (g *KavenegarGateway) SendOTPSMS(otp *OTP) (string, error) {
	return g.SendSMSMessage(&Message{
		Text:   otp.OTP,
		Number: otp.Number,
	})
}

func (g *KavenegarGateway) SendOTPCallout(otp *OTP) (string, error) {
	return g.SendCalloutMessage(&Message{
		Text:   otp.OTP,
		Number: otp.Number,
	})
}

func (g *KavenegarGateway) SendSMSMessage(message *Message) (string, error) {
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

func (g *KavenegarGateway) SendCalloutMessage(message *Message) (string, error) {
	api := kavenegar.New(g.APIKey)
	tts, err := api.Call.MakeTTS(message.Number, message.Text, &kavenegar.CallParam{})
	if err != nil {
		return err.Error(), err
	}

	return tts.StatusText, nil
}
