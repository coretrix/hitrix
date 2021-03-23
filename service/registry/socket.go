package registry

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/socket"
	"github.com/sarulabs/di"
)

func ServiceSocketRegistry(registerHandler, unregisterHandler func(s *socket.Socket)) *service.Definition {
	return &service.Definition{
		Name:   service.SocketRegistryService,
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return socket.BuildAndRunSocketRegistry(registerHandler, unregisterHandler), nil
		},
	}
}
