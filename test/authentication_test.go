package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/service"
	generatorMock "github.com/coretrix/hitrix/service/component/generator/mocks"
	"github.com/coretrix/hitrix/service/component/password"
	smsMock "github.com/coretrix/hitrix/service/component/sms/mocks"
	"github.com/coretrix/hitrix/service/registry"
	"github.com/coretrix/hitrix/service/registry/mocks"
)

func createUser(input map[string]interface{}) *entity.DevPanelUserEntity {
	ormService := service.DI().OrmEngine()
	devPanelUserEntity := &entity.DevPanelUserEntity{}

	for field, val := range input {
		switch field {
		case "Email":
			devPanelUserEntity.Email = val.(string)
		case "Password":
			devPanelUserEntity.Password = val.(string)
		}
	}

	ormService.Flush(devPanelUserEntity)

	return devPanelUserEntity
}

func TestAuthenticate(t *testing.T) {
	t.Run("simple successful", func(t *testing.T) {
		fakeSMS := &smsMock.FakeSMSSender{}
		fakeGenerator := &generatorMock.FakeGenerator{}

		fakeGenerator.On("GenerateUUID").Return("randomid")

		createContextMyApp(t, "my-app", nil,
			[]*service.DefinitionGlobal{
				registry.ServiceProviderErrorLogger(),
				registry.ServiceProviderJWT(),
				registry.ServiceProviderPassword(password.NewSimpleManager),
				registry.ServiceProviderUUID(),
				registry.ServiceProviderAuthentication(),
				registry.ServiceProviderClock(),
				mocks.ServiceProviderMockSMS(fakeSMS),
				mocks.ServiceProviderMockGenerator(fakeGenerator),
			},
			nil,
		)

		passwordService := service.DI().Password()
		hashedPassword, _ := passwordService.HashPassword("1234")
		ormService := service.DI().OrmEngine()

		createUser(map[string]interface{}{
			"Email":    "test@test.com",
			"Password": hashedPassword,
		})

		authenticationService := service.DI().Authentication()
		fetchedAdminEntity := &entity.DevPanelUserEntity{}
		_, _, err := authenticationService.Authenticate(ormService, "test@test.com", "1234", fetchedAdminEntity)
		assert.Nil(t, err)
	})

	t.Run("simple successful by id", func(t *testing.T) {
		fakeSMS := &smsMock.FakeSMSSender{}
		fakeGenerator := &generatorMock.FakeGenerator{}

		fakeGenerator.On("GenerateUUID").Return("randomid")

		createContextMyApp(t, "my-app", nil,
			[]*service.DefinitionGlobal{
				registry.ServiceProviderErrorLogger(),
				registry.ServiceProviderJWT(),
				registry.ServiceProviderPassword(password.NewSimpleManager),
				registry.ServiceProviderUUID(),
				registry.ServiceProviderAuthentication(),
				registry.ServiceProviderClock(),
				mocks.ServiceProviderMockSMS(fakeSMS),
				mocks.ServiceProviderMockGenerator(fakeGenerator),
			},
			nil,
		)

		passwordService := service.DI().Password()
		hashedPassword, _ := passwordService.HashPassword("1234")
		ormService := service.DI().OrmEngine()

		userEntity := createUser(map[string]interface{}{
			"Email":    "test@test.com",
			"Password": hashedPassword,
		})

		authenticationService := service.DI().Authentication()
		fetchedAdminEntity := &entity.DevPanelUserEntity{}
		_, _, err := authenticationService.AuthenticateByID(ormService, userEntity.GetID(), fetchedAdminEntity)
		assert.Nil(t, err)
		assert.Equal(t, fetchedAdminEntity.GetID(), userEntity.GetID())
	})

	t.Run("wrong email", func(t *testing.T) {
		fakeSMS := &smsMock.FakeSMSSender{}
		fakeGenerator := &generatorMock.FakeGenerator{}
		createContextMyApp(t, "my-app", nil,
			[]*service.DefinitionGlobal{
				registry.ServiceProviderErrorLogger(),
				registry.ServiceProviderJWT(),
				registry.ServiceProviderPassword(password.NewSimpleManager),
				registry.ServiceProviderUUID(),
				registry.ServiceProviderAuthentication(),
				registry.ServiceProviderClock(),
				mocks.ServiceProviderMockSMS(fakeSMS),
				mocks.ServiceProviderMockGenerator(fakeGenerator),
			},
			nil,
		)

		passwordService := service.DI().Password()
		hashedPassword, _ := passwordService.HashPassword("1234")
		ormService := service.DI().OrmEngine()

		createUser(map[string]interface{}{
			"Email":    "test@test.com",
			"Password": hashedPassword,
		})

		authenticationService := service.DI().Authentication()
		fetchedAdminEntity := &entity.DevPanelUserEntity{}
		_, _, err := authenticationService.Authenticate(ormService, "test@tesat.com", "1234", fetchedAdminEntity)
		assert.NotNil(t, err)
	})
}

