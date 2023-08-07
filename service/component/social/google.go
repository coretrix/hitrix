package social

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type googleUserData struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	LastName  string `json:"family_name"`
	FirstName string `json:"given_name"`
	Picture   string `json:"picture"`
}

type Google struct {
}

func (p *Google) GetUserData(_ context.Context, token string, _ bool) (*UserData, error) {
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	googleUser := &googleUserData{}

	err = json.Unmarshal(body, googleUser)
	if err != nil {
		panic(err.Error())
	}

	return &UserData{
		ID:        googleUser.ID,
		FirstName: googleUser.FirstName,
		LastName:  googleUser.LastName,
		Avatar:    googleUser.Picture,
		Email:     googleUser.Email,
	}, nil
}
