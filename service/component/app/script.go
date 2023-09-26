package app

import (
	"context"
	"time"
)

type IScript interface {
	Description() string
	Run(ctx context.Context, exit IExit)
	Unique() bool
}

type IExit interface {
	Valid()
	Error()
	Custom(exitCode int)
}

type Infinity interface {
	Infinity() bool
}

type Interval interface {
	Interval() time.Duration
}

type IntervalOptional interface {
	IntervalActive() bool
}

type Intermediate interface {
	IsIntermediate() bool
}

type Optional interface {
	Active() bool
}