func TestVerifyAccessToken(t *testing.T) {
	fakeSMS := &smsMock.FakeSMSSender{}
	fakeGenerator := &generatorMock.FakeGenerator{}

	t.Run("simple success", func(t *testing.T) {
		createContextMyApp(t, "my-app", nil,
			[]*service.DefinitionGlobal{
				registry.ServiceProviderErrorLogger(),
				registry.ServiceProviderJWT(),
				registry.ServiceProviderPassword(password.NewSimpleManager),
				registry.ServiceProviderUUID(),
				registry.ServiceProviderAuthentication(),
				registry.ServiceProviderClock(),
				mocks.ServiceProviderMockSMS(fakeSMS),
				mocks.ServiceProviderMockGenerator(fakeGenerator),
			},
			nil,
		)

		passwordService := service.DI().Password()
		hashedPassword, _ := passwordService.HashPassword("1234")

		currentUser := createUser(map[string]interface{}{
			"Email":    "test@test.com",
			"Password": hashedPassword,
		})

		accessKey := fmt.Sprintf("ACCESS:%d:%s", currentUser.ID, service.DI().UUID().Generate())
		ormService := service.DI().OrmEngine()
		ormService.GetRedis().Set(accessKey, "", 10)

		authenticationService := service.DI().Authentication()
		token, err := authenticationService.GenerateTokenPair(currentUser.ID, accessKey, 10)
		assert.Nil(t, err)
		fetchedAdminEntity := &entity.DevPanelUserEntity{}
		payload, err := authenticationService.VerifyAccessToken(ormService, token, fetchedAdminEntity)
		assert.Nil(t, err)
		assert.Equal(t, accessKey, payload["jti"])
	})

	t.Run("wrong token", func(t *testing.T) {
		createContextMyApp(t, "my-app", nil,
			[]*service.DefinitionGlobal{
				registry.ServiceProviderErrorLogger(),
				registry.ServiceProviderJWT(),
				registry.ServiceProviderPassword(password.NewSimpleManager),
				registry.ServiceProviderUUID(),
				registry.ServiceProviderAuthentication(),
				registry.ServiceProviderClock(),
				mocks.ServiceProviderMockSMS(fakeSMS),
				mocks.ServiceProviderMockGenerator(fakeGenerator),
			},
			nil,
		)

		passwordService := service.DI().Password()
		hashedPassword, _ := passwordService.HashPassword("1234")

		currentUser := createUser(map[string]interface{}{
			"Email":    "test@test.com",
			"Password": hashedPassword,
		})

		accessKey := fmt.Sprintf("ACCESS:%d:%s", currentUser.ID, service.DI().UUID().Generate())
		ormService := service.DI().OrmEngine()
		ormService.GetRedis().Set(accessKey, "", 10)

		authenticationService := service.DI().Authentication()
		token, err := authenticationService.GenerateTokenPair(currentUser.ID, accessKey, 10)
		assert.Nil(t, err)
		fetchedAdminEntity := &entity.DevPanelUserEntity{}
		_, err = authenticationService.VerifyAccessToken(ormService, " wef"+token, fetchedAdminEntity)
		assert.NotNil(t, err)
	})
}

