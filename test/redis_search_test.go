package main

import (
	"testing"

	redisearch "github.com/coretrix/beeorm-redisearch-plugin"
	"github.com/stretchr/testify/assert"

	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/service"
)

func TestRedisSearch(t *testing.T) {
	createContextMyApp(t, "server", nil, nil, nil)

	ormService := service.DI().OrmEngine()

	query := &redisearch.RedisSearchQuery{}
	query.FilterString("Email", "test@coretrix.com")

	newDevPanelUserEntity := &entity.DevPanelUserEntity{
		Email: "test@coretrix.com",
	}
	ormService.Flush(newDevPanelUserEntity)

	devPanelUserEntity := &entity.DevPanelUserEntity{}
	found := ormService.RedisSearchOne(devPanelUserEntity, query)

	assert.True(t, found)
}
