package app

import (
	"context"
	"time"

	"github.com/latolukasz/beeorm"
)

type IScript interface {
	Description() string
	Run(ctx context.Context, exit IExit, ormService *beeorm.Engine)
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
