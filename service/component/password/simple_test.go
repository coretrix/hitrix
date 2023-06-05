package password_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coretrix/hitrix/service/component/password"
)

func TestHashPasswordTrue(t *testing.T) {
	passwordService := &password.SimpleManager{}

	hash, err := passwordService.HashPassword("Str0NGPa$$W0rD!")

	assert.NoError(t, err, "Cannot create hash")
	assert.Equal(t, "eh71ZMSd5oCpYTaazon8jc53bo0sMiWSPmPVuMVB9mU=", hash, "Hash is not valid")
}

func TestHashPasswordFalse(t *testing.T) {
	passwordService := &password.SimpleManager{}

	hash, err := passwordService.HashPassword("Str0NGPa$$W0rD!1")

	assert.NoError(t, err, "Cannot create hash")
	assert.NotEqual(t, "eh71ZMSd5oCpYTaazon8jc53bo0sMiWSPmPVuMVB9mU=", hash)
}

func TestVerifyPasswordTrue(t *testing.T) {
	passwordService := &password.SimpleManager{}

	assert.True(t, passwordService.VerifyPassword("Str0NGPa$$W0rD!", "eh71ZMSd5oCpYTaazon8jc53bo0sMiWSPmPVuMVB9mU="))
}

func TestVerifyPasswordFalse(t *testing.T) {
	passwordService := &password.SimpleManager{}

	assert.False(t, passwordService.VerifyPassword("Str0NGPa$$W0rD!", ""))
}
