package registry

import (
	"github.com/sarulabs/di"
	tusd "github.com/tus/tusd/pkg/handler"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/oss"
	"github.com/coretrix/hitrix/service/component/uploader"
	"github.com/coretrix/hitrix/service/component/uploader/locker"
)

func ServiceProviderUploader(c tusd.Config, getLockerFunc locker.GetLockerFunc) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.UploaderService,
		Build: func(ctn di.Container) (interface{}, error) {
			osService := ctn.Get(service.OSService).(oss.IProvider)

			uploaderBucketName := osService.GetBucketConfig(oss.BucketPublic).Name

			composer := tusd.NewStoreComposer()

			store := uploader.GetStore(osService.GetClient(), uploaderBucketName)
			store.UseIn(composer)

			if getLockerFunc != nil {
				composer.UseLocker(getLockerFunc(ctn))
			}

			c.StoreComposer = composer

			return uploader.NewTUSDUploader(c, uploaderBucketName), nil
		},
	}
}
