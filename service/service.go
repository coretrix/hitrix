package service

import (
	"context"

	"github.com/coretrix/hitrix/service/component/generator"

	"github.com/coretrix/hitrix/service/component/sms"

	"github.com/coretrix/hitrix/service/component/authentication"

	"github.com/coretrix/hitrix/service/component/clock"

	apilogger "github.com/coretrix/hitrix/service/component/api_logger"

	"github.com/coretrix/hitrix/service/component/socket"

	"github.com/coretrix/hitrix/service/component/oss"

	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	slackapi "github.com/coretrix/hitrix/service/component/slack_api"

	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/jwt"
	"github.com/coretrix/hitrix/service/component/password"

	"github.com/coretrix/hitrix/service/component/app"

	"github.com/latolukasz/orm"
)

const (
	AppService              = "app"
	ConfigService           = "config"
	ErrorLoggerService      = "error_logger"
	JWTService              = "jwt"
	DDOSService             = "ddos"
	ORMConfigService        = "orm_config"
	ORMEngineGlobalService  = "orm_engine_global"
	ORMEngineRequestService = "orm_engine_request"
	OSSGoogleService        = "oss_google"
	PasswordService         = "password"
	SlackAPIService         = "slack_api"
	SocketRegistryService   = "socket_registry"
	APILoggerService        = "api_logger"
	AuthenticationService   = "authentication"
	ClockService            = "clock"
	SMSService              = "sms"
	GeneratorService        = "generator_service"
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
	OSSGoogle() (oss.Client, bool)
	SocketRegistry() (*socket.Registry, bool)
	APILoggerService() (apilogger.APILogger, bool)
	AuthenticationService() (*authentication.Authentication, bool)
	SMSService() (sms.ISender, bool)
	GeneratorService() (generator.Generator, bool)
}

type diContainer struct {
}

var dicInstance = &diContainer{}

func DI() DIInterface {
	return dicInstance
}

func (d *diContainer) App() *app.App {
	return GetServiceRequired(AppService).(*app.App)
}

func (d *diContainer) Config() *config.Config {
	return GetServiceRequired(ConfigService).(*config.Config)
}

func (d *diContainer) OrmConfig() (orm.ValidatedRegistry, bool) {
	v, has := GetServiceOptional(ORMConfigService)
	if has {
		return v.(orm.ValidatedRegistry), true
	}
	return nil, false
}

func (d *diContainer) OrmEngine() (*orm.Engine, bool) {
	v, has := GetServiceOptional(ORMEngineGlobalService)
	if has {
		return v.(*orm.Engine), true
	}
	return nil, false
}

func (d *diContainer) OrmEngineForContext(ctx context.Context) (*orm.Engine, bool) {
	v, has := GetServiceForRequestOptional(ctx, ORMEngineRequestService)
	if has {
		return v.(*orm.Engine), true
	}
	return nil, false
}

func (d *diContainer) JWT() (*jwt.JWT, bool) {
	v, has := GetServiceOptional(JWTService)
	if has {
		return v.(*jwt.JWT), true
	}
	return nil, false
}

func (d *diContainer) SMSService() (sms.ISender, bool) {
	v, has := GetServiceOptional(SMSService)
	if has {
		return v.(sms.ISender), true
	}
	return nil, false
}

func (d *diContainer) GeneratorService() (generator.Generator, bool) {
	v, has := GetServiceOptional(GeneratorService)
	if has {
		return v.(generator.Generator), true
	}
	return nil, false
}

func (d *diContainer) Password() (*password.Password, bool) {
	v, has := GetServiceOptional(PasswordService)
	if has {
		return v.(*password.Password), true
	}
	return nil, false
}

func (d *diContainer) SlackAPI() (*slackapi.SlackAPI, bool) {
	v, has := GetServiceOptional(SlackAPIService)
	if has {
		return v.(*slackapi.SlackAPI), true
	}
	return nil, false
}

func (d *diContainer) ErrorLogger() (errorlogger.ErrorLogger, bool) {
	v, has := GetServiceOptional(ErrorLoggerService)
	if has {
		return v.(errorlogger.ErrorLogger), true
	}
	return nil, false
}

func (d *diContainer) OSSGoogle() (oss.Client, bool) {
	v, has := GetServiceOptional(OSSGoogleService)
	if has {
		return v.(oss.Client), true
	}
	return nil, false
}

func (d *diContainer) SocketRegistry() (*socket.Registry, bool) {
	v, has := GetServiceOptional(SocketRegistryService)
	if has {
		return v.(*socket.Registry), true
	}
	return nil, false
}

func (d *diContainer) APILoggerService() (apilogger.APILogger, bool) {
	v, has := GetServiceOptional(APILoggerService)
	if has {
		return v.(apilogger.APILogger), true
	}
	return nil, false
}

func (d *diContainer) ClockService() (clock.Clock, bool) {
	v, has := GetServiceOptional(ClockService)
	if has {
		return v.(clock.Clock), true
	}
	return nil, false
}

func (d *diContainer) AuthenticationService() (*authentication.Authentication, bool) {
	v, has := GetServiceOptional(AuthenticationService)
	if has {
		return v.(*authentication.Authentication), true
	}
	return nil, false
}
