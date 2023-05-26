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

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
)

type AmazonOSS struct {
	client       *s3.S3
	clockService clock.IClock
	ctx          context.Context
	buckets      bucketsConfig
	namespaces   namespacesConfig
}

func NewAmazonOSS(configService config.IConfig, clockService clock.IClock, namespaces Namespaces) (IProvider, error) {
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

	bucketsConfiguration, namespacesConfiguration := loadConfig(configService, namespaces)

	return &AmazonOSS{
		client:       s3.New(newSession),
		clockService: clockService,
		ctx:          context.Background(),
		buckets:      bucketsConfiguration,
		namespaces:   namespacesConfiguration,
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

func (ossStorage *AmazonOSS) GetObjectURL(namespace Namespace, object *entity.FileObject) (string, error) {
	cdnURL, err := ossStorage.GetObjectCDNURL(namespace, object)
	if err != nil {
		return "", err
	}

	if cdnURL != "" {
		return cdnURL, nil
	}

	return ossStorage.GetObjectOSSURL(namespace, object)
}

func (ossStorage *AmazonOSS) GetObjectOSSURL(_ Namespace, _ *entity.FileObject) (string, error) {
	panic("not implemented")
}

func (ossStorage *AmazonOSS) GetObjectCDNURL(namespace Namespace, object *entity.FileObject) (string, error) {
	if object == nil {
		return "", errors.New("nil file object")
	}

	bucketConfig, err := ossStorage.namespaces.getBucketConfig(namespace)
	if err != nil {
		return "", err
	}

	return getObjectCDNURL(bucketConfig, object.StorageKey), nil
}

func (ossStorage *AmazonOSS) GetObjectSignedURL(namespace Namespace, object *entity.FileObject, expires time.Time) (string, error) {
	if object == nil {
		return "", errors.New("nil file object")
	}

	now := ossStorage.clockService.Now()

	if now.After(expires) {
		return "", errors.New("expire time is before now")
	}

	bucketConfig, err := ossStorage.namespaces.getBucketConfig(namespace)
	if err != nil {
		return "", err
	}

	req, _ := ossStorage.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketConfig.Name),
		Key:    aws.String(object.StorageKey)},
	)

	return req.Presign(expires.Sub(now))
}

func (ossStorage *AmazonOSS) GetObjectBase64Content(_ Namespace, _ *entity.FileObject) (string, error) {
	panic("not implemented")
}

func (ossStorage *AmazonOSS) UploadObjectFromByte(
	ormService *beeorm.Engine,
	namespace Namespace,
	objectContent []byte,
	extension string,
) (entity.FileObject, error) {
	bucketConfig, err := ossStorage.namespaces.getBucketConfig(namespace)
	if err != nil {
		return entity.FileObject{}, err
	}

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

	_, err = ossStorage.client.PutObject(putObjectInput)
	if err != nil {
		return entity.FileObject{}, err
	}

	return entity.FileObject{
		ID:         storageCounter,
		StorageKey: objectKey,
	}, nil
}

func (ossStorage *AmazonOSS) UploadObjectFromFile(ormService *beeorm.Engine, namespace Namespace, localFile string) (entity.FileObject, error) {
	fileContent, ext, err := readContentFile(localFile)
	if err != nil {
		return entity.FileObject{}, err
	}

	return ossStorage.UploadObjectFromByte(ormService, namespace, fileContent, ext)
}

func (ossStorage *AmazonOSS) UploadObjectFromBase64(
	ormService *beeorm.Engine,
	namespace Namespace,
	content string,
	extension string,
) (entity.FileObject, error) {
	byteData, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return entity.FileObject{}, err
	}

	return ossStorage.UploadObjectFromByte(ormService, namespace, byteData, extension)
}

func (ossStorage *AmazonOSS) GetBucketConfigNamespace(namespace Namespace) (*BucketConfig, error) {
	bucketConfig, err := ossStorage.namespaces.getBucketConfig(namespace)

	if err != nil {
		return nil, err
	}

	return bucketConfig, nil
}

func (ossStorage *AmazonOSS) UploadImageFromFile(ormService *beeorm.Engine, namespace Namespace, localFile string) (entity.FileObject, error) {
	return ossStorage.UploadObjectFromFile(ormService, namespace, localFile)
}

func (ossStorage *AmazonOSS) UploadImageFromBase64(
	ormService *beeorm.Engine,
	namespace Namespace,
	image string,
	extension string,
) (entity.FileObject, error) {
	byteData, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		return entity.FileObject{}, err
	}

	return ossStorage.UploadObjectFromByte(ormService, namespace, byteData, extension)
}

func (ossStorage *AmazonOSS) DeleteObject(namespace Namespace, object *entity.FileObject) error {
	bucketConfig, err := ossStorage.namespaces.getBucketConfig(namespace)
	if err != nil {
		return err
	}

	_, err = ossStorage.client.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(bucketConfig.Name),
		Delete: &s3.Delete{Objects: []*s3.ObjectIdentifier{{
			Key: aws.String(object.StorageKey),
		}}},
	})

	return err
}

func (ossStorage *AmazonOSS) getObjectKey(namespace Namespace, storageCounter uint64, fileExtension string) string {
	return namespace.String() + "/" + strconv.FormatUint(storageCounter, 10) + fileExtension
}
