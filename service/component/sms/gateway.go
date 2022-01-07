package sms

import "github.com/dongri/phonenumber"

const (
	Sinch     = "sinch"
	Twilio    = "twilio"
	Kavenegar = "kavenegar"

	success = "sent successfully"
	failure = "sent unsuccessfully"

	timeoutInSeconds = 5
)

type Gateway interface {
	SendOTPSMS(*OTP) (string, error)
	SendOTPCallout(*OTP) (string, error)
	SendSMSMessage(*Message) (string, error)
	SendCalloutMessage(*Message) (string, error)
	SendVerificationSMS(*OTP) (string, error)
	SendVerificationCallout(*OTP) (string, error)
	VerifyCode(*OTP) (string, error)
}

type Phone struct {
	Number  string
	ISO3166 phonenumber.ISO3166
}

type OTP struct {
	OTP      string
	Phone    *Phone
	Provider *Provider
	Template string
}

type Message struct {
	Text     string
	Number   string
	Provider *Provider
}

type Provider struct {
	Primary   string
	Secondary string
}
