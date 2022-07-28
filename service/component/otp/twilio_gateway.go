package otp

import (
	"encoding/json"

	"github.com/twilio/twilio-go"
	"github.com/twilio/twilio-go/client"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
)

const SMSOTPProviderTwilio = "Twilio"

type Twilio struct {
	Client          *twilio.RestClient
	VerificationSID string
}

func NewTwilioSMSOTPProvider(accountSid, authToken, verificationSid string) *Twilio {
	return &Twilio{
		Client: twilio.NewRestClientWithParams(twilio.RestClientParams{
			Username: accountSid,
			Password: authToken,
		}),
		VerificationSID: verificationSid,
	}
}

func (t *Twilio) GetName() string {
	return SMSOTPProviderTwilio
}

func (t *Twilio) GetCode() string {
	return ""
}

func (t *Twilio) GetPhonePrefixes() []string {
	return nil
}

func (t *Twilio) SendOTP(phone *Phone, _ string) (string, string, error) {
	createVerificationParams := &openapi.CreateVerificationParams{}

	createVerificationParams.SetChannel("sms")
	createVerificationParams.SetTo(phone.Number)

	request, jsonError := t.toJSON(createVerificationParams)

	if jsonError != nil {
		return "", "", jsonError
	}

	verifyV2Verification, err := t.Client.VerifyV2.CreateVerification(t.VerificationSID, createVerificationParams)

	response, jsonError := t.toJSON(verifyV2Verification)

	if jsonError != nil {
		return request, "", jsonError
	}

	if err != nil {
		return request, response, err
	}

	//TODO check error codes

	return request, response, nil
}

func (t *Twilio) Call(phone *Phone, _ string, customMessage string) (string, string, error) {
	createVerificationParams := &openapi.CreateVerificationParams{}

	createVerificationParams.SetChannel("call")
	createVerificationParams.SetTo(phone.Number)

	if customMessage != "" {
		createVerificationParams.SetCustomMessage(customMessage)
	}

	createVerificationParams.SetTo(phone.Number)

	request, jsonError := t.toJSON(createVerificationParams)

	if jsonError != nil {
		return "", "", jsonError
	}

	verifyV2Verification, err := t.Client.VerifyV2.CreateVerification(t.VerificationSID, createVerificationParams)

	response, jsonError := t.toJSON(verifyV2Verification)

	if jsonError != nil {
		return request, "", jsonError
	}

	if err != nil {
		return request, response, err
	}

	//TODO check error codes

	return request, response, nil
}

func (t *Twilio) VerifyOTP(phone *Phone, code, _ string) (string, string, bool, bool, error) {
	createVerificationCheckParams := &openapi.CreateVerificationCheckParams{}

	createVerificationCheckParams.SetTo(phone.Number)
	createVerificationCheckParams.SetCode(code)

	request, jsonError := t.toJSON(createVerificationCheckParams)

	if jsonError != nil {
		return "", "", false, false, jsonError
	}

	verifyV2VerificationCheck, err := t.Client.VerifyV2.CreateVerificationCheck(t.VerificationSID, createVerificationCheckParams)

	response, jsonError := t.toJSON(verifyV2VerificationCheck)

	if jsonError != nil {
		return request, "", false, false, jsonError
	}

	if err != nil {
		twilioRestError, ok := err.(*client.TwilioRestError)

		if !ok {
			return request, response, false, false, err
		}

		if twilioRestError.Status == 404 && twilioRestError.Code == 20404 {
			//already confirmed or expired session - invalid request
			return request, twilioRestError.Error(), false, false, nil
		}

		return request, twilioRestError.Error(), false, false, err
	}

	return request, response, true, *verifyV2VerificationCheck.Valid, nil
}

func (t *Twilio) toJSON(data interface{}) (string, error) {
	dataJSON, err := json.Marshal(data)

	if err != nil {
		return "", err
	}

	return string(dataJSON), nil
}
