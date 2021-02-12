package hitrix

import (
	"context"

	apexLog "github.com/apex/log"
	"github.com/summer-solutions/orm"
)

type DICInterface interface {
	App() *AppDefinition
	Log() apexLog.Interface
	Config() *Config
	OrmConfig() (orm.ValidatedRegistry, bool)
	OrmEngine() (*orm.Engine, bool)
	LogForContext(ctx context.Context) *RequestLog
	OrmEngineForContext(ctx context.Context) (*orm.Engine, bool)
	JWT() (*JWT, bool)
	Password() (*Password, bool)
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

func (d *dic) Log() apexLog.Interface {
	return GetServiceRequired("log").(apexLog.Interface)
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

func (d *dic) LogForContext(ctx context.Context) *RequestLog {
	return GetServiceForRequestRequired(ctx, "log_request").(*RequestLog)
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
