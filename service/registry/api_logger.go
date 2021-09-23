package registry

import (
	"github.com/coretrix/hitrix/service"
	apilogger "github.com/coretrix/hitrix/service/component/api_logger"
	"github.com/sarulabs/di"
)

func ServiceProviderAPILogger(entity apilogger.ILogEntity) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.APILoggerService,
		Build: func(ctn di.Container) (interface{}, error) {
			return apilogger.NewMysqlAPILogger(entity), nil
		},
	}
}
