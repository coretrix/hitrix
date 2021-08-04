package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/jwt"
	"github.com/sarulabs/di"
)

func ServiceProviderJWT() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.JWTService,
		Build: func(ctn di.Container) (interface{}, error) {
			return &jwt.JWT{}, nil
		},
	}
}
