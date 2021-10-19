package mocks

import (
	"time"

	"github.com/coretrix/hitrix/service/component/oss"

	"github.com/latolukasz/beeorm"
	"github.com/stretchr/testify/mock"
)

type FakeOSSClient struct {
	mock.Mock
}

func (t *FakeOSSClient) GetObjectURL(bucket string, object *oss.Object) (string, error) {
	return t.Called(bucket, object).Get(0).(string), t.Called(bucket, object).Get(0).(error)
}

func (t *FakeOSSClient) GetObjectOSSURL(bucket string, object *oss.Object) (string, error) {
	return t.Called(bucket, object).Get(0).(string), t.Called(bucket, object).Get(0).(error)
}

func (t *FakeOSSClient) GetObjectCDNURL(bucket string, object *oss.Object) (string, error) {
	return t.Called(bucket, object).Get(0).(string), t.Called(bucket, object).Get(0).(error)
}

func (t *FakeOSSClient) GetObjectSignedURL(bucket string, object *oss.Object, expires time.Time) (string, error) {
	return t.Called(bucket, object, expires).Get(0).(string), t.Called(bucket, object).Get(0).(error)
}

func (t *FakeOSSClient) GetObjectBase64Content(bucket string, object *oss.Object) (string, error) {
	return t.Called(bucket, object).Get(0).(string), t.Called(bucket, object).Get(0).(error)
}

func (t *FakeOSSClient) UploadObjectFromFile(_ *beeorm.Engine, bucket, localFile string) (oss.Object, error) {
	return t.Called(bucket, localFile).Get(0).(oss.Object), t.Called(bucket, localFile).Get(0).(error)
}

func (t *FakeOSSClient) UploadObjectFromBase64(_ *beeorm.Engine, bucket, content, extension string) (oss.Object, error) {
	return t.Called(bucket, content, extension).Get(0).(oss.Object), t.Called(bucket, content, extension).Get(0).(error)
}

func (t *FakeOSSClient) UploadObjectFromByte(_ *beeorm.Engine, bucket string, content []byte, extension string) (oss.Object, error) {
	return t.Called(bucket, content, extension).Get(0).(oss.Object), t.Called(bucket, content, extension).Get(0).(error)
}

func (t *FakeOSSClient) UploadImageFromFile(_ *beeorm.Engine, bucket, localFile string) (oss.Object, error) {
	return t.Called(bucket, localFile).Get(0).(oss.Object), t.Called(bucket, localFile).Get(0).(error)
}

func (t *FakeOSSClient) UploadImageFromBase64(_ *beeorm.Engine, bucket, image, extension string) (oss.Object, error) {
	return t.Called(bucket, image, extension).Get(0).(oss.Object), t.Called(bucket, image, extension).Get(0).(error)
}

func (t *FakeOSSClient) DeleteObject(_ string, _ *oss.Object) error {
	return nil
}
