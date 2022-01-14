package registry

import (
	"errors"

	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	pdf "github.com/coretrix/hitrix/service/component/pdf"
)

func ServiceProviderPDF() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.PDFService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)
			chromeWebSocketURL, ok := configService.String("chrome_headless.web_socket_url")
			if !ok {
				return nil, errors.New("missing chrome_headless.web_socket_url")
			}
			return pdf.NewPDFService(chromeWebSocketURL), nil
		},
	}
}
