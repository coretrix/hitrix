package uploader

import (
	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/tus/tusd/pkg/s3store"
)

type AmazonStore struct {
	store s3store.S3Store
}

func (s *AmazonStore) UseIn(composer *tusd.StoreComposer) {
	s.store.UseIn(composer)
}
