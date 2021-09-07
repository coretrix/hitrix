package service

import (
	"context"

	"github.com/coretrix/hitrix/service/component/crud"

	"github.com/coretrix/hitrix/service/component/checkout"

	s3 "github.com/coretrix/hitrix/service/component/amazon/storage"
	apilogger "github.com/coretrix/hitrix/service/component/api_logger"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/authentication"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	fileextractor "github.com/coretrix/hitrix/service/component/file_extractor"
	"github.com/coretrix/hitrix/service/component/generator"
	"github.com/coretrix/hitrix/service/component/jwt"
	"github.com/coretrix/hitrix/service/component/localizer"
	"github.com/coretrix/hitrix/service/component/mail"
	"github.com/coretrix/hitrix/service/component/oss"
	"github.com/coretrix/hitrix/service/component/password"
	slackapi "github.com/coretrix/hitrix/service/component/slack_api"
	"github.com/coretrix/hitrix/service/component/sms"
	"github.com/coretrix/hitrix/service/component/social"
	"github.com/coretrix/hitrix/service/component/socket"
	"github.com/coretrix/hitrix/service/component/stripe"
	"github.com/coretrix/hitrix/service/component/uploader"

	"github.com/latolukasz/beeorm"
)

const (
	AppService              = "app"
	ConfigService           = "config"
	ErrorLoggerService      = "error_logger"
	LocalizerService        = "localizer"
	ExtractorService        = "extractor"
	JWTService              = "jwt"
	DDOSService             = "ddos"
	FCMService              = "fcm"
	ORMConfigService        = "orm_config"
	ORMEngineGlobalService  = "orm_engine_global"
	ORMEngineRequestService = "orm_engine_request"
	OSSGoogleService        = "oss_google"
	PasswordService         = "password"
	SlackAPIService         = "slack_api"
	AmazonS3Service         = "amazon_s3"
	UploaderService         = "uploader"
	StripeService           = "stripe"
	CheckoutService         = "checkout"
	DynamicLinkService      = "dynamic_link"
	SocketRegistryService   = "socket_registry"
	APILoggerService        = "api_logger"
	AuthenticationService   = "authentication"
	ClockService            = "clock"
	SMSService              = "sms"
	GeneratorService        = "generator_service"
	MailMandrill            = "mail_mandrill"
	GoogleService           = "google"
	FacebookService         = "facebook"
	CrudService             = "crud"
)

type DIInterface interface {
	App() *app.App
	Config() config.IConfig
	OrmConfig() (beeorm.ValidatedRegistry, bool)
	OrmEngine() (*beeorm.Engine, bool)
	OrmEngineForContext(ctx context.Context) *beeorm.Engine
	JWT() (*jwt.JWT, bool)
	Password() (*password.Password, bool)
	SlackAPI() (*slackapi.SlackAPI, bool)
	ErrorLogger() (errorlogger.ErrorLogger, bool)
	OSSGoogle() (oss.Client, bool)
	AmazonS3() (s3.Client, bool)
	SocketRegistry() (*socket.Registry, bool)
	APILoggerService() (apilogger.APILogger, bool)
	AuthenticationService() (*authentication.Authentication, bool)
	SMSService() (sms.ISender, bool)
	GeneratorService() (generator.Generator, bool)
	MailMandrillService() mail.Sender
	Stripe() (stripe.IStripe, bool)
	GoogleService() *social.Google
	Checkout() (checkout.ICheckout, bool)
	ClockService() clock.Clock
	UploaderService() (uploader.Uploader, bool)
	CrudService() *crud.Crud
	LocalizerService() localizer.Localizer
	FileExtractorService() *fileextractor.FileExtractor
}

type diContainer struct {
}

func (d *diContainer) Checkout() (checkout.ICheckout, bool) {
	v, has := GetServiceOptional(CheckoutService)
	if has {
		return v.(checkout.ICheckout), true
	}
	return nil, false
}
func (d *diContainer) AmazonS3() (s3.Client, bool) {
	v, has := GetServiceOptional(AmazonS3Service)
	if has {
		return v.(s3.Client), true
	}
	return nil, false
}

func (d *diContainer) Stripe() (stripe.IStripe, bool) {
	v, has := GetServiceOptional(StripeService)
	if has {
		return v.(stripe.IStripe), true
	}
	return nil, false
}

var dicInstance = &diContainer{}

func DI() DIInterface {
	return dicInstance
}

func (d *diContainer) App() *app.App {
	return GetServiceRequired(AppService).(*app.App)
}

func (d *diContainer) Config() config.IConfig {
	return GetServiceRequired(ConfigService).(config.IConfig)
}

func (d *diContainer) OrmConfig() (beeorm.ValidatedRegistry, bool) {
	v, has := GetServiceOptional(ORMConfigService)
	if has {
		return v.(beeorm.ValidatedRegistry), true
	}
	return nil, false
}

func (d *diContainer) OrmEngine() (*beeorm.Engine, bool) {
	v, has := GetServiceOptional(ORMEngineGlobalService)
	if has {
		return v.(*beeorm.Engine), true
	}
	return nil, false
}

func (d *diContainer) OrmEngineForContext(ctx context.Context) *beeorm.Engine {
	return GetServiceForRequestRequired(ctx, ORMEngineRequestService).(*beeorm.Engine)
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

func (d *diContainer) ClockService() clock.Clock {
	return GetServiceRequired(ClockService).(clock.Clock)
}

func (d *diContainer) AuthenticationService() (*authentication.Authentication, bool) {
	v, has := GetServiceOptional(AuthenticationService)
	if has {
		return v.(*authentication.Authentication), true
	}
	return nil, false
}

func (d *diContainer) MailMandrillService() mail.Sender {
	return GetServiceRequired(MailMandrill).(mail.Sender)
}

func (d *diContainer) GoogleService() *social.Google {
	return GetServiceRequired(GoogleService).(*social.Google)
}

func (d *diContainer) UploaderService() (uploader.Uploader, bool) {
	v, has := GetServiceOptional(UploaderService)
	if has {
		return v.(uploader.Uploader), true
	}
	return nil, false
}

func (d *diContainer) CrudService() *crud.Crud {
	return GetServiceRequired(CrudService).(*crud.Crud)
}

func (d *diContainer) LocalizerService() localizer.Localizer {
	return GetServiceRequired(LocalizerService).(localizer.Localizer)
}

func (d *diContainer) FileExtractorService() *fileextractor.FileExtractor {
	return GetServiceRequired(ExtractorService).(*fileextractor.FileExtractor)
}
