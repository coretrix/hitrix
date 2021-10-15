# WebSocket
This service add support of websockets. It manage the connections and provide you easy way to read and write messages

Register the service into your `main.go` file:
```go
registry.ServiceProviderSocketRegistry(registerHandler, unregisterHandler func(s *socket.Socket))
```

Access the service:
```go
service.DI().SocketRegistry()
```

To be able to handle new connections you should create your own route and create a handler for it.
Your handler should looks like that:
```go
type WebsocketController struct {
}

func (controller *WebsocketController) InitConnection(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	socketRegistryService, has := service.DI().SocketRegistry()
	if !has {
		panic("Socket Registry is not registered")
	}

	errorLoggerService, has := service.DI().ErrorLogger()
	if !has {
		panic("Socket Registry is not registered")
	}

	connection := &socket.Connection{Send: make(chan []byte, 256), Ws: ws}
	socketHolder := &socket.Socket{
		ErrorLogger: errorLoggerService,
		Connection:  connection,
		ID:          "unique connection hash based on userID, deviceID and timestamp",
		Namespace:   model.DefaultNamespace,
	}

	socketRegistryService.Register <- socketHolder

	go socketHolder.WritePump()
	go socketHolder.ReadPump(socketRegistryService, func(rawData []byte) {
		s, _ := socketRegistryService.Sockets.Load(socketHolder.ID)
		
        dto := &DTOMessage{}
        err = json.Unmarshal(rawData, dto)
        if err != nil {
            errorLoggerService.LogError(err)
            retrun
        }
        //handle business logic here
        s.(*socket.Socket).Emit(dto)
	})
}

```
This handler initializes the new coming connections and have 2 go routines - one for writing messages and the second one for reading messages
If you want to send message you should use ```socketRegistryService.Emit```

If you want to read coming messages you should do it in the function we are passing as second parameter of ```ReadPump``` method

If you want to select certain connection you can do it by the ID and this method 
```go 
s, err := socketRegistryService.Sockets.Load(ID)
```

Also websocket service provide you hooks for registering new connections and for unregistering already existing connections.
You can define those handlers when you register the service based on namespace of socket.
