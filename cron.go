package hitrix

import (
	"log"
	"time"
)

func (h *Hitrix) StartCron(cron func(), beforeExitFunc func()) {
	go h.cronProcess(cron)

	<-h.done

	if beforeExitFunc != nil {
		beforeExitFunc()
	}

	log.Println("killing...")
	h.cancel()
	time.Sleep(time.Millisecond * 300)
}

func (h *Hitrix) cronProcess(cron func()) {
	defer func() {
		if r := recover(); r != nil {
			errorLogger, has := DIC().ErrorLogger()
			if has {
				errorLogger.LogRecover(r)
			} else {
				log.Println(r.(string))
			}
			time.Sleep(500)
			h.cronProcess(cron)
		}
	}()

	cron()
}
