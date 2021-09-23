package goroutine

import (
	"time"

	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
)

type Manager struct {
	ErrorLogger errorlogger.ErrorLogger
}

func NewGoroutineManager(errorLoggerService errorlogger.ErrorLogger) IGoroutine {
	return &Manager{ErrorLogger: errorLoggerService}
}

func (g *Manager) Goroutine(fn func()) {
	go g.routine(fn, false)
}

func (g *Manager) GoroutineWithRestart(fn func()) {
	go g.routine(fn, true)
}

func (g *Manager) routine(fn func(), autoRestart bool) {
	defer func() {
		if r := recover(); r != nil {
			g.ErrorLogger.LogError(r)

			if autoRestart {
				time.Sleep(time.Second)

				go g.routine(fn, true)
			}
		}
	}()

	fn()
}
