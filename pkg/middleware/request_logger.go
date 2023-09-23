package middleware

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/pkg/response"
	"github.com/coretrix/hitrix/service"
	requestlogger "github.com/coretrix/hitrix/service/component/request_logger"
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

func RequestLogger(ginEngine *gin.Engine, extender func(context *gin.Context, requestEntity *entity.RequestLoggerEntity)) {
	ormConfig := service.DI().OrmConfig()

	entities := ormConfig.GetEntities()
	if _, ok := entities["entity.RequestLoggerEntity"]; !ok {
		panic("you should register RequestLoggerEntity")
	}

	ginEngine.Use(func(context *gin.Context) {
		requestStart := time.Now()

		ormService := service.DI().OrmEngineForContext(context.Request.Context())
		requestLoggerService := service.DI().RequestLogger()

		requestLoggerEntity := requestLoggerService.LogRequest(
			ormService,
			service.DI().App().Name,
			context.Request.URL.String(),
			context.Request,
			context.ContentType(),
		)

		context.Set(requestlogger.ID, requestLoggerEntity.ID)
		context.Header("X-Request-ID", strconv.FormatUint(requestLoggerEntity.ID, 10))

		logger := &dbLogger{}
		ormService.RegisterQueryLogger(logger, true, true, true)

		context.Next()

		requestLoggerEntity.RequestDuration = time.Now().Sub(requestStart).Milliseconds()

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
