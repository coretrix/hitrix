package oss

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/latolukasz/beeorm"
)

const ProviderGoogleOSS = 1
const ProviderAmazonOSS = 2

type IProvider interface {
	GetObjectURL(bucket string, object *Object) (string, error)
	GetObjectOSSURL(bucket string, object *Object) (string, error)
	GetObjectCDNURL(bucket string, object *Object) (string, error)
	GetObjectSignedURL(bucket string, object *Object, expires time.Time) (string, error)
	GetObjectBase64Content(bucket string, object *Object) (string, error)
	UploadObjectFromFile(ormService *beeorm.Engine, bucket, localFile string) (Object, error)
	UploadObjectFromBase64(ormService *beeorm.Engine, bucket, content, extension string) (Object, error)
	UploadObjectFromByte(ormService *beeorm.Engine, bucket string, content []byte, extension string) (Object, error)
	UploadImageFromFile(ormService *beeorm.Engine, bucket, localFile string) (Object, error)
	UploadImageFromBase64(ormService *beeorm.Engine, bucket, image, extension string) (Object, error)
	DeleteObject(bucket string, object *Object) error
	//TODO Remove
	CreateObjectFromKey(ormService *beeorm.Engine, bucket, key string) Object
}

type Object struct {
	ID         uint64
	StorageKey string
	Data       interface{}
}

type Buckets struct {
	Mapping map[string]uint64
	Configs map[string]map[string]*BucketConfig
}

type BucketConfig struct {
	Name   string
	CDNURL string
}

func loadBucketsConfig(configService config.IConfig, bucketsMapping map[string]uint64) *Buckets {
	bucketsConfigDefinitions, ok := configService.Get("oss.buckets")

	if !ok {
		panic("oss: missing bucket configuration")
	}

	buckets := &Buckets{
		Mapping: bucketsMapping,
		Configs: map[string]map[string]*BucketConfig{},
	}

	for bucket, envsBucketConfig := range bucketsConfigDefinitions.(map[interface{}]interface{}) {
		for env, bucketConfig := range envsBucketConfig.(map[interface{}]interface{}) {
			bucketConfigMap := map[string]string{}

			for key, value := range bucketConfig.(map[interface{}]interface{}) {
				bucketConfigMap[key.(string)] = value.(string)
			}

			name, has := bucketConfigMap["name"]

			if !has {
				panic("oss: missing bucket name for bucket: " + bucket.(string) + " and env: " + env.(string))
			}

			cdnUrl, has := bucketConfigMap["cdn_url"]

			_, has = buckets.Configs[bucket.(string)]

			if !has {
				buckets.Configs[bucket.(string)] = map[string]*BucketConfig{}
			}

			buckets.Configs[bucket.(string)][env.(string)] = &BucketConfig{
				Name:   name,
				CDNURL: cdnUrl,
			}
		}
	}

	return buckets
}

func getBucketConfig(buckets *Buckets, bucket string) map[string]*BucketConfig {
	bucketExists(buckets, bucket)

	return buckets.Configs[bucket]
}

func getBucketEnvConfig(buckets *Buckets, bucket string, env string) *BucketConfig {
	return getBucketConfig(buckets, bucket)[env]
}

func getBucketName(buckets *Buckets, bucket, env string) string {
	return getBucketEnvConfig(buckets, bucket, env).Name
}

func getBucketCDNURL(buckets *Buckets, bucket, env string) string {
	return getBucketEnvConfig(buckets, bucket, env).CDNURL
}

func getObjectCDNURL(buckets *Buckets, bucket, env, storageKey string) string {
	cdnURL := getBucketCDNURL(buckets, bucket, env)

	if cdnURL == "" {
		return ""
	}

	replacer := strings.NewReplacer("{StorageKey}", storageKey, "{Env}", env, "{Bucket}", getBucketName(buckets, bucket, env))

	return replacer.Replace(cdnURL)
}

func getStorageCounter(ormService *beeorm.Engine, buckets *Buckets, bucket string) uint64 {
	bucketExists(buckets, bucket)

	bucketID := buckets.Mapping[bucket]

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

func bucketExists(buckets *Buckets, bucket string) {
	_, has := buckets.Mapping[bucket]

	if !has {
		panic("bucket [" + bucket + "] not found")
	}
}

func readContentFile(localFile string) ([]byte, string, error) {
	fileContent, err := ioutil.ReadFile(localFile)

	if err != nil {
		return nil, "", err
	}

	return fileContent, filepath.Ext(localFile), nil
}
