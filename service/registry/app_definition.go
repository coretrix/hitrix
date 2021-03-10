package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/sarulabs/di"
)

func ServiceApp(app *app.App) *service.Definition {
	return &service.Definition{
		Name:   service.AppService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return app, nil
		},
	}
}
