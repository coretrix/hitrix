package socket

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	"github.com/coretrix/hitrix/service/component/goroutine"
	componentmetrics "github.com/coretrix/hitrix/service/component/metrics"
)

const (
	defaultWriteWait       = 10 * time.Second
	defaultPongWait        = 60 * time.Second
	defaultPingPeriod      = 30 * time.Second
	defaultMaxMessageSize  = 512
	defaultReadBufferSize  = 1024
	defaultWriteBufferSize = 1024
	defaultSendBufferSize  = 256

	MetricActiveConnections = "SocketActiveConnections"
	MetricConnectionsTotal  = "SocketConnectionsTotal"
	MetricDisconnectsTotal  = "SocketDisconnectsTotal"
)

var (
	ErrMissingSocketID  = errors.New("socket ID is required")
	ErrMissingNamespace = errors.New("socket namespace is required")
	ErrUnknownNamespace = errors.New("socket namespace is not registered")
)

type Config struct {
	WriteWait         time.Duration
	PongWait          time.Duration
	PingPeriod        time.Duration
	MaxMessageSize    int64
	ReadBufferSize    int
	WriteBufferSize   int
	SendBufferSize    int
	EnableCompression bool
	CheckOrigin       func(request *http.Request) bool
}

func DefaultConfig() Config {
	return Config{
		WriteWait:         defaultWriteWait,
		PongWait:          defaultPongWait,
		PingPeriod:        defaultPingPeriod,
		MaxMessageSize:    defaultMaxMessageSize,
		ReadBufferSize:    defaultReadBufferSize,
		WriteBufferSize:   defaultWriteBufferSize,
		SendBufferSize:    defaultSendBufferSize,
		EnableCompression: true,
		CheckOrigin:       checkSameOrigin,
	}
}

type ConnectionOptions struct {
	ID                 string
	Namespace          string
	Context            context.Context
	ErrorLogger        errorlogger.ErrorLogger
	ReadMessageHandler ReadMessageHandler
	ResponseHeader     http.Header
}

type ReadMessageHandler func(socket *Socket, rawData []byte)

type Registry struct {
	Sockets          *sync.Map
	ServiceGoroutine goroutine.IGoroutine

	eventHandlersMap NamespaceEventHandlerMap
	config           Config
	lifecycleMu      sync.Mutex
}

func NewSocketRegistry(eventHandlersMap NamespaceEventHandlerMap, serviceGoroutine goroutine.IGoroutine) *Registry {
	return NewSocketRegistryWithConfig(eventHandlersMap, serviceGoroutine, DefaultConfig())
}

func NewSocketRegistryWithConfig(
	eventHandlersMap NamespaceEventHandlerMap,
	serviceGoroutine goroutine.IGoroutine,
	config Config,
) *Registry {
	config.applyDefaults()

	return &Registry{
		Sockets:          &sync.Map{},
		ServiceGoroutine: serviceGoroutine,
		eventHandlersMap: eventHandlersMap,
		config:           config,
	}
}

// ServeHTTP upgrades and owns the connection until it is closed. Authentication
// must be completed by the caller before this method is invoked.
func (registry *Registry) ServeHTTP(
	writer http.ResponseWriter,
	request *http.Request,
	options ConnectionOptions,
) error {
	if options.ID == "" {
		return ErrMissingSocketID
	}
	if options.Namespace == "" {
		return ErrMissingNamespace
	}
	eventHandlers, ok := registry.eventHandlersMap[options.Namespace]
	if !ok || eventHandlers == nil {
		return fmt.Errorf("%w: %s", ErrUnknownNamespace, options.Namespace)
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:    registry.config.ReadBufferSize,
		WriteBufferSize:   registry.config.WriteBufferSize,
		EnableCompression: registry.config.EnableCompression,
		CheckOrigin:       registry.config.CheckOrigin,
	}

	websocketConnection, err := upgrader.Upgrade(writer, request, options.ResponseHeader)
	if err != nil {
		return err
	}

	ctx := options.Context
	if ctx == nil {
		ctx = context.Background()
	}

	socketHolder := newSocket(
		ctx,
		options.ErrorLogger,
		&Connection{
			Ws:   websocketConnection,
			Send: make(chan []byte, registry.config.SendBufferSize),
		},
		options.ID,
		options.Namespace,
		registry.config,
	)

	registry.register(socketHolder)
	defer registry.unregister(socketHolder)

	registry.ServiceGoroutine.Goroutine(func() {
		socketHolder.writePump()
		socketHolder.Close()
	})
	registry.ServiceGoroutine.Goroutine(func() {
		select {
		case <-ctx.Done():
			socketHolder.Close()
		case <-socketHolder.Done():
		}
	})

	socketHolder.readPump(options.ReadMessageHandler)
	socketHolder.Close()

	return nil
}

