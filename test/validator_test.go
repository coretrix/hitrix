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

func TestTimestampGteTodayValidation(t *testing.T) {
	now := time.Now()
	valueBeforeNow := time.Now().AddDate(0, 0, -1)
	err := binding.NewValidator().Validate(helper.GetTimestamp(&valueBeforeNow), "timestamp_gte_today")
	assert.NotNil(t, err)

	err = binding.NewValidator().Validate(helper.GetTimestamp(&now), "timestamp_gte_today")
	assert.Nil(t, err)
}
