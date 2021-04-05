package clock

import (
	"time"

	"github.com/xorcare/pointer"
)

// Clock provides an abstraction over the time package
type Clock interface {
	Now() time.Time
	NowPointer() *time.Time
}

// SysClock is wrapper over the standard time package
type SysClock struct{}

// Now returns time.Now
func (c *SysClock) Now() time.Time {
	return time.Now().UTC()
}

// Now returns *time.Now
func (c *SysClock) NowPointer() *time.Time {
	return pointer.Time(time.Now().UTC())
}
