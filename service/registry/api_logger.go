package registry

import (
	"github.com/coretrix/hitrix/service"
	apilogger "github.com/coretrix/hitrix/service/component/api_logger"
	"github.com/sarulabs/di"
)

func APILogger(entity apilogger.LogEntity) *service.Definition {
	return &service.Definition{
		Name:   service.APILoggerService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return apilogger.NewAPILog(entity), nil
		},
	}
}
