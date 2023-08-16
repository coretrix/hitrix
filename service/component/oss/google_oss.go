package oss

import (
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/pkg/helper"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
)

type GoogleOSS struct {
	client       *storage.Client
	clockService clock.IClock
	ctx          context.Context
	buckets      bucketsConfig
	namespaces   namespacesConfig
	accessID     string
	privateKey   []byte
}

func NewGoogleOSS(configService config.IConfig, clockService clock.IClock, namespaces Namespaces) (IProvider, error) {
	ctx := context.Background()

	if !helper.ExistsInDir(".oss.json", configService.GetFolderPath()) {
		return nil, errors.New(configService.GetFolderPath() + "/.oss.json does not exists")
	}

	credentialsFile := configService.GetFolderPath() + "/.oss.json"

	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, err
	}

	jwtCredentialsJSONString, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, err
	}

	jwtConfig, _ := google.JWTConfigFromJSON(jwtCredentialsJSONString)

	bucketsConfiguration, namespacesConfiguration := loadConfig(configService, namespaces)

	return &GoogleOSS{
		client:       client,
		clockService: clockService,
		ctx:          ctx,
		buckets:      bucketsConfiguration,
		namespaces:   namespacesConfiguration,
		accessID:     jwtConfig.Email,
		privateKey:   jwtConfig.PrivateKey,
	}, nil
}

func (ossStorage *GoogleOSS) GetBucketConfig(bucket Bucket) *BucketConfig {
	return getBucketConfig(ossStorage.buckets[bucket])
}

func (ossStorage *GoogleOSS) GetClient() interface{} {
	return ossStorage.client
}

func (ossStorage *GoogleOSS) GetObjectURL(namespace Namespace, object *entity.FileObject) (string, error) {
	cdnURL, err := ossStorage.GetObjectCDNURL(namespace, object)
	if err != nil {
		return "", err
	}

	if cdnURL != "" {
		return cdnURL, nil
	}

	return ossStorage.GetObjectOSSURL(namespace, object)
}

func (ossStorage *GoogleOSS) GetObjectOSSURL(namespace Namespace, object *entity.FileObject) (string, error) {
	if object == nil {
		return "", errors.New("nil file object")
	}

	bucketConfig, err := ossStorage.namespaces.getBucketConfig(namespace)
	if err != nil {
		return "", err
	}

	ossBucketObjectAttributes, err := ossStorage.client.Bucket(bucketConfig.Name).Object(object.StorageKey).Attrs(ossStorage.ctx)
	if err != nil {
		return "", err
	}

	return ossBucketObjectAttributes.MediaLink, nil
}

func (ossStorage *GoogleOSS) GetObjectCDNURL(namespace Namespace, object *entity.FileObject) (string, error) {
	if object == nil {
		return "", errors.New("nil file object")
	}

	bucketConfig, err := ossStorage.namespaces.getBucketConfig(namespace)
	if err != nil {
		return "", err
	}

	return getObjectCDNURL(bucketConfig, object.StorageKey), nil
}

func (ossStorage *GoogleOSS) GetNamespaceBucketConfig(namespace Namespace) (*BucketConfig, error) {
	bucketConfig, err := ossStorage.namespaces.getBucketConfig(namespace)

	if err != nil {
		return nil, err
	}

	return bucketConfig, nil
}

func (ossStorage *GoogleOSS) GetObjectSignedURL(namespace Namespace, object *entity.FileObject, expires time.Time) (string, error) {
	if object == nil {
		return "", errors.New("nil file object")
	}

	bucketConfig, err := ossStorage.namespaces.getBucketConfig(namespace)
	if err != nil {
		return "", err
	}

	return storage.SignedURL(
		bucketConfig.Name,
		object.StorageKey,
		&storage.SignedURLOptions{
			GoogleAccessID: ossStorage.accessID,
			PrivateKey:     ossStorage.privateKey,
			Method:         http.MethodGet,
			Expires:        expires,
		})
}

