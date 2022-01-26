package hitrix

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/latolukasz/beeorm"
	"github.com/ryanuber/columnize"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
)

type BackgroundProcessor struct {
	Server *Hitrix
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

func (processor *BackgroundProcessor) RunScript(s app.IScript) {
	options, isOptional := s.(app.Optional)

	if isOptional {
		if !options.Active() {
			log.Print("IScript not active. Exiting.")
			return
		}
	}
	interval, isInterval := s.(app.Interval)
	_, isInfinity := s.(app.Infinity)

	go func() {
		for {
			log.Println("Run script - " + s.Description())

			valid := processor.runScript(s)

			if valid {
				log.Println("Successfully executed script - " + s.Description())
			} else {
				log.Println("Failed script - " + s.Description())
				time.Sleep(time.Second * 10)
				continue
			}

			if isInfinity {
				log.Println("Infinity - " + s.Description())
				select {}
			}

			if !isInterval {
				log.Println("Finished - " + s.Description())
				processor.Server.done <- true
				break
			}

			log.Println("Sleep for " + fmt.Sprint(interval.Interval()) + " seconds - " + s.Description())

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
			def := service.GetServiceRequired(defCode).(app.IScript)
			options := make([]string, 0)
			interval, is := def.(app.Interval)
			if is {
				options = append(options, "interval")
				duration := "every " + interval.Interval().String()
				_, is := def.(app.IntervalOptional)
				if is {
					duration += " with condition"
				}
				options = append(options, duration)
			}

			if def.Unique() {
				options = append(options, "unique")
			}
			optional, is := def.(app.Optional)
			if is {
				options = append(options, "optional")
				if optional.Active() {
					options = append(options, "active")
				} else {
					options = append(options, "inactive")
				}
			}
			intermediate, is := def.(app.Intermediate)
			if is && intermediate.IsIntermediate() {
				options = append(options, "intermediate")
			}
			output = append(output, strings.Join([]string{defCode, strings.Join(options, ","), def.Description()}, " | "))
		}
		_, _ = os.Stdout.WriteString(columnize.SimpleFormat(output) + "\n")
	}
}

func (processor *BackgroundProcessor) runScript(s app.IScript) bool {
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

				service.DI().ErrorLogger().LogError(message)

				valid = false
			}
		}()

		appService := service.DI().App()

		s.Run(appService.GlobalContext, &exit{s: processor.Server})

		return valid
	}()
}

func (processor *BackgroundProcessor) RunAsyncOrmConsumer() {
	ormService := service.DI().OrmEngine()
	appService := service.DI().App()

	GoroutineWithRestart(func() {
		asyncConsumer := beeorm.NewBackgroundConsumer(ormService)
		for {
			if asyncConsumer.Digest(appService.GlobalContext) {
				break
			}
			time.Sleep(time.Second * 30)
		}
	})
}
