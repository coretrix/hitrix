package oss

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
)

type AmazonOSS struct {
	client       *s3.S3
	clockService clock.IClock
	ctx          context.Context
	buckets      map[Bucket]*BucketConfig
}

func NewAmazonOSS(configService config.IConfig, clockService clock.IClock, publicNamespaces, privateNamespaces []Namespace) (IProvider, error) {
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

	return &AmazonOSS{
		client:       s3.New(newSession),
		clockService: clockService,
		ctx:          context.Background(),
		buckets:      loadBucketsConfig(configService, publicNamespaces, privateNamespaces),
	}, nil
}

type CachedObjectURLTemplate struct {
	Environment string
	BucketName  string
	StorageKey  string
	CounterID   string
}

func (ossStorage *AmazonOSS) GetBucketConfig(bucket Bucket) *BucketConfig {
	return getBucketConfig(ossStorage.buckets[bucket])
}

func (ossStorage *AmazonOSS) GetClient() interface{} {
	return ossStorage.client
}

func (ossStorage *AmazonOSS) GetObjectURL(bucket Bucket, object *Object) (string, error) {
	cdnURL, err := ossStorage.GetObjectCDNURL(bucket, object)

	if err != nil {
		return "", err
	}

	if cdnURL != "" {
		return cdnURL, nil
	}

	return ossStorage.GetObjectOSSURL(bucket, object)
}

func (ossStorage *AmazonOSS) GetObjectOSSURL(_ Bucket, _ *Object) (string, error) {
	panic("not implemented")
}

func (ossStorage *AmazonOSS) GetObjectCDNURL(bucket Bucket, object *Object) (string, error) {
	return getObjectCDNURL(ossStorage.buckets[bucket], object.StorageKey), nil
}

func (ossStorage *AmazonOSS) GetObjectSignedURL(bucket Bucket, object *Object, expires time.Time) (string, error) {
	now := ossStorage.clockService.Now()

	if now.After(expires) {
		return "", errors.New("expire time is before now")
	}

	req, _ := ossStorage.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(ossStorage.buckets[bucket].Name),
		Key:    aws.String(object.StorageKey)},
	)

	return req.Presign(expires.Sub(now))
}

func (ossStorage *AmazonOSS) GetObjectBase64Content(_ Bucket, _ *Object) (string, error) {
	panic("not implemented")
}

func (ossStorage *AmazonOSS) UploadObjectFromByte(ormService *beeorm.Engine, bucket Bucket, namespace Namespace, objectContent []byte, extension string) (Object, error) {
	bucketConfig := ossStorage.buckets[bucket]

	bucketConfig.validateNamespace(namespace)

	storageCounter := getStorageCounter(ormService, bucketConfig)

	objectKey := ossStorage.getObjectKey(namespace, storageCounter, extension)

	putObjectInput := &s3.PutObjectInput{
		Body:   bytes.NewReader(objectContent),
		Bucket: aws.String(bucketConfig.Name),
		Key:    aws.String(objectKey),
	}

	if bucketConfig.ACL != "" {
		putObjectInput.ACL = aws.String(bucketConfig.ACL)
	}

	_, err := ossStorage.client.PutObject(putObjectInput)

	if err != nil {
		return Object{}, err
	}

	return Object{
		ID:         storageCounter,
		StorageKey: objectKey,
	}, nil
}

func (ossStorage *AmazonOSS) UploadObjectFromFile(ormService *beeorm.Engine, bucket Bucket, namespace Namespace, localFile string) (Object, error) {
	fileContent, ext, err := readContentFile(localFile)

	if err != nil {
		return Object{}, err
	}

	return ossStorage.UploadObjectFromByte(ormService, bucket, namespace, fileContent, ext)
}

func (ossStorage *AmazonOSS) UploadObjectFromBase64(ormService *beeorm.Engine, bucket Bucket, namespace Namespace, content, extension string) (Object, error) {
	byteData, err := base64.StdEncoding.DecodeString(content)

	if err != nil {
		return Object{}, err
	}

	return ossStorage.UploadObjectFromByte(ormService, bucket, namespace, byteData, extension)
}

func (ossStorage *AmazonOSS) UploadImageFromFile(ormService *beeorm.Engine, bucket Bucket, namespace Namespace, localFile string) (Object, error) {
	return ossStorage.UploadObjectFromFile(ormService, bucket, namespace, localFile)
}

func (ossStorage *AmazonOSS) UploadImageFromBase64(ormService *beeorm.Engine, bucket Bucket, namespace Namespace, image, extension string) (Object, error) {
	byteData, err := base64.StdEncoding.DecodeString(image)

	if err != nil {
		return Object{}, err
	}

	return ossStorage.UploadObjectFromByte(ormService, bucket, namespace, byteData, extension)
}

func (ossStorage *AmazonOSS) DeleteObject(bucket Bucket, object *Object) error {
	_, err := ossStorage.client.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(ossStorage.buckets[bucket].Name),
		Delete: &s3.Delete{Objects: []*s3.ObjectIdentifier{{
			Key: aws.String(object.StorageKey),
		}}},
	})

	return err
}

func (ossStorage *AmazonOSS) CreateObjectFromKey(ormService *beeorm.Engine, bucket Bucket, key string) Object {
	//TODO remove
	return Object{
		ID:         getStorageCounter(ormService, ossStorage.buckets[bucket]),
		StorageKey: key,
	}
}

func (ossStorage *AmazonOSS) getObjectKey(namespace Namespace, storageCounter uint64, fileExtension string) string {
	if namespace != "" {
		return namespace.String() + "/" + strconv.FormatUint(storageCounter, 10) + fileExtension
	}

	return strconv.FormatUint(storageCounter, 10) + fileExtension
}
