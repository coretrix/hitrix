package main

import (
	"fmt"
	"testing"

	"github.com/coretrix/hitrix/pkg/helper"

	"github.com/coretrix/hitrix/service/registry"
	"github.com/stretchr/testify/assert"

	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/service"
)

func createUser(input map[string]interface{}) *entity.AdminUserEntity {
	ormService, _ := service.DI().OrmEngine()
	adminEntity := &entity.AdminUserEntity{}
	for field, val := range input {
		switch field {
		case "Email":
			adminEntity.Email = val.(string)
		case "Password":
			adminEntity.Password = val.(string)
		}
	}
	ormService.Flush(adminEntity)
	return adminEntity
}

func TestAuthenticate(t *testing.T) {
	t.Run("simple successful", func(t *testing.T) {
		createContextMyApp(t, "my-app", nil,
			registry.ServiceProviderJWT(),
			registry.ServiceProviderPassword(),
			registry.ServiceProviderAuthentication(),
		)

		passwordService, _ := service.DI().Password()
		hashedPassword, _ := passwordService.HashPassword("1234")

		createUser(map[string]interface{}{
			"Email":    "test@test.com",
			"Password": hashedPassword,
		})

		authenticationService, _ := service.DI().AuthenticationService()
		fetchedAdminEntity := &entity.AdminUserEntity{}
		_, _, err := authenticationService.Authenticate("test@test.com", "1234", fetchedAdminEntity)
		assert.Nil(t, err)
	})

	t.Run("wrong email", func(t *testing.T) {
		createContextMyApp(t, "my-app", nil,
			registry.ServiceProviderJWT(),
			registry.ServiceProviderPassword(),
			registry.ServiceProviderAuthentication(),
		)

		passwordService, _ := service.DI().Password()
		hashedPassword, _ := passwordService.HashPassword("1234")

		createUser(map[string]interface{}{
			"Email":    "test@test.com",
			"Password": hashedPassword,
		})

		authenticationService, _ := service.DI().AuthenticationService()
		fetchedAdminEntity := &entity.AdminUserEntity{}
		_, _, err := authenticationService.Authenticate("test@tesat.com", "1234", fetchedAdminEntity)
		assert.NotNil(t, err)
	})
}

func TestVerifyAccessToken(t *testing.T) {
	t.Run("simple success", func(t *testing.T) {
		createContextMyApp(t, "my-app", nil,
			registry.ServiceProviderJWT(),
			registry.ServiceProviderPassword(),
			registry.ServiceProviderAuthentication(),
		)

		passwordService, _ := service.DI().Password()
		hashedPassword, _ := passwordService.HashPassword("1234")

		currentUser := createUser(map[string]interface{}{
			"Email":    "test@test.com",
			"Password": hashedPassword,
		})

		accessKey := fmt.Sprintf("ACCESS:%d:%s", currentUser.ID, helper.GenerateUUID())
		ormService, _ := service.DI().OrmEngine()
		ormService.GetRedis().Set(accessKey, "", 10)

		authenticationService, _ := service.DI().AuthenticationService()
		token, err := authenticationService.GenerateTokenPair(currentUser.ID, accessKey, 10)
		assert.Nil(t, err)
		fetchedAdminEntity := &entity.AdminUserEntity{}
		payload, err := authenticationService.VerifyAccessToken(token, fetchedAdminEntity)
		assert.Nil(t, err)
		assert.Equal(t, accessKey, payload["jti"])
	})

	t.Run("wrong token", func(t *testing.T) {
		createContextMyApp(t, "my-app", nil,
			registry.ServiceProviderJWT(),
			registry.ServiceProviderPassword(),
			registry.ServiceProviderAuthentication(),
		)

		passwordService, _ := service.DI().Password()
		hashedPassword, _ := passwordService.HashPassword("1234")

		currentUser := createUser(map[string]interface{}{
			"Email":    "test@test.com",
			"Password": hashedPassword,
		})

		accessKey := fmt.Sprintf("ACCESS:%d:%s", currentUser.ID, helper.GenerateUUID())
		ormService, _ := service.DI().OrmEngine()
		ormService.GetRedis().Set(accessKey, "", 10)

		authenticationService, _ := service.DI().AuthenticationService()
		token, err := authenticationService.GenerateTokenPair(currentUser.ID, accessKey, 10)
		assert.Nil(t, err)
		fetchedAdminEntity := &entity.AdminUserEntity{}
		_, err = authenticationService.VerifyAccessToken(" wef"+token, fetchedAdminEntity)
		assert.NotNil(t, err)
	})
}

