package registry

import (
	"errors"

	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/oss"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/uploader"
	"github.com/coretrix/hitrix/service/component/uploader/locker"

	"github.com/sarulabs/di"
	tusd "github.com/tus/tusd/pkg/handler"
)

func ServiceProviderUploader(c tusd.Config, getLockerFunc locker.GetLockerFunc) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.UploaderService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)

			bucket, ok := configService.String("oss.uploader.bucket")

			if !ok {
				panic(errors.New("missing uploader bucket"))
			}

			osService := ctn.Get(service.OSService).(oss.IProvider)

			composer := tusd.NewStoreComposer()

			store := uploader.GetStore(osService.GetClient(), osService.GetUploaderBucketConfig().Name)
			store.UseIn(composer)

			if getLockerFunc != nil {
				composer.UseLocker(getLockerFunc(ctn))
			}

			c.StoreComposer = composer

			return uploader.NewTUSDUploader(c, bucket), nil
		},
	}
}
