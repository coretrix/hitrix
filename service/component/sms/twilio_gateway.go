package sms

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/coretrix/hitrix/pkg/helper"

	"github.com/kevinburke/twilio-go"
	twilioGo "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
)

type TwilioGateway struct {
	SID         string
	Token       string
	FromNumber  string
	AuthyURL    string
	AuthyAPIKey string
	VerifySID   string
}

func (g *TwilioGateway) SendOTPSMS(otp *OTP) (string, error) {
	data := url.Values{}
	data.Set("api_key", g.AuthyAPIKey)
	data.Set("via", "sms")
	data.Set("phone_number", otp.Number)
	data.Set("country_code", otp.CC)
	data.Set("custom_code", otp.OTP)
	data.Set("locale", "en")
	data.Set("code_length", "4")

	baseURL, err := url.Parse(g.AuthyURL)
	if err != nil {
		return err.Error(), err
	}

	baseURL.RawQuery = data.Encode()

	headers := map[string]string{
		"Content-Type":   "application/x-www-form-urlencoded",
		"Content-Length": strconv.Itoa(len(data.Encode())),
	}

	_, _, code, err := helper.Call(
		context.Background(),
		"POST",
		baseURL.String(),
		headers,
		time.Duration(timeoutInSeconds)*time.Second,
		nil,
		nil)

	if err != nil {
		return err.Error(), err
	}

	if code != http.StatusOK {
		e := fmt.Errorf("expected status code OK, but got %v", code)
		return e.Error(), e
	}

	// TODO: find out the format of response

	return "success", nil
}

func (g *TwilioGateway) SendOTPCallout(otp *OTP) (string, error) {
	data := url.Values{}
	data.Set("api_key", g.AuthyAPIKey)
	data.Set("via", "call")
	data.Set("phone_number", otp.Number)
	data.Set("country_code", otp.CC)
	data.Set("custom_code", otp.OTP)
	data.Set("locale", "en")
	data.Set("code_length", "4")

	baseURL, err := url.Parse(g.AuthyURL)
	if err != nil {
		return failure, err
	}

	baseURL.RawQuery = data.Encode()

	headers := map[string]string{
		"Content-Type":   "application/x-www-form-urlencoded",
		"Content-Length": strconv.Itoa(len(data.Encode())),
	}

	_, _, code, err := helper.Call(
		context.Background(),
		baseURL.String(),
		"POST",
		headers,
		time.Duration(timeoutInSeconds)*time.Second,
		nil,
		nil)

	if err != nil {
		return failure, err
	}

	if code != http.StatusOK {
		e := fmt.Errorf("expected status code OK, but got %v", code)
		return e.Error(), e
	}
	// TODO: find out the format of response

	return success, nil
}

func (g *TwilioGateway) SendSMSMessage(message *Message) (string, error) {
	api := twilio.NewClient(g.SID, g.Token, nil)

	msg, err := api.Messages.SendMessage(g.FromNumber, message.Number, message.Text, nil)
	if err != nil {
		return err.Error(), err
	}

	return msg.Status.Friendly(), nil
}

func (g *TwilioGateway) SendCalloutMessage(message *Message) (string, error) {
	// not supported for now
	return "", nil
}

func (g *TwilioGateway) SendVerificationSMS(otp *OTP) (string, error) {
	client := twilioGo.NewRestClientWithParams(twilioGo.RestClientParams{
		Username: g.SID,
		Password: g.Token,
	})
	channel := "sms"
	locale := "en"
	params := &openapi.CreateVerificationParams{
		To: &otp.Number,
		// CustomCode: &otp.OTP,
		Channel: &channel,
		Locale:  &locale,
	}
	_, err := client.VerifyV2.CreateVerification(g.VerifySID, params)
	if err != nil {
		return err.Error(), err
	}

	return success, nil
}

func (g *TwilioGateway) SendVerificationCallout(otp *OTP) (string, error) {
	client := twilioGo.NewRestClientWithParams(twilioGo.RestClientParams{
		Username: g.SID,
		Password: g.Token,
	})
	channel := "call"
	locale := "en"
	params := &openapi.CreateVerificationParams{
		To: &otp.Number,
		// CustomCode: &otp.OTP,
		Channel: &channel,
		Locale:  &locale,
	}
	_, err := client.VerifyV2.CreateVerification(g.VerifySID, params)
	if err != nil {
		return err.Error(), err
	}

	return success, nil
}

// func (g *TwilioGateway) VerifyCode(code *string, number *string) (string, error) {
func (g *TwilioGateway) VerifyCode(opt *OTP) (string, error) {
	client := twilioGo.NewRestClientWithParams(twilioGo.RestClientParams{
		Username: g.SID,
		Password: g.Token,
	})
	params := &openapi.CreateVerificationCheckParams{
		Code: &opt.OTP,
		To:   &opt.Number,
	}
	_, err := client.VerifyV2.CreateVerificationCheck(g.VerifySID, params)
	if err != nil {
		return err.Error(), err
	}

	return success, nil
}
