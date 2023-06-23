package registry

import (
	"errors"

	"github.com/latolukasz/beeorm"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/translation"
)

func ServiceProviderTranslation() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.TranslationService,
		Build: func(ctn di.Container) (interface{}, error) {
			ormConfig := ctn.Get(service.ORMConfigService).(beeorm.ValidatedRegistry)
			entities := ormConfig.GetEntities()

			if _, ok := entities["entity.TranslationTextEntity"]; !ok {
				return nil, errors.New("you should register TranslationTextEntity")
			}

			return translation.NewTranslationService(ctn.Get(service.AppService).(*app.App)), nil
		},
	}
}
