package hitrix

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/latolukasz/orm"

	"github.com/coretrix/hitrix/service"

	"github.com/ryanuber/columnize"
)

type Script interface {
	Description() string
	Run(ctx context.Context, exit Exit)
	Unique() bool
}

type BackgroundProcessor struct {
	Server *Hitrix
}

type Exit interface {
	Valid()
	Error()
	Custom(exitCode int)
}

type exit struct {
	s *Hitrix
}

func (e *exit) Custom(exitCode int) {
	e.s.exit <- exitCode
}

func (e *exit) Valid() {
	e.Custom(0)
}

func (e *exit) Error() {
	e.Custom(1)
}

type ScriptInfinity interface {
	Infinity() bool
}

type ScriptInterval interface {
	Interval() time.Duration
}

type ScriptIntervalOptional interface {
	IntervalActive() bool
}

type ScriptIntermediate interface {
	IsIntermediate() bool
}

type ScriptOptional interface {
	Active() bool
}

func (processor *BackgroundProcessor) RunScript(script Script) {
	options, isOptional := script.(ScriptOptional)

	if isOptional {
		if !options.Active() {
			log.Print("Script not active. Exiting.")
			return
		}
	}
	interval, isInterval := script.(ScriptInterval)
	_, isInfinity := script.(ScriptInfinity)

	go func() {
		for {
			valid := processor.runScript(script)
			if isInfinity {
				select {}
			}

			if !isInterval {
				processor.Server.done <- true
				break
			}

			//TODO
			if !valid {
				log.Print("Error in last run.")
			}

			time.Sleep(interval.Interval())
		}
	}()
	processor.Server.await()
}

func listScrips() {
	scripts := service.DI().App().Scripts
	if len(scripts) > 0 {
		output := []string{
			"NAME | OPTIONS | DESCRIPTION ",
		}
		for _, defCode := range scripts {
			def := service.GetServiceRequired(defCode).(Script)
			options := make([]string, 0)
			interval, is := def.(ScriptInterval)
			if is {
				options = append(options, "interval")
				duration := "every " + interval.Interval().String()
				_, is := def.(ScriptIntervalOptional)
				if is {
					duration += " with condition"
				}
				options = append(options, duration)
			}

			if def.Unique() {
				options = append(options, "unique")
			}
			optional, is := def.(ScriptOptional)
			if is {
				options = append(options, "optional")
				if optional.Active() {
					options = append(options, "active")
				} else {
					options = append(options, "inactive")
				}
			}
			intermediate, is := def.(ScriptIntermediate)
			if is && intermediate.IsIntermediate() {
				options = append(options, "intermediate")
			}
			output = append(output, strings.Join([]string{defCode, strings.Join(options, ","), def.Description()}, " | "))
		}
		_, _ = os.Stdout.WriteString(columnize.SimpleFormat(output) + "\n")
	}
}

func (processor *BackgroundProcessor) runScript(script Script) bool {
	return func() bool {
		valid := true
		defer func() {
			if err := recover(); err != nil {
				var message string
				asErr, is := err.(error)
				if is {
					message = asErr.Error()
				} else {
					message = fmt.Sprint(err)
				}

				errorLogger, has := service.DI().ErrorLogger()
				if has {
					errorLogger.LogError(message + "\n" + string(debug.Stack()))
				}
				valid = false
			}
		}()
		script.Run(processor.Server.ctx, &exit{s: processor.Server})
		return valid
	}()
}

func (processor *BackgroundProcessor) RunAsyncOrmConsumer() {
	ormService, has := service.DI().OrmEngine()
	if !has {
		panic("Orm is not registered")
	}

	go func() {
		asyncConsumer := orm.NewBackgroundConsumer(ormService)
		asyncConsumer.Digest(processor.Server.ctx)
	}()
}
