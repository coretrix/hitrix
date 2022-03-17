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

func (t *FakeOSSClient) GetObjectURL(bucket oss.Bucket, object *oss.Object) (string, error) {
	return t.Called(bucket, object).Get(0).(string), nil
}

func (t *FakeOSSClient) GetObjectOSSURL(bucket oss.Bucket, object *oss.Object) (string, error) {
	return t.Called(bucket, object).Get(0).(string), nil
}

func (t *FakeOSSClient) GetObjectCDNURL(bucket oss.Bucket, object *oss.Object) (string, error) {
	return t.Called(bucket, object).Get(0).(string), nil
}

func (t *FakeOSSClient) GetObjectSignedURL(bucket oss.Bucket, object *oss.Object, expires time.Time) (string, error) {
	return t.Called(bucket, object, expires).Get(0).(string), nil
}

func (t *FakeOSSClient) GetObjectBase64Content(bucket oss.Bucket, object *oss.Object) (string, error) {
	return t.Called(bucket, object).Get(0).(string), nil
}

func (t *FakeOSSClient) UploadObjectFromFile(_ *beeorm.Engine, bucket oss.Bucket, path oss.Namespace, localFile string) (oss.Object, error) {
	return t.Called(bucket, path, localFile).Get(0).(oss.Object), nil
}

func (t *FakeOSSClient) UploadObjectFromBase64(_ *beeorm.Engine, bucket oss.Bucket, path oss.Namespace, content, extension string) (oss.Object, error) {
	return t.Called(bucket, path, content, extension).Get(0).(oss.Object), nil
}

func (t *FakeOSSClient) UploadObjectFromByte(_ *beeorm.Engine, bucket oss.Bucket, path oss.Namespace, content []byte, fileName *string, extension string) (oss.Object, error) {
	return t.Called(bucket, path, content, fileName, extension).Get(0).(oss.Object), nil
}

func (t *FakeOSSClient) UploadImageFromFile(_ *beeorm.Engine, bucket oss.Bucket, path oss.Namespace, localFile string) (oss.Object, error) {
	return t.Called(bucket, path, localFile).Get(0).(oss.Object), nil
}

func (t *FakeOSSClient) UploadImageFromBase64(_ *beeorm.Engine, bucket oss.Bucket, path oss.Namespace, image, extension string) (oss.Object, error) {
	return t.Called(bucket, path, image, extension).Get(0).(oss.Object), nil
}

func (t *FakeOSSClient) DeleteObject(_ oss.Bucket, _ *oss.Object) error {
	return nil
}

func (t *FakeOSSClient) CreateObjectFromKey(_ *beeorm.Engine, _ oss.Bucket, _ string) oss.Object {
	return oss.Object{}
}

func (t *FakeOSSClient) GetBucketConfig(_ oss.Bucket) *oss.BucketConfig {
	return nil
}
