package main

import (
	"testing"

	"github.com/coretrix/hitrix/service/component/password"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	passwordService := &password.Password{}

	hash, err := passwordService.HashPassword("Str0NGPa$$W0rD!")

	assert.NoError(t, err, "Cannot create hash")
	assert.Equal(t, "eh71ZMSd5oCpYTaazon8jc53bo0sMiWSPmPVuMVB9mU=", hash, "Hash is not valid")
}

func TestVerifyPassword(t *testing.T) {
	passwordService := &password.Password{}

	passwordInput := "Str0NGPa$$W0rD!"
	valid := passwordService.VerifyPassword(passwordInput, "eh71ZMSd5oCpYTaazon8jc53bo0sMiWSPmPVuMVB9mU=")

	assert.True(t, valid)

	valid = passwordService.VerifyPassword(passwordInput, "")
	assert.False(t, valid)
}
