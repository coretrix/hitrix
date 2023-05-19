package sms

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coretrix/hitrix/pkg/helper"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
	"net/http"
	"time"
)

const LinkMobility = "link-mobility"

type LinkMobilityProvider struct {
	Service   string
	Key       string
	Secret    string
	Endpoint  string
	Shortcode string
}

func NewLinkMobilityProvider(configService config.IConfig, _ clock.IClock) (IProvider, error) {
	service, ok := configService.String("sms.link_mobility.service")
	if !ok {
		return nil, errors.New("missing sms.link_mobility.service")
	}

	key, ok := configService.String("sms.link_mobility.key")
	if !ok {
		return nil, errors.New("missing sms.link_mobility.key")
	}

	secret, ok := configService.String("sms.link_mobility.secret")
	if !ok {
		return nil, errors.New("missing sms.link_mobility.secret")
	}

	endpoint, ok := configService.String("sms.link_mobility.endpoint")
	if !ok {
		return nil, errors.New("missing sms.link_mobility.endpoint")
	}

	shortcode, ok := configService.String("sms.link_mobility.shortcode")
	if !ok {
		return nil, errors.New("missing sms.link_mobility.shortcode")
	}

	return &LinkMobilityProvider{
		Service:   service,
		Key:       key,
		Secret:    secret,
		Endpoint:  endpoint,
		Shortcode: shortcode,
	}, nil
}

func NewLinkMobilityProviderNoConfig(_ config.IConfig, _ clock.IClock) (IProvider, error) {
	return &LinkMobilityProvider{}, nil
}

func (g *LinkMobilityProvider) GetName() string {
	return LinkMobility
}

type linkMobilityMsg struct {
	From      string `json:"sc"`
	Message   string `json:"text"`
	To        string `json:"msisdn"`
	ServiceID string `json:"service_id"`
}

func (g *LinkMobilityProvider) SendSMSMessage(message *Message) (string, error) {
	row := &linkMobilityMsg{
		ServiceID: g.Service,
		From:      g.Shortcode,
		To:        message.Number,
		Message:   message.Text,
	}

	body := []*linkMobilityMsg{row}

	headers := g.getHeaders(body)
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

func (g *LinkMobilityProvider) getHeaders(body []*linkMobilityMsg) map[string]string {
	bodyByte, _ := json.Marshal(body)
	hmacSignature := genHMAC256(bodyByte, []byte(g.Secret))
	signature := base64.StdEncoding.EncodeToString(hmacSignature)

	return map[string]string{
		"Content-Type": "application/json",
		"x-api-key":    g.Key,
		"x-api-sign":   signature,
		"Expect":       "",
	}
}

func genHMAC256(ciphertext, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(ciphertext)
	hmacSum := mac.Sum(nil)
	return hmacSum
}
