package main

import (
	"testing"

	"github.com/coretrix/hitrix/pkg/helper"
	"github.com/stretchr/testify/assert"
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
