package oss

import (
	"io/ioutil"
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
	bucketPublicDBID  = 1
	bucketPrivateDBID = 2
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

	bucketsConfig := bucketsConfigDefinitions.(map[string]map[string]interface{})

	publicBucketConfig, hasPublicBucketConfig := bucketsConfig[BucketPublic.String()]
	privateBucketConfig, hasPrivateBucketConfig := bucketsConfig[BucketPrivate.String()]

	if !hasPublicBucketConfig && !hasPrivateBucketConfig {
		panic("oss: invalid bucket configuration. no buckets defined")
	}

	if hasPublicBucketConfig {
		bucketName, has := publicBucketConfig["name"]
		if !has || bucketName == nil {
			panic("oss: missing bucket name for public bucket")
		}

		if publicNamespaces == nil {
			panic("oss: missing namespaces for public bucket")
		}

		bucketConfig := &BucketConfig{
			StorageCounterDatabaseID: bucketPublicDBID,
			Name:                     bucketName.(string),
			Namespaces:               map[Namespace]Namespace{},
		}

		for _, publicNamespace := range publicNamespaces {
			bucketConfig.Namespaces[publicNamespace] = publicNamespace
		}

		bucketCDNURL, has := publicBucketConfig["cdn_url"]

		if has && bucketCDNURL != nil {
			bucketConfig.CDNURL = bucketCDNURL.(string)
		}

		buckets[BucketPublic] = bucketConfig
	}

	if hasPrivateBucketConfig {
		bucketName, has := privateBucketConfig["name"]
		if !has || bucketName == nil {
			panic("oss: missing bucket name for private bucket")
		}

		if privateNamespaces == nil {
			panic("oss: missing namespaces for private bucket")
		}

		bucketConfig := &BucketConfig{
			StorageCounterDatabaseID: bucketPrivateDBID,
			Name:                     bucketName.(string),
			Namespaces:               map[Namespace]Namespace{},
		}

		for _, privateNamespace := range privateNamespaces {
			bucketConfig.Namespaces[privateNamespace] = privateNamespace
		}

		bucketCDNURL, has := privateBucketConfig["cdn_url"]

		if has && bucketCDNURL != nil {
			bucketConfig.CDNURL = bucketCDNURL.(string)
		}

		buckets[BucketPrivate] = bucketConfig
	}

	return buckets
}

func getObjectCDNURL(bucketConfig *BucketConfig, storageKey string) string {
	cdnURL := bucketConfig.CDNURL

	if cdnURL == "" {
		return ""
	}

	replacer := strings.NewReplacer("{StorageKey}", storageKey, "{Bucket}", bucketConfig.Name)

	return replacer.Replace(cdnURL)
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
	fileContent, err := ioutil.ReadFile(localFile)

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
	}
}
