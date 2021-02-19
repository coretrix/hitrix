package service

import (
	"context"

	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	slackapi "github.com/coretrix/hitrix/service/component/slack_api"

	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/jwt"
	"github.com/coretrix/hitrix/service/component/password"

	"github.com/coretrix/hitrix/service/component/app"

	"github.com/summer-solutions/orm"
)

type DIInterface interface {
	App() *app.App
	Config() *config.Config
	OrmConfig() (orm.ValidatedRegistry, bool)
	OrmEngine() (*orm.Engine, bool)
	OrmEngineForContext(ctx context.Context) (*orm.Engine, bool)
	JWT() (*jwt.JWT, bool)
	Password() (*password.Password, bool)
	SlackAPI() (*slackapi.SlackAPI, bool)
	ErrorLogger() (errorlogger.ErrorLogger, bool)
}

type diContainer struct {
}

var dicInstance = &diContainer{}

func DI() DIInterface {
	return dicInstance
}

func (d *diContainer) App() *app.App {
	return GetServiceRequired("app").(*app.App)
}

func (d *diContainer) Config() *config.Config {
	return GetServiceRequired("config").(*config.Config)
}

func (d *diContainer) OrmConfig() (orm.ValidatedRegistry, bool) {
	v, has := GetServiceOptional("orm_config")
	if has {
		return v.(orm.ValidatedRegistry), true
	}
	return nil, false
}

func (d *diContainer) OrmEngine() (*orm.Engine, bool) {
	v, has := GetServiceOptional("orm_engine_global")
	if has {
		return v.(*orm.Engine), true
	}
	return nil, false
}

func (d *diContainer) OrmEngineForContext(ctx context.Context) (*orm.Engine, bool) {
	v, has := GetServiceForRequestOptional(ctx, "orm_engine_request")
	if has {
		return v.(*orm.Engine), true
	}
	return nil, false
}

func (d *diContainer) JWT() (*jwt.JWT, bool) {
	v, has := GetServiceOptional("jwt")
	if has {
		return v.(*jwt.JWT), true
	}
	return nil, false
}

func (d *diContainer) Password() (*password.Password, bool) {
	v, has := GetServiceOptional("password")
	if has {
		return v.(*password.Password), true
	}
	return nil, false
}

func (d *diContainer) SlackAPI() (*slackapi.SlackAPI, bool) {
	v, has := GetServiceOptional("slack_api")
	if has {
		return v.(*slackapi.SlackAPI), true
	}
	return nil, false
}

func (d *diContainer) ErrorLogger() (errorlogger.ErrorLogger, bool) {
	v, has := GetServiceOptional("error_logger")
	if has {
		return v.(errorlogger.ErrorLogger), true
	}
	return nil, false
}
