package clock

import (
	"time"

	"github.com/xorcare/pointer"
)

type SysClock struct{}

func (c *SysClock) Now() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}

func (c *SysClock) NowPointer() *time.Time {
	return pointer.Time(time.Now().UTC().Truncate(time.Second))
}
