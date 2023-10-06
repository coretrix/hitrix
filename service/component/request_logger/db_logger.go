package requestlogger

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/clock"
)

type DBLogger struct {
	clockService clock.IClock
}

func NewDBLogger(clockService clock.IClock) IRequestLogger {
	return &DBLogger{clockService}
}

func (g *DBLogger) LogRequest(ormService *beeorm.Engine, appName, url string, request *http.Request, contentType string) *entity.RequestLoggerEntity {
	headers, err := httputil.DumpRequest(request, false)

	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	content := append(headers, body...)

	// And now set a new body, which will simulate the same data we read:
	request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	requestLoggerEntity := &entity.RequestLoggerEntity{
		URL:       url,
		AppName:   appName,
		CreatedAt: g.clockService.Now(),
	}

	if isText(contentType) && len(content)*4 <= 64000 {
		requestLoggerEntity.Request = content
	} else {
		if len(body)*4 <= 64000 {
			requestLoggerEntity.Request = []byte("Skipped HEADERS \n\n")
			requestLoggerEntity.Request = append(requestLoggerEntity.Request, body...)
		} else {
			requestLoggerEntity.Request = []byte("__TOO_LARGE__")
		}
	}

	ormService.Flush(requestLoggerEntity)

	return requestLoggerEntity
}

func (g *DBLogger) LogResponse(ormService *beeorm.Engine, requestLoggerEntity *entity.RequestLoggerEntity, responseBody []byte, status int) {
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
