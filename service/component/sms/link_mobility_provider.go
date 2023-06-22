package sms

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/coretrix/hitrix/pkg/helper"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
)

const LinkMobility = "link-mobility"

type LinkMobilityProvider struct {
	Service   int
	Key       string
	Secret    string
	Endpoint  string
	Shortcode int
}

func NewLinkMobilityProvider(configService config.IConfig, _ clock.IClock) (IProvider, error) {
	service, ok := configService.Int("sms.link_mobility.service")
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

	shortcode, ok := configService.Int("sms.link_mobility.shortcode")
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
		ServiceID: strconv.Itoa(g.Service),
		From:      strconv.Itoa(g.Shortcode),
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

	responseBodyJSON := &struct {
		Meta struct {
			Code int `json:"code"`
		} `json:"meta"`
	}{}

	err = json.Unmarshal(responseBody, responseBodyJSON)
	if err != nil {
		return failure, fmt.Errorf("cannot unmarshal response Response: %s", string(responseBody))
	}

	if responseBodyJSON.Meta.Code != 200 {
		return failure, fmt.Errorf("unexpected status code Response: %s", string(responseBody))
	}

	return success, nil
}

func (g *LinkMobilityProvider) getHeaders(body []*linkMobilityMsg) map[string]string {
	bodyByte, _ := json.Marshal(body)
	hmacSignature := genHMAC512(bodyByte, []byte(g.Secret))
	signature := hex.EncodeToString(hmacSignature)

	return map[string]string{
		"Content-Type": "application/json",
		"x-api-key":    g.Key,
		"x-api-sign":   signature,
		"Expect":       "",
	}
}

func genHMAC512(ciphertext, secret []byte) []byte {
	mac := hmac.New(sha512.New, secret)
	mac.Write(ciphertext)
	hmacSum := mac.Sum(nil)

	return hmacSum
}
