package datastore

import (
	"errors"

	"github.com/coretrix/hitrix/service"
	s3hitrix "github.com/coretrix/hitrix/service/component/amazon/storage"
	"github.com/coretrix/hitrix/service/component/config"

	"github.com/sarulabs/di"
	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/tus/tusd/pkg/s3store"
)

type Store interface {
	UseIn(composer *tusd.StoreComposer)
}

type GetStoreFunc func(ctn di.Container) Store

func GetAmazonS3Store(ctn di.Container) Store {
	s3Client := ctn.Get(service.AmazonS3Service).(*s3hitrix.AmazonS3)
	configService := ctn.Get(service.ConfigService).(config.IConfig)

	bucket, ok := configService.String("uploader.bucket")
	if !ok {
		panic(errors.New("missing bucket"))
	}

	store := s3store.New(s3Client.GetBucketName(bucket), s3Client.GetClient().(s3store.S3API))
	return &AmazonS3Store{s3: store}
}
