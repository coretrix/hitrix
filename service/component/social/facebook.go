package social

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type facebookUserData struct {
	FirstName string
	LastName  string
	Avatar    string
	Email     string
}

type Facebook struct {
}

func (p *Facebook) GetUserData(_ context.Context, token string) (*UserData, error) {
	resp, err := http.Get("https://graph.facebook.com/me?access_token=" +
		url.QueryEscape(token))
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(resp.Body)

	if err != nil {
		return nil, err
	}

	// read all response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	facebookUser := &facebookUserData{}

	err = json.Unmarshal(body, facebookUser)
	if err != nil {
		panic(err.Error())
	}

	return &UserData{
		FirstName: facebookUser.FirstName,
		LastName:  facebookUser.LastName,
		Avatar:    facebookUser.Avatar,
		Email:     facebookUser.Email,
	}, nil
}

func (p *Facebook) SetIsAndroid(_ bool) {}
