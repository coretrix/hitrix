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

	newAdminUserEntity := &entity.AdminUserEntity{
		Email: "test@coretrix.com",
	}
	ormService.Flush(newAdminUserEntity)

	adminUserEntity := &entity.AdminUserEntity{}
	found := ormService.RedisSearchOne(adminUserEntity, query)

	assert.True(t, found)
}
