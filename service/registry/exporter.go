package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/exporter"
	"github.com/sarulabs/di"
)

func ServiceProviderExporter() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.ExporterService,
		Build: func(ctn di.Container) (interface{}, error) {
			return exporter.NewExportService(exporter.NewXLSXExportService(), exporter.NewCSVExportService()), nil
		},
	}
}
