package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/sarulabs/di"
)

func ServiceSetting(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.SettingService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
