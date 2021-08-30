package registry

import (
	"github.com/coretrix/hitrix/service"
	fileextractor "github.com/coretrix/hitrix/service/component/file_extractor"
	"github.com/sarulabs/di"
)

func ServiceProviderExtractor() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.ExtractorService,
		Build: func(ctn di.Container) (interface{}, error) {
			return fileextractor.NewFileExtractor(), nil
		},
	}
}
