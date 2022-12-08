package requestlogger

import (
	"net/http"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
)

type IRequestLogger interface {
	LogRequest(ormService *beeorm.Engine, appName, url string, request *http.Request, contentType string) *entity.RequestLoggerEntity
	LogResponse(ormService *beeorm.Engine, requestLoggerEntity *entity.RequestLoggerEntity, responseBody []byte, status int)
}
