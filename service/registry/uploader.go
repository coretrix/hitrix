package registry

import (
	"errors"

	"github.com/coretrix/hitrix/service"
	s3hitrix "github.com/coretrix/hitrix/service/component/amazon/storage"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/uploader"
	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/tus/tusd/pkg/memorylocker"
	"github.com/tus/tusd/pkg/s3store"

	"github.com/sarulabs/di"
)

func ServiceDefinitionUploader(c tusd.Config) *service.Definition {
	return &service.Definition{
		Name:   service.UploaderService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)
			s3Client := ctn.Get(service.AmazonS3Service).(*s3hitrix.AmazonS3)

			bucket, ok := configService.String("uploader.bucket")
			if !ok {
				return nil, errors.New("missing bucket")
			}

			store := s3store.New(s3Client.GetBucketName(bucket), s3Client.GetClient().(s3store.S3API))

			composer := tusd.NewStoreComposer()
			store.UseIn(composer)

			// implement redis locker
			composer.UseLocker(memorylocker.New())

			c.StoreComposer = composer
			return uploader.NewTUSDUploader(c), nil
		},
	}
}
