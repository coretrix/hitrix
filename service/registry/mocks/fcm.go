package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/fcm"
)

func ServiceProviderMockFCM(mock fcm.FCM) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.FCMService,
		Build: func(ctn di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
