package s3

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/latolukasz/beeorm"
)

type AmazonS3 struct {
	client                   *s3.S3
	ctx                      context.Context
	environment              string
	bucketsMapping           map[string]uint64
	bucketsConfigDefinitions map[string]map[string]string
	urlPrefix                string
	domain                   string
}

func NewAmazonS3(endpoint string, accessKeyID string, secretAccessKey string, allowedBuckets map[string]uint64, bucketsConfigDefinitions map[string]map[string]string,
	region string, disableSSL bool, urlPrefix string, domain string, environment string) *AmazonS3 {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(region),
		DisableSSL:       aws.Bool(disableSSL),
	}
	newSession, _ := session.NewSession(s3Config)
	s3Client := s3.New(newSession)

	return &AmazonS3{
		client:                   s3Client,
		ctx:                      context.Background(),
		bucketsMapping:           allowedBuckets,
		bucketsConfigDefinitions: bucketsConfigDefinitions,
		environment:              environment,
		urlPrefix:                urlPrefix,
		domain:                   domain,
	}
}

func (amazonS3 *AmazonS3) getCounter(ormService *beeorm.Engine, bucket string) uint64 {
	amazonS3.checkBucket(bucket)

	bucketID, has := amazonS3.bucketsMapping[bucket]

	if !has {
		panic("s3 bucket [" + bucket + "] id not found")
	}

	amazonS3BucketCounterEntity := &entity.S3BucketCounterEntity{}

	locker := ormService.GetRedis().GetLocker()
	lock, hasLock := locker.Obtain("locker_amazon_s3_counters_bucket_"+bucket, 2*time.Second, 5*time.Second)
	defer lock.Release()

	if !hasLock {
		panic("Failed to obtain lock for locker_amazon_s3_counters_bucket_" + bucket)
	}

	has = ormService.LoadByID(bucketID, amazonS3BucketCounterEntity)
	if !has {
		amazonS3BucketCounterEntity.ID = bucketID
		amazonS3BucketCounterEntity.Counter = 1
	} else {
		amazonS3BucketCounterEntity.Counter = amazonS3BucketCounterEntity.Counter + 1
	}
	ormService.Flush(amazonS3BucketCounterEntity)

	ttl := lock.TTL()
	if ttl == 0 {
		panic("lock lost")
	}

	return amazonS3BucketCounterEntity.Counter
}

func (amazonS3 *AmazonS3) checkBucket(bucketName string) {
	_, ok := amazonS3.bucketsMapping[bucketName]

	if !ok {
		panic("bucket [" + bucketName + "] not found")
	}
}

func (amazonS3 *AmazonS3) getBucketName(bucketName string) string {
	if bucketConfig, ok := amazonS3.bucketsConfigDefinitions[bucketName]; ok {
		if bucket, ok := bucketConfig[amazonS3.environment]; ok {
			return bucket
		}
	}

	return ""
}

func (amazonS3 *AmazonS3) DeleteObject(bucket string, objects ...*Object) bool {
	objectIds := make([]*s3.ObjectIdentifier, len(objects))

	for i, file := range objects {
		objectIds[i] = &s3.ObjectIdentifier{
			Key: aws.String(file.StorageKey),
		}
	}

	input := s3.DeleteObjectsInput{
		Bucket: aws.String(amazonS3.getBucketName(bucket)),
		Delete: &s3.Delete{Objects: objectIds},
	}
	deletedObjects, err := amazonS3.client.DeleteObjects(&input)
	if err != nil {
		panic("s3BucketObjectsDelete:" + err.Error())
	}
	return len(deletedObjects.Deleted) == len(objects)
}

func (amazonS3 *AmazonS3) putObject(ormService *beeorm.Engine, bucket string, objectContent []byte, extension string) Object {
	storageCounter := amazonS3.getCounter(ormService, bucket)

	objectKey := amazonS3.getObjectKey(storageCounter, extension)

	bucket = amazonS3.getBucketName(bucket)

	_, err := amazonS3.client.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(objectContent),
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		panic("s3BucketObjectPut:" + err.Error())
	}

	return Object{
		ID:         storageCounter,
		StorageKey: objectKey,
	}
}

