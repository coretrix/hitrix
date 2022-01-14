package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/jwt"
)

func ServiceProviderJWT() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.JWTService,
		Build: func(ctn di.Container) (interface{}, error) {
			return &jwt.JWT{}, nil
		},
	}
}
