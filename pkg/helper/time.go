package helper

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/icza/gox/timex"
)

const (
	TimeLayoutY                 = "2006"
	TimeLayoutYM                = "2006-01"
	TimeLayoutYMD               = "2006-01-02"
	TimeLayoutYMDHM             = "2006-01-02 15:04"
	TimeLayoutYMDHMS            = "2006-01-02 15:04:05"
	TimeLayoutHM                = "15:04"
	TimeLayoutMDYYYYHMMSSPM     = "1/2/2006 3:04:05 PM"
	TimeLayoutNoSepYYYYMMDDHHMM = "200601021504"
	TimeLayoutTextMD            = "January 2"
)

const Second = 1
const Minute = 60
const Hour = 3600
const Day = 86400
const Week = 604800
const Month = 2592000
const Year = 31104000

type TimeDifference struct {
	Years, Months, Days, Hours, Minutes, Seconds int
}

func GetTimeDifference(from, to time.Time) *TimeDifference {
	years, months, days, hours, minutes, seconds := timex.Diff(from, to)

	return &TimeDifference{
		Years:   years,
		Months:  months,
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
		Seconds: seconds,
	}
}

func GetSecondsDifference(from, to time.Time) float64 {
	var delta time.Duration
	if from.After(to) {
		delta = from.Sub(to)
	} else {
		delta = to.Sub(from)
	}

	return delta.Seconds()
}

func GetWeekDay(dateTime time.Time) uint8 {
	weekday := dateTime.Weekday()

	if weekday == 0 {
		weekday = 7
	}

	return uint8(weekday)
}

func ValidateAndParseTimeRange(startsAt, endsAt, fromName, toName string, startsAtGTENow bool, c *gin.Context) (*time.Time, *time.Time, error) {
	var from *time.Time
	var to *time.Time

	startDate, err := time.Parse(TimeLayoutYMDHM, startsAt)
	if err == nil {
		from = &startDate
	}

	endDate, err := time.Parse(TimeLayoutYMDHM, endsAt)
	if err == nil {
		to = &endDate
	}

	if startsAtGTENow && from != nil && from.Before(time.Now()) {
		return nil, nil, errors.New("it should be greater or equal to now")
	}

	if from != nil && to != nil && from.After(*to) {
		return nil, nil, fmt.Errorf("it should be greater than %v", fromName)
	}

	return from, to, nil
}

func GetTimeDifferenceHumanBySeconds(seconds float64) string {
	return GetTimeDifferenceHuman(time.Now(), time.Now().Add(time.Duration(seconds)*time.Second))
}

func GetTimeDifferenceHuman(startDate, endDate time.Time) string {
	timeDifference := GetTimeDifference(startDate, endDate)

	duration := ""
	if timeDifference.Years > 0 {
		duration += strconv.Itoa(timeDifference.Years) + "y "
	}

	if timeDifference.Months > 0 || timeDifference.Years > 0 {
		duration += strconv.Itoa(timeDifference.Months) + "m "
	}

	if timeDifference.Days > 0 || timeDifference.Months > 0 || timeDifference.Years > 0 {
		duration += strconv.Itoa(timeDifference.Days) + "d "
	}

	if timeDifference.Hours > 0 || timeDifference.Days > 0 || timeDifference.Months > 0 || timeDifference.Years > 0 {
		duration += strconv.Itoa(timeDifference.Hours) + "h "
	}

	if timeDifference.Minutes > 0 || timeDifference.Hours > 0 || timeDifference.Days > 0 || timeDifference.Months > 0 || timeDifference.Years > 0 {
		duration += strconv.Itoa(timeDifference.Minutes) + "min "
	}

	duration += strconv.Itoa(timeDifference.Seconds) + "s"

	return duration
}

func GetTimestamp(t *time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
