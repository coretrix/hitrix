package main

import (
	"testing"

	"github.com/latolukasz/beeorm"
	"github.com/stretchr/testify/assert"

	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/service"
)

func TestRedisSearch(t *testing.T) {
	createContextMyApp(t, "my-app", nil, nil, nil)

	ormService, _ := service.DI().OrmEngine()

	query := &beeorm.RedisSearchQuery{}
	query.FilterString("Email", "test@coretrix.com")

	newDevPanelUserEntity := &entity.DevPanelUserEntity{
		Email: "test@coretrix.com",
	}
	ormService.Flush(newDevPanelUserEntity)

	devPanelUserEntity := &entity.DevPanelUserEntity{}
	found := ormService.RedisSearchOne(devPanelUserEntity, query)

	assert.True(t, found)
}