func TestRefreshToken(t *testing.T) {
	fakeSMS := &smsMock.FakeSMSSender{}
	fakeGenerator := &generatorMock.FakeGenerator{}
	fakeGenerator.On("GenerateUUID").Return("randomid")
	t.Run("success refresh", func(t *testing.T) {
		createContextMyApp(t, "my-app", nil,
			[]*service.DefinitionGlobal{
				registry.ServiceProviderErrorLogger(),
				registry.ServiceProviderJWT(),
				registry.ServiceProviderPassword(password.NewSimpleManager),
				registry.ServiceProviderUUID(),
				registry.ServiceProviderAuthentication(),
				registry.ServiceProviderClock(),
				mocks.ServiceProviderMockSMS(fakeSMS),
				mocks.ServiceProviderMockGenerator(fakeGenerator),
			},
			nil,
		)

		passwordService := service.DI().Password()
		hashedPassword, _ := passwordService.HashPassword("1234")

		currentUser := createUser(map[string]interface{}{
			"Email":    "test@test.com",
			"Password": hashedPassword,
		})

		accessKey := fmt.Sprintf("ACCESS:%d:%s", currentUser.ID, service.DI().UUID().Generate())
		ormService := service.DI().OrmEngine()
		ormService.GetRedis().Set(accessKey, "", 10)

		authenticationService := service.DI().Authentication()
		refresh, err := authenticationService.GenerateTokenPair(currentUser.ID, accessKey, 10)
		assert.Nil(t, err)
		_, _, err = authenticationService.RefreshToken(ormService, refresh)
		assert.Nil(t, err)
	})

	t.Run("wrong refresh", func(t *testing.T) {
		createContextMyApp(t, "my-app", nil,
			[]*service.DefinitionGlobal{
				registry.ServiceProviderErrorLogger(),
				registry.ServiceProviderJWT(),
				registry.ServiceProviderPassword(password.NewSimpleManager),
				registry.ServiceProviderUUID(),
				registry.ServiceProviderAuthentication(),
				registry.ServiceProviderClock(),
				mocks.ServiceProviderMockSMS(fakeSMS),
				mocks.ServiceProviderMockGenerator(fakeGenerator),
			},
			nil,
		)

		passwordService := service.DI().Password()
		hashedPassword, _ := passwordService.HashPassword("1234")

		currentUser := createUser(map[string]interface{}{
			"Email":    "test@test.com",
			"Password": hashedPassword,
		})

		accessKey := fmt.Sprintf("ACCESS:%d:%s", currentUser.ID, service.DI().UUID().Generate())
		ormService := service.DI().OrmEngine()
		ormService.GetRedis().Set(accessKey, "", 10)

		authenticationService := service.DI().Authentication()
		refresh, err := authenticationService.GenerateTokenPair(currentUser.ID, accessKey, 10)
		assert.Nil(t, err)
		_, _, err = authenticationService.RefreshToken(ormService, "ef"+refresh)
		assert.NotNil(t, err)
	})
}

func TestLogoutCurrentSession(t *testing.T) {
	fakeSMS := &smsMock.FakeSMSSender{}
	fakeGenerator := &generatorMock.FakeGenerator{}

	t.Run("simple logout", func(t *testing.T) {
		createContextMyApp(t, "my-app", nil,
			[]*service.DefinitionGlobal{
				registry.ServiceProviderErrorLogger(),
				registry.ServiceProviderJWT(),
				registry.ServiceProviderPassword(password.NewSimpleManager),
				registry.ServiceProviderUUID(),
				registry.ServiceProviderAuthentication(),
				registry.ServiceProviderClock(),
				mocks.ServiceProviderMockSMS(fakeSMS),
				mocks.ServiceProviderMockGenerator(fakeGenerator),
			},
			nil,
		)

		passwordService := service.DI().Password()
		hashedPassword, _ := passwordService.HashPassword("1234")

		currentUser := createUser(map[string]interface{}{
			"Email":    "test@test.com",
			"Password": hashedPassword,
		})

		accessKey := fmt.Sprintf("ACCESS:%d:%s", currentUser.ID, service.DI().UUID().Generate())
		ormService := service.DI().OrmEngine()
		ormService.GetRedis().Set(accessKey, "", 10)

		authenticationService := service.DI().Authentication()
		accessToken, err := authenticationService.GenerateTokenPair(currentUser.ID, accessKey, 10)
		assert.Nil(t, err)
		fetchedAdminEntity := &entity.DevPanelUserEntity{}
		_, err = authenticationService.VerifyAccessToken(ormService, accessToken, fetchedAdminEntity)
		assert.Nil(t, err)
		authenticationService.LogoutCurrentSession(ormService, accessKey)
		_, err = authenticationService.VerifyAccessToken(ormService, accessToken, fetchedAdminEntity)
		assert.NotNil(t, err)
	})
}

