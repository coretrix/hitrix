package mocks

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/fcm"
	"github.com/sarulabs/di"
)

func FakeFCMService(fakeFCMService fcm.FCM) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.FCMService,
		Build: func(ctn di.Container) (interface{}, error) {
			return fakeFCMService, nil
		},
	}
}
