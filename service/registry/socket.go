package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/goroutine"
	"github.com/coretrix/hitrix/service/component/socket"
)

func ServiceProviderSocketRegistry(eventHandlersMap socket.NamespaceEventHandlerMap) *service.DefinitionGlobal {
	return ServiceProviderSocketRegistryWithConfig(eventHandlersMap, socket.DefaultConfig())
}

func ServiceProviderSocketRegistryWithConfig(
	eventHandlersMap socket.NamespaceEventHandlerMap,
	config socket.Config,
) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.SocketRegistryService,
		Build: func(ctn di.Container) (interface{}, error) {
			return socket.NewSocketRegistryWithConfig(
					eventHandlersMap,
					ctn.Get(service.GoroutineService).(goroutine.IGoroutine),
					config,
				),
				nil
		},
	}
}
