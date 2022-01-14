package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/exporter"
)

func ServiceProviderExporter() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.ExporterService,
		Build: func(ctn di.Container) (interface{}, error) {
			return exporter.NewExportService(exporter.NewXLSXExportService(), exporter.NewCSVExportService()), nil
		},
	}
}
