package main

//
//import (
//	"net/http"
//	"os"
//	"testing"
//	"time"
//
//	clock "github.com/coretrix/hitrix/service/component/clock/mocks"
//	"github.com/coretrix/hitrix/service/component/oss"
//	ossMock "github.com/coretrix/hitrix/service/component/oss/mocks"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//
//	userMock "github.com/coretrix/mobzzo-backend/pkg/ioc/service/user/mocks"
//	"github.com/coretrix/mobzzo-backend/tests"
//)
//
//func TestPostUploadImageAction(t *testing.T) {
//	t.Run("simple create", func(t *testing.T) {
//		link := "https://somelink.com"
//		storageKey := "test.jpg"
//		fileID := uint64(1)
//
//		fakeFile := oss.Object{
//			ID:         fileID,
//			StorageKey: storageKey,
//		}
//
//		fakeOSService := &ossMock.FakeOSSClient{}
//		fakeOSService.On("UploadImageFromFile", mock.Anything, mock.Anything, mock.Anything).Return(fakeFile)
//		fakeOSService.On("GetObjectURL", mock.Anything, mock.Anything).Return(link)
//
//		fakeUserService := &userMock.FakeUserService{}
//		fakeUserService.On("GetSession").Return(&userService.Session{
//			AccessKey: "access key",
//			User:      &entity.UserEntity{},
//		}, true)
//
//		fakeUserService.On("MustGetSession").Return(&userService.Session{
//			AccessKey: "access key",
//			User:      &entity.UserEntity{ID: 1},
//		})
//
//		now := time.Unix(1, 0)
//		fakeClock := &clock.FakeSysClock{}
//		fakeClock.On("Now").Return(now)
//
//		mockServices := &tests.IoCMocks{
//			OSService:    fakeOSService,
//			ClockService: fakeClock,
//			UserService:  fakeUserService,
//		}
//
//		e := tests.CreateContextWebAPI(t, mockServices)
//
//		imageFile, err := os.Open("../../../fixtures/test.jpeg")
//		assert.Nil(t, err)
//
//		body := map[string]interface{}{
//			"file":      imageFile,
//			"namespace": entity.FileNamespaceSubscriptionDocuments.String(),
//		}
//
//		got := &file.File{}
//
//		err = tests.SendHTTPRequestWithMultipartBody(e, http.MethodPost, "/v1/file/upload/", body, true, got)
//
//		assert.Nil(t, err)
//		assert.NotNil(t, got)
//		assert.Equal(t, got.URL, link)
//		assert.Equal(t, got.ID, fileID)
//	})
//}
