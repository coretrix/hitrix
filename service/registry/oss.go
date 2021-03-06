package registry

import (
	"errors"

	"github.com/latolukasz/orm"

	"github.com/coretrix/hitrix/pkg/helper"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/oss/storage"
	"github.com/sarulabs/di"
)

// OSSGoogle Be sure that you registered entity OSSBucketCounterEntity
func OSSGoogle(buckets map[string]uint64) *service.Definition {
	return &service.Definition{
		Name:   service.OSSGoogleService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)
			appService := ctn.Get(service.AppService).(*app.App)

			if !helper.ExistsInDir(".oss.json", configService.GetFolderPath()) {
				return nil, errors.New(configService.GetFolderPath() + "/.oss.json does not exists")
			}

			ormConfig := ctn.Get(service.ORMConfigService).(orm.ValidatedRegistry)
			entities := ormConfig.GetEntities()
			if _, ok := entities["entity.OSSBucketCounterEntity"]; !ok {
				return nil, errors.New("you should register OSSBucketCounterEntity")
			}

			if len(buckets) == 0 {
				return nil, errors.New("please define buckets")
			}

			return storage.NewGoogleOSS(configService.GetFolderPath()+"/.oss.json", appService.Mode, buckets, configService), nil
		},
	}
}
