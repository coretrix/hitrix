package requestlogger

import (
	"net/http"

	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/pkg/entity"
)

const ID = "request_logger_id"

type IRequestLogger interface {
	LogRequest(ormService *datalayer.ORM, appName, url string, request *http.Request, contentType string) *entity.RequestLoggerEntity
	LogResponse(ormService *datalayer.ORM, requestLoggerEntity *entity.RequestLoggerEntity, responseBody []byte, status int)
}
