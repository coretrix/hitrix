package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/socket"
	"github.com/sarulabs/di"
)

func ServiceSocketRegistry(eventHandlersMap socket.NamespaceEventHandlerMap) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.SocketRegistryService,
		Build: func(ctn di.Container) (interface{}, error) {
			return socket.BuildAndRunSocketRegistry(eventHandlersMap), nil
		},
	}
}
