package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/setting"
)

func ServiceProviderSetting() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.SettingService,
		Build: func(ctn di.Container) (interface{}, error) {
			return setting.NewSettingService(), nil
		},
	}
}
