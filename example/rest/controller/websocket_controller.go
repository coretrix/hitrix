package controller

import (
	"encoding/json"
	"net/http"

	"github.com/coretrix/hitrix"

	model "github.com/coretrix/hitrix/example/model/socket"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/socket"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type DTOMessage struct {
	Type     string
	SocketID string
	Data     interface{}
}

type WebsocketController struct {
}

//if u want to test with websocket-demo.html u should add it to CORS policy

func (controller *WebsocketController) InitConnection(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	socketRegistryService := service.DI().SocketRegistry()

	errorLogger := service.DI().ErrorLogger()

	connection := &socket.Connection{Send: make(chan []byte, 256), Ws: ws}
	socketHolder := &socket.Socket{
		ErrorLogger: errorLogger,
		Connection:  connection,
		ID:          "unique connection hash based on userID, deviceID and timestamp",
		Namespace:   model.DefaultWebsocketNamespace,
	}

	socketRegistryService.Register <- socketHolder

	hitrix.Goroutine(func() {
		socketHolder.WritePump()
	})

	hitrix.Goroutine(func() {
		socketHolder.ReadPump(socketRegistryService, func(rawData []byte) {
			dto := &DTOMessage{}
			err = json.Unmarshal(rawData, dto)
			if err != nil {
				errorLogger.LogError(err)
				return
			}
			//return back the received message
			s, _ := socketRegistryService.Sockets.Load(socketHolder.ID)
			s.(*socket.Socket).Emit(dto)
		})
	})
}
