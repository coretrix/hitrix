package helper_test

import (
	"testing"
	"time"

	"github.com/coretrix/hitrix/pkg/helper"

	"github.com/stretchr/testify/assert"
)

func TestGetTimeDifference(t *testing.T) {
	from := time.Date(1986, 3, 7, 17, 30, 15, 0, time.UTC)
	to := time.Date(2020, 5, 2, 19, 35, 10, 0, time.UTC)

	timeDifference := helper.GetTimeDifference(from, to)

	assert.Equal(t, 34, timeDifference.Years)
	assert.Equal(t, 1, timeDifference.Months)
	assert.Equal(t, 26, timeDifference.Days)
	assert.Equal(t, 2, timeDifference.Hours)
	assert.Equal(t, 4, timeDifference.Minutes)
	assert.Equal(t, 55, timeDifference.Seconds)

	timeDifference = helper.GetTimeDifference(to, from)

	assert.Equal(t, 34, timeDifference.Years)
	assert.Equal(t, 1, timeDifference.Months)
	assert.Equal(t, 26, timeDifference.Days)
	assert.Equal(t, 2, timeDifference.Hours)
	assert.Equal(t, 4, timeDifference.Minutes)
	assert.Equal(t, 55, timeDifference.Seconds)
}

func TestGetWeekDay(t *testing.T) {
	location := time.Now().Location()

	//Use 2020-06-01 as base because Monday is first day of the month
	for day := 1; day <= 7; day++ {
		weekDay := helper.GetWeekDay(time.Date(2020, 6, day, 0, 0, 0, 0, location))
		assert.Equal(t, uint8(day), weekDay)
	}
}

func TestGetSecondsDifference(t *testing.T) {
	location := time.Now().Location()

	from := time.Date(2020, 6, 1, 0, 0, 0, 0, location)
	to := time.Date(2020, 6, 1, 1, 0, 0, 0, location)

	assert.Equal(t, float64(3600), helper.GetSecondsDifference(to, from))
}