func (ossStorage *GoogleOSS) GetObjectBase64Content(namespace Namespace, object *entity.FileObject) (string, error) {
	if object == nil {
		return "", errors.New("nil file object")
	}

	bucketConfig, err := ossStorage.namespaces.getBucketConfig(namespace)
	if err != nil {
		return "", err
	}

	reader, err := ossStorage.client.Bucket(bucketConfig.Name).Object(object.StorageKey).NewReader(context.Background())
	if err != nil {
		return "", err
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(content), nil
}

func (ossStorage *GoogleOSS) UploadObjectFromFile(ormService *datalayer.DataLayer, namespace Namespace, localFile string) (entity.FileObject, error) {
	fileContent, ext, err := readContentFile(localFile)
	if err != nil {
		return entity.FileObject{}, err
	}

	return ossStorage.UploadObjectFromByte(ormService, namespace, fileContent, ext)
}

func (ossStorage *GoogleOSS) UploadObjectFromBase64(
	ormService *datalayer.DataLayer,
	namespace Namespace,
	base64content,
	extension string,
) (entity.FileObject, error) {
	byteData, err := base64.StdEncoding.DecodeString(base64content)
	if err != nil {
		return entity.FileObject{}, err
	}

	return ossStorage.UploadObjectFromByte(ormService, namespace, byteData, extension)
}

func (ossStorage *GoogleOSS) UploadImageFromBase64(
	ormService *datalayer.DataLayer,
	namespace Namespace,
	base64image string,
	extension string,
) (entity.FileObject, error) {
	byteData, err := base64.StdEncoding.DecodeString(base64image)
	if err != nil {
		return entity.FileObject{}, err
	}

	return ossStorage.UploadObjectFromByte(ormService, namespace, byteData, extension)
}

func (ossStorage *GoogleOSS) UploadImageFromFile(ormService *datalayer.DataLayer, namespace Namespace, localFile string) (entity.FileObject, error) {
	return ossStorage.UploadObjectFromFile(ormService, namespace, localFile)
}

func (ossStorage *GoogleOSS) UploadObjectFromByte(
	ormService *datalayer.DataLayer,
	namespace Namespace,
	objectContent []byte,
	extension string,
) (entity.FileObject, error) {
	bucketConfig, err := ossStorage.namespaces.getBucketConfig(namespace)
	if err != nil {
		return entity.FileObject{}, err
	}

	storageCounter := getStorageCounter(ossStorage.ctx, ormService, bucketConfig)

	objectKey := ossStorage.getObjectKey(namespace, storageCounter, extension)

	ossBucketObject := ossStorage.client.Bucket(bucketConfig.Name).Object(objectKey).NewWriter(ossStorage.ctx)

	//TODO Remove
	ossStorage.setObjectContentType(ossBucketObject, extension)

	_, err = ossBucketObject.Write(objectContent)
	if err != nil {
		return entity.FileObject{}, err
	}

	err = ossBucketObject.Close()
	if err != nil {
		return entity.FileObject{}, err
	}

	return entity.FileObject{
		ID:         storageCounter,
		StorageKey: objectKey,
	}, nil
}

func (ossStorage *GoogleOSS) DeleteObject(namespace Namespace, object *entity.FileObject) error {
	bucketConfig, err := ossStorage.namespaces.getBucketConfig(namespace)
	if err != nil {
		return err
	}

	return ossStorage.client.Bucket(bucketConfig.Name).Object(object.StorageKey).Delete(ossStorage.ctx)
}

func (ossStorage *GoogleOSS) getObjectKey(namespace Namespace, storageCounter uint64, fileExtension string) string {
	return namespace.String() + "/" + strconv.FormatUint(storageCounter, 10) + fileExtension
}

// TODO Remove
func (ossStorage *GoogleOSS) setObjectContentType(writer *storage.Writer, extension string) {
	if writer == nil {
		return
	}

	if extension == ".svg" && writer.ObjectAttrs.ContentType == "" {
		writer.ObjectAttrs.ContentType = "image/svg+xml"
	}
}
