package datastore

import (
	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/tus/tusd/pkg/s3store"
)

type AmazonS3Store struct {
	s3 s3store.S3Store
}

func (s *AmazonS3Store) UseIn(composer *tusd.StoreComposer) {
	s.s3.UseIn(composer)
}
