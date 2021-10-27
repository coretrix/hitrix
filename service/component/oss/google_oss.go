package oss

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/coretrix/hitrix/service/component/clock"

	"github.com/coretrix/hitrix/pkg/helper"
	"github.com/coretrix/hitrix/service/component/config"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"cloud.google.com/go/storage"
	"github.com/latolukasz/beeorm"
	"golang.org/x/net/context"
)

type GoogleOSS struct {
	client         *storage.Client
	clockService   clock.IClock
	ctx            context.Context
	env            string
	buckets        *Buckets
	uploaderBucket string
	accessID       string
	privateKey     []byte
}

func NewGoogleOSS(configService config.IConfig, clockService clock.IClock, bucketsMapping map[string]uint64, env string) (*GoogleOSS, error) {
	ctx := context.Background()

	if !helper.ExistsInDir(".oss.json", configService.GetFolderPath()) {
		return nil, errors.New(configService.GetFolderPath() + "/.oss.json does not exists")
	}

	credentialsFile := configService.GetFolderPath() + "/.oss.json"

	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))

	if err != nil {
		return nil, err
	}

	jwtCredentialsJSONString, err := ioutil.ReadFile(credentialsFile)

	if err != nil {
		return nil, err
	}

	jwtConfig, _ := google.JWTConfigFromJSON(jwtCredentialsJSONString)

	bucket, _ := configService.String("oss.uploader.bucket")

	return &GoogleOSS{
		client:         client,
		clockService:   clockService,
		ctx:            ctx,
		env:            env,
		buckets:        loadBucketsConfig(configService, bucketsMapping),
		uploaderBucket: bucket,
		accessID:       jwtConfig.Email,
		privateKey:     jwtConfig.PrivateKey,
	}, nil
}

func (ossStorage *GoogleOSS) GetClient() interface{} {
	return ossStorage.client
}

func (ossStorage *GoogleOSS) GetObjectURL(bucket string, object *Object) (string, error) {
	cdnURL, err := ossStorage.GetObjectCDNURL(bucket, object)

	if err != nil {
		return "", err
	}

	if cdnURL != "" {
		return cdnURL, nil
	}

	return ossStorage.GetObjectOSSURL(bucket, object)
}

func (ossStorage *GoogleOSS) GetObjectOSSURL(bucket string, object *Object) (string, error) {
	bucketName := getBucketName(ossStorage.buckets, bucket, ossStorage.env)

	ossBucketObjectAttributes, err := ossStorage.client.Bucket(bucketName).Object(object.StorageKey).Attrs(ossStorage.ctx)

	if err != nil {
		return "", err
	}

	return ossBucketObjectAttributes.MediaLink, nil
}

func (ossStorage *GoogleOSS) GetObjectCDNURL(bucket string, object *Object) (string, error) {
	return getObjectCDNURL(
		ossStorage.buckets,
		bucket,
		ossStorage.env,
		object.StorageKey), nil
}

func (ossStorage *GoogleOSS) GetObjectSignedURL(bucket string, object *Object, expires time.Time) (string, error) {
	return storage.SignedURL(getBucketName(ossStorage.buckets, bucket, ossStorage.env),
		object.StorageKey,
		&storage.SignedURLOptions{
			GoogleAccessID: ossStorage.accessID,
			PrivateKey:     ossStorage.privateKey,
			Method:         http.MethodGet,
			Expires:        expires,
		})
}

func (ossStorage *GoogleOSS) GetObjectBase64Content(bucket string, object *Object) (string, error) {
	bucketName := getBucketName(ossStorage.buckets, bucket, ossStorage.env)

	reader, err := ossStorage.client.Bucket(bucketName).Object(object.StorageKey).NewReader(context.Background())

	if err != nil {
		return "", err
	}

	content, err := ioutil.ReadAll(reader)

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(content), nil
}

func (ossStorage *GoogleOSS) UploadObjectFromFile(ormService *beeorm.Engine, bucket, localFile string) (Object, error) {
	fileContent, ext, err := readContentFile(localFile)

	if err != nil {
		return Object{}, err
	}

	return ossStorage.UploadObjectFromByte(ormService, bucket, fileContent, ext)
}

func (ossStorage *GoogleOSS) UploadObjectFromBase64(ormService *beeorm.Engine, bucket, base64content, extension string) (Object, error) {
	byteData, err := base64.StdEncoding.DecodeString(base64content)

	if err != nil {
		return Object{}, err
	}

	return ossStorage.UploadObjectFromByte(ormService, bucket, byteData, extension)
}

func (ossStorage *GoogleOSS) UploadImageFromBase64(ormService *beeorm.Engine, bucket, base64image, extension string) (Object, error) {
	byteData, err := base64.StdEncoding.DecodeString(base64image)

	if err != nil {
		return Object{}, err
	}

	return ossStorage.UploadObjectFromByte(ormService, bucket, byteData, extension)
}

func (ossStorage *GoogleOSS) UploadImageFromFile(ormService *beeorm.Engine, bucket, localFile string) (Object, error) {
	return ossStorage.UploadObjectFromFile(ormService, bucket, localFile)
}

func (ossStorage *GoogleOSS) UploadObjectFromByte(ormService *beeorm.Engine, bucket string, objectContent []byte, extension string) (Object, error) {
	storageCounter := getStorageCounter(ormService, ossStorage.buckets, bucket)

	objectKey := ossStorage.getObjectKey(storageCounter, extension)

	bucketName := getBucketName(ossStorage.buckets, bucket, ossStorage.env)

	ossBucketObject := ossStorage.client.Bucket(bucketName).Object(objectKey).NewWriter(ossStorage.ctx)

	//TODO Remove
	ossStorage.setObjectContentType(ossBucketObject, extension)

	_, err := ossBucketObject.Write(objectContent)

	if err != nil {
		return Object{}, err
	}

	err = ossBucketObject.Close()

	if err != nil {
		return Object{}, err
	}

	return Object{
		ID:         storageCounter,
		StorageKey: objectKey,
	}, nil
}

func (ossStorage *GoogleOSS) DeleteObject(_ string, _ *Object) error {
	panic("not implemented")
}

func (ossStorage *GoogleOSS) CreateObjectFromKey(ormService *beeorm.Engine, bucket, key string) Object {
	return Object{
		ID:         getStorageCounter(ormService, ossStorage.buckets, bucket),
		StorageKey: key,
	}
}

func (ossStorage *GoogleOSS) getObjectKey(storageCounter uint64, fileExtension string) string {
	return strconv.FormatUint(storageCounter, 10) + fileExtension
}

//TODO Remove
func (ossStorage *GoogleOSS) setObjectContentType(writer *storage.Writer, extension string) {
	if writer == nil {
		return
	}
	if extension == ".svg" && writer.ObjectAttrs.ContentType == "" {
		writer.ObjectAttrs.ContentType = "image/svg+xml"
	}
}

func (ossStorage *GoogleOSS) GetUploaderBucketConfig() *BucketConfig {
	return getBucketEnvConfig(ossStorage.buckets, ossStorage.uploaderBucket, ossStorage.env)
}
