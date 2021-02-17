package hitrix

import (
	"context"

	"github.com/summer-solutions/orm"
)

type DICInterface interface {
	App() *AppDefinition
	Config() *Config
	OrmConfig() (orm.ValidatedRegistry, bool)
	OrmEngine() (*orm.Engine, bool)
	OrmEngineForContext(ctx context.Context) (*orm.Engine, bool)
	JWT() (*JWT, bool)
	Password() (*Password, bool)
	SlackAPI() (*SlackAPI, bool)
	ErrorLogger() (ErrorLogger, bool)
}

type dic struct {
}

var dicInstance = &dic{}

func DIC() DICInterface {
	return dicInstance
}

func (d *dic) App() *AppDefinition {
	return GetServiceRequired("app").(*AppDefinition)
}

func (d *dic) Config() *Config {
	return GetServiceRequired("config").(*Config)
}

func (d *dic) OrmConfig() (orm.ValidatedRegistry, bool) {
	v, has := GetServiceOptional("orm_config")
	if has {
		return v.(orm.ValidatedRegistry), true
	}
	return nil, false
}

func (d *dic) OrmEngine() (*orm.Engine, bool) {
	v, has := GetServiceOptional("orm_engine_global")
	if has {
		return v.(*orm.Engine), true
	}
	return nil, false
}

func (d *dic) OrmEngineForContext(ctx context.Context) (*orm.Engine, bool) {
	v, has := GetServiceForRequestOptional(ctx, "orm_engine_request")
	if has {
		return v.(*orm.Engine), true
	}
	return nil, false
}

func (d *dic) JWT() (*JWT, bool) {
	v, has := GetServiceOptional("jwt")
	if has {
		return v.(*JWT), true
	}
	return nil, false
}

func (d *dic) Password() (*Password, bool) {
	v, has := GetServiceOptional("password")
	if has {
		return v.(*Password), true
	}
	return nil, false
}

func (d *dic) SlackAPI() (*SlackAPI, bool) {
	v, has := GetServiceOptional("slack_api")
	if has {
		return v.(*SlackAPI), true
	}
	return nil, false
}

func (d *dic) ErrorLogger() (ErrorLogger, bool) {
	v, has := GetServiceOptional("error_logger")
	if has {
		return v.(ErrorLogger), true
	}
	return nil, false
}
