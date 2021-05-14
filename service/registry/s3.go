package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/amazon/storage"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/sarulabs/di"
)

func ServiceDefinitionAmazonS3(buckets map[string]uint64) *service.Definition {
	return &service.Definition{
		Name:   service.AmazonS3Service,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			conf := ctn.Get(service.ConfigService).(*config.Config).GetStringMap("amazon_s3")
			appService := ctn.Get(service.AppService).(*app.App)
			disableSSL := false

			if val, ok := conf["disable_ssl"]; ok {
				if val == "true" {
					disableSSL = true
				}
			}
			return s3.NewAmazonS3(conf["endpoint"].(string), conf["access_key_id"].(string),
				conf["secret_access_key"].(string), buckets, conf["region"].(string), disableSSL,
				conf["url_prefix"].(string), conf["domain"].(string), appService.Mode, conf), nil
		},
	}
}
