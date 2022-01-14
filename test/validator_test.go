package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coretrix/hitrix/pkg/helper"
)

func TestBasicValidation(t *testing.T) {
	err := helper.NewValidator().Validate("no-email-string", "email")
	assert.NotNil(t, err)

	err = helper.NewValidator().Validate("awesome-dude@awesome-com.com", "email")
	assert.Nil(t, err)
}

func TestCountryCodeValidation(t *testing.T) {
	err := helper.NewValidator().Validate("ABCDEFG", "country_code_custom")
	assert.NotNil(t, err)

	err = helper.NewValidator().Validate("SE", "country_code_custom")
	assert.Nil(t, err)
}
