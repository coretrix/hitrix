package oss

import (
	"time"

	"github.com/latolukasz/beeorm"
)

type Client interface {
	GetObjectURL(bucket string, object *Object) string
	GetObjectCachedURL(bucket string, object *Object) string
	GetObjectSignedURL(bucket string, object *Object, expires time.Time) string
	UploadObjectFromFile(ormService *beeorm.Engine, bucket, localFile string) Object
	UploadObjectFromBase64(ormService *beeorm.Engine, bucket, content, extension string) Object
	UploadImageFromFile(ormService *beeorm.Engine, bucket, localFile string) Object
	UploadImageFromBase64(ormService *beeorm.Engine, bucket, image, extension string) Object
	GetObjectBase64Content(bucket string, object *Object) (string, error)
}

type Object struct {
	ID         uint64
	StorageKey string
	CachedURL  string
	Data       interface{}
}
