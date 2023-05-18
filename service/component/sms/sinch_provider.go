package sms

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/coretrix/hitrix/service/component/config"
	"net/http"
	"time"

	"github.com/coretrix/hitrix/pkg/helper"
	"github.com/coretrix/hitrix/service/component/clock"
)

const (
	javascriptISOString = `2006-01-02T15:04:05.999Z07:00`
	Sinch               = "sinch"
)

type SinchProvider struct {
	Clock      clock.IClock
	AppID      string
	AppSecret  string
	MsgURL     string
	FromNumber string
}

func NewSinchProvider(configService config.IConfig, clockService clock.IClock) (IProvider, error) {
	appID, ok := configService.String("sms.sinch.app_id")
	if !ok {
		return nil, errors.New("missing sms.sinch.app_id")
	}
	appSecret, ok := configService.String("sms.sinch.app_secret")
	if !ok {
		return nil, errors.New("missing sms.sinch.app_secret")
	}
	msgURL, ok := configService.String("sms.sinch.msg_url")
	if !ok {
		return nil, errors.New("missing sms.sinch.msg_url")
	}
	fromNumberSinch, ok := configService.String("sms.sinch.from_number")
	if !ok {
		return nil, errors.New("missing sms.sinch.from_number")
	}

	return &SinchProvider{
		Clock:      clockService,
		AppID:      appID,
		AppSecret:  appSecret,
		MsgURL:     msgURL,
		FromNumber: fromNumberSinch,
	}, nil
}

func (g *SinchProvider) GetName() string {
	return Sinch
}

func (g *SinchProvider) SendSMSMessage(message *Message) (string, error) {
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
	responseBody, _, code, err := helper.Call(
		context.Background(),
		"POST",
		g.MsgURL+"/"+message.Number,
		headers,
		time.Duration(timeoutInSeconds)*time.Second,
		body,
		nil)

	if err != nil {
		return failure, err
	}

	if code != http.StatusOK {
		return failure, fmt.Errorf("expected status code OK, but got %v Response: %s", code, string(responseBody))
	}

	return success, nil
}

func (g *SinchProvider) getSinchHeaders() map[string]string {
	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(g.AppID+":"+g.AppSecret)),
		"X-Timestamp":   g.Clock.Now().Format(javascriptISOString),
	}
}
