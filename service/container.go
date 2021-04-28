package service

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/service/component/app"

	"github.com/sarulabs/di"
)

type key int

const (
	GinKey key = iota
)

var container di.Container

type Definition struct {
	Name   string
	Global bool
	Script bool
	Build  func(ctn di.Container) (interface{}, error)
	Close  func(obj interface{}) error
	Flags  func(registry *app.FlagsRegistry)
}

func SetContainer(c di.Container) {
	container = c
}

func HasService(key string) bool {
	_, has := container.Definitions()[key]
	return has
}

func GetServiceSafe(key string) (service interface{}, has bool, err error) {
	return getServiceSafe(container, key)
}

func GetServiceOptional(key string) (service interface{}, has bool) {
	return getServiceOptional(container, key)
}

func GetServiceRequired(key string) interface{} {
	return getServiceRequired(container, key)
}

func GetServiceForRequestSafe(ctx context.Context, key string) (service interface{}, has bool, err error) {
	return getServiceSafe(GetContainerFromRequest(ctx), key)
}

func GetServiceForRequestOptional(ctx context.Context, key string) (service interface{}, has bool) {
	return getServiceOptional(GetContainerFromRequest(ctx), key)
}

func GetServiceForRequestRequired(ctx context.Context, key string) interface{} {
	return getServiceRequired(GetContainerFromRequest(ctx), key)
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

func GetContainerFromRequest(ctx context.Context) (ctn di.Container) {
	c := GinFromContext(ctx)
	requestContainer, has := c.Get("RequestContainer")

	if !has {
		var err error
		ctn, err = container.SubContainer()
		if err != nil {
			panic(err)
		}
		c.Set("RequestContainer", ctn)
	} else {
		ctn = requestContainer.(di.Container)
	}
	return ctn
}

func GinFromContext(ctx context.Context) *gin.Context {
	return ctx.Value(GinKey).(*gin.Context)
}
