package storage

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/coretrix/hitrix/service/component/config"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"

	"google.golang.org/api/option"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/oss"

	"github.com/coretrix/hitrix/service/component/app"

	"cloud.google.com/go/storage"
	"github.com/latolukasz/beeorm"
	"golang.org/x/net/context"
)

type GoogleOSS struct {
	client      *storage.Client
	ctx         context.Context
	environment string
	buckets     map[string]uint64
	domain      string
	urlPrefix   string
	jwtConfig   *jwt.Config
}

func NewGoogleOSS(credentialsFile string, environment string, buckets map[string]uint64, configService config.IConfig) *GoogleOSS {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))

	if err != nil {
		panic(err)
	}

	domain, ok := configService.String("oss.domain")
	if !ok {
		panic("oss domain is not set in config")
	}

	urlPrefix := "static-"
	if configURLPrefix, ok := configService.String("oss.url_prefix"); ok {
		urlPrefix = configURLPrefix
	}

	jwtCredentialsJSONString, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		panic("could not read and parse " + credentialsFile)
	}
	jwtConfig, _ := google.JWTConfigFromJSON(jwtCredentialsJSONString)

	return &GoogleOSS{
		client:      client,
		ctx:         ctx,
		environment: environment,
		buckets:     buckets,
		domain:      domain,
		urlPrefix:   urlPrefix,
		jwtConfig:   jwtConfig,
	}
}

func (ossStorage *GoogleOSS) GetObjectURL(bucket string, object *oss.Object) string {
	ossStorage.checkBucket(bucket)

	bucketByEnv := bucket

	if ossStorage.environment != app.ModeProd {
		bucketByEnv += "-" + ossStorage.environment
	}

	envInUrl := ""
	if ossStorage.environment != app.ModeProd {
		envInUrl = ossStorage.environment + "."
	}

	ossBucketObjectAttributes, err := ossStorage.client.Bucket(bucketByEnv).Object(object.StorageKey).Attrs(ossStorage.ctx)

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("https://%s%s%s/%s/%s", ossStorage.urlPrefix, envInUrl, ossStorage.domain, bucketByEnv, ossBucketObjectAttributes.Name)
}

func (ossStorage *GoogleOSS) GetObjectCachedURL(bucket string, object *oss.Object) string {
	ossStorage.checkBucket(bucket)

	bucketByEnv := bucket

	if ossStorage.environment != app.ModeProd {
		bucketByEnv += "-" + ossStorage.environment
	}

	envInUrl := ""
	if ossStorage.environment != app.ModeProd {
		envInUrl = ossStorage.environment + "."
	}

	return fmt.Sprintf("https://%s%s%s/%s/%s", ossStorage.urlPrefix, envInUrl, ossStorage.domain, bucketByEnv, object.CachedURL)
}

func (ossStorage *GoogleOSS) GetObjectSignedURL(bucket string, object *oss.Object, expires time.Time) string {
	ossStorage.checkBucket(bucket)

	bucketByEnv := bucket

	if ossStorage.environment != app.ModeProd {
		bucketByEnv += "-" + ossStorage.environment
	}

	signedURL, err := storage.SignedURL(bucketByEnv, object.StorageKey, &storage.SignedURLOptions{
		GoogleAccessID: ossStorage.jwtConfig.Email,
		PrivateKey:     ossStorage.jwtConfig.PrivateKey,
		Method:         http.MethodGet,
		Expires:        expires,
	})

	if err != nil {
		panic(err)
	}

	return signedURL
}

func (ossStorage *GoogleOSS) UploadObjectFromFile(ormService *beeorm.Engine, bucket, localFile string) oss.Object {
	ossStorage.checkBucket(bucket)

	fileContent, ext := ossStorage.ReadFile(localFile)

	return ossStorage.putObject(ormService, bucket, fileContent, ext)
}

func (ossStorage *GoogleOSS) UploadObjectFromBase64(ormService *beeorm.Engine, bucket, base64content, extension string) oss.Object {
	byteData, err := base64.StdEncoding.DecodeString(base64content)

	if err != nil {
		panic(err)
	}

	return ossStorage.putObject(ormService, bucket, byteData, extension)
}

