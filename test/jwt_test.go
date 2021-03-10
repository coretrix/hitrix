package main

import (
	"fmt"
	"testing"
	"time"

	jwt2 "github.com/coretrix/hitrix/service/component/jwt"

	"github.com/stretchr/testify/assert"
)

func TestJWTCreation(t *testing.T) {
	expectedJWT := "eyJhbGdvIjoiSFMyNTYiLCJ0eXBlIjoiSldUIn0=" +
		".eyJleHAiOiIxNTc5NTE0NDAwIiwiaXNzIjoiYmx1ZWxvZyIsInN1YiI6IlVzZXIifQ==.vZPl0FGz1opXMrrcEC6pFdFIQ3I10WHxsToB12CQ1nA="

	headers := map[string]string{
		"algo": "HS256",
		"type": "JWT",
	}

	payload := map[string]string{
		"exp": "1579514400",
		"iss": "bluelog",
		"sub": "User",
	}

	jwt := &jwt2.JWT{}
	jwtToken, err := jwt.EncodeJWT("mynewsecret", headers, payload)

	assert.NoError(t, err, "An error occured while creating JWT!")
	assert.Equal(t, expectedJWT, jwtToken, "Cannot create jwt token!")
}

func TestVerifyJWT(t *testing.T) {
	headers := map[string]string{
		"algo": "HS256",
		"type": "JWT",
	}

	payload := map[string]string{
		"exp": fmt.Sprintf("%v", time.Now().Unix()),
		"iss": "bluelog",
		"sub": "User",
	}

	jwt := &jwt2.JWT{}

	jwtToken, _ := jwt.EncodeJWT("mynewsecret", headers, payload)

	err := jwt.VerifyJWT("mynewsecret", jwtToken, 72000)

	var msg string
	if err != nil {
		msg = err.Error()
	}
	assert.NoError(t, err, msg)
}

func TestVerifyJWTExpired(t *testing.T) {
	headers := map[string]string{
		"algo": "HS256",
		"type": "JWT",
	}

	payload := map[string]string{
		"exp": fmt.Sprintf("%v", time.Now().Unix()-1000),
		"iss": "bluelog",
		"sub": "User",
	}

	jwt := &jwt2.JWT{}

	jwtToken, _ := jwt.EncodeJWT("mynewsecret", headers, payload)

	err := jwt.VerifyJWT("mynewsecret", jwtToken, 72000)
	assert.NoError(t, err, "Cannot verify expired jwt token")
}

func TestExtractPayload(t *testing.T) {
	expectedTokenSubject := "User"
	headers := map[string]string{
		"algo": "HS256",
		"type": "JWT",
	}

	payload := map[string]string{
		"exp": fmt.Sprintf("%v", time.Now().Unix()),
		"iss": "bluelog",
		"sub": expectedTokenSubject,
	}

	jwt := &jwt2.JWT{}

	jwtToken, _ := jwt.EncodeJWT("mynewsecret", headers, payload)

	extractedPayload, err := jwt.VerifyJWTAndGetPayload("mynewsecret", jwtToken, 72000)
	assert.NoError(t, err, "Cannot verify expired jwt token")
	assert.Equal(t, expectedTokenSubject, extractedPayload["sub"], "Cannot create jwt token!")
}
