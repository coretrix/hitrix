package sms

const (
	Sinch  = "sinch"
	Twilio = "twilio"

	success = "sent successfully"
	failure = "sent unsuccessfully"

	timeoutInSeconds = 5

	Kavenegar = "kavenegar"
)

type Gateway interface {
	SendOTPSMS(*OTP) (string, error)
	SendOTPCallout(*OTP) (string, error)
	SendSMSMessage(*Message) (string, error)
	SendCalloutMessage(*Message) (string, error)
	SendVerificationSMS(*OTP) (string, error)
	SendVerificationCallout(*OTP) (string, error)
	VerifyCode(*OTP) (string, error)
	// VerifyCode(code *string, number *string) (string, error)
}

type OTP struct {
	OTP      string
	Number   string
	CC       string
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
