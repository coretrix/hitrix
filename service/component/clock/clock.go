package clock

import (
	"time"
)

type IClock interface {
	Now() time.Time
	NowPointer() *time.Time
}
