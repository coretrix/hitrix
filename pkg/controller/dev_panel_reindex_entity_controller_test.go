package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestReindexEntityRejectsInvalidIDsBeforeAccessingServices(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, body := range []string{
		`{}`, `{"ID":"0"}`, `{"ID":"-1"}`, `{"ID":"1.5"}`,
		`{"ID":"18446744073709551616"}`, `{"ID":"+42"}`,
		`{"ID":"042"}`, `{"ID":42}`, `{"ID":"42 OR 1=1"}`, `not json`,
	} {
		t.Run(body, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(recorder)
			context.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
			context.Request.Header.Set("Content-Type", "application/json")
			controller := &DevPanelController{}
			controller.PostRedisSearchReindexEntity(context)
			if recorder.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d: %s", recorder.Code, recorder.Body.String())
			}
		})
	}
}
