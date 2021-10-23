package service

import (
	"context"

	"github.com/coretrix/hitrix/service/component/exporter"

	"github.com/coretrix/hitrix/service/component/fcm"
	"github.com/coretrix/hitrix/service/component/pdf"

	"github.com/coretrix/hitrix/service/component/ddos"
	dynamiclink "github.com/coretrix/hitrix/service/component/dynamic_link"

	"github.com/coretrix/hitrix/service/component/otp"

	"github.com/coretrix/hitrix/service/component/uuid"

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
	OSService               = "oss"
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
	ExporterService         = "exporter"
)

type diContainer struct {
}

var dicInstance = &diContainer{}

func DI() *diContainer {
	return dicInstance
}

func (d *diContainer) Checkout() checkout.ICheckout {
	return GetServiceRequired(CheckoutService).(checkout.ICheckout)
}

func (d *diContainer) Exporter() exporter.IExporter {
	return GetServiceRequired(ExporterService).(exporter.IExporter)
}

func (d *diContainer) AmazonS3() s3.Client {
	return GetServiceRequired(AmazonS3Service).(s3.Client)
}

func (d *diContainer) Stripe() stripe.IStripe {
	return GetServiceRequired(StripeService).(stripe.IStripe)
}

func (d *diContainer) App() *app.App {
	return GetServiceRequired(AppService).(*app.App)
}

func (d *diContainer) Config() config.IConfig {
	return GetServiceRequired(ConfigService).(config.IConfig)
}

func (d *diContainer) OrmConfig() beeorm.ValidatedRegistry {
	return GetServiceRequired(ORMConfigService).(beeorm.ValidatedRegistry)
}

func (d *diContainer) OrmEngine() *beeorm.Engine {
	return GetServiceRequired(ORMEngineGlobalService).(*beeorm.Engine)
}

func (d *diContainer) OrmEngineForContext(ctx context.Context) *beeorm.Engine {
	return GetServiceForRequestRequired(ctx, ORMEngineRequestService).(*beeorm.Engine)
}

func (d *diContainer) JWT() *jwt.JWT {
	return GetServiceRequired(JWTService).(*jwt.JWT)
}

func (d *diContainer) SMS() sms.ISender {
	return GetServiceRequired(SMSService).(sms.ISender)
}

func (d *diContainer) Generator() generator.IGenerator {
	return GetServiceRequired(GeneratorService).(generator.IGenerator)
}

func (d *diContainer) Password() password.IPassword {
	return GetServiceRequired(PasswordService).(password.IPassword)
}

func (d *diContainer) Slack() slack.Slack {
	return GetServiceRequired(SlackService).(slack.Slack)
}

func (d *diContainer) ErrorLogger() errorlogger.ErrorLogger {
	return GetServiceRequired(ErrorLoggerService).(errorlogger.ErrorLogger)
}

func (d *diContainer) OSService() oss.IProvider {
	return GetServiceRequired(OSService).(oss.IProvider)
}

func (d *diContainer) SocketRegistry() *socket.Registry {
	return GetServiceRequired(SocketRegistryService).(*socket.Registry)
}

func (d *diContainer) APILogger() apilogger.IAPILogger {
	return GetServiceRequired(APILoggerService).(apilogger.IAPILogger)
}

func (d *diContainer) Clock() clock.IClock {
	return GetServiceRequired(ClockService).(clock.IClock)
}

func (d *diContainer) Authentication() *authentication.Authentication {
	return GetServiceRequired(AuthenticationService).(*authentication.Authentication)
}

func (d *diContainer) MailMandrill() mail.Sender {
	return GetServiceRequired(MailMandrillService).(mail.Sender)
}

func (d *diContainer) Google() *social.Google {
	return GetServiceRequired(GoogleService).(*social.Google)
}

func (d *diContainer) Uploader() uploader.Uploader {
	return GetServiceRequired(UploaderService).(uploader.Uploader)
}

func (d *diContainer) Crud() *crud.Crud {
	return GetServiceRequired(CrudService).(*crud.Crud)
}

func (d *diContainer) Localize() localize.ILocalizer {
	return GetServiceRequired(LocalizeService).(localize.ILocalizer)
}

func (d *diContainer) FileExtractor() *fileextractor.FileExtractor {
	return GetServiceRequired(ExtractorService).(*fileextractor.FileExtractor)
}

func (d *diContainer) UUID() uuid.IUUID {
	return GetServiceRequired(UUIDService).(uuid.IUUID)
}

func (d *diContainer) OTP() otp.IOTP {
	return GetServiceRequired(OTPService).(otp.IOTP)
}

func (d *diContainer) DDOS() ddos.IDDOS {
	return GetServiceRequired(DDOSService).(ddos.IDDOS)
}

func (d *diContainer) DynamicLink() dynamiclink.IGenerator {
	return GetServiceRequired(DynamicLinkService).(dynamiclink.IGenerator)
}

func (d *diContainer) FCM() fcm.FCM {
	return GetServiceRequired(FCMService).(fcm.FCM)
}

func (d *diContainer) PDF() pdf.ServiceInterface {
	return GetServiceRequired(PDFService).(pdf.ServiceInterface)
}
