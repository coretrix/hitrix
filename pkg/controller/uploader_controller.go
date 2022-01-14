package controller

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/service"
)

type UploaderController struct {
}

func (controller *UploaderController) PostFileAction(c *gin.Context) {
	uploaderService := service.DI().Uploader()

	appService := service.DI().App()
	if !appService.IsInLocalMode() {
		c.Request.Header.Set("X-Forwarded-Proto", "https")
	}

	isPartial := c.Request.Header.Get("Upload-Concat") == "partial"

	var rec http.ResponseWriter
	switch isPartial {
	case true:
		rec = c.Writer
	default:
		rec = &httptest.ResponseRecorder{}
	}

	uploaderService.PostFile(rec, c.Request)
	if isPartial {
		return
	}

	for name, values := range rec.Header() {
		c.Writer.Header()[name] = values
	}

	c.Writer.WriteHeader(rec.(*httptest.ResponseRecorder).Code)
}

func (controller *UploaderController) GetFileAction(c *gin.Context) {
	uploaderService := service.DI().Uploader()

	appService := service.DI().App()
	if !appService.IsInLocalMode() {
		c.Request.Header.Set("X-Forwarded-Proto", "https")
	}

	uploaderService.GetFile(c.Writer, c.Request)
}

func (controller *UploaderController) HeadFile(c *gin.Context) {
	uploaderService := service.DI().Uploader()

	appService := service.DI().App()
	if !appService.IsInLocalMode() {
		c.Request.Header.Set("X-Forwarded-Proto", "https")
	}

	uploaderService.HeadFile(c.Writer, c.Request)
}

func (controller *UploaderController) PatchFile(c *gin.Context) {
	uploaderService := service.DI().Uploader()

	appService := service.DI().App()
	if !appService.IsInLocalMode() {
		c.Request.Header.Set("X-Forwarded-Proto", "https")
	}

	uploaderService.PatchFile(c.Writer, c.Request)
}

func (controller *UploaderController) DeleteFile(c *gin.Context) {
	uploaderService := service.DI().Uploader()

	appService := service.DI().App()
	if !appService.IsInLocalMode() {
		c.Request.Header.Set("X-Forwarded-Proto", "https")
	}

	uploaderService.DelFile(c.Writer, c.Request)
}
