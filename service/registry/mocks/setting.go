package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
)

func ServiceSetting(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.SettingService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
