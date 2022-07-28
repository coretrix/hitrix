package registry

import (
	"errors"

	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/html2pdf"
)

func ServiceProviderHTML2PDF() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.HTML2PDFService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)
			chromeWebSocketURL, ok := configService.String("chrome_headless.web_socket_url")
			if !ok {
				return nil, errors.New("missing chrome_headless.web_socket_url")
			}

			return html2pdf.NewHTML2PDFService(chromeWebSocketURL), nil
		},
	}
}
