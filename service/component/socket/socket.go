package socket

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
)

var (
	ErrSocketClosed   = errors.New("socket is closed")
	ErrSendBufferFull = errors.New("socket send buffer is full")
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

	config         Config
	done           chan struct{}
	closeOnce      sync.Once
	unregisterOnce sync.Once
}

func newSocket(
	ctx context.Context,
	errorLogger errorlogger.ErrorLogger,
	connection *Connection,
	id string,
	namespace string,
	config Config,
) *Socket {
	return &Socket{
		Ctx:         ctx,
		ErrorLogger: errorLogger,
		Connection:  connection,
		ID:          id,
		Namespace:   namespace,
		config:      config,
		done:        make(chan struct{}),
	}
}

func (c *Connection) write(messageType int, payload []byte, writeWait time.Duration) error {
	if err := c.Ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return err
	}

	return c.Ws.WriteMessage(messageType, payload)
}

func (s *Socket) readPump(readMessageHandler ReadMessageHandler) {
	s.Connection.Ws.SetReadLimit(s.config.MaxMessageSize)
	_ = s.Connection.Ws.SetReadDeadline(time.Now().Add(s.config.PongWait))
	s.Connection.Ws.SetPongHandler(func(string) error {
		return s.Connection.Ws.SetReadDeadline(time.Now().Add(s.config.PongWait))
	})

	for {
		_, rawData, err := s.Connection.Ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				s.logError(err)
			}

			return
		}

		if readMessageHandler != nil {
			readMessageHandler(s, rawData)
		}
	}
}

func (s *Socket) writePump() {
	ticker := time.NewTicker(s.config.PingPeriod)
	defer ticker.Stop()

	for {
		select {
		case message := <-s.Connection.Send:
			if err := s.Connection.write(websocket.TextMessage, message, s.config.WriteWait); err != nil {
				s.logError(err)

				return
			}
		case <-ticker.C:
			if err := s.Connection.write(websocket.PingMessage, nil, s.config.WriteWait); err != nil {
				s.logError(err)

				return
			}
		case <-s.Ctx.Done():
			return
		case <-s.done:
			return
		}
	}
}

func (s *Socket) Emit(dto interface{}) error {
	data, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	select {
	case <-s.Ctx.Done():
		return ErrSocketClosed
	case <-s.done:
		return ErrSocketClosed
	case s.Connection.Send <- data:
		return nil
	default:
		return ErrSendBufferFull
	}
}

func (s *Socket) Done() <-chan struct{} {
	return s.done
}

func (s *Socket) Close() {
	s.closeOnce.Do(func() {
		close(s.done)
		_ = s.Connection.Ws.Close()
	})
}

func (s *Socket) logError(err error) {
	if s.ErrorLogger != nil {
		s.ErrorLogger.LogError(err)
	}
}
