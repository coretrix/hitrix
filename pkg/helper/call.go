package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Call helper for api calls
// TODO : make service and use debug service to
func Call(ctx context.Context,
	method,
	url string,
	headers map[string]string,
	timeout time.Duration,
	payload interface{},
	cookies []*http.Cookie) ([]byte, http.Header, int, error) {
	d, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	var b io.Reader
	b = bytes.NewReader(d)
	method = strings.ToUpper(method)
	if StringInArray(method, "GET", "DELETE") {
		b = nil
	}

	r, err := http.NewRequest(method, url, b)
	if err != nil {
		return nil, nil, 0, errors.New("error while creating request")
	}

	for i := range headers {
		r.Header.Set(i, headers[i])
	}

	for i := range cookies {
		r.AddCookie(cookies[i])
	}

	nCtx, cnl := context.WithTimeout(ctx, timeout)
	defer cnl()

	resp, err := http.DefaultClient.Do(r.WithContext(nCtx))
	if err != nil {
		return nil, nil, 0, errors.New("error in return response")
	}

	data, err := ioutil.ReadAll(resp.Body)

	defer func() {
		_ = resp.Body.Close()
	}()

	if err != nil {
		return nil, nil, resp.StatusCode, errors.New("error in reading response")
	}
	return data, resp.Header, resp.StatusCode, nil
}
