package middleware

import (
	"encoding/json"
	"net/http/httputil"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/pkg/response"
	"github.com/coretrix/hitrix/service"
	"github.com/gin-gonic/gin"
)

type dbLogger struct {
	disabled bool
	logs     []map[string]interface{}
}

func (l *dbLogger) Disable() {
	l.disabled = true
}

func (l *dbLogger) Handle(data map[string]interface{}) {
	if l.disabled {
		return
	}

	l.logs = append(l.logs, data)
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

func RequestLogger(ginEngine *gin.Engine, extender func(context *gin.Context, requestEntity *entity.RequestLoggerEntity)) {
	ormConfig := service.DI().OrmConfig()
	entities := ormConfig.GetEntities()
	if _, ok := entities["entity.RequestLoggerEntity"]; !ok {
		panic("you should register RequestLoggerEntity")
	}

	ginEngine.Use(func(context *gin.Context) {
		ormService := service.DI().OrmEngineForContext(context.Request.Context())

		requestEntity := &entity.RequestLoggerEntity{
			AppName:   service.DI().App().Name,
			CreatedAt: service.DI().Clock().Now(),
		}

		content, err := httputil.DumpRequest(context.Request, true)

		if err != nil {
			return
		}

		if isText(context.ContentType()) && len(content) <= 16_000_000 {
			requestEntity.Text = string(content)
		} else {
			requestEntity.Content = content
		}

		ormService.Flush(requestEntity)
		logger := &dbLogger{}
		ormService.RegisterQueryLogger(logger, true, true, true)
		context.Next()

		logger.Disable()
		encoded, err := json.Marshal(logger.logs)
		if err != nil {
			return
		}

		extender(context, requestEntity)

		responseBody, has := context.Get(response.ResponseBody)

		if has {
			responseBodyByte, _ := json.Marshal(responseBody)
			if len(responseBodyByte) <= 16_000_000 {
				requestEntity.ResponseText = string(responseBodyByte)
			} else {
				requestEntity.ResponseContent = responseBodyByte
			}
		}

		requestEntity.Log = encoded
		requestEntity.URL = context.Request.URL.String()
		requestEntity.Status = context.Writer.Status()
		ormService.Flush(requestEntity)
	})
}
