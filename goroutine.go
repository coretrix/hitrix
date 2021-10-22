package hitrix

import (
	"time"

	"github.com/coretrix/hitrix/service"
)

func Goroutine(fn func()) {
	go routine(fn, false)
}

func GoroutineWithRestart(fn func()) {
	go routine(fn, true)
}

func routine(fn func(), autoRestart bool) {
	defer func() {
		if r := recover(); r != nil {
			service.DI().ErrorLogger().LogError(r)

			if autoRestart {
				time.Sleep(time.Second)

				go routine(fn, true)
			}
		}
	}()

	fn()
}
