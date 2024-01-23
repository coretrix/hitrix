package hitrix

import (
	"encoding/json"
	"expvar"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/TwiN/go-color"
	"github.com/latolukasz/beeorm"
	"github.com/ryanuber/columnize"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/pkg/helper"
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
	interval, isInterval := s.(app.Interval)
	_, isInfinity := s.(app.Infinity)

	if !isInterval && !isInfinity {
		log.Println("Failed script - "+s.Description(), " - it must implement either Interval or Infinity interface")

		return
	}

	options, isOptional := s.(app.Optional)

	if isOptional {
		if !options.Active() {
			log.Print(s.Description() + "  script not active. Exiting.")

			return
		}
	}

	var ticker *time.Ticker
	if isInterval {
		ticker = time.NewTicker(interval.Interval())
	}

	go func() {
		if isInfinity {
			log.Println("Infinity - " + s.Description())

			for {
				valid := processor.run(s)
				if valid {
					processor.Server.done <- true

					break
				}

				time.Sleep(time.Second * 10)
			}
		}

		if !isInterval {
			service.DI().App().Add(1)
			defer service.DI().App().Done()

			processor.run(s)
			log.Println(color.InGreen("Finished script") + " - " + s.Description())
			processor.Server.done <- true

			return
		}

		if isInterval {
			service.DI().App().Add(1)
			processor.run(s)
			service.DI().App().Done()

			//nolint //a
			for {
				select {
				case <-ticker.C:
					service.DI().App().Add(1)
					processor.run(s)
					service.DI().App().Done()
				}
			}
		}
	}()
}

func (processor *BackgroundProcessor) run(s app.IScript) bool {
	log.Println(color.InGreen("Start script") + " - " + s.Description())

	valid := processor.runScript(s)

	if valid {
		_, isInfinity := s.(app.Infinity)
		if !isInfinity {
			log.Println(color.InGreen("Running script") + " - " + s.Description())
		}
	} else {
		log.Println(color.InRed("Failed script") + " - " + s.Description())
	}

	return valid
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

		ormService := service.DI().OrmEngine().Clone()
		ormService.SetLogMetaData("Script", s.Description())

		s.Run(appService.GlobalContext, &exit{s: processor.Server}, ormService)

		return valid
	}()
}

func (processor *BackgroundProcessor) RunAsyncOrmConsumer() {
	ormService := service.DI().OrmEngine().Clone()
	appService := service.DI().App()

	GoroutineWithRestart(func() {
		log.Println("starting orm background consumer")

		asyncConsumer := beeorm.NewBackgroundConsumer(ormService)
		for {
			if asyncConsumer.Digest(appService.GlobalContext) {
				log.Println("orm background consumer exited successfully")

				break
			}

			log.Println("orm background consumer can not obtain lock, sleeping for 30 seconds")
			time.Sleep(time.Second * 30)
		}
	})
}

type FieldProcessor map[string]func(value interface{}) float64

func BytesToMB(value interface{}) float64 {
	return value.(float64) / float64(1000000)
}

func NanoToMilli(value interface{}) float64 {
	return value.(float64) / float64(1000000)
}

func (processor *BackgroundProcessor) RunAsyncMetricsCollector(fieldProcessor FieldProcessor) {
	ormConfig := service.DI().OrmConfig()
	entities := ormConfig.GetEntities()

	if _, ok := entities["entity.MetricsEntity"]; !ok {
		panic("you should register MetricsEntity")
	}

	clockService := service.DI().Clock()
	configService := service.DI().Config()

	fields, hasFields := configService.Strings("metrics.fields")

	if !hasFields {
		panic("Metrics fields are required")
	}

	fieldsMap := map[string]struct{}{}
	for _, field := range fields {
		fieldsMap[field] = struct{}{}
	}

	intervalCollectorInMilli, hasIntervalCollector := configService.Int("metrics.interval_collector")

	if !hasIntervalCollector {
		intervalCollectorInMilli = 1000
	}

	countFlusher, hasCountFlusher := configService.Int("metrics.count_flusher")

	if !hasCountFlusher {
		countFlusher = 60
	}

	ticker := time.NewTicker(time.Duration(intervalCollectorInMilli) * time.Millisecond)
	appName := service.DI().App().Name
	counter := 0

	GoroutineWithRestart(func() {
		ormService := service.DI().OrmEngine().Clone()
		flusher := ormService.NewFlusher()

		log.Println("starting metrics collector")

		for range ticker.C {
			data := "{"

			expvar.Do(func(kv expvar.KeyValue) {
				if kv.Key == "memstats" {
					memStats := &map[string]interface{}{}

					err := json.Unmarshal([]byte(kv.Value.String()), memStats)
					if err != nil {
						panic(err)
					}

					for k, v := range *memStats {
						if _, ok := fieldsMap[k]; ok {
							if _, hasFieldProcessor := fieldProcessor[k]; hasFieldProcessor {
								data += fmt.Sprintf("%q: %f,", k, fieldProcessor[k](v))
							} else {
								data += fmt.Sprintf("%q: %f,", k, v)
							}
						}
					}

					if _, ok := fieldsMap["NumGoroutine"]; ok {
						data += fmt.Sprintf("\"NumGoroutine\": %d,", runtime.NumGoroutine())
					}
					data = strings.TrimSuffix(data, ",")
				}
			})

			data += "}"

			flusher.Track(&entity.MetricsEntity{
				AppName:   appName,
				Metrics:   data,
				CreatedAt: clockService.Now(),
			})

			counter++

			if counter == countFlusher {
				flusher.Flush()
				counter = 0
			}
		}
	})
}

