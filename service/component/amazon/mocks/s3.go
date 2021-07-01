package mocks

import (
	"time"

	s3 "github.com/coretrix/hitrix/service/component/amazon/storage"

	s3sdk "github.com/aws/aws-sdk-go/service/s3"
	"github.com/latolukasz/orm"
	"github.com/stretchr/testify/mock"
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

func (t *FakeS3Client) UploadObjectFromFile(_ *orm.Engine, bucket, localFile string) s3.Object {
	return t.Called(bucket, localFile).Get(0).(s3.Object)
}

func (t *FakeS3Client) UploadObjectFromBase64(_ *orm.Engine, bucket, content, extension string) s3.Object {
	return t.Called(bucket, content, extension).Get(0).(s3.Object)
}

func (t *FakeS3Client) UploadImageFromFile(_ *orm.Engine, bucket, localFile string) s3.Object {
	return t.Called(bucket, localFile).Get(0).(s3.Object)
}

func (t *FakeS3Client) UploadImageFromBase64(_ *orm.Engine, bucket, image, extension string) s3.Object {
	return t.Called(bucket, image, extension).Get(0).(s3.Object)
}

func (t *FakeS3Client) DeleteObject(bucket string, objects ...*s3.Object) bool {
	return t.Called(bucket, objects).Get(0).(bool)
}

func (t *FakeS3Client) GetClient() *s3sdk.S3 {
	return t.Called().Get(0).(*s3sdk.S3)
}
