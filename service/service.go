package service

import (
	"context"

	"github.com/coretrix/clockwork"
	"github.com/latolukasz/beeorm"

	s3 "github.com/coretrix/hitrix/service/component/amazon/storage"
	apilogger "github.com/coretrix/hitrix/service/component/api_logger"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/authentication"
	"github.com/coretrix/hitrix/service/component/checkout"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/crud"
	"github.com/coretrix/hitrix/service/component/ddos"
	dynamiclink "github.com/coretrix/hitrix/service/component/dynamic_link"
	"github.com/coretrix/hitrix/service/component/elorus"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	"github.com/coretrix/hitrix/service/component/exporter"
	"github.com/coretrix/hitrix/service/component/fcm"
	featureflag "github.com/coretrix/hitrix/service/component/feature_flag"
	fileextractor "github.com/coretrix/hitrix/service/component/file_extractor"
	"github.com/coretrix/hitrix/service/component/generator"
	"github.com/coretrix/hitrix/service/component/gql"
	"github.com/coretrix/hitrix/service/component/instagram"
	"github.com/coretrix/hitrix/service/component/jwt"
	"github.com/coretrix/hitrix/service/component/localize"
	"github.com/coretrix/hitrix/service/component/mail"
	"github.com/coretrix/hitrix/service/component/oss"
	"github.com/coretrix/hitrix/service/component/otp"
	"github.com/coretrix/hitrix/service/component/password"
	"github.com/coretrix/hitrix/service/component/pdf"
	"github.com/coretrix/hitrix/service/component/setting"
	"github.com/coretrix/hitrix/service/component/slack"
	"github.com/coretrix/hitrix/service/component/sms"
	"github.com/coretrix/hitrix/service/component/social"
	"github.com/coretrix/hitrix/service/component/socket"
	"github.com/coretrix/hitrix/service/component/stripe"
	"github.com/coretrix/hitrix/service/component/template"
	"github.com/coretrix/hitrix/service/component/uploader"
	"github.com/coretrix/hitrix/service/component/uuid"
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
	ClockWorkRequestService = "clockwork_request"
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
	MailService             = "mail"
	GoogleService           = "google"
	FacebookService         = "facebook"
	CrudService             = "crud"
	UUIDService             = "uuid"
	OTPService              = "otp"
	ExporterService         = "exporter"
	SettingService          = "setting"
	FeatureFlagService      = "feature_flag"
	TemplateService         = "template"
	GQLService              = "gql"
	ElorusService           = "elorus"
	InstagramService        = "instagram"
)

type DIContainer struct {
}

var dicInstance = &DIContainer{}

func DI() *DIContainer {
	return dicInstance
}

func (d *DIContainer) Checkout() checkout.ICheckout {
	return GetServiceRequired(CheckoutService).(checkout.ICheckout)
}

func (d *DIContainer) Exporter() exporter.IExporter {
	return GetServiceRequired(ExporterService).(exporter.IExporter)
}

func (d *DIContainer) AmazonS3() s3.Client {
	return GetServiceRequired(AmazonS3Service).(s3.Client)
}

func (d *DIContainer) Stripe() stripe.IStripe {
	return GetServiceRequired(StripeService).(stripe.IStripe)
}

func (d *DIContainer) App() *app.App {
	return GetServiceRequired(AppService).(*app.App)
}

func (d *DIContainer) Config() config.IConfig {
	return GetServiceRequired(ConfigService).(config.IConfig)
}

func (d *DIContainer) OrmConfig() beeorm.ValidatedRegistry {
	return GetServiceRequired(ORMConfigService).(beeorm.ValidatedRegistry)
}

func (d *DIContainer) OrmEngine() *beeorm.Engine {
	return GetServiceRequired(ORMEngineGlobalService).(*beeorm.Engine)
}

func (d *DIContainer) OrmEngineForContext(ctx context.Context) *beeorm.Engine {
	return GetServiceForRequestRequired(ctx, ORMEngineRequestService).(*beeorm.Engine)
}

func (d *DIContainer) ClockWorkForContext(ctx context.Context) *clockwork.Clockwork {
	return GetServiceForRequestRequired(ctx, ClockWorkRequestService).(*clockwork.Clockwork)
}

