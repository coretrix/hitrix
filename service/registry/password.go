package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/password"
	"github.com/sarulabs/di"
)

func ServiceProviderPassword() *service.Definition {
	return &service.Definition{
		Name:   service.PasswordService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return &password.Password{}, nil
		},
	}
}
