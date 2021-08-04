package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/password"
	"github.com/sarulabs/di"
)

func ServiceProviderPassword() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.PasswordService,
		Build: func(ctn di.Container) (interface{}, error) {
			return &password.Password{}, nil
		},
	}
}
