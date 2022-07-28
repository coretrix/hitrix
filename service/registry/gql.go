package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/gql"
	"github.com/coretrix/hitrix/service/component/localize"
)

func ServiceProviderGql() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.GQLService,
		Build: func(ctn di.Container) (interface{}, error) {
			localizeService, err := ctn.SafeGet(service.LocalizeService)
			if err == nil {
				return gql.NewGqlService(localizeService.(localize.ILocalizer)), nil
			}

			return gql.NewGqlService(nil), nil
		},
	}
}
