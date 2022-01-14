package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/password"
)

func ServiceProviderPassword() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.PasswordService,
		Build: func(ctn di.Container) (interface{}, error) {
			return &password.Manager{}, nil
		},
	}
}
