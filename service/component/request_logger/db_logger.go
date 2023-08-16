package requestlogger

import (
	"net/http"
	"net/http/httputil"

	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/clock"
)

type DBLogger struct {
	clockService clock.IClock
}

func NewDBLogger(clockService clock.IClock) IRequestLogger {
	return &DBLogger{clockService}
}

//nolint // info
func (g *DBLogger) LogRequest(ormService *datalayer.DataLayer, appName, url string, request *http.Request, contentType string) *entity.RequestLoggerEntity {
	content, err := httputil.DumpRequest(request, true)

	if err != nil {
		panic(err)
	}

	requestLoggerEntity := &entity.RequestLoggerEntity{
		URL:       url,
		AppName:   appName,
		CreatedAt: g.clockService.Now(),
	}

	if isText(contentType) && len(content)*4 <= 64000 {
		requestLoggerEntity.Request = content
	} else {
		requestLoggerEntity.Request = []byte("__TOO_LARGE__")
	}

	ormService.Flush(requestLoggerEntity)

	return requestLoggerEntity
}

// nolint //info
func (g *DBLogger) LogResponse(ormService *datalayer.DataLayer, requestLoggerEntity *entity.RequestLoggerEntity, responseBody []byte, status int) {
	requestLoggerEntity.Status = status

	if len(responseBody) > 0 && len(string(responseBody))*4 <= 64000 {
		requestLoggerEntity.Response = responseBody
	} else {
		requestLoggerEntity.Response = []byte("__TOO_LARGE__")
	}

	if len(string(requestLoggerEntity.Log))*4 > 64000 {
		requestLoggerEntity.Log = []byte("__TOO_LARGE__")
	}

	ormService.Flush(requestLoggerEntity)
}

func isText(contentType string) bool {
	allowedHeaders := map[string]struct{}{
		"application/json": {},
		"text/plain":       {},
		"text/html":        {},
	}

	_, ok := allowedHeaders[contentType]

	return ok
}
