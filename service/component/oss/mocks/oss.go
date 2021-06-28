package mocks

import (
	"time"

	"github.com/coretrix/hitrix/service/component/oss"

	"github.com/latolukasz/orm"
	"github.com/stretchr/testify/mock"
)

type FakeOSSClient struct {
	mock.Mock
}

func (t *FakeOSSClient) GetObjectURL(bucket string, object *oss.Object) string {
	return t.Called(bucket, object).String(0)
}

func (t *FakeOSSClient) GetObjectCachedURL(bucket string, object *oss.Object) string {
	return t.Called(bucket, object).String(0)
}

func (t *FakeOSSClient) GetObjectSignedURL(bucket string, object *oss.Object, expires time.Time) string {
	return t.Called(bucket, object, expires).String(0)
}

func (t *FakeOSSClient) UploadObjectFromFile(_ *orm.Engine, bucket, localFile string) oss.Object {
	return t.Called(bucket, localFile).Get(0).(oss.Object)
}

func (t *FakeOSSClient) UploadObjectFromBase64(_ *orm.Engine, bucket, content, extension string) oss.Object {
	return t.Called(bucket, content, extension).Get(0).(oss.Object)
}

func (t *FakeOSSClient) UploadImageFromFile(_ *orm.Engine, bucket, localFile string) oss.Object {
	return t.Called(bucket, localFile).Get(0).(oss.Object)
}

func (t *FakeOSSClient) UploadImageFromBase64(_ *orm.Engine, bucket, image, extension string) oss.Object {
	return t.Called(bucket, image, extension).Get(0).(oss.Object)
}

func (t *FakeOSSClient) GetObjectBase64Content(bucket string, object *oss.Object) (string, error) {
	return t.Called(bucket, object).Get(0).(string), t.Called(bucket, object).Get(0).(error)
}
