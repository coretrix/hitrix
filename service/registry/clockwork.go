package registry

import (
	"github.com/coretrix/clockwork"
	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/service"
)

func ServiceProviderClockWorkForContext() *service.DefinitionRequest {
	return &service.DefinitionRequest{
		Name: service.ClockWorkRequestService,
		Build: func(c *gin.Context) (interface{}, error) {
			return &clockwork.Clockwork{}, nil
		},
	}
}