func TestRefreshToken(t *testing.T) {
	t.Run("success refresh", func(t *testing.T) {
		createContextMyApp(t, "my-app", nil,
			registry.ServiceProviderJWT(),
			registry.ServiceProviderPassword(),
			registry.ServiceProviderAuthentication(),
		)

		passwordService, _ := service.DI().Password()
		hashedPassword, _ := passwordService.HashPassword("1234")

		currentUser := createUser(map[string]interface{}{
			"Email":    "test@test.com",
			"Password": hashedPassword,
		})

		accessKey := fmt.Sprintf("ACCESS:%d:%s", currentUser.ID, helper.GenerateUUID())
		ormService, _ := service.DI().OrmEngine()
		ormService.GetRedis().Set(accessKey, "", 10)

		authenticationService, _ := service.DI().AuthenticationService()
		refresh, err := authenticationService.GenerateTokenPair(currentUser.ID, accessKey, 10)
		assert.Nil(t, err)
		_, _, err = authenticationService.RefreshToken(refresh)
		assert.Nil(t, err)
	})

	t.Run("wrong refresh", func(t *testing.T) {
		createContextMyApp(t, "my-app", nil,
			registry.ServiceProviderJWT(),
			registry.ServiceProviderPassword(),
			registry.ServiceProviderAuthentication(),
		)

		passwordService, _ := service.DI().Password()
		hashedPassword, _ := passwordService.HashPassword("1234")

		currentUser := createUser(map[string]interface{}{
			"Email":    "test@test.com",
			"Password": hashedPassword,
		})

		accessKey := fmt.Sprintf("ACCESS:%d:%s", currentUser.ID, helper.GenerateUUID())
		ormService, _ := service.DI().OrmEngine()
		ormService.GetRedis().Set(accessKey, "", 10)

		authenticationService, _ := service.DI().AuthenticationService()
		refresh, err := authenticationService.GenerateTokenPair(currentUser.ID, accessKey, 10)
		assert.Nil(t, err)
		_, _, err = authenticationService.RefreshToken("ef" + refresh)
		assert.NotNil(t, err)
	})
}

func TestLogoutCurrentSession(t *testing.T) {
	t.Run("simple logout", func(t *testing.T) {
		createContextMyApp(t, "my-app", nil,
			registry.ServiceProviderJWT(),
			registry.ServiceProviderPassword(),
			registry.ServiceProviderAuthentication(),
		)

		passwordService, _ := service.DI().Password()
		hashedPassword, _ := passwordService.HashPassword("1234")

		currentUser := createUser(map[string]interface{}{
			"Email":    "test@test.com",
			"Password": hashedPassword,
		})

		accessKey := fmt.Sprintf("ACCESS:%d:%s", currentUser.ID, helper.GenerateUUID())
		ormService, _ := service.DI().OrmEngine()
		ormService.GetRedis().Set(accessKey, "", 10)

		authenticationService, _ := service.DI().AuthenticationService()
		accessToken, err := authenticationService.GenerateTokenPair(currentUser.ID, accessKey, 10)
		assert.Nil(t, err)
		fetchedAdminEntity := &entity.AdminUserEntity{}
		_, err = authenticationService.VerifyAccessToken(accessToken, fetchedAdminEntity)
		assert.Nil(t, err)
		authenticationService.LogoutCurrentSession(accessKey)
		_, err = authenticationService.VerifyAccessToken(accessToken, fetchedAdminEntity)
		assert.NotNil(t, err)
	})
}

