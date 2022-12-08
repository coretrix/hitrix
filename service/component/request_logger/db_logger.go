package requestlogger

import (
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/latolukasz/beeorm"
	"net/http"
	"net/http/httputil"
)

type DbLogger struct {
	clockService clock.IClock
}

func NewDbLogger(clockService clock.IClock) IRequestLogger {
	return &DbLogger{clockService}
}

func (g *DbLogger) LogRequest(ormService *beeorm.Engine, appName, url string, request *http.Request, contentType string) *entity.RequestLoggerEntity {
	content, err := httputil.DumpRequest(request, true)

	if err != nil {
		panic(err)
	}

	requestLoggerEntity := &entity.RequestLoggerEntity{
		URL:       url,
		AppName:   appName,
		Request:   content,
		CreatedAt: g.clockService.Now(),
	}

	if isText(contentType) && len(content) <= 16_000_000 {
		requestLoggerEntity.RequestText = string(content)
	} else {
		requestLoggerEntity.Request = content
	}

	ormService.Flush(requestLoggerEntity)

	return requestLoggerEntity
}

func (g *DbLogger) LogResponse(ormService *beeorm.Engine, requestLoggerEntity *entity.RequestLoggerEntity, responseBody []byte, status int) {
	requestLoggerEntity.Status = status

	if len(string(responseBody)) <= 16_000_000 {
		requestLoggerEntity.ResponseText = string(responseBody)
	} else {
		requestLoggerEntity.Response = responseBody
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
