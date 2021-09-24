package trustpilot

import (
	"bytes"
	"net/http"
	"net/url"
	"time"
)

func makeTrustpilotAuthenticatedRequest(accessToken string, method string, endpoint string, params url.Values, body []byte) (*http.Request, error) {
	requestUrl, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	if params != nil {
		requestUrl.RawQuery = params.Encode()
	}

	req, err := http.NewRequest(method, requestUrl.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	setRequestHeaders(req, accessToken)

	return req, nil
}

func makeTrustpilotAuthenticatedClient(accessToken string) (*http.Client, error) {
	client := http.Client{
		Timeout: time.Second * 10,

		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			setRequestHeaders(req, accessToken)
			return nil
		},
	}
	return &client, nil
}

func setRequestHeaders(req *http.Request, accessToken string) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
}
