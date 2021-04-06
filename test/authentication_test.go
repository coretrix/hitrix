package main

import (
	"github.com/coretrix/hitrix/service/registry"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/service"
)

func TestGenerateToken(t *testing.T) {
	createContextMyApp(t, "my-app", nil, []*service.Definition{
		registry.ServiceProviderJWT(),
		registry.ServiceProviderPassword(),
		registry.ServiceProviderAuthentication(&entity.AdminUserEntity{}),
	})
	ormService, _ := service.DI().OrmEngine()
	passwordService, _ := service.DI().Password()
	hashedPassword, _ := passwordService.HashPassword("123")

	adminEntity := &entity.AdminUserEntity{
		Email:    "test@test.com",
		Password: hashedPassword,
	}
	ormService.Flush(adminEntity)

	authenticationService, _ := service.DI().AuthenticationService()
	accessToken, refreshToken, err := authenticationService.Authenticate("test@test.com", "1234")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid_password")
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)

	accessToken, refreshToken, err = authenticationService.Authenticate("test1@test.com", "123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user_not_found")
	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)

	accessToken, refreshToken, err = authenticationService.Authenticate("test@test.com", "123")
	assert.Nil(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	fetchedAdminEntity := &entity.AdminUserEntity{}
	err = authenticationService.VerifyAccessToken("something-invalid", fetchedAdminEntity)
	assert.NotNil(t, err)

	err = authenticationService.VerifyAccessToken(accessToken, fetchedAdminEntity)
	assert.Nil(t, err)
	assert.Equal(t, fetchedAdminEntity.ID, uint64(1))

	newAccessToken, newRefreshToken, err := authenticationService.RefreshToken("some-invalid-refresh-token")
	assert.NotNil(t, err)
	assert.Empty(t, newAccessToken)
	assert.Empty(t, newRefreshToken)

	newAccessToken, newRefreshToken, err = authenticationService.RefreshToken(refreshToken)
	assert.Nil(t, err)
	assert.NotEmpty(t, newAccessToken)
	assert.NotEmpty(t, newRefreshToken)

	fetchedAdminEntityUsingRefreshToken := &entity.AdminUserEntity{}
	err = authenticationService.VerifyAccessToken(newAccessToken, fetchedAdminEntityUsingRefreshToken)
	assert.Nil(t, err)
	assert.Equal(t, fetchedAdminEntityUsingRefreshToken.ID, uint64(1))
}
