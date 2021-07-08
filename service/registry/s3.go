package registry

import (
	"errors"

	"github.com/coretrix/hitrix/service"
	s3 "github.com/coretrix/hitrix/service/component/amazon/storage"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/latolukasz/orm"
	"github.com/sarulabs/di"
)

// ServiceDefinitionAmazonS3 Be sure that you registered entity S3BucketCounterEntity
func ServiceDefinitionAmazonS3(bucketsMapping map[string]uint64) *service.Definition {
	return &service.Definition{
		Name:   service.AmazonS3Service,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			ormConfig := ctn.Get(service.ORMConfigService).(orm.ValidatedRegistry)
			entities := ormConfig.GetEntities()
			if _, ok := entities["entity.S3BucketCounterEntity"]; !ok {
				return nil, errors.New("you should register S3BucketCounterEntity")
			}

			configService := ctn.Get(service.ConfigService).(config.IConfig)
			appService := ctn.Get(service.AppService).(*app.App)
			disableSSL := false

			if val, ok := configService.Bool("amazon_s3.disable_ssl"); ok && val {
				disableSSL = true
			}

			endpoint, ok := configService.String("amazon_s3.endpoint")
			if !ok {
				return nil, errors.New("missing endpoint")
			}

			accessKeyID, ok := configService.String("amazon_s3.access_key_id")
			if !ok {
				return nil, errors.New("missing access_key_id")
			}

			secretAccessKey, ok := configService.String("amazon_s3.secret_access_key")
			if !ok {
				return nil, errors.New("missing secret_access_key")
			}

			region, ok := configService.String("amazon_s3.region")
			if !ok {
				return nil, errors.New("missing region")
			}

			urlPrefix, ok := configService.String("amazon_s3.url_prefix")
			if !ok {
				return nil, errors.New("missing url_prefix")
			}

			domain, ok := configService.String("amazon_s3.domain")
			if !ok {
				return nil, errors.New("missing domain")
			}

			bucketsConfigDefinitions, ok := configService.Get("amazon_s3.buckets")
			if !ok {
				return nil, errors.New("missing buckets")
			}

			bucketsConfigDefinitionsMap := map[string]map[string]string{}
			for k, v := range bucketsConfigDefinitions.(map[interface{}]interface{}) {
				bucketsConfigDefinitionsInnerMap := map[string]string{}
				for k1, v1 := range v.(map[interface{}]interface{}) {
					bucketsConfigDefinitionsInnerMap[k1.(string)] = v1.(string)
				}

				bucketsConfigDefinitionsMap[k.(string)] = bucketsConfigDefinitionsInnerMap
			}

			return s3.NewAmazonS3(
				endpoint,
				accessKeyID,
				secretAccessKey,
				bucketsMapping,
				bucketsConfigDefinitionsMap,
				region,
				disableSSL,
				urlPrefix,
				domain,
				appService.Mode), nil
		},
	}
}
