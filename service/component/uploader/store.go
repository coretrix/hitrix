package uploader

import (
	tusd "github.com/tus/tusd/pkg/handler"
)

type Store interface {
	UseIn(composer *tusd.StoreComposer)
}
