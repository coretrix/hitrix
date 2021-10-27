package uploader

import (
	"github.com/tus/tusd/pkg/gcsstore"
	tusd "github.com/tus/tusd/pkg/handler"
)

type GoogleStore struct {
	store gcsstore.GCSStore
}

func (s *GoogleStore) UseIn(composer *tusd.StoreComposer) {
	s.store.UseIn(composer)
}
