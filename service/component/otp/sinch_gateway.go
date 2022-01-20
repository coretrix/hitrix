package otp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/coretrix/hitrix/pkg/helper"
)

const SMSOTPProviderSinch = "Sinch"

type Sinch struct {
	AppID           string
	AppSecret       string
	VerificationURL string
}

func NewSinchSMSOTPProvider(appID, appSecret, verificationURL string) *Sinch {
	return &Sinch{
		AppID:           appID,
		AppSecret:       appSecret,
		VerificationURL: verificationURL,
	}
}

func (s *Sinch) GetName() string {
	return SMSOTPProviderSinch
}

func (s *Sinch) GetCode() string {
	return ""
}

func (s *Sinch) GetPhonePrefixes() []string {
	return nil
}

func (s *Sinch) SendOTP(phone *Phone, _ string) (string, string, error) {
	body := &struct {
		Identity *struct {
			Type     string `json:"type"`
			Endpoint string `json:"endpoint"`
		} `json:"identity"`
		Method string `json:"method"`
	}{
		Identity: &struct {
			Type     string `json:"type"`
			Endpoint string `json:"endpoint"`
		}{
			Type:     "number",
			Endpoint: phone.Number,
		},
		Method: "sms",
	}

	request, err := s.toJSON(body)
	if err != nil {
		return request, "", err
	}

	headers := s.getSinchHeaders()
	responseBody, _, code, err := helper.Call(
		context.Background(),
		"POST",
		s.VerificationURL,
		headers,
		time.Duration(5)*time.Second,
		body,
		nil)

	if err != nil {
		return request, string(responseBody), err
	}

	if code != http.StatusOK {
		return request, string(responseBody), fmt.Errorf("expected status code 200, but got %v", code)
	}

	return request, string(responseBody), nil
}

func (s *Sinch) Call(_ *Phone, _ string, _ string) (string, string, error) {
	// not implemented
	return "", "", nil
}

func (s *Sinch) VerifyOTP(phone *Phone, code, _ string) (string, string, bool, bool, error) {
	body := &struct {
		Method string `json:"method"`
		SMS    *struct {
			Code string `json:"code"`
		} `json:"sms"`
	}{
		Method: "sms",
		SMS: &struct {
			Code string `json:"code"`
		}{Code: code},
	}

	request, err := s.toJSON(body)
	if err != nil {
		return request, "", false, false, err
	}

	headers := s.getSinchHeaders()
	responseBody, _, respCode, err := helper.Call(
		context.Background(),
		"PUT",
		s.VerificationURL+"/number/"+phone.Number,
		headers,
		time.Duration(5)*time.Second,
		body,
		nil)

	if err != nil {
		return request, string(responseBody), false, false, err
	}

	if respCode != http.StatusOK && strings.Contains(string(responseBody), "Invalid identity or code") {
		return request, string(responseBody), false, false, nil
	} else if respCode != http.StatusOK {
		return request, string(responseBody), false, false, fmt.Errorf("expected status code 200, but got %v", respCode)
	}

	respStruct := &struct {
		ID        string `json:"id"`
		Method    string `json:"method"`
		Status    string `json:"status"`
		Reason    string `json:"reason"`
		Reference string `json:"reference"`
		Source    string `json:"source"`
	}{}

	if err := json.Unmarshal(responseBody, respStruct); err != nil {
		return request, string(responseBody), false, false, err
	}

	return request, string(responseBody), true, respStruct.Status == "SUCCESSFUL", nil
}

func (s *Sinch) getSinchHeaders() map[string]string {
	base64Credentials := base64.StdEncoding.EncodeToString([]byte(s.AppID + ":" + s.AppSecret))
	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Basic " + base64Credentials,
	}
}

func (s *Sinch) toJSON(data interface{}) (string, error) {
	dataJSON, err := json.Marshal(data)

	if err != nil {
		return "", err
	}

	return string(dataJSON), nil
}