func (registry *Registry) Emit(socketID string, dto interface{}) error {
	value, ok := registry.Sockets.Load(socketID)
	if !ok {
		return ErrSocketClosed
	}

	return value.(*Socket).Emit(dto)
}

func (registry *Registry) Close(socketID string) bool {
	value, ok := registry.Sockets.Load(socketID)
	if !ok {
		return false
	}

	value.(*Socket).Close()

	return true
}

func (registry *Registry) register(socketHolder *Socket) {
	registry.lifecycleMu.Lock()
	defer registry.lifecycleMu.Unlock()

	if current, loaded := registry.Sockets.Load(socketHolder.ID); loaded {
		currentSocket := current.(*Socket)
		currentSocket.Close()
		registry.unregisterLocked(currentSocket)
	}
	registry.Sockets.Store(socketHolder.ID, socketHolder)
	componentmetrics.Add(MetricActiveConnections, 1)
	componentmetrics.Add(MetricConnectionsTotal, 1)

	if handler := registry.eventHandlersMap[socketHolder.Namespace].RegisterHandler; handler != nil {
		handler(socketHolder)
	}
}

func (registry *Registry) unregister(socketHolder *Socket) {
	registry.lifecycleMu.Lock()
	defer registry.lifecycleMu.Unlock()

	registry.unregisterLocked(socketHolder)
}

func (registry *Registry) unregisterLocked(socketHolder *Socket) {
	socketHolder.unregisterOnce.Do(func() {
		registry.Sockets.CompareAndDelete(socketHolder.ID, socketHolder)
		componentmetrics.Add(MetricActiveConnections, -1)
		componentmetrics.Add(MetricDisconnectsTotal, 1)

		if handler := registry.eventHandlersMap[socketHolder.Namespace].UnregisterHandler; handler != nil {
			handler(socketHolder)
		}
	})
}

func (config *Config) applyDefaults() {
	defaults := DefaultConfig()

	if config.WriteWait <= 0 {
		config.WriteWait = defaults.WriteWait
	}
	if config.PongWait <= 0 {
		config.PongWait = defaults.PongWait
	}
	if config.PingPeriod <= 0 || config.PingPeriod >= config.PongWait {
		config.PingPeriod = defaults.PingPeriod
	}
	if config.MaxMessageSize <= 0 {
		config.MaxMessageSize = defaults.MaxMessageSize
	}
	if config.ReadBufferSize <= 0 {
		config.ReadBufferSize = defaults.ReadBufferSize
	}
	if config.WriteBufferSize <= 0 {
		config.WriteBufferSize = defaults.WriteBufferSize
	}
	if config.SendBufferSize <= 0 {
		config.SendBufferSize = defaults.SendBufferSize
	}
	if config.CheckOrigin == nil {
		config.CheckOrigin = defaults.CheckOrigin
	}
}

func checkSameOrigin(request *http.Request) bool {
	origin := request.Header.Get("Origin")
	if origin == "" {
		return true
	}

	parsedOrigin, err := url.Parse(origin)
	if err != nil {
		return false
	}

	return strings.EqualFold(parsedOrigin.Host, request.Host)
}

type NamespaceEventHandlerMap map[string]*EventHandlers

type EventHandlers struct {
	RegisterHandler, UnregisterHandler func(s *Socket)
}
