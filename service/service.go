package service

import (
	"context"

	"github.com/coretrix/hitrix/service/component/otp"

	"github.com/coretrix/hitrix/service/component/uuid"

	"github.com/coretrix/hitrix/service/component/goroutine"
	"github.com/coretrix/hitrix/service/component/localize"

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
	"github.com/coretrix/hitrix/service/component/mail"
	"github.com/coretrix/hitrix/service/component/oss"
	"github.com/coretrix/hitrix/service/component/password"
	"github.com/coretrix/hitrix/service/component/slack"
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
	LocalizeService         = "localize"
	PDFService              = "pdf"
	ExtractorService        = "extractor"
	JWTService              = "jwt"
	DDOSService             = "ddos"
	FCMService              = "fcm"
	ORMConfigService        = "orm_config"
	ORMEngineGlobalService  = "orm_engine_global"
	ORMEngineRequestService = "orm_engine_request"
	OSSGoogleService        = "oss_google"
	PasswordService         = "password"
	SlackService            = "slack"
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
	GoroutineService        = "goroutine"
	GeneratorService        = "generator"
	MailMandrillService     = "mail_mandrill"
	GoogleService           = "google"
	FacebookService         = "facebook"
	CrudService             = "crud"
	UUIDService             = "uuid"
	OTPService              = "otp"
)

type DIInterface interface {
	App() *app.App
	Config() config.IConfig
	OrmConfig() (beeorm.ValidatedRegistry, bool)
	OrmEngine() (*beeorm.Engine, bool)
	OrmEngineForContext(ctx context.Context) *beeorm.Engine
	JWT() (*jwt.JWT, bool)
	Password() (password.IPassword, bool)
	Slack() (slack.Slack, bool)
	ErrorLogger() (errorlogger.ErrorLogger, bool)
	OSSGoogle() (oss.Client, bool)
	AmazonS3() (s3.Client, bool)
	SocketRegistry() (*socket.Registry, bool)
	APILogger() (apilogger.IAPILogger, bool)
	Authentication() (*authentication.Authentication, bool)
	SMS() (sms.ISender, bool)
	Generator() (generator.IGenerator, bool)
	MailMandrill() mail.Sender
	Stripe() (stripe.IStripe, bool)
	Google() *social.Google
	Checkout() (checkout.ICheckout, bool)
	Clock() clock.IClock
	Uploader() (uploader.Uploader, bool)
	CrudService() *crud.Crud
	Localize() localize.ILocalizer
	FileExtractor() *fileextractor.FileExtractor
	Goroutine() goroutine.IGoroutine
	UUID() uuid.IUUID
	OTP() otp.IOTP
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

func (d *diContainer) SMS() (sms.ISender, bool) {
	v, has := GetServiceOptional(SMSService)
	if has {
		return v.(sms.ISender), true
	}
	return nil, false
}

func (d *diContainer) Generator() (generator.IGenerator, bool) {
	v, has := GetServiceOptional(GeneratorService)
	if has {
		return v.(generator.IGenerator), true
	}
	return nil, false
}

func (d *diContainer) Password() (password.IPassword, bool) {
	v, has := GetServiceOptional(PasswordService)
	if has {
		return v.(password.IPassword), true
	}
	return nil, false
}

func (d *diContainer) Slack() (slack.Slack, bool) {
	v, has := GetServiceOptional(SlackService)
	if has {
		return v.(slack.Slack), true
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

func (d *diContainer) APILogger() (apilogger.IAPILogger, bool) {
	v, has := GetServiceOptional(APILoggerService)
	if has {
		return v.(apilogger.IAPILogger), true
	}
	return nil, false
}

func (d *diContainer) Clock() clock.IClock {
	return GetServiceRequired(ClockService).(clock.IClock)
}

func (d *diContainer) Authentication() (*authentication.Authentication, bool) {
	v, has := GetServiceOptional(AuthenticationService)
	if has {
		return v.(*authentication.Authentication), true
	}
	return nil, false
}

func (d *diContainer) MailMandrill() mail.Sender {
	return GetServiceRequired(MailMandrillService).(mail.Sender)
}

func (d *diContainer) Google() *social.Google {
	return GetServiceRequired(GoogleService).(*social.Google)
}

func (d *diContainer) Uploader() (uploader.Uploader, bool) {
	v, has := GetServiceOptional(UploaderService)
	if has {
		return v.(uploader.Uploader), true
	}
	return nil, false
}

func (d *diContainer) CrudService() *crud.Crud {
	return GetServiceRequired(CrudService).(*crud.Crud)
}

func (d *diContainer) Localize() localize.ILocalizer {
	return GetServiceRequired(LocalizeService).(localize.ILocalizer)
}

func (d *diContainer) FileExtractor() *fileextractor.FileExtractor {
	return GetServiceRequired(ExtractorService).(*fileextractor.FileExtractor)
}

func (d *diContainer) Goroutine() goroutine.IGoroutine {
	return GetServiceRequired(GoroutineService).(goroutine.IGoroutine)
}

func (d *diContainer) UUID() uuid.IUUID {
	return GetServiceRequired(UUIDService).(uuid.IUUID)
}

func (d *diContainer) OTP() otp.IOTP {
	return GetServiceRequired(OTPService).(otp.IOTP)
}
