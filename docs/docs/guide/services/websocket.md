# WebSocket

The socket registry owns the complete WebSocket lifecycle: HTTP upgrade,
registration, ping/pong keepalive, reads, buffered writes and cleanup.

Register the service into your `main.go` file:
```go
eventHandlersMap := socket.NamespaceEventHandlerMap{
	"default": {
		RegisterHandler:   registerHandler,
		UnregisterHandler: unregisterHandler,
	},
}

registry.ServiceProviderSocketRegistry(eventHandlersMap)
```

The default origin policy accepts requests without an `Origin` header and
same-origin browser requests. Cross-origin applications must provide an
explicit policy:

```go
socketConfig := socket.DefaultConfig()
socketConfig.CheckOrigin = func(request *http.Request) bool {
	return request.Header.Get("Origin") == "https://business.example.com"
}

registry.ServiceProviderSocketRegistryWithConfig(eventHandlersMap, socketConfig)
```

Access the service:
```go
service.DI().SocketRegistry()
```

Authenticate and authorize the request before calling `ServeHTTP`. For example,
an application using short-lived, single-use socket tickets must consume and
validate the ticket first. Invalid requests should never be upgraded.

`ServeHTTP` blocks until the connection closes and manages its goroutines:

```go
type WebsocketController struct {
}

func (controller *WebsocketController) InitConnection(c *gin.Context) {
	// Validate and consume the application's socket ticket here.
	socketRegistryService := service.DI().SocketRegistry()
	errorLoggerService := service.DI().ErrorLogger()

	err := socketRegistryService.ServeHTTP(c.Writer, c.Request, socket.ConnectionOptions{
		ID:          "unique authenticated connection ID",
		Namespace:   model.DefaultNamespace,
		Context:     c.Request.Context(),
		ErrorLogger: errorLoggerService,
		ReadMessageHandler: func(socketHolder *socket.Socket, rawData []byte) {
			// Handle application messages here.
		},
	})
	if err != nil {
		errorLoggerService.LogErrorWithRequest(c, err)
	}
}
```

Send a message to one connection by ID:

```go
err := socketRegistryService.Emit(socketID, dto)
```

`Emit` is non-blocking. It returns `socket.ErrSendBufferFull` when a slow
consumer fills its outbound buffer and `socket.ErrSocketClosed` after the
connection is gone. Register and unregister hooks are configured per namespace.
