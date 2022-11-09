package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/coretrix/hitrix/pkg/binding"
	"github.com/coretrix/hitrix/pkg/helper"
)

func TestBasicValidation(t *testing.T) {
	err := binding.NewValidator().Validate("no-email-string", "email")
	assert.NotNil(t, err)

	err = binding.NewValidator().Validate("awesome-dude@awesome-com.com", "email")
	assert.Nil(t, err)
}

func TestCountryCodeValidation(t *testing.T) {
	err := binding.NewValidator().Validate("ABCDEFG", "country_code_custom")
	assert.NotNil(t, err)

	err = binding.NewValidator().Validate("SE", "country_code_custom")
	assert.Nil(t, err)
}

func TestTimestampGteNowValidation(t *testing.T) {
	valueAfterNow := time.Now().AddDate(0, 0, 1)
	valueBeforeNow := time.Now().AddDate(0, 0, -1)
	err := binding.NewValidator().Validate(helper.GetTimestamp(&valueBeforeNow), "timestamp_gte_now")
	assert.NotNil(t, err)

	err = binding.NewValidator().Validate(helper.GetTimestamp(&valueAfterNow), "timestamp_gte_now")
	assert.Nil(t, err)
}
