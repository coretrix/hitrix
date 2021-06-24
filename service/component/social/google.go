package social

import (
	"encoding/json"
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

	if err != nil {
		return nil, err
	}

	// read all response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	googleUserData := &googleUserData{}
	err = json.Unmarshal(body, googleUserData)
	if err != nil {
		panic(err.Error())
	}

	return &UserData{
		FirstName: googleUserData.FirstName,
		LastName:  googleUserData.LastName,
		Avatar:    googleUserData.Picture,
		Email:     googleUserData.Email,
	}, nil
}
