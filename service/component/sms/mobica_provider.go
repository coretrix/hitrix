package sms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/coretrix/hitrix/pkg/helper"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
)

const Mobica = "mobica"

type MobicaProvider struct {
	Email    string
	Password string
	Route    string
	From     string
	Endpoint string
}

func NewMobicaProvider(configService config.IConfig, _ clock.IClock) (IProvider, error) {
	email, ok := configService.String("sms.mobica.email")
	if !ok {
		return nil, errors.New("missing sms.mobica.email")
	}

	password, ok := configService.String("sms.mobica.password")
	if !ok {
		return nil, errors.New("missing sms.mobica.password")
	}

	route, _ := configService.String("sms.mobica.route")

	from, _ := configService.String("sms.mobica.from")

	endpoint, ok := configService.String("sms.mobica.endpoint")
	if !ok {
		return nil, errors.New("missing sms.mobica.endpoint")
	}

	return &MobicaProvider{
		Email:    email,
		Password: password,
		Route:    route,
		From:     from,
		Endpoint: endpoint,
	}, nil
}

func NewMobicaProviderNoConfig(_ config.IConfig, _ clock.IClock) (IProvider, error) {
	return &MobicaProvider{}, nil
}

func (g *MobicaProvider) GetName() string {
	return Mobica
}

type sms struct {
	Route    string `json:"route,omitempty"`
	SmartCut uint8  `json:"smartCut"`
	Message  string `json:"message"`
	From     string `json:"from,omitempty"`
}

type mobicaMsg struct {
	Phone string `json:"phone"` //it supports comma separated phones
	Sms   sms    `json:"sms"`
	User  string `json:"user"`
	Pass  string `json:"pass"`
}

func (g *MobicaProvider) SendSMSMessage(message *Message) (string, error) {
	body := &mobicaMsg{
		Phone: message.Number,
		Sms: sms{
			Route:    g.Route,
			SmartCut: 1,
			Message:  message.Text,
			From:     g.From,
		},
		User: g.Email,
		Pass: g.Password,
	}

	headers := g.getHeaders()
	responseBody, _, code, err := helper.Call(
		context.Background(),
		"POST",
		g.Endpoint,
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

	responseBodyJSON := &struct {
		Status int    `json:"status"`
		Desc   string `json:"desc"`
	}{}

	err = json.Unmarshal(responseBody, responseBodyJSON)
	if err != nil {
		return failure, fmt.Errorf("cannot unmarshal response Response: %s", string(responseBody))
	}

	if responseBodyJSON.Status != 1004 {
		return failure, fmt.Errorf("unexpected status code Response: %s", string(responseBody))
	}

	return success, nil
}

func (g *MobicaProvider) getHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
}
