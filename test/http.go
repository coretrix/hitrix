package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/pkg/errors"

	"github.com/coretrix/hitrix/pkg/test"
)

func SendHTTPRequest(env *test.Environment, method string, pathAndQuery string, authenticate bool, response interface{}) error {
	ts := httptest.NewServer(env.GinEngine)
	defer ts.Close()

	r, err := http.NewRequestWithContext(env.Cxt, method, ts.URL+"/v1"+pathAndQuery, nil)
	if err != nil {
		return err
	}

	if authenticate {
		r.Header.Set("Authorization", "Bearer accessToken")
	}

	return sendHTTPRequest(r, response)
}

func SendHTTPRequestWithBody(
	env *test.Environment,
	method string,
	pathAndQuery string,
	input interface{},
	authenticate bool,
	response interface{},
) error {
	return SendHTTPRequestWithBodyAndHeaders(env, method, pathAndQuery, input, authenticate, response, map[string]string{})
}

func SendHTTPRequestWithBodyAndHeaders(
	env *test.Environment,
	method string,
	pathAndQuery string,
	input interface{},
	authenticate bool,
	response interface{},
	headers map[string]string,
) error {
	body, err := json.Marshal(input)
	if err != nil {
		return errors.Wrap(err, "could not marshal post params")
	}

	ts := httptest.NewServer(env.GinEngine)
	defer ts.Close()

	r, err := http.NewRequestWithContext(env.Cxt, method, ts.URL+"/v1"+pathAndQuery, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	for key, value := range headers {
		r.Header.Set(key, value)
	}

	if authenticate {
		r.Header.Set("Authorization", "Bearer accessToken")
	}

	return sendHTTPRequest(r, response)
}

func sendHTTPRequest(request *http.Request, response interface{}) error {
	w, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	defer w.Body.Close()

	body, err := io.ReadAll(w.Body)
	if err != nil {
		return err
	}

	if w.StatusCode != http.StatusOK {
		return fmt.Errorf("got http status %d : %s", w.StatusCode, string(body))
	}

	res := response

	err = json.Unmarshal(body, &res)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal response")
	}

	return nil
}