func (d *DIContainer) JWT() *jwt.JWT {
	return GetServiceRequired(JWTService).(*jwt.JWT)
}

func (d *DIContainer) SMS() sms.ISender {
	return GetServiceRequired(SMSService).(sms.ISender)
}

func (d *DIContainer) Generator() generator.IGenerator {
	return GetServiceRequired(GeneratorService).(generator.IGenerator)
}

func (d *DIContainer) Password() password.IPassword {
	return GetServiceRequired(PasswordService).(password.IPassword)
}

func (d *DIContainer) Slack() slack.Slack {
	return GetServiceRequired(SlackService).(slack.Slack)
}

func (d *DIContainer) ErrorLogger() errorlogger.ErrorLogger {
	return GetServiceRequired(ErrorLoggerService).(errorlogger.ErrorLogger)
}

func (d *DIContainer) OSService() oss.IProvider {
	return GetServiceRequired(OSService).(oss.IProvider)
}

func (d *DIContainer) SocketRegistry() *socket.Registry {
	return GetServiceRequired(SocketRegistryService).(*socket.Registry)
}

func (d *DIContainer) APILogger() apilogger.IAPILogger {
	return GetServiceRequired(APILoggerService).(apilogger.IAPILogger)
}

func (d *DIContainer) Clock() clock.IClock {
	return GetServiceRequired(ClockService).(clock.IClock)
}

func (d *DIContainer) Setting() setting.ServiceSettingInterface {
	return GetServiceRequired(SettingService).(setting.ServiceSettingInterface)
}

func (d *DIContainer) Authentication() *authentication.Authentication {
	return GetServiceRequired(AuthenticationService).(*authentication.Authentication)
}

func (d *DIContainer) Mail() mail.Sender {
	return GetServiceRequired(MailService).(mail.Sender)
}

func (d *DIContainer) Google() *social.Google {
	return GetServiceRequired(GoogleService).(*social.Google)
}

func (d *DIContainer) Uploader() uploader.Uploader {
	return GetServiceRequired(UploaderService).(uploader.Uploader)
}

func (d *DIContainer) Crud() *crud.Crud {
	return GetServiceRequired(CrudService).(*crud.Crud)
}

func (d *DIContainer) Localize() localize.ILocalizer {
	return GetServiceRequired(LocalizeService).(localize.ILocalizer)
}

func (d *DIContainer) FileExtractor() *fileextractor.FileExtractor {
	return GetServiceRequired(ExtractorService).(*fileextractor.FileExtractor)
}

func (d *DIContainer) UUID() uuid.IUUID {
	return GetServiceRequired(UUIDService).(uuid.IUUID)
}

func (d *DIContainer) OTP() otp.IOTP {
	return GetServiceRequired(OTPService).(otp.IOTP)
}

func (d *DIContainer) DDOS() ddos.IDDOS {
	return GetServiceRequired(DDOSService).(ddos.IDDOS)
}

func (d *DIContainer) DynamicLink() dynamiclink.IGenerator {
	return GetServiceRequired(DynamicLinkService).(dynamiclink.IGenerator)
}

func (d *DIContainer) FCM() fcm.FCM {
	return GetServiceRequired(FCMService).(fcm.FCM)
}

func (d *DIContainer) PDF() pdf.ServiceInterface {
	return GetServiceRequired(PDFService).(pdf.ServiceInterface)
}

func (d *DIContainer) FeatureFlag() featureflag.ServiceFeatureFlagInterface {
	return GetServiceRequired(FeatureFlagService).(featureflag.ServiceFeatureFlagInterface)
}

func (d *DIContainer) Template() template.ITemplateInterface {
	return GetServiceRequired(TemplateService).(template.ITemplateInterface)
}

func (d *DIContainer) Gql() gql.IGQLInterface {
	return GetServiceRequired(GQLService).(gql.IGQLInterface)
}

func (d *DIContainer) Elorus() elorus.IProvider {
	return GetServiceRequired(ElorusService).(elorus.IProvider)
}

func (d *DIContainer) Instagram() instagram.IAPIManager {
	return GetServiceRequired(InstagramService).(instagram.IAPIManager)
}
