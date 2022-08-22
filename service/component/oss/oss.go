package oss

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
)

type Bucket string

const (
	BucketPrivate Bucket = "private"
	BucketPublic  Bucket = "public"
)

func (b Bucket) String() string {
	return string(b)
}

type Namespace string

func (n Namespace) String() string {
	return string(n)
}

const (
	bucketPublicStorageCounterDatabaseID  = 1
	bucketPrivateStorageCounterDatabaseID = 2
)

type NewProviderFunc func(configService config.IConfig, clockService clock.IClock, publicNamespaces, privateNamespaces []Namespace) (IProvider, error)

type IProvider interface {
	GetBucketConfig(bucket Bucket) *BucketConfig
	GetClient() interface{}
	GetObjectURL(bucket Bucket, object *Object) (string, error)
	GetObjectOSSURL(bucket Bucket, object *Object) (string, error)
	GetObjectCDNURL(bucket Bucket, object *Object) (string, error)
	GetObjectSignedURL(bucket Bucket, object *Object, expires time.Time) (string, error)
	GetObjectBase64Content(bucket Bucket, object *Object) (string, error)
	UploadObjectFromFile(ormService *beeorm.Engine, bucket Bucket, namespace Namespace, localFile string) (Object, error)
	UploadObjectFromBase64(ormService *beeorm.Engine, bucket Bucket, namespace Namespace, content, extension string) (Object, error)
	UploadObjectFromByte(ormService *beeorm.Engine, bucket Bucket, namespace Namespace, content []byte, extension string) (Object, error)
	UploadImageFromFile(ormService *beeorm.Engine, bucket Bucket, namespace Namespace, localFile string) (Object, error)
	UploadImageFromBase64(ormService *beeorm.Engine, bucket Bucket, namespace Namespace, image, extension string) (Object, error)
	DeleteObject(bucket Bucket, object *Object) error
	// CreateObjectFromKey TODO Remove
	CreateObjectFromKey(ormService *beeorm.Engine, bucket Bucket, key string) Object
}

type Object struct {
	ID         uint64
	StorageKey string
	Data       interface{}
}

type BucketConfig struct {
	StorageCounterDatabaseID uint64
	Type                     Bucket
	Name                     string
	CDNURL                   string
	Namespaces               map[Namespace]Namespace
	ACL                      string
}

func (b *BucketConfig) validateNamespace(namespace Namespace) {
	_, has := b.Namespaces[namespace]

	if !has {
		panic("Namespace [" + namespace.String() + "] not found in " + b.Type.String())
	}
}

func loadBucketsConfig(configService config.IConfig, publicNamespaces, privateNamespaces []Namespace) map[Bucket]*BucketConfig {
	bucketsConfigDefinitions, ok := configService.Get("oss.buckets")

	if !ok {
		panic("oss: missing bucket configuration")
	}

	buckets := map[Bucket]*BucketConfig{}

	bucketsConfigs := bucketsConfigDefinitions.(map[interface{}]interface{})

	for bucket, bucketConfigsI := range bucketsConfigs {
		bucketConfigs := bucketConfigsI.(map[interface{}]interface{})

		bucketConfig := &BucketConfig{
			Type:       Bucket(bucket.(string)),
			Namespaces: map[Namespace]Namespace{},
		}

		if bucketConfig.Type == BucketPublic {
			bucketConfig.StorageCounterDatabaseID = bucketPublicStorageCounterDatabaseID

			if publicNamespaces == nil {
				panic("oss: missing namespaces for public bucket")
			}

			for _, publicNamespace := range publicNamespaces {
				bucketConfig.Namespaces[publicNamespace] = publicNamespace
			}
		}

		if bucketConfig.Type == BucketPrivate {
			bucketConfig.StorageCounterDatabaseID = bucketPrivateStorageCounterDatabaseID

			if privateNamespaces == nil {
				panic("oss: missing namespaces for private bucket")
			}

			for _, privateNamespace := range privateNamespaces {
				bucketConfig.Namespaces[privateNamespace] = privateNamespace
			}
		}

		for keyI, valueI := range bucketConfigs {
			key := keyI.(string)
			if key == "name" {
				if valueI == nil {
					panic("value is nil for key name for bucket: " + bucketConfig.Type)
				}

				bucketConfig.Name = valueI.(string)

				continue
			}

			if key == "cdn_url" {
				if valueI != nil {
					bucketConfig.CDNURL = valueI.(string)
				}

				continue
			}

			if key == "ACL" {
				if valueI == nil {
					panic("value is nil for key ACL for bucket: " + bucketConfig.Type)
				}

				bucketConfig.ACL = valueI.(string)

				continue
			}

			panic("invalid key " + key + " for bucket: " + bucketConfig.Type.String())
		}

		buckets[bucketConfig.Type] = bucketConfig
	}

	if len(buckets) == 0 {
		panic("missing buckets configuration")
	}

	return buckets
}

func getObjectCDNURL(bucketConfig *BucketConfig, storageKey string) string {
	if bucketConfig.CDNURL == "" {
		return ""
	}

	replacer := strings.NewReplacer("{StorageKey}", storageKey, "{Bucket}", bucketConfig.Name)

	return replacer.Replace(bucketConfig.CDNURL)
}

func getStorageCounter(ormService *beeorm.Engine, bucketConfig *BucketConfig) uint64 {
	bucketID := bucketConfig.StorageCounterDatabaseID

	ossBucketCounterEntity := &entity.OSSBucketCounterEntity{}

	locker := ormService.GetRedis().GetLocker()
	lockerKey := "locker_oss_counters_bucket_" + strconv.FormatUint(bucketID, 10)

	lock, hasLock := locker.Obtain(lockerKey, 2*time.Second, 5*time.Second)
	defer lock.Release()

	if !hasLock {
		panic("Failed to obtain lock for :" + lockerKey)
	}

	has := ormService.LoadByID(bucketID, ossBucketCounterEntity)

	if !has {
		ossBucketCounterEntity.ID = bucketID
		ossBucketCounterEntity.Counter = 1
	} else {
		ossBucketCounterEntity.Counter = ossBucketCounterEntity.Counter + 1
	}

	ormService.Flush(ossBucketCounterEntity)

	if lock.TTL() == 0 {
		panic("lock lost for :" + lockerKey)
	}

	return ossBucketCounterEntity.Counter
}

func readContentFile(localFile string) ([]byte, string, error) {
	fileContent, err := os.ReadFile(localFile)

	if err != nil {
		return nil, "", err
	}

	return fileContent, filepath.Ext(localFile), nil
}

func getBucketConfig(bucketConfig *BucketConfig) *BucketConfig {
	return &BucketConfig{
		StorageCounterDatabaseID: bucketConfig.StorageCounterDatabaseID,
		Type:                     bucketConfig.Type,
		Name:                     bucketConfig.Name,
		CDNURL:                   bucketConfig.CDNURL,
		Namespaces:               bucketConfig.Namespaces,
		ACL:                      bucketConfig.ACL,
	}
}
