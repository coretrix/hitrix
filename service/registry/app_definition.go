package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/sarulabs/di"
)

func ServiceProviderApp(app *app.App) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.AppService,
		Build: func(ctn di.Container) (interface{}, error) {
			return app, nil
		},
	}
}
