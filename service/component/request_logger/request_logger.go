package requestlogger

import (
	"net/http"

	"github.com/coretrix/trixorm"

	"github.com/coretrix/hitrix/pkg/entity"
)

const ID = "request_logger_id"

type IRequestLogger interface {
	LogRequest(ormService *trixorm.Engine, appName, url string, request *http.Request, contentType string) *entity.RequestLoggerEntity
	LogResponse(ormService *trixorm.Engine, requestLoggerEntity *entity.RequestLoggerEntity, responseBody []byte, status int)
}