func (processor *BackgroundProcessor) RunAsyncRequestLoggerCleaner() {
	ormConfig := service.DI().OrmConfig()
	entities := ormConfig.GetEntities()

	if _, ok := entities["entity.RequestLoggerEntity"]; !ok {
		panic("you should register RequestLoggerEntity")
	}

	configService := service.DI().Config()
	ttlInDays, has := configService.Int("metrics.ttl_in_days")

	if !has {
		ttlInDays = 30
	}

	GoroutineWithRestart(func() {
		ormService := service.DI().OrmEngine().Clone()

		log.Println("starting request logger cleaner")

		for {
			removeAllOldRequestLoggerRows(ormService, ttlInDays)

			log.Println("sleeping request logger cleaner")
			time.Sleep(time.Minute * 30)
		}
	})
}

func removeAllOldRequestLoggerRows(ormService *beeorm.Engine, ttlInDays int) {
	pager := beeorm.NewPager(1, 1000)

	for {
		where := beeorm.NewWhere("CreatedAt < ?", service.DI().Clock().Now().AddDate(0, 0, -ttlInDays).Format(helper.TimeLayoutYMDHMS))

		var requestLoggerEntities []*entity.RequestLoggerEntity
		ormService.Search(where, pager, &requestLoggerEntities)

		flusher := ormService.NewFlusher()

		for _, requestLoggerEntity := range requestLoggerEntities {
			flusher.Delete(requestLoggerEntity)
		}

		flusher.Flush()

		log.Printf("%d rows was removed", len(requestLoggerEntities))

		if len(requestLoggerEntities) < pager.PageSize {
			break
		}
	}
}

func (processor *BackgroundProcessor) RunAsyncMetricsCleaner() {
	ormConfig := service.DI().OrmConfig()
	entities := ormConfig.GetEntities()

	if _, ok := entities["entity.MetricsEntity"]; !ok {
		panic("you should register MetricsEntity")
	}

	configService := service.DI().Config()
	ttlInDays, has := configService.Int("metrics.ttl_in_days")

	if !has {
		ttlInDays = 30
	}

	GoroutineWithRestart(func() {
		ormService := service.DI().OrmEngine().Clone()

		log.Println("starting metrics cleaner")

		for {
			removeAllOldMetricsRows(ormService, ttlInDays)

			log.Println("sleeping metrics cleaner")
			time.Sleep(time.Minute * 30)
		}
	})
}

func removeAllOldMetricsRows(ormService *beeorm.Engine, ttlInDays int) {
	pager := beeorm.NewPager(1, 1000)

	for {
		where := beeorm.NewWhere("CreatedAt < ?", service.DI().Clock().Now().AddDate(0, 0, -ttlInDays).Format(helper.TimeLayoutYMDHMS))

		var metricsEntities []*entity.MetricsEntity
		ormService.Search(where, pager, &metricsEntities)

		flusher := ormService.NewFlusher()

		for _, requestLoggerEntity := range metricsEntities {
			flusher.Delete(requestLoggerEntity)
		}

		flusher.Flush()

		log.Printf("%d rows was removed", len(metricsEntities))

		if len(metricsEntities) < pager.PageSize {
			break
		}
	}
}
