package s3

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/latolukasz/orm"
)

type AmazonS3 struct {
	client      *s3.S3
	ctx         context.Context
	environment string
	buckets     map[string]uint64
	s3Config    map[string]interface{}
	urlPrefix   string
	domain      string
}

func NewAmazonS3(endpoint string, accessKeyID string, secretAccessKey string, allowedBuckets map[string]uint64,
	region string, disableSSL bool, urlPrefix string, domain string, environment string, config map[string]interface{}) *AmazonS3 {
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
		client:      s3Client,
		ctx:         context.Background(),
		buckets:     allowedBuckets,
		environment: environment,
		s3Config:    config,
		urlPrefix:   urlPrefix,
		domain:      domain,
	}
}

func (amazonS3 *AmazonS3) getCounter(ormService *orm.Engine, bucket string) uint64 {
	amazonS3.checkBucket(bucket)

	bucketID, has := amazonS3.buckets[bucket]

	if !has {
		panic("s3 bucket [" + bucket + "] id not found")
	}

	amazonS3BucketCounterEntity := &entity.S3BucketCounterEntity{}

	locker := ormService.GetRedis().GetLocker()
	lock, hasLock := locker.Obtain(amazonS3.ctx, "locker_amazon_s3_counters_bucket_"+bucket, 2*time.Second, 5*time.Second)
	defer lock.Release()

	if !hasLock {
		panic("Failed to obtain lock for locker_google_oss_counters_bucket_" + bucket)
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

func (amazonS3 *AmazonS3) createBucket(bucketName string) {
	params := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}
	_, err := amazonS3.client.CreateBucket(params)
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if !(ok && (aerr.Code() == s3.ErrCodeBucketAlreadyOwnedByYou || aerr.Code() == s3.ErrCodeBucketAlreadyExists)) {
			log.Panic(err)
		}
	}
}

func (amazonS3 *AmazonS3) checkBucket(bucketName string) {
	_, ok := amazonS3.buckets[bucketName]

	if !ok {
		panic("bucket [" + bucketName + "] not found")
	}

	amazonS3.createBucket(bucketName)
}

func (amazonS3 *AmazonS3) putObject(ormService *orm.Engine, bucket string, objectContent []byte, extension string) Object {
	storageCounter := amazonS3.getCounter(ormService, bucket)

	objectKey := amazonS3.getObjectKey(storageCounter, extension)

	bucketByEnv := bucket

	if amazonS3.environment != app.ModeProd {
		bucketByEnv += "-" + amazonS3.environment
		amazonS3.createBucket(bucketByEnv)
	}

	_, err := amazonS3.client.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(objectContent),
		Bucket: aws.String(bucketByEnv),
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

func (amazonS3 *AmazonS3) UploadObjectFromFile(ormService *orm.Engine, bucket, localFile string) Object {
	amazonS3.checkBucket(bucket)

	fileContent, ext := amazonS3.ReadFile(localFile)

	return amazonS3.putObject(ormService, bucket, fileContent, ext)
}

func (amazonS3 *AmazonS3) UploadObjectFromBase64(ormService *orm.Engine, bucket, base64content, extension string) Object {
	byteData, err := base64.StdEncoding.DecodeString(base64content)

	if err != nil {
		panic(err)
	}

	return amazonS3.putObject(ormService, bucket, byteData, extension)
}

func (amazonS3 *AmazonS3) UploadImageFromBase64(ormService *orm.Engine, bucket, base64image, extension string) Object {
	byteData, err := base64.StdEncoding.DecodeString(base64image)

	if err != nil {
		panic(err)
	}

	return amazonS3.putObject(ormService, bucket, byteData, extension)
}

func (amazonS3 *AmazonS3) UploadImageFromFile(ormService *orm.Engine, bucket, localFile string) Object {
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

	bucketByEnv := bucket

	if amazonS3.environment != app.ModeProd {
		bucketByEnv += "-" + amazonS3.environment
	}

	return fmt.Sprintf("https://%s%s.%s/%s/%s", amazonS3.urlPrefix, amazonS3.environment, amazonS3.domain,
		bucketByEnv, object.StorageKey)
}

func (amazonS3 *AmazonS3) GetObjectSignedURL(bucket string, object *Object, expiresIn time.Duration) string {
	amazonS3.checkBucket(bucket)

	bucketByEnv := bucket

	if amazonS3.environment != app.ModeProd {
		bucketByEnv += "-" + amazonS3.environment
	}

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

type Object struct {
	ID         uint64
	StorageKey string
	CachedURL  string
	Data       interface{}
}
