package social

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type googleUserData struct {
	Email     string `json:"email"`
	LastName  string `json:"family_name"`
	FirstName string `json:"given_name"`
	Picture   string `json:"picture"`
}

type Google struct {
}

func (p *Google) GetUserData(token string) (*UserData, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=" +
		url.QueryEscape(token))
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Status: " + resp.Status)
	}

	// read all response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	googleUser := &googleUserData{}

	err = json.Unmarshal(body, googleUser)
	if err != nil {
		panic(err.Error())
	}

	return &UserData{
		FirstName: googleUser.FirstName,
		LastName:  googleUser.LastName,
		Avatar:    googleUser.Picture,
		Email:     googleUser.Email,
	}, nil
}
