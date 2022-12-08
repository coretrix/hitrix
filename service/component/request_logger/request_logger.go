package requestlogger

import (
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/latolukasz/beeorm"
	"net/http"
)

type IRequestLogger interface {
	LogRequest(ormService *beeorm.Engine, appName, url string, request *http.Request, contentType string) *entity.RequestLoggerEntity
	LogResponse(ormService *beeorm.Engine, requestLoggerEntity *entity.RequestLoggerEntity, responseBody []byte, status int)
}
