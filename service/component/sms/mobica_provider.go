package sms

import (
	"context"
	"errors"
	"fmt"
	"github.com/coretrix/hitrix/pkg/helper"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
	"net/http"
	"time"
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

	route, ok := configService.String("sms.mobica.route")
	if !ok {
		return nil, errors.New("missing sms.mobica.route")
	}

	from, ok := configService.String("sms.mobica.from")
	if !ok {
		return nil, errors.New("missing sms.mobica.from")
	}

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
	route    string
	smartCut uint8
	message  string
	from     string
}

type mobicaMsg struct {
	phone string
	sms   sms
	user  string
	pass  string
}

func (g *MobicaProvider) SendSMSMessage(message *Message) (string, error) {
	body := &mobicaMsg{
		phone: message.Number,
		sms: sms{
			route:    g.Route,
			smartCut: 1,
			message:  message.Text,
			from:     g.From,
		},
		user: g.Email,
		pass: g.Password,
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

	return success, nil
}

func (g *MobicaProvider) getHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
}
