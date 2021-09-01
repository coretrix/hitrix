package helper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

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
func TestTruncateTime(t *testing.T) {
	t.Run("truncate utc", func(t *testing.T) {
		now := time.Date(2021, 6, 10, 11, 12, 13, 14, time.UTC)
		truncated := helper.TruncateTime(now)
		assert.Equal(t, 2021, truncated.Year())
		assert.Equal(t, time.Month(6), truncated.Month())
		assert.Equal(t, 10, truncated.Day())
		assert.Equal(t, 0, truncated.Hour())
		assert.Equal(t, 0, truncated.Minute())
		assert.Equal(t, 0, truncated.Minute())
		assert.Equal(t, 0, truncated.Second())
		assert.Equal(t, 0, truncated.Nanosecond())
		assert.Equal(t, time.UTC, truncated.Location())
	})
	t.Run("truncate tehran", func(t *testing.T) {
		tehran, err := time.LoadLocation("Asia/Tehran")
		require.Nil(t, err)

		now := time.Date(2021, 6, 10, 11, 12, 13, 14, tehran)
		truncated := helper.TruncateTime(now)
		assert.Equal(t, 2021, truncated.Year())
		assert.Equal(t, time.Month(6), truncated.Month())
		assert.Equal(t, 10, truncated.Day())
		assert.Equal(t, 0, truncated.Hour())
		assert.Equal(t, 0, truncated.Minute())
		assert.Equal(t, 0, truncated.Minute())
		assert.Equal(t, 0, truncated.Second())
		assert.Equal(t, 0, truncated.Nanosecond())
		assert.Equal(t, tehran, truncated.Location())
	})
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
