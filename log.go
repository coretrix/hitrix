package hitrix

import (
	apexLog "github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di"
)

type LogRequestFieldProvider func(ctx *gin.Context) apexLog.Fielder
type LogFieldProvider func() apexLog.Fielder

type RequestLog struct {
	providers []LogRequestFieldProvider
	entry     apexLog.Interface
}

func serviceLogGlobal() *ServiceDefinition {
	return &ServiceDefinition{
		Name:   "log",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			l := apexLog.WithFields(&apexLog.Fields{"app": DIC().App().Name()})
			key := "_log_providers"
			_, has := ctn.Definitions()[key]
			if has {
				providers := ctn.Get(key).([]LogFieldProvider)
				for _, fields := range providers {
					l = l.WithFields(fields())
				}
			}
			return l, nil
		},
	}
}

func serviceLogForRequest() *ServiceDefinition {
	return &ServiceDefinition{
		Name:   "log_request",
		Global: false,
		Build: func(ctn di.Container) (interface{}, error) {
			key := "_log_request_providers"
			_, has := ctn.Definitions()[key]
			var l *RequestLog
			if has {
				providers := ctn.Get(key).([]LogRequestFieldProvider)
				l = &RequestLog{providers: providers}
			} else {
				l = &RequestLog{}
			}
			return l, nil
		},
	}
}

func (l *RequestLog) Log(ctx *gin.Context) apexLog.Interface {
	if l.entry == nil {
		entry := apexLog.WithFields(&apexLog.Fields{"app": DIC().App().Name()})
		for _, p := range l.providers {
			fields := p(ctx)
			if fields != nil {
				entry = entry.WithFields(fields)
			}
		}
		l.entry = entry
	}
	return l.entry
}

func (r *Registry) initializeLog() {
	if DIC().App().IsInProdMode() {
		h, has := GetServiceOptional("log_handler")
		if has {
			apexLog.SetHandler(h.(apexLog.Handler))
		} else {
			apexLog.SetHandler(json.Default)
		}
		apexLog.SetLevel(apexLog.WarnLevel)
	} else {
		apexLog.SetHandler(text.Default)
		apexLog.SetLevel(apexLog.DebugLevel)
	}
}
