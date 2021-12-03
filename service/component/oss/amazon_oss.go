package oss

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"strconv"
	"time"

	"github.com/coretrix/hitrix/service/component/clock"

	"github.com/coretrix/hitrix/service/component/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/latolukasz/beeorm"
)

type AmazonOSS struct {
	client         *s3.S3
	clockService   clock.IClock
	ctx            context.Context
	env            string
	buckets        *Buckets
	uploaderBucket string
}

func NewAmazonOSS(configService config.IConfig, clockService clock.IClock, bucketsMapping map[string]*Bucket, env string) (*AmazonOSS, error) {
	disableSSL := false

	if val, ok := configService.Bool("oss.amazon.disable_ssl"); ok && val {
		disableSSL = true
	}

	endpoint, ok := configService.String("oss.amazon.endpoint")
	if !ok {
		return nil, errors.New("missing endpoint")
	}

	accessKeyID, ok := configService.String("oss.amazon.access_key_id")
	if !ok {
		return nil, errors.New("missing access_key_id")
	}

	secretAccessKey, ok := configService.String("oss.amazon.secret_access_key")
	if !ok {
		return nil, errors.New("missing secret_access_key")
	}

	region, ok := configService.String("oss.amazon.region")
	if !ok {
		return nil, errors.New("missing region")
	}

	newSession, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(region),
		DisableSSL:       aws.Bool(disableSSL),
	})

	if err != nil {
		return nil, err
	}

	bucket, _ := configService.String("oss.uploader.bucket")

	return &AmazonOSS{
		client:         s3.New(newSession),
		clockService:   clockService,
		ctx:            context.Background(),
		env:            env,
		buckets:        loadBucketsConfig(configService, bucketsMapping),
		uploaderBucket: bucket,
	}, nil
}

type CachedObjectURLTemplate struct {
	Environment string
	BucketName  string
	StorageKey  string
	CounterID   string
}

func (ossStorage *AmazonOSS) GetClient() interface{} {
	return ossStorage.client
}

func (ossStorage *AmazonOSS) GetObjectURL(bucket string, object *Object) (string, error) {
	cdnURL, err := ossStorage.GetObjectCDNURL(bucket, object)

	if err != nil {
		return "", err
	}

	if cdnURL != "" {
		return cdnURL, nil
	}

	return ossStorage.GetObjectOSSURL(bucket, object)
}

func (ossStorage *AmazonOSS) GetObjectOSSURL(_ string, _ *Object) (string, error) {
	panic("not implemented")
}

func (ossStorage *AmazonOSS) GetObjectCDNURL(bucket string, object *Object) (string, error) {
	return getObjectCDNURL(
		ossStorage.buckets,
		bucket,
		ossStorage.env,
		object.StorageKey), nil
}

func (ossStorage *AmazonOSS) GetObjectSignedURL(bucket string, object *Object, expires time.Time) (string, error) {
	now := ossStorage.clockService.Now()

	if now.After(expires) {
		return "", errors.New("expire time is before now")
	}

	req, _ := ossStorage.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(getBucketName(ossStorage.buckets, bucket, ossStorage.env)),
		Key:    aws.String(object.StorageKey)},
	)

	return req.Presign(expires.Sub(now))
}

func (ossStorage *AmazonOSS) GetObjectBase64Content(_ string, _ *Object) (string, error) {
	panic("not implemented")
}

func (ossStorage *AmazonOSS) UploadObjectFromByte(ormService *beeorm.Engine, bucket, path string, objectContent []byte, extension string) (Object, error) {
	bucketExists(ossStorage.buckets, bucket)
	pathExists(ossStorage.buckets, bucket, path)

	storageCounter := getStorageCounter(ormService, ossStorage.buckets, bucket)

	objectKey := ossStorage.getObjectKey(path, storageCounter, extension)

	_, err := ossStorage.client.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(objectContent),
		Bucket: aws.String(getBucketName(ossStorage.buckets, bucket, ossStorage.env)),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		return Object{}, err
	}

	return Object{
		ID:         storageCounter,
		StorageKey: objectKey,
	}, nil
}

func (ossStorage *AmazonOSS) UploadObjectFromFile(ormService *beeorm.Engine, bucket, path, localFile string) (Object, error) {
	fileContent, ext, err := readContentFile(localFile)

	if err != nil {
		return Object{}, err
	}

	return ossStorage.UploadObjectFromByte(ormService, bucket, path, fileContent, ext)
}

func (ossStorage *AmazonOSS) UploadObjectFromBase64(ormService *beeorm.Engine, bucket, path, content, extension string) (Object, error) {
	byteData, err := base64.StdEncoding.DecodeString(content)

	if err != nil {
		return Object{}, err
	}

	return ossStorage.UploadObjectFromByte(ormService, bucket, path, byteData, extension)
}

func (ossStorage *AmazonOSS) UploadImageFromFile(ormService *beeorm.Engine, bucket, path, localFile string) (Object, error) {
	return ossStorage.UploadObjectFromFile(ormService, bucket, path, localFile)
}

func (ossStorage *AmazonOSS) UploadImageFromBase64(ormService *beeorm.Engine, bucket, path, image, extension string) (Object, error) {
	byteData, err := base64.StdEncoding.DecodeString(image)

	if err != nil {
		return Object{}, err
	}

	return ossStorage.UploadObjectFromByte(ormService, bucket, path, byteData, extension)
}

func (ossStorage *AmazonOSS) DeleteObject(bucket string, object *Object) error {
	_, err := ossStorage.client.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(getBucketName(ossStorage.buckets, bucket, ossStorage.env)),
		Delete: &s3.Delete{Objects: []*s3.ObjectIdentifier{{
			Key: aws.String(object.StorageKey),
		}}},
	})

	return err
}

func (ossStorage *AmazonOSS) getObjectKey(path string, storageCounter uint64, fileExtension string) string {
	if path != "" {
		return path + "/" + strconv.FormatUint(storageCounter, 10) + fileExtension
	} else {
		return strconv.FormatUint(storageCounter, 10) + fileExtension
	}
}

func (ossStorage *AmazonOSS) CreateObjectFromKey(ormService *beeorm.Engine, bucket, key string) Object {
	//TODO remove
	return Object{
		ID:         getStorageCounter(ormService, ossStorage.buckets, bucket),
		StorageKey: key,
	}
}

func (ossStorage *AmazonOSS) GetUploaderBucketConfig() *BucketConfig {
	return getBucketEnvConfig(ossStorage.buckets, ossStorage.uploaderBucket, ossStorage.env)
}
