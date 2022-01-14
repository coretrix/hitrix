package socket

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"

	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
)

type Connection struct {
	Ws   *websocket.Conn
	Send chan []byte
}

type Socket struct {
	Ctx         context.Context
	ErrorLogger errorlogger.ErrorLogger
	Connection  *Connection
	ID          string
	Namespace   string
}

func (c *Connection) write(mt int, payload []byte) error {
	_ = c.Ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Ws.WriteMessage(mt, payload)
}

func (s *Socket) ReadPump(registry *Registry, readMessageHandler func(rawData []byte)) {
	defer func() {
		registry.Unregister <- s
		s.Connection.Ws.Close()
	}()

	s.Connection.Ws.SetReadLimit(maxMessageSize)
	_ = s.Connection.Ws.SetReadDeadline(time.Now().Add(pongWait))
	s.Connection.Ws.SetPongHandler(func(string) error { _ = s.Connection.Ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, rawData, err := s.Connection.Ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				s.ErrorLogger.LogError(err)
			}
			break
		}

		readMessageHandler(rawData)
	}
}

func (s *Socket) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		s.Connection.Ws.Close()
	}()
	for {
		select {
		case message, ok := <-s.Connection.Send:
			if !ok {
				err := s.Connection.write(websocket.CloseMessage, []byte{})
				if err != nil {
					s.ErrorLogger.LogError(err)
				}
				return
			}
			if err := s.Connection.write(websocket.TextMessage, message); err != nil {
				s.ErrorLogger.LogError(err)
				return
			}
		case <-ticker.C:
			if err := s.Connection.write(websocket.PingMessage, []byte{}); err != nil {
				s.ErrorLogger.LogError(err)
				return
			}
		}
	}
}

func (s *Socket) Emit(dto interface{}) {
	data, err := json.Marshal(dto)
	if err != nil {
		s.ErrorLogger.LogError(err)
		return
	}

	s.Connection.Send <- data
}
