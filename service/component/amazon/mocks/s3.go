package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/coretrix/hitrix/datalayer"
	s3 "github.com/coretrix/hitrix/service/component/amazon/storage"
)

type FakeS3Client struct {
	mock.Mock
}

func (t *FakeS3Client) GetObjectCachedURL(bucket string, object *s3.Object) string {
	return t.Called(bucket, object).String(0)
}

func (t *FakeS3Client) GetObjectSignedURL(bucket string, object *s3.Object, expires time.Duration) string {
	return t.Called(bucket, object, expires).String(0)
}

func (t *FakeS3Client) UploadObjectFromFile(_ *datalayer.DataLayer, bucket, localFile string) s3.Object {
	return t.Called(bucket, localFile).Get(0).(s3.Object)
}

func (t *FakeS3Client) UploadObjectFromBase64(_ *datalayer.DataLayer, bucket, content, extension string) s3.Object {
	return t.Called(bucket, content, extension).Get(0).(s3.Object)
}

func (t *FakeS3Client) UploadObjectFromByte(_ *datalayer.DataLayer, bucket string, byteData []byte, extension string) s3.Object {
	return t.Called(bucket, byteData, extension).Get(0).(s3.Object)
}

func (t *FakeS3Client) UploadImageFromFile(_ *datalayer.DataLayer, bucket, localFile string) s3.Object {
	return t.Called(bucket, localFile).Get(0).(s3.Object)
}

func (t *FakeS3Client) UploadImageFromBase64(_ *datalayer.DataLayer, bucket, image, extension string) s3.Object {
	return t.Called(bucket, image, extension).Get(0).(s3.Object)
}

func (t *FakeS3Client) DeleteObject(bucket string, objects ...*s3.Object) bool {
	return t.Called(bucket, objects).Get(0).(bool)
}

func (t *FakeS3Client) GetClient() interface{} {
	return t.Called().Get(0)
}

func (t *FakeS3Client) CreateObjectFromKey(_ *datalayer.DataLayer, bucket, key string) s3.Object {
	return t.Called(bucket, key).Get(0).(s3.Object)
}

func (t *FakeS3Client) GetBucketName(bucket string) string {
	return t.Called(bucket).Get(0).(string)
}