func (ossStorage *GoogleOSS) UploadImageFromBase64(ormService *beeorm.Engine, bucket, base64image, extension string) oss.Object {
	byteData, err := base64.StdEncoding.DecodeString(base64image)

	if err != nil {
		panic(err)
	}

	return ossStorage.putObject(ormService, bucket, byteData, extension)
}

func (ossStorage *GoogleOSS) UploadImageFromFile(ormService *beeorm.Engine, bucket, localFile string) oss.Object {
	return ossStorage.UploadObjectFromFile(ormService, bucket, localFile)
}

func (ossStorage *GoogleOSS) checkBucket(bucket string) {
	_, ok := ossStorage.buckets[bucket]

	if !ok {
		panic("bucket [" + bucket + "] not found")
	}
}

func (ossStorage *GoogleOSS) putObject(ormService *beeorm.Engine, bucket string, objectContent []byte, extension string) oss.Object {
	storageCounter := ossStorage.getStorageCounter(ormService, bucket)

	objectKey := ossStorage.getObjectKey(storageCounter, extension)

	bucketByEnv := bucket

	if ossStorage.environment != app.ModeProd {
		bucketByEnv += "-" + ossStorage.environment
	}

	ossBucketObject := ossStorage.client.Bucket(bucketByEnv).Object(objectKey).NewWriter(ossStorage.ctx)
	ossStorage.setObjectContentTypeIfNecessary(ossBucketObject, extension)

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

func (ossStorage *GoogleOSS) getStorageCounter(ormService *beeorm.Engine, bucket string) uint64 {
	ossStorage.checkBucket(bucket)

	bucketID, has := ossStorage.buckets[bucket]

	if !has {
		panic("oss bucket [" + bucket + "] id not found")
	}

	googleOSSBucketCounterEntity := &entity.OSSBucketCounterEntity{}

	locker := ormService.GetRedis().GetLocker()
	lock, hasLock := locker.Obtain("locker_google_oss_counters_bucket_"+bucket, 2*time.Second, 5*time.Second)
	defer lock.Release()

	if !hasLock {
		panic("Failed to obtain lock for locker_google_oss_counters_bucket_" + bucket)
	}

	has = ormService.LoadByID(bucketID, googleOSSBucketCounterEntity)
	if !has {
		googleOSSBucketCounterEntity.ID = bucketID
		googleOSSBucketCounterEntity.Counter = 1
	} else {
		googleOSSBucketCounterEntity.Counter = googleOSSBucketCounterEntity.Counter + 1
	}
	ormService.Flush(googleOSSBucketCounterEntity)

	ttl := lock.TTL()
	if ttl == 0 {
		panic("lock lost")
	}

	return googleOSSBucketCounterEntity.Counter
}

func (ossStorage *GoogleOSS) getObjectKey(storageCounter uint64, fileExtension string) string {
	return strconv.FormatUint(storageCounter, 10) + fileExtension
}

func (ossStorage *GoogleOSS) GetObjectBase64Content(bucket string, object *oss.Object) (string, error) {
	ossStorage.checkBucket(bucket)

	bucketByEnv := bucket

	if ossStorage.environment != app.ModeProd {
		bucketByEnv += "-" + ossStorage.environment
	}

	reader, err := ossStorage.client.Bucket(bucketByEnv).Object(object.StorageKey).NewReader(context.Background())
	if err != nil {
		return "", err
	}

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(content), nil
}

func (ossStorage *GoogleOSS) ReadFile(localFile string) ([]byte, string) {
	fileContent, err := ioutil.ReadFile(localFile)

	if err != nil {
		panic(err)
	}

	return fileContent, filepath.Ext(localFile)
}

func (ossStorage *GoogleOSS) setObjectContentTypeIfNecessary(writer *storage.Writer, extension string) {
	if writer == nil {
		return
	}
	if extension == ".svg" && writer.ObjectAttrs.ContentType == "" {
		writer.ObjectAttrs.ContentType = "image/svg+xml"
	}
}
