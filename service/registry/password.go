package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/password"
)

func ServiceProviderPassword(newFunc password.NewPasswordManagerFunc) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.PasswordService,
		Build: func(ctn di.Container) (interface{}, error) {
			return newFunc(ctn.Get(service.ConfigService).(config.IConfig)), nil
		},
	}
}
