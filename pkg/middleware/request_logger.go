package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/pkg/response"
	"github.com/coretrix/hitrix/service"
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
		requestLoggerService := service.DI().RequestLogger()

		requestLoggerEntity := requestLoggerService.LogRequest(
			ormService,
			service.DI().App().Name,
			context.Request.URL.String(),
			context.Request,
			context.ContentType(),
		)

		logger := &dbLogger{}
		ormService.RegisterQueryLogger(logger, true, true, true)
		context.Next()

		logger.Disable()
		encoded, err := json.Marshal(logger.logs)
		if err != nil {
			return
		}

		extender(context, requestLoggerEntity)

		responseBody, _ := context.Get(response.ResponseBody)
		responseBodyByte, _ := json.Marshal(responseBody)

		requestLoggerEntity.Log = encoded
		requestLoggerService.LogResponse(ormService, requestLoggerEntity, responseBodyByte, context.Writer.Status())
	})
}
