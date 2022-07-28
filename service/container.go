package service

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service/component/app"
)

type key int

const (
	GinKey key = iota
	RequestBodyKey
)

var container di.Container
var servicesDefinitionsRequestList []*DefinitionRequest

type DefinitionGlobal struct {
	Name   string
	Script bool
	Build  func(ctn di.Container) (interface{}, error)
	Close  func(obj interface{}) error
	Flags  func(registry *app.FlagsRegistry)
}

type DefinitionRequest struct {
	Name  string
	Build func(ctn *gin.Context) (interface{}, error)
}

func SetContainer(c di.Container) {
	if container != nil {
		_ = container.Delete()
	}

	container = c
}

func SetRequestServices(servicesDefinitionsRequest []*DefinitionRequest) {
	servicesDefinitionsRequestList = servicesDefinitionsRequest
}

func HasService(key string) bool {
	_, has := container.Definitions()[key]

	return has
}

func GetServiceOptional(key string) (service interface{}, has bool) {
	return getServiceOptional(container, key)
}

func GetServiceRequired(key string) interface{} {
	return getServiceRequired(container, key)
}

func GetServiceForRequestRequired(ctx context.Context, key string) interface{} {
	return GetServiceFromRequest(ctx, key)
}

func getServiceSafe(ctn di.Container, key string) (service interface{}, has bool, err error) {
	service, err = ctn.SafeGet(key)
	if err == nil {
		return service, true, nil
	}

	_, has = ctn.Definitions()[key]
	if !has {
		return nil, false, nil
	}

	return nil, true, err
}

func getServiceOptional(ctn di.Container, key string) (service interface{}, has bool) {
	service, has, err := getServiceSafe(ctn, key)
	if err != nil {
		panic(err)
	}

	return service, has
}

func getServiceRequired(ctn di.Container, key string) interface{} {
	service, has, err := getServiceSafe(ctn, key)
	if err != nil {
		panic(err)
	} else if !has {
		panic(fmt.Errorf("missing service " + key))
	}

	return service
}

func GetServiceFromRequest(ctx context.Context, key string) interface{} {
	c := GinFromContext(ctx)
	requestService, has := c.Get(key)

	if !has {
		var err error

		for _, s := range servicesDefinitionsRequestList {
			if s.Name == key {
				requestService, err = s.Build(c)
				if err != nil {
					panic(err)
				}
			}
		}

		c.Set(key, requestService)
	}

	if requestService == nil {
		panic("not defined service " + key)
	}

	return requestService
}

func GinFromContext(ctx context.Context) *gin.Context {
	return ctx.Value(GinKey).(*gin.Context)
}
