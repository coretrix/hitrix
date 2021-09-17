package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/goroutine"
	"github.com/coretrix/hitrix/service/component/socket"
	"github.com/sarulabs/di"
)

func ServiceProviderSocketRegistry(eventHandlersMap socket.NamespaceEventHandlerMap) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.SocketRegistryService,
		Build: func(ctn di.Container) (interface{}, error) {
			return socket.NewSocketRegistry(
					eventHandlersMap,
					ctn.Get(service.GoroutineService).(goroutine.IGoroutine)),
				nil
		},
	}
}
