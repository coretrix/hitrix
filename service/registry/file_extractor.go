package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	fileextractor "github.com/coretrix/hitrix/service/component/file_extractor"
)

func ServiceProviderExtractor() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.ExtractorService,
		Build: func(ctn di.Container) (interface{}, error) {
			return fileextractor.NewFileExtractor(), nil
		},
	}
}
