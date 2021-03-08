package storage

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/oss"

	"github.com/coretrix/hitrix/service/component/app"

	"google.golang.org/api/option"

	"cloud.google.com/go/storage"
	"github.com/summer-solutions/orm"
	"golang.org/x/net/context"
)

type GoogleOSS struct {
	client      *storage.Client
	ctx         context.Context
	environment string
	buckets     map[string]uint64
}

func NewGoogleOSS(credentialsFile string, environment string, buckets map[string]uint64) *GoogleOSS {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))

	if err != nil {
		panic(err)
	}

	return &GoogleOSS{
		client:      client,
		ctx:         ctx,
		environment: environment,
		buckets:     buckets,
	}
}

func (ossStorage *GoogleOSS) GetObjectURL(bucket string, object *oss.Object) string {
	ossStorage.checkBucket(bucket)

	bucketByEnv := bucket

	if ossStorage.environment != app.ModeProd {
		bucketByEnv += "-" + ossStorage.environment
	}

	ossBucketObjectAttributes, err := ossStorage.client.Bucket(bucketByEnv).Object(object.StorageKey).Attrs(ossStorage.ctx)

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("https://static-%s.hymn.tv/%s/%s", ossStorage.environment, bucketByEnv, ossBucketObjectAttributes.Name)
}

func (ossStorage *GoogleOSS) GetObjectCachedURL(bucket string, object *oss.Object) string {
	ossStorage.checkBucket(bucket)

	bucketByEnv := bucket

	if ossStorage.environment != app.ModeProd {
		bucketByEnv += "-" + ossStorage.environment
	}

	return fmt.Sprintf("https://static-%s.hymn.tv/%s/%s", ossStorage.environment, bucketByEnv, object.CachedURL)
}

func (ossStorage *GoogleOSS) GetObjectSignedURL(bucket string, object *oss.Object, expires time.Time) string {
	ossStorage.checkBucket(bucket)

	bucketByEnv := bucket

	if ossStorage.environment != app.ModeProd {
		bucketByEnv += "-" + ossStorage.environment
	}

	signedURL, err := storage.SignedURL(bucketByEnv, object.StorageKey, &storage.SignedURLOptions{
		GoogleAccessID: "", //todo anton
		PrivateKey:     nil,
		Method:         http.MethodGet,
		Expires:        expires,
	})

	if err != nil {
		panic(err)
	}

	return signedURL
}

func (ossStorage *GoogleOSS) UploadObjectFromFile(ormService *orm.Engine, bucket, localFile string) oss.Object {
	ossStorage.checkBucket(bucket)

	fileContent, ext := ossStorage.ReadFile(localFile)

	return ossStorage.putObject(ormService, bucket, fileContent, ext)
}

func (ossStorage *GoogleOSS) UploadObjectFromBase64(ormService *orm.Engine, bucket, base64content, extension string) oss.Object {
	byteData, err := base64.StdEncoding.DecodeString(base64content)

	if err != nil {
		panic(err)
	}

	return ossStorage.putObject(ormService, bucket, byteData, extension)
}

func (ossStorage *GoogleOSS) UploadImageFromBase64(ormService *orm.Engine, bucket, base64image, extension string) oss.Object {
	byteData, err := base64.StdEncoding.DecodeString(base64image)

	if err != nil {
		panic(err)
	}

	return ossStorage.putObject(ormService, bucket, byteData, extension)
}

func (ossStorage *GoogleOSS) UploadImageFromFile(ormService *orm.Engine, bucket, localFile string) oss.Object {
	return ossStorage.UploadObjectFromFile(ormService, bucket, localFile)
}

func (ossStorage *GoogleOSS) checkBucket(bucket string) {
	_, ok := ossStorage.buckets[bucket]

	if !ok {
		panic("bucket [" + bucket + "] not found")
	}
}

func (ossStorage *GoogleOSS) putObject(ormService *orm.Engine, bucket string, objectContent []byte, extension string) oss.Object {
	storageCounter := ossStorage.getStorageCounter(ormService, bucket)

	objectKey := ossStorage.getObjectKey(storageCounter, extension)

	bucketByEnv := bucket

	if ossStorage.environment != app.ModeProd {
		bucketByEnv += "-" + ossStorage.environment
	}

	ossBucketObject := ossStorage.client.Bucket(bucketByEnv).Object(objectKey).NewWriter(ossStorage.ctx)

	_, err := ossBucketObject.Write(objectContent)

	if err != nil {
		panic("ossBucketObjectWrite:" + err.Error())
	}

	err = ossBucketObject.Close()

	if err != nil {
		panic("ossBucketObjectClose:" + err.Error())
	}

	return oss.Object{
		ID:         storageCounter,
		StorageKey: objectKey,
		CachedURL:  ossBucketObject.Name,
	}
}

func (ossStorage *GoogleOSS) getStorageCounter(ormService *orm.Engine, bucket string) uint64 {
	ossStorage.checkBucket(bucket)

	bucketID, has := ossStorage.buckets[bucket]

	if !has {
		panic("oss bucket [" + bucket + "] id not found")
	}

	googleOSSBucketCounterEntity := &entity.OSSBucketCounterEntity{}

	//todo waiting for lukasz to improve orm
	has = ormService.LoadByID(bucketID, googleOSSBucketCounterEntity)
	if !has {
		googleOSSBucketCounterEntity.ID = bucketID
		googleOSSBucketCounterEntity.Counter = 1
	} else {
		googleOSSBucketCounterEntity.Counter = googleOSSBucketCounterEntity.Counter + 1
	}
	flusher := ormService.NewFlusher()
	flusher.Track(googleOSSBucketCounterEntity)
	flusher.FlushWithLock("default", "locker_google_oss_counters_bucket_"+bucket, 2*time.Second, 5*time.Second)

	return googleOSSBucketCounterEntity.Counter
}

func (ossStorage *GoogleOSS) getObjectKey(storageCounter uint64, fileExtension string) string {
	return strconv.FormatUint(storageCounter, 10) + fileExtension
}

func (ossStorage *GoogleOSS) ReadFile(localFile string) ([]byte, string) {
	fileContent, err := ioutil.ReadFile(localFile)

	if err != nil {
		panic(err)
	}

	return fileContent, filepath.Ext(localFile)
}
