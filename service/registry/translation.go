package registry

import (
	"errors"

	"github.com/latolukasz/beeorm/v2"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
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

			return translation.NewTranslationService(ctn.Get(service.ErrorLoggerService).(errorlogger.ErrorLogger)), nil
		},
	}
}
