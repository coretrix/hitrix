package requestlogger

import (
	"bytes"
	"io"
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

	var body = []byte("")

	if request.Body != nil {
		body, err = io.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}
	}

	content := append(headers, body...)

	// And now set a new body, which will simulate the same data we read:
	request.Body = io.NopCloser(bytes.NewBuffer(body))

	requestLoggerEntity := &entity.RequestLoggerEntity{
		URL:       url,
		AppName:   appName,
		CreatedAt: g.clockService.Now(),
	}

	if isText(contentType) && len(content)*4 <= 64000 {
		requestLoggerEntity.Request = content
	} else if len(body) > 0 {
		if len(body)*4 <= 64000 {
			requestLoggerEntity.Request = []byte("Skipped HEADERS \n\n")
			requestLoggerEntity.Request = append(requestLoggerEntity.Request, body...)
		} else {
			requestLoggerEntity.Request = []byte("Skipped HEADERS \n\n Partial BODY \n\n")
			requestLoggerEntity.Request = append(requestLoggerEntity.Request, body[0:16000]...)
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
		requestLoggerEntity.Response = []byte("Partial BODY \n\n")
		requestLoggerEntity.Response = append(requestLoggerEntity.Response, responseBody[0:16000]...)
	}

	if len(string(requestLoggerEntity.Log))*4 > 64000 {
		log := requestLoggerEntity.Log
		requestLoggerEntity.Log = []byte("Partial LOG \n\n")
		requestLoggerEntity.Log = append(requestLoggerEntity.Log, log[0:16000]...)
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
