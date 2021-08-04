package registry

import (
	"errors"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/uploader"
	datastore "github.com/coretrix/hitrix/service/component/uploader/data_store"
	"github.com/coretrix/hitrix/service/component/uploader/locker"

	"github.com/sarulabs/di"
	tusd "github.com/tus/tusd/pkg/handler"
)

func ServiceDefinitionUploader(c tusd.Config, getStoreFunc datastore.GetStoreFunc, getLockerFunc locker.GetLockerFunc) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.UploaderService,
		Build: func(ctn di.Container) (interface{}, error) {
			composer := tusd.NewStoreComposer()

			if getStoreFunc == nil {
				panic(errors.New("missing get data store func"))
			}

			store := getStoreFunc(ctn)
			store.UseIn(composer)

			if getLockerFunc != nil {
				composer.UseLocker(getLockerFunc(ctn))
			}

			c.StoreComposer = composer
			return uploader.NewTUSDUploader(c), nil
		},
	}
}