func TestLogoutAllSessions(t *testing.T) {
	fakeSMS := &smsMock.FakeSMSSender{}
	fakeGenerator := &generatorMock.FakeGenerator{}

	t.Run("logout from one session", func(t *testing.T) {
		createContextMyApp(t, "my-app", nil,
			[]*service.DefinitionGlobal{
				registry.ServiceProviderErrorLogger(),
				registry.ServiceProviderJWT(),
				registry.ServiceProviderPassword(password.NewSimpleManager),
				registry.ServiceProviderUUID(),
				registry.ServiceProviderAuthentication(),
				registry.ServiceProviderClock(),
				mocks.ServiceProviderMockSMS(fakeSMS),
				mocks.ServiceProviderMockGenerator(fakeGenerator),
			},
			nil,
		)

		passwordService := service.DI().Password()
		hashedPassword, _ := passwordService.HashPassword("1234")

		currentUser := createUser(map[string]interface{}{
			"Email":    "test@test.com",
			"Password": hashedPassword,
		})

		accessKey := fmt.Sprintf("ACCESS:%d:%s", currentUser.ID, service.DI().UUID().Generate())
		ormService := service.DI().OrmEngine()
		ormService.GetRedis().Set(accessKey, "", 10)

		accessListKey := fmt.Sprintf("USER_KEYS:%d", currentUser.ID)
		ormService.GetRedis().Set(accessListKey, accessKey, 10)

		authenticationService := service.DI().Authentication()
		accessToken, err := authenticationService.GenerateTokenPair(currentUser.ID, accessKey, 10)
		assert.Nil(t, err)
		fetchedAdminEntity := &entity.DevPanelUserEntity{}
		_, err = authenticationService.VerifyAccessToken(ormService, accessToken, fetchedAdminEntity)
		assert.Nil(t, err)
		authenticationService.LogoutAllSessions(ormService, currentUser.ID)
		_, err = authenticationService.VerifyAccessToken(ormService, accessToken, fetchedAdminEntity)
		assert.NotNil(t, err)
	})

	t.Run("logout from both sessions", func(t *testing.T) {
		createContextMyApp(t, "my-app", nil,
			[]*service.DefinitionGlobal{
				registry.ServiceProviderErrorLogger(),
				registry.ServiceProviderJWT(),
				registry.ServiceProviderPassword(password.NewSimpleManager),
				registry.ServiceProviderUUID(),
				registry.ServiceProviderAuthentication(),
				registry.ServiceProviderClock(),
				mocks.ServiceProviderMockSMS(fakeSMS),
				mocks.ServiceProviderMockGenerator(fakeGenerator),
			},
			nil,
		)

		passwordService := service.DI().Password()
		hashedPassword, _ := passwordService.HashPassword("1234")

		currentUser := createUser(map[string]interface{}{
			"Email":    "test@test.com",
			"Password": hashedPassword,
		})

		accessKey1 := fmt.Sprintf("ACCESS:%d:%s", currentUser.ID, service.DI().UUID().Generate())
		accessKey2 := fmt.Sprintf("ACCESS:%d:%s", currentUser.ID, service.DI().UUID().Generate())
		ormService := service.DI().OrmEngine()
		ormService.GetRedis().Set(accessKey1, "", 10)
		ormService.GetRedis().Set(accessKey2, "", 10)

		accessListKey := fmt.Sprintf("USER_KEYS:%d", currentUser.ID)
		ormService.GetRedis().Set(accessListKey, accessKey1+";"+accessKey2, 10)

		authenticationService := service.DI().Authentication()
		accessToken1, err := authenticationService.GenerateTokenPair(currentUser.ID, accessKey1, 10)
		assert.Nil(t, err)
		accessToken2, err := authenticationService.GenerateTokenPair(currentUser.ID, accessKey2, 10)
		assert.Nil(t, err)
		fetchedAdminEntity := &entity.DevPanelUserEntity{}
		_, err = authenticationService.VerifyAccessToken(ormService, accessToken1, fetchedAdminEntity)
		assert.Nil(t, err)
		_, err = authenticationService.VerifyAccessToken(ormService, accessToken2, fetchedAdminEntity)
		assert.Nil(t, err)
		authenticationService.LogoutAllSessions(ormService, currentUser.ID)
		_, err = authenticationService.VerifyAccessToken(ormService, accessToken1, fetchedAdminEntity)
		assert.NotNil(t, err)
		_, err = authenticationService.VerifyAccessToken(ormService, accessToken2, fetchedAdminEntity)
		assert.NotNil(t, err)
	})
}

func TestGenerateTokenPair(t *testing.T) {
	fakeSMS := &smsMock.FakeSMSSender{}
	fakeGenerator := &generatorMock.FakeGenerator{}
	createContextMyApp(t, "my-app", nil,
		[]*service.DefinitionGlobal{
			registry.ServiceProviderErrorLogger(),
			registry.ServiceProviderJWT(),
			registry.ServiceProviderPassword(password.NewSimpleManager),
			registry.ServiceProviderUUID(),
			registry.ServiceProviderAuthentication(),
			registry.ServiceProviderClock(),
			mocks.ServiceProviderMockSMS(fakeSMS),
			mocks.ServiceProviderMockGenerator(fakeGenerator),
		},
		nil,
	)

	passwordService := service.DI().Password()
	hashedPassword, _ := passwordService.HashPassword("1234")

	currentUser := createUser(map[string]interface{}{
		"Email":    "test@test.com",
		"Password": hashedPassword,
	})

	authenticationService := service.DI().Authentication()

	ormService := service.DI().OrmEngine()
	ormService.GetRedis().Set("test_key", "", 10)

	accessToken, err := authenticationService.GenerateTokenPair(currentUser.ID, "test_key", 10)

	assert.Nil(t, err)

	fetchedAdminEntity := &entity.DevPanelUserEntity{}
	payload, err := authenticationService.VerifyAccessToken(ormService, accessToken, fetchedAdminEntity)

	assert.Nil(t, err)
	assert.Equal(t, "test_key", payload["jti"])
	assert.Equal(t, fmt.Sprint(currentUser.ID), payload["sub"])
}
