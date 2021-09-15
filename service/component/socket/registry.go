package socket

import (
	"fmt"
	"sync"
	"time"
)

const (
	// Time allowed to write a broadcastMessage to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong broadcastMessage from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum broadcastMessage size allowed from peer.
	maxMessageSize = 512
)

type Registry struct {
	Register   chan *Socket
	Unregister chan *Socket
	Sockets    *sync.Map
}

func BuildAndRunSocketRegistry(eventHandlersMap NamespaceEventHandlerMap) *Registry {
	registry := &Registry{
		Register:   make(chan *Socket),
		Unregister: make(chan *Socket),
		Sockets:    &sync.Map{},
	}

	//TODO use hitrix
	//hitrix.GoroutineWithRestart(func() {
	go registry.run(eventHandlersMap)
	//})

	return registry
}

func (registry *Registry) run(eventHandlersMap NamespaceEventHandlerMap) {
	for {
		select {
		case s := <-registry.Register: //new connection
			registry.Sockets.Store(s.ID, s)

			eventHandlers, ok := eventHandlersMap[s.Namespace]
			if !ok {
				panic(fmt.Errorf("register handler for namespace %v not found", s.Namespace))
			}

			eventHandlers.RegisterHandler(s)
		case s := <-registry.Unregister:
			registry.Sockets.Delete(s.ID)

			eventHandlers, ok := eventHandlersMap[s.Namespace]
			if !ok {
				panic(fmt.Errorf("unregister handler for namespace %v not found", s.Namespace))
			}

			eventHandlers.UnregisterHandler(s)
		}
	}
}

type NamespaceEventHandlerMap map[string]*EventHandlers

type EventHandlers struct {
	RegisterHandler, UnregisterHandler func(s *Socket)
}
