package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/jwt"
	"github.com/sarulabs/di"
)

func ServiceProviderJWT() *service.Definition {
	return &service.Definition{
		Name:   "jwt",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return &jwt.JWT{}, nil
		},
	}
}
