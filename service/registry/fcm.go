package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/fcm"
	"github.com/sarulabs/di"
)

func ServiceProviderFCM() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.FCMService,
		Build: func(ctn di.Container) (interface{}, error) {
			appService := ctn.Get(service.AppService).(*app.App)
			return fcm.NewFCM(appService.GlobalContext)
		},
	}
}
