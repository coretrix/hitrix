package oss

import (
	"time"

	"github.com/latolukasz/orm"
)

type Client interface {
	GetObjectURL(bucket string, object *Object) string
	GetObjectCachedURL(bucket string, object *Object) string
	GetObjectSignedURL(bucket string, object *Object, expires time.Time) string
	UploadObjectFromFile(ormService *orm.Engine, bucket, localFile string) Object
	UploadObjectFromBase64(ormService *orm.Engine, bucket, content, extension string) Object
	UploadImageFromFile(ormService *orm.Engine, bucket, localFile string) Object
	UploadImageFromBase64(ormService *orm.Engine, bucket, image, extension string) Object
	GetObjectBase64Content(bucket string, object *Object) (string, error)
}

type Object struct {
	ID         uint64
	StorageKey string
	CachedURL  string
	Data       interface{}
}
