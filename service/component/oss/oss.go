package oss

import (
	"errors"
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

type Namespace string

func (n Namespace) String() string {
	return string(n)
}

type Namespaces map[Namespace]Bucket

type namespacesConfig map[Namespace]*BucketConfig

func (n namespacesConfig) getBucketConfig(namespace Namespace) (*BucketConfig, error) {
	bucketConfig, has := n[namespace]

	if !has {
		return nil, errors.New("Namespace [" + namespace.String() + "] not found!")
	}

	return bucketConfig, nil
}

type BucketConfig struct {
	StorageCounterDatabaseID uint64
	Type                     Bucket
	Name                     string
	CDNURL                   string
	ACL                      string
}

type bucketsConfig map[Bucket]*BucketConfig

type Bucket string

const (
	BucketPrivate Bucket = "private"
	BucketPublic  Bucket = "public"
)

func (b Bucket) String() string {
	return string(b)
}

const (
	bucketPublicStorageCounterDatabaseID  = 1
	bucketPrivateStorageCounterDatabaseID = 2
)

type NewProviderFunc func(configService config.IConfig, clockService clock.IClock, namespaces Namespaces) (IProvider, error)

type IProvider interface {
	GetBucketConfig(bucket Bucket) *BucketConfig
	GetClient() interface{}
	GetObjectURL(namespace Namespace, object *entity.FileObject) (string, error)
	GetObjectOSSURL(namespace Namespace, object *entity.FileObject) (string, error)
	GetObjectCDNURL(namespace Namespace, object *entity.FileObject) (string, error)
	GetObjectSignedURL(namespace Namespace, object *entity.FileObject, expires time.Time) (string, error)
	GetObjectBase64Content(namespace Namespace, object *entity.FileObject) (string, error)
	UploadObjectFromFile(ormService *beeorm.Engine, namespace Namespace, localFile string) (entity.FileObject, error)
	UploadObjectFromBase64(ormService *beeorm.Engine, namespace Namespace, content, extension string) (entity.FileObject, error)
	UploadObjectFromByte(ormService *beeorm.Engine, namespace Namespace, content []byte, extension string) (entity.FileObject, error)
	UploadImageFromFile(ormService *beeorm.Engine, namespace Namespace, localFile string) (entity.FileObject, error)
	UploadImageFromBase64(ormService *beeorm.Engine, namespace Namespace, image, extension string) (entity.FileObject, error)
	DeleteObject(namespace Namespace, object *entity.FileObject) error
}

func loadConfig(configService config.IConfig, namespaces Namespaces) (bucketsConfig, namespacesConfig) {
	bucketsConfigDefinitions, ok := configService.Get("oss.buckets")

	if !ok {
		panic("oss: missing bucket configuration")
	}

	namespacesConfiguration := namespacesConfig{}
	bucketsConfiguration := bucketsConfig{}

	bucketsConfigs := bucketsConfigDefinitions.(map[interface{}]interface{})

	for bucket, bucketConfigsI := range bucketsConfigs {
		bucketConfigs := bucketConfigsI.(map[interface{}]interface{})

		bucketConfig := &BucketConfig{
			Type: Bucket(bucket.(string)),
		}

		if bucketConfig.Type == BucketPublic {
			bucketConfig.StorageCounterDatabaseID = bucketPublicStorageCounterDatabaseID

			for namespace, bucket := range namespaces {
				if bucket != BucketPublic {
					continue
				}

				namespacesConfiguration[namespace] = bucketConfig
			}
		} else if bucketConfig.Type == BucketPrivate {
			bucketConfig.StorageCounterDatabaseID = bucketPrivateStorageCounterDatabaseID

			for namespace, bucket := range namespaces {
				if bucket != BucketPrivate {
					continue
				}

				namespacesConfiguration[namespace] = bucketConfig
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

		bucketsConfiguration[bucketConfig.Type] = bucketConfig
	}

	if len(bucketsConfiguration) == 0 {
		panic("missing buckets configuration")
	}

	return bucketsConfiguration, namespacesConfiguration
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
		ACL:                      bucketConfig.ACL,
	}
}
