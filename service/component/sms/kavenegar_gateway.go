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
	api := kavenegar.New(g.APIKey)

	res, err := api.Verify.Lookup(otp.Phone.Number, otp.Template, otp.OTP, nil)
	if err != nil {
		return err.Error(), err
	}

	if res.Status != kavenegar.Type_MessageStatus_Sent && res.Status != kavenegar.Type_MessageStatus_Delivered {
		e := fmt.Errorf("there was a problem sending sms")

		return e.Error(), e
	}

	return res.Message, nil
}

func (g *KavenegarGateway) SendOTPCallout(otp *OTP) (string, error) {
	return g.SendCalloutMessage(&Message{
		Text:   fmt.Sprintf(otp.Template, otp.OTP),
		Number: otp.Phone.Number,
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

func (g *KavenegarGateway) SendVerificationSMS(_ *OTP) (string, error) {
	panic("not implemented.")
}

func (g *KavenegarGateway) SendVerificationCallout(_ *OTP) (string, error) {
	panic("not implemented.")
}

func (g *KavenegarGateway) VerifyCode(_ *OTP) (string, error) {
	panic("not implemented.")
}
