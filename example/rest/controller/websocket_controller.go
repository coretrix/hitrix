package controller

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	model "github.com/coretrix/hitrix/example/model/socket"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/socket"
)

type DTOMessage struct {
	Type     string
	SocketID string
	Data     interface{}
}

type WebsocketController struct {
}

func (controller *WebsocketController) InitConnection(c *gin.Context) {
	socketRegistryService := service.DI().SocketRegistry()
	errorLogger := service.DI().ErrorLogger()

	err := socketRegistryService.ServeHTTP(
		c.Writer,
		c.Request,
		socket.ConnectionOptions{
			ID:          "unique connection hash based on userID, deviceID and timestamp",
			Namespace:   model.DefaultWebsocketNamespace,
			Context:     c.Request.Context(),
			ErrorLogger: errorLogger,
			ReadMessageHandler: func(socketHolder *socket.Socket, rawData []byte) {
				dto := &DTOMessage{}
				if unmarshalErr := json.Unmarshal(rawData, dto); unmarshalErr != nil {
					errorLogger.LogError(unmarshalErr)

					return
				}

				//return the received message
				if emitErr := socketHolder.Emit(dto); emitErr != nil {
					errorLogger.LogError(emitErr)
				}
			},
		})
	if err != nil {
		errorLogger.LogErrorWithRequest(c, err)
	}
}
