package sms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/coretrix/hitrix/pkg/helper"

	"github.com/kevinburke/twilio-go"
)

type TwilioGateway struct {
	SID         string
	Token       string
	FromNumber  string
	AuthyURL    string
	AuthyAPIKey string
	VerifyURL   string
	VerifySID   string
}

type twilioResponse struct {
	Valid bool
}

func (g *TwilioGateway) SendOTPSMS(otp *OTP) (string, error) {
	data := url.Values{}
	data.Set("via", "sms")
	data.Set("phone_number", otp.Phone.Number)
	data.Set("country_code", otp.Phone.ISO3166.CountryCode)
	data.Set("custom_code", otp.OTP)
	data.Set("locale", "en")
	data.Set("code_length", "4")

	baseURL, err := url.Parse(g.AuthyURL)
	if err != nil {
		return err.Error(), err
	}

	baseURL.RawQuery = data.Encode()

	headers := map[string]string{
		"Content-Type":    "application/x-www-form-urlencoded",
		"Content-Length":  strconv.Itoa(len(data.Encode())),
		"X-Authy-API-Key": g.AuthyAPIKey,
	}

	responseBody, _, code, err := helper.Call(
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
		return failure, fmt.Errorf("expected status code OK, but got %v Response: %s", code, string(responseBody))
	}

	// TODO: find out the format of response

	return success, nil
}

func (g *TwilioGateway) SendOTPCallout(otp *OTP) (string, error) {
	data := url.Values{}
	data.Set("api_key", g.AuthyAPIKey)
	data.Set("via", "call")
	data.Set("phone_number", otp.Phone.Number)
	data.Set("country_code", otp.Phone.ISO3166.CountryCode)
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

	responseBody, _, code, err := helper.Call(
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
		return failure, fmt.Errorf("expected status code OK, but got %v Response: %s", code, string(responseBody))
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
	data := url.Values{}
	data.Set("To", otp.Phone.Number)
	data.Set("Channel", "sms")

	endpoint := strings.Join([]string{
		g.VerifyURL,
		"Services", "/",
		g.VerifySID, "/",
		"Verifications",
	}, "")
	baseURL, err := url.Parse(endpoint)
	if err != nil {
		return err.Error(), err
	}

	headers := map[string]string{
		"Authorization": "Basic " + helper.BasicAuth(g.SID, g.Token),
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	responseBody, _, code, err := helper.Call(
		context.Background(),
		"POST",
		baseURL.String(),
		headers,
		time.Duration(timeoutInSeconds)*time.Second,
		data.Encode(),
		nil)

	if err != nil {
		return failure, err
	}

	if code < 200 && code >= 300 {
		return failure, fmt.Errorf("expected status code OK, but got %v Response: %s", code, string(responseBody))
	}

	return success, nil
}

func (g *TwilioGateway) SendVerificationCallout(otp *OTP) (string, error) {
	data := url.Values{}
	data.Set("To", otp.Phone.Number)
	data.Set("Channel", "call")

	endpoint := strings.Join([]string{
		g.VerifyURL,
		"Services", "/",
		g.VerifySID, "/",
		"Verifications",
	}, "")
	baseURL, err := url.Parse(endpoint)
	if err != nil {
		return err.Error(), err
	}

	headers := map[string]string{
		"Authorization": "Basic " + helper.BasicAuth(g.SID, g.Token),
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	responseBody, _, code, err := helper.Call(
		context.Background(),
		"POST",
		baseURL.String(),
		headers,
		time.Duration(timeoutInSeconds)*time.Second,
		data.Encode(),
		nil)

	if err != nil {
		return failure, err
	}

	if code < 200 && code >= 300 {
		return failure, fmt.Errorf("expected status code OK, but got %v Response: %s", code, string(responseBody))
	}

	return success, nil
}

func (g *TwilioGateway) VerifyCode(otp *OTP) (string, error) {
	data := url.Values{}
	data.Set("To", otp.Phone.Number)
	data.Set("Code", otp.OTP)

	endpoint := strings.Join([]string{
		g.VerifyURL,
		"Services", "/",
		g.VerifySID, "/",
		"VerificationCheck",
	}, "")
	baseURL, err := url.Parse(endpoint)
	if err != nil {
		return err.Error(), err
	}

	headers := map[string]string{
		"Authorization": "Basic " + helper.BasicAuth(g.SID, g.Token),
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	responseBody, _, code, err := helper.Call(
		context.Background(),
		"POST",
		baseURL.String(),
		headers,
		time.Duration(timeoutInSeconds)*time.Second,
		data.Encode(),
		nil)

	if err != nil {
		return failure, err
	}

	if code < 200 || code >= 300 {
		return failure, fmt.Errorf("expected status code OK, but got %v Response: %s", code, string(responseBody))
	}

	response := twilioResponse{}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return failure, err
	}

	if !response.Valid {
		return failure, errors.New("the code is not valid")
	}

	return success, nil
}
