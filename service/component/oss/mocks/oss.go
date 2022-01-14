package mocks

import (
	"time"

	"github.com/latolukasz/beeorm"
	"github.com/stretchr/testify/mock"

	"github.com/coretrix/hitrix/service/component/oss"
)

type FakeOSSClient struct {
	mock.Mock
}

func (t *FakeOSSClient) GetClient() interface{} {
	return nil
}

func (t *FakeOSSClient) GetObjectURL(bucket string, object *oss.Object) (string, error) {
	return t.Called(bucket, object).Get(0).(string), nil
}

func (t *FakeOSSClient) GetObjectOSSURL(bucket string, object *oss.Object) (string, error) {
	return t.Called(bucket, object).Get(0).(string), nil
}

func (t *FakeOSSClient) GetObjectCDNURL(bucket string, object *oss.Object) (string, error) {
	return t.Called(bucket, object).Get(0).(string), nil
}

func (t *FakeOSSClient) GetObjectSignedURL(bucket string, object *oss.Object, expires time.Time) (string, error) {
	return t.Called(bucket, object, expires).Get(0).(string), nil
}

func (t *FakeOSSClient) GetObjectBase64Content(bucket string, object *oss.Object) (string, error) {
	return t.Called(bucket, object).Get(0).(string), nil
}

func (t *FakeOSSClient) UploadObjectFromFile(_ *beeorm.Engine, bucket, path, localFile string) (oss.Object, error) {
	return t.Called(bucket, path, localFile).Get(0).(oss.Object), nil
}

func (t *FakeOSSClient) UploadObjectFromBase64(_ *beeorm.Engine, bucket, path, content, extension string) (oss.Object, error) {
	return t.Called(bucket, path, content, extension).Get(0).(oss.Object), nil
}

func (t *FakeOSSClient) UploadObjectFromByte(_ *beeorm.Engine, bucket, path string, content []byte, extension string) (oss.Object, error) {
	return t.Called(bucket, path, content, extension).Get(0).(oss.Object), nil
}

func (t *FakeOSSClient) UploadImageFromFile(_ *beeorm.Engine, bucket, path, localFile string) (oss.Object, error) {
	return t.Called(bucket, path, localFile).Get(0).(oss.Object), nil
}

func (t *FakeOSSClient) UploadImageFromBase64(_ *beeorm.Engine, bucket, path, image, extension string) (oss.Object, error) {
	return t.Called(bucket, path, image, extension).Get(0).(oss.Object), nil
}

func (t *FakeOSSClient) DeleteObject(_ string, _ *oss.Object) error {
	return nil
}

func (t *FakeOSSClient) CreateObjectFromKey(_ *beeorm.Engine, _, _ string) oss.Object {
	return oss.Object{}
}

func (t *FakeOSSClient) GetUploaderBucketConfig() *oss.BucketConfig {
	return &oss.BucketConfig{}
}
