package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/socket"
	"github.com/sarulabs/di"
)

func ServiceSocketRegistry(eventHandlersMap socket.NamespaceEventHandlerMap) *service.Definition {
	return &service.Definition{
		Name:   service.SocketRegistryService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return socket.BuildAndRunSocketRegistry(eventHandlersMap), nil
		},
	}
}