func (amazonS3 *AmazonS3) getObjectKey(storageCounter uint64, fileExtension string) string {
	return strconv.FormatUint(storageCounter, 10) + fileExtension
}

func (amazonS3 *AmazonS3) UploadObjectFromFile(ormService *beeorm.Engine, bucket, localFile string) Object {
	amazonS3.checkBucket(bucket)

	fileContent, ext := amazonS3.ReadFile(localFile)

	return amazonS3.putObject(ormService, bucket, fileContent, ext)
}

func (amazonS3 *AmazonS3) UploadObjectFromBase64(ormService *beeorm.Engine, bucket, base64content, extension string) Object {
	byteData, err := base64.StdEncoding.DecodeString(base64content)

	if err != nil {
		panic(err)
	}

	return amazonS3.putObject(ormService, bucket, byteData, extension)
}

func (amazonS3 *AmazonS3) UploadImageFromBase64(ormService *beeorm.Engine, bucket, base64image, extension string) Object {
	byteData, err := base64.StdEncoding.DecodeString(base64image)

	if err != nil {
		panic(err)
	}

	return amazonS3.putObject(ormService, bucket, byteData, extension)
}

func (amazonS3 *AmazonS3) UploadImageFromFile(ormService *beeorm.Engine, bucket, localFile string) Object {
	return amazonS3.UploadObjectFromFile(ormService, bucket, localFile)
}

func (amazonS3 *AmazonS3) ReadFile(localFile string) ([]byte, string) {
	fileContent, err := ioutil.ReadFile(localFile)

	if err != nil {
		panic(err)
	}

	return fileContent, filepath.Ext(localFile)
}

func (amazonS3 *AmazonS3) GetObjectCachedURL(bucket string, object *Object) string {
	amazonS3.checkBucket(bucket)

	bucketByEnv := amazonS3.getBucketName(bucket)

	return fmt.Sprintf("https://%s%s.%s/%s/%s", amazonS3.urlPrefix, amazonS3.environment, amazonS3.domain,
		bucketByEnv, object.StorageKey)
}

func (amazonS3 *AmazonS3) GetObjectSignedURL(bucket string, object *Object, expiresIn time.Duration) string {
	amazonS3.checkBucket(bucket)

	bucketByEnv := amazonS3.getBucketName(bucket)

	req, _ := amazonS3.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketByEnv),
		Key:    aws.String(object.StorageKey)},
	)

	url, err := req.Presign(expiresIn) // Set link expiration time

	if err != nil {
		panic(err)
	}

	return url
}

func (amazonS3 *AmazonS3) CreateObjectFromKey(ormService *beeorm.Engine, bucket, key string) Object {
	return Object{
		ID:         amazonS3.getCounter(ormService, bucket),
		StorageKey: key,
	}
}

func (amazonS3 *AmazonS3) GetClient() interface{} {
	return amazonS3.client
}

func (amazonS3 *AmazonS3) GetBucketName(bucket string) string {
	return amazonS3.bucketsConfigDefinitions[bucket][amazonS3.environment]
}

type Object struct {
	ID         uint64
	StorageKey string
	CachedURL  string
	Data       interface{}
}

type Client interface {
	GetClient() interface{}
	GetBucketName(bucket string) string
	CreateObjectFromKey(ormService *beeorm.Engine, bucket, key string) Object
	GetObjectCachedURL(bucket string, object *Object) string
	GetObjectSignedURL(bucket string, object *Object, expires time.Duration) string
	UploadObjectFromFile(ormService *beeorm.Engine, bucket, localFile string) Object
	UploadObjectFromBase64(ormService *beeorm.Engine, bucket, content, extension string) Object
	UploadImageFromFile(ormService *beeorm.Engine, bucket, localFile string) Object
	UploadImageFromBase64(ormService *beeorm.Engine, bucket, image, extension string) Object
	DeleteObject(bucket string, objects ...*Object) bool
}
