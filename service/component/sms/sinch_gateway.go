package sms

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/coretrix/hitrix/pkg/helper"

	"github.com/coretrix/hitrix/service/component/clock"
)

const (
	javascriptISOString = `2006-01-02T15:04:05.999Z07:00`
)

type SinchGateway struct {
	Clock        clock.Clock
	AppID        string
	AppSecret    string
	MsgURL       string
	FromNumber   string
	CallURL      string
	CallerNumber string
}

func (g *SinchGateway) SendOTPSMS(otp *OTP) (string, error) {
	return g.SendSMSMessage(&Message{
		Text:   otp.OTP,
		Number: otp.Number,
	})
}

func (g *SinchGateway) SendOTPCallout(otp *OTP) (string, error) {
	return g.SendCalloutMessage(&Message{
		Text:   otp.OTP,
		Number: otp.Number,
	})
}

func (g *SinchGateway) SendSMSMessage(message *Message) (string, error) {
	body := struct {
		From    string `json:"from"`
		Message string `json:"message"`
		Caller  string `json:"caller"`
	}{
		From:    g.FromNumber,
		Message: message.Text,
		Caller:  message.Number,
	}

	headers := g.getSinchHeaders()
	_, _, code, err := helper.Call(
		context.Background(),
		g.MsgURL+"/"+message.Number,
		"POST",
		headers,
		time.Duration(timeoutInSeconds)*time.Second,
		body,
		nil)

	if err != nil {
		return failure, err
	}

	if code != http.StatusOK {
		return failure, fmt.Errorf("expected status code OK, but got %v", code)
	}

	return success, nil
}

type ttsCallOut struct {
	CLI         string       `json:"cli"`
	Domain      string       `json:"domain"`
	Locale      string       `json:"locale"`
	Text        string       `json:"text"`
	Destination *destination `json:"destination"`
}

type destination struct {
	Type     string `json:"type"`
	Endpoint string `json:"endpoint"`
}

func (g *SinchGateway) SendCalloutMessage(message *Message) (string, error) {
	body := struct {
		Method     string      `json:"method"`
		TTSCallOut *ttsCallOut `json:"ttsCallOut"`
	}{
		Method: "ttsCallout",
		TTSCallOut: &ttsCallOut{
			CLI:    g.CallerNumber,
			Domain: "pstn",
			Locale: "en-US",
			Text:   message.Text + "...." + message.Text + "...." + message.Text,
			Destination: &destination{
				Type:     "number",
				Endpoint: message.Number,
			},
		},
	}

	headers := g.getSinchHeaders()
	_, _, code, err := helper.Call(
		context.Background(),
		g.CallURL,
		"POST",
		headers,
		time.Duration(timeoutInSeconds)*time.Second,
		body,
		nil)

	if err != nil {
		return failure, err
	}

	if code != http.StatusOK {
		return failure, fmt.Errorf("expected status code OK, but got %v", code)
	}

	return success, nil
}

func (g *SinchGateway) getSinchHeaders() map[string]string {
	base64Credentials := base64.StdEncoding.EncodeToString([]byte(g.AppID + ":" + g.AppSecret))
	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Basic" + base64Credentials,
		"X-Timestamp":   g.Clock.Now().Format(javascriptISOString),
	}
}
