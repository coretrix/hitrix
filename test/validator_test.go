package main

import (
	"testing"
	"time"

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

func TestTimestampGteNowValidation(t *testing.T) {
	valueAfterNow := time.Now().Add(time.Duration(12312312312))
	valueNow := time.Now()
	err := helper.NewValidator().Validate(helper.GetTimestamp(&valueNow), "timestamp_gte_now")
	assert.NotNil(t, err)

	err = helper.NewValidator().Validate(helper.GetTimestamp(&valueAfterNow), "timestamp_gte_now")
	assert.Nil(t, err)
}
