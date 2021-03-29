package registry

import (
	"github.com/coretrix/hitrix/service"
	apilogger "github.com/coretrix/hitrix/service/component/api_logger"
	"github.com/latolukasz/orm"
	"github.com/sarulabs/di"
)

func APILogger(entity apilogger.LogEntity) *service.Definition {
	//if reflect.ValueOf(entity).Kind() == reflect.Ptr {
	//	panic("passed entity should not be reference")
	//}

	return &service.Definition{
		Name:   service.APILoggerService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			subContainer, err := ctn.SubContainer()
			if err != nil {
				return nil, err
			}

			ormService := subContainer.Get("orm_engine_request").(*orm.Engine)
			return apilogger.NewAPILog(ormService, entity), nil
		},
	}
}