func TestLogoutAllSessions(t *testing.T) {
	t.Run("logout from one session", func(t *testing.T) {
		createContextMyApp(t, "my-app", nil,
			registry.ServiceProviderJWT(),
			registry.ServiceProviderPassword(),
			registry.ServiceProviderAuthentication(),
		)

		passwordService, _ := service.DI().Password()
		hashedPassword, _ := passwordService.HashPassword("1234")

		currentUser := createUser(map[string]interface{}{
			"Email":    "test@test.com",
			"Password": hashedPassword,
		})

		accessKey := fmt.Sprintf("ACCESS:%d:%s", currentUser.ID, helper.GenerateUUID())
		ormService, _ := service.DI().OrmEngine()
		ormService.GetRedis().Set(accessKey, "", 10)

		accessListKey := fmt.Sprintf("USER_KEYS:%d", currentUser.ID)
		ormService.GetRedis().Set(accessListKey, accessKey, 10)

		authenticationService, _ := service.DI().AuthenticationService()
		accessToken, err := authenticationService.GenerateTokenPair(currentUser.ID, accessKey, 10)
		assert.Nil(t, err)
		fetchedAdminEntity := &entity.AdminUserEntity{}
		_, err = authenticationService.VerifyAccessToken(accessToken, fetchedAdminEntity)
		assert.Nil(t, err)
		authenticationService.LogoutAllSessions(currentUser.ID)
		_, err = authenticationService.VerifyAccessToken(accessToken, fetchedAdminEntity)
		assert.NotNil(t, err)
	})

	t.Run("logout from both sessions", func(t *testing.T) {
		createContextMyApp(t, "my-app", nil,
			registry.ServiceProviderJWT(),
			registry.ServiceProviderPassword(),
			registry.ServiceProviderAuthentication(),
		)

		passwordService, _ := service.DI().Password()
		hashedPassword, _ := passwordService.HashPassword("1234")

		currentUser := createUser(map[string]interface{}{
			"Email":    "test@test.com",
			"Password": hashedPassword,
		})

		accessKey1 := fmt.Sprintf("ACCESS:%d:%s", currentUser.ID, helper.GenerateUUID())
		accessKey2 := fmt.Sprintf("ACCESS:%d:%s", currentUser.ID, helper.GenerateUUID())
		ormService, _ := service.DI().OrmEngine()
		ormService.GetRedis().Set(accessKey1, "", 10)
		ormService.GetRedis().Set(accessKey2, "", 10)

		accessListKey := fmt.Sprintf("USER_KEYS:%d", currentUser.ID)
		ormService.GetRedis().Set(accessListKey, accessKey1+";"+accessKey2, 10)

		authenticationService, _ := service.DI().AuthenticationService()
		accessToken1, err := authenticationService.GenerateTokenPair(currentUser.ID, accessKey1, 10)
		assert.Nil(t, err)
		accessToken2, err := authenticationService.GenerateTokenPair(currentUser.ID, accessKey2, 10)
		assert.Nil(t, err)
		fetchedAdminEntity := &entity.AdminUserEntity{}
		_, err = authenticationService.VerifyAccessToken(accessToken1, fetchedAdminEntity)
		assert.Nil(t, err)
		_, err = authenticationService.VerifyAccessToken(accessToken2, fetchedAdminEntity)
		assert.Nil(t, err)
		authenticationService.LogoutAllSessions(currentUser.ID)
		_, err = authenticationService.VerifyAccessToken(accessToken1, fetchedAdminEntity)
		assert.NotNil(t, err)
		_, err = authenticationService.VerifyAccessToken(accessToken2, fetchedAdminEntity)
		assert.NotNil(t, err)
	})
}

func TestGenerateTokenPair(t *testing.T) {
	createContextMyApp(t, "my-app", nil,
		registry.ServiceProviderJWT(),
		registry.ServiceProviderPassword(),
		registry.ServiceProviderAuthentication(),
	)

	passwordService, _ := service.DI().Password()
	hashedPassword, _ := passwordService.HashPassword("1234")

	currentUser := createUser(map[string]interface{}{
		"Email":    "test@test.com",
		"Password": hashedPassword,
	})

	authenticationService, _ := service.DI().AuthenticationService()

	ormService, _ := service.DI().OrmEngine()
	ormService.GetRedis().Set("test_key", "", 10)
	accessToken, err := authenticationService.GenerateTokenPair(currentUser.ID, "test_key", 10)
	assert.Nil(t, err)
	fetchedAdminEntity := &entity.AdminUserEntity{}
	payload, err := authenticationService.VerifyAccessToken(accessToken, fetchedAdminEntity)
	assert.Nil(t, err)
	assert.Equal(t, "test_key", payload["jti"])
	assert.Equal(t, fmt.Sprint(currentUser.ID), payload["sub"])
}
