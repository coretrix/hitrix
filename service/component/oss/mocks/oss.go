package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/oss"
)

type FakeOSSClient struct {
	mock.Mock
}

func (t *FakeOSSClient) GetClient() interface{} {
	return nil
}

func (t *FakeOSSClient) GetObjectURL(namespace oss.Namespace, object *entity.FileObject) (string, error) {
	return t.Called(namespace, object).Get(0).(string), nil
}

func (t *FakeOSSClient) GetObjectOSSURL(namespace oss.Namespace, object *entity.FileObject) (string, error) {
	return t.Called(namespace, object).Get(0).(string), nil
}

func (t *FakeOSSClient) GetObjectCDNURL(namespace oss.Namespace, object *entity.FileObject) (string, error) {
	return t.Called(namespace, object).Get(0).(string), nil
}

func (t *FakeOSSClient) GetObjectSignedURL(namespace oss.Namespace, object *entity.FileObject, expires time.Time) (string, error) {
	return t.Called(namespace, object, expires).Get(0).(string), nil
}

func (t *FakeOSSClient) GetObjectBase64Content(namespace oss.Namespace, object *entity.FileObject) (string, error) {
	return t.Called(namespace, object).Get(0).(string), nil
}

func (t *FakeOSSClient) UploadObjectFromFile(_ *datalayer.DataLayer, namespace oss.Namespace, localFile string) (entity.FileObject, error) {
	return t.Called(namespace, localFile).Get(0).(entity.FileObject), nil
}

func (t *FakeOSSClient) UploadObjectFromBase64(
	_ *datalayer.DataLayer,
	namespace oss.Namespace,
	content string,
	extension string,
) (entity.FileObject, error) {
	return t.Called(namespace, content, extension).Get(0).(entity.FileObject), nil
}

func (t *FakeOSSClient) UploadObjectFromByte(
	_ *datalayer.DataLayer,
	namespace oss.Namespace,
	content []byte,
	extension string,
) (entity.FileObject, error) {
	return t.Called(namespace, content, extension).Get(0).(entity.FileObject), nil
}

func (t *FakeOSSClient) UploadImageFromFile(_ *datalayer.DataLayer, namespace oss.Namespace, localFile string) (entity.FileObject, error) {
	return t.Called(namespace, localFile).Get(0).(entity.FileObject), nil
}

func (t *FakeOSSClient) UploadImageFromBase64(_ *datalayer.DataLayer, namespace oss.Namespace, image, extension string) (entity.FileObject, error) {
	return t.Called(namespace, image, extension).Get(0).(entity.FileObject), nil
}

func (t *FakeOSSClient) DeleteObject(namespace oss.Namespace, object *entity.FileObject) error {
	return t.Called(namespace, object).Error(0)
}

func (t *FakeOSSClient) CreateObjectFromKey(_ *datalayer.DataLayer, namespace oss.Namespace, key string) entity.FileObject {
	return t.Called(namespace, key).Get(0).(entity.FileObject)
}

func (t *FakeOSSClient) GetBucketConfig(bucket oss.Bucket) *oss.BucketConfig {
	return t.Called(bucket).Get(0).(*oss.BucketConfig)
}

func (t *FakeOSSClient) GetNamespaceBucketConfig(namespace oss.Namespace) (*oss.BucketConfig, error) {
	return t.Called(namespace).Get(0).(*oss.BucketConfig), nil
}
