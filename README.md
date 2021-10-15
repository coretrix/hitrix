![Check & test](https://github.com/coretrix/hitrix/workflows/Check%20&%20test/badge.svg)
[![codecov](https://codecov.io/gh/coretrix/hitrix/branch/main/graph/badge.svg)](https://codecov.io/gh/coretrix/hitrix)
[![Go Report Card](https://goreportcard.com/badge/github.com/coretrix/hitrix)](https://goreportcard.com/report/github.com/coretrix/hitrix)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)



# Hitrix

Hitrix is a web framework written in Go (Golang) and support Graphql and REST api.
Hitrix is based on top of [Gqlgen](https://gqlgen.com/]) and [Gin Framework](https://github.com/gin-gonic/gin) and it's high performance and easy to use

### Built-in features:

 * It supports all features of [Gqlgen](https://gqlgen.com/]) and [Gin Framework](https://github.com/gin-gonic/gin)
 * Integrated with [ORM](https://github.com/latolukasz/beeorm)
 * Follows [Dependency injection](https://en.wikipedia.org/wiki/Dependency_injection) pattern
 * Provides many DI services that makes your live easier. You can read more about them [here](https://github.com/coretrix/hitrix#built-in-services)
 * Provides [Dev panel](https://github.com/coretrix/dev-frontend) where you can monitor and manage your application(monitoring, error log, db alters redis status and so on)
 * Other Features
    * [Database Seeding](#seeding)

## Installation

```
go get -u github.com/coretrix/hitrix
``` 
 
 
## Quick start
1. Run next command into your project's main folder and the graph structure will be created
```
go run github.com/99designs/gqlgen init
```


2. Create `cmd` folder into your project and file called `main.go`

Put the next code into the file:
```go
package main

import (
	"github.com/coretrix/hitrix"
	"github.com/gin-gonic/gin"
	
	"your-project/graph" //path you your graph
	"your-project/graph/generated" //path you your graph generated folder
)

func main() {
	s, deferFunc := hitrix.New(
		"app-name", "your secret",
	).Build()
    defer deferFunc()
	s.RunServer(9999, generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}),  func(ginEngine *gin.Engine) {
		//here you can register all your middlewares
	})
}

```

You are able register DI services in your `main.go` file in that way:

```go
package main

import (
	"github.com/coretrix/hitrix"
	"github.com/coretrix/hitrix/service/registry"
	"your-project/entity"
	"your-project/graph"
	"your-project/graph/generated"
	"github.com/coretrix/hitrix/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	s, deferFunc := hitrix.New(
		"app-name", "your secret",
	).RegisterDIGlobalService(
		registry.ServiceProviderErrorLogger(), //register redis error logger
		registry.ServiceProviderConfigDirectory("../config"), //register config service. As param you should point to the folder of your config file
		registry.ServiceProviderOrmRegistry(entity.Init), //register our ORM and pass function where we set some configurations 
		registry.ServiceProviderOrmEngine(), //register our ORM engine for background processes
		registry.ServiceProviderJWT(), //register JWT DI service
		registry.ServiceProviderPassword(), //register pasword DI service
	).RegisterDIRequestService(
		registry.ServiceProviderOrmEngineForContext(), //register our ORM engine per context used in foreground processes 
	).Build()
    defer deferFunc()

	s.RunServer(9999, generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}),  func(ginEngine *gin.Engine) {
		middleware.Cors(ginEngine)
	})
}

```
Now I will explain the main.go file line by line
 1. We create **New** instance of Hitrix and pass app name and a secret that is used from our security services 
 2. We register some DI services
    
    2.1. Global DI service for error logger. It will be used for error handler as well in case of panic
        If you register SlackApi error logger also it will send messages to slack channel
    
    2.2. Global DI service that loads config file
    
    2.3. Global DI service that initialize our ORM registry
    
    2.4. Global DI ORM engine used in background processes
    
    2.5. Request DI ORM engine used in foreground processes
    
    2.6. Global DI JWT service used by dev panel
    
    2.7. Global DI Password service used by dev-panel
 4. We run the server on port `9999`, pass graphql resolver and as third param we pass all middlewares we need.   
As you can see in our example we register only Cors middleware
    

### Register [Dev Panel](https://github.com/coretrix/dev-frontend)
If you want to use our dev panel and to be able to manage alters, error log, redis monitoring, redis stream and so  on you should execute next steps:

#### Create DevPanelUserEntity
```go
package entity

import (
	"github.com/latolukasz/beeorm"
)

type DevPanelUserEntity struct {
	beeorm.ORM   `orm:"table=admin_users;redisCache"`
	ID        uint64
	Email     string `orm:"unique=Email"`
	Password  string

	UserEmailIndex *beeorm.CachedQuery `queryOne:":Email = ?"`
}

func (e *DevPanelUserEntity) GetUsername() string {
	return e.Email
}

func (e *DevPanelUserEntity) GetPassword() string {
	return e.Password
}

```

After that you should register it to the `entity.Init` function
```go
package entity

import "github.com/latolukasz/beeorm"

func Init(registry *beeorm.Registry) {
	registry.RegisterEntity(
		&DevPanelUserEntity{},
	)
}

```

Please execute this alter into your database
```sql

create table dev_panel_users
(
    ID       bigint unsigned auto_increment primary key,
    Email    varchar(255) null,
    Password varchar(255) null,
    constraint Email unique (Email)
) charset = utf8mb4;


```

After that you can make GET request to http://localhost:9999/dev/create-dev-panel-user/?username=contact@coretrix.com&password=coretrix
This will generate sql query that should be executed into your database to create new user for dev panel

#### Register dev panel when you make new instance of hitrix framework in your `main.go` file
```go
s, deferFunc := hitrix.New(
		"app-name", "your secret",
	).RegisterDIGlobalService(
		registry.ServiceProviderErrorLogger(), //register redis error logger
		//...
	).
    RegisterDevPanel(&entity.DevPanelUserEntity{}, middleware.Router, nil). //register our dev-panel and pass the entity where we save admin users, the router and the third param is used for the redis stream pool if its used
    Build()
```



### Defining DI services
We have two types of DI services - Global and Request services
Global services are singletons created once for the whole application
Request services are singletons created once per request

### Calling DI services
If you want to access the registered DI services you can do in in that way:
```go
service.DI().App() //access the app
service.DI().Config() //access config
service.DI().OrmEngine() //access global orm engine
service.DI().OrmEngineForContext() //access reqeust orm engine
service.DI().JWT() //access JWT
service.DI().Password() //access JWT
//...and so on
```

#### Register new DI global service
```go

func ServiceProviderMyService() *ServiceProvider {
	return &ServiceProvider{
		Name:   "my_service",
		Build: func(ctn di.Container) (interface{}, error) {
			return &yourService{}, nil
		},
	}
}

```

And you have to register `ServiceProviderMyService()` in your `main.go` file


Now you can access this service in your code using:

```go
import (
    "github.com/coretrix/hitrix"
)

func SomeResolver(ctx context.Context) {

    service.HasService("my_service") // return true
    
    // return error if Build function returned error
    myService, has, err := service.GetServiceSafe("my_service") 
    // will panic if Build function returns error
    myService, has := service.GetServiceOptional("my_service") 
    // will panic if service is not registered or Build function returned errors
    myService := service.GetServiceRequired("my_service") 

    // if you registered service with field "Global" set to false (request service)

    myContextService, has, err := hitrix.GetServiceForRequestSafe(ctx).Get("my_service_request")
    myContextService, has := hitrix.GetServiceForRequestOptional(ctx).Get("my_service_request") 
    myContextService := hitrix.GetServiceForRequestRequired(ctx).Get("my_service_request") 
}

```

It's a good practice to define one object to return all available services:

```go
package my_package
import (
    "github.com/coretrix/hitrix"
)



func MyService() MyService {
    return service.GetServiceRequired("service_key").(*MyService)
}


```

### Setting mode

#### APP_MODE environment variable
You can define hitrix mode using special environment variable "**APP_MODE**".

Hitrix provides by default four modes:

 * **hitrix.ModeLocal - local**
   * should be used on local development machine (developer laptop)
   * errors and stack trace is printed directly to system console
   * log level is set to Debug level
   * log is formatted using human friendly console text formatter
   * Gin Framework is running in GinDebug mode
 * **hitrix.ModeTest - test**
   * should be used when you run your application tests
 * **hitrix.ModeDev - dev**
   * should be used on your dev server
 * **hitrix.ModeDemo - demo**
   * should be used on your demo server
 * **hitrix.ModeProd - prod**
   * errors and stack trace is printed only using Log
   * log level is set to Warn level
   * log is formatted using json formatter   
    
Mode is just a string. You can define any name you want. Remember that every mode that you create
follows **hitrix.ModeProd** rules explained above.
    
    
In code you can easly check current mode using one of these methods:    

```go
service.DI().App().Mode()
service.DI().App().IsInLocalMode()
service.DI().App().IsInProdMode()
service.DI().App().IsInMode("my_mode")
```

#### APP_CONFIG_FOLDER environment variable
There are another important environment variable called `APP_CONFIG_FOLDER`
You can set path to your config folder for your demo, prod or any other environment

#### Environment variables in config file
Its good practice to keep your secrets like database credentials and so on out of the repository.
Our advice is to keep them like environment variables and call them into config.yaml file
For example your config can looks like this:
```yaml
orm:
  default:
    mysql: ENV[DEFAULT_MYSQL]
    redis: ENV[DEFAULT_REDIS]
    locker: default
    local_cache: 1000
```
where `DEFAULT_MYSQL` and `DEFAULT_REDIS` are env variables and our framework will automatically replace `ENV[DEFAULT_MYSQL]` and `ENV[DEFAULT_REDIS]` with the right values

If you want to define array of values you should split them by `;` and they will be presented into the yaml file in that way:
```yaml
cors:
    - test1
    - test2
```

If you want to enable the debug for orm you can add this tag `orm_debug: true` on the main level of your config

Also we check if there is .env.XXX file in main config folder where XXX is the value of the APP_MODE.
If there is for example .env.local we are reading those env variables and merge them with config.yaml how we presented above

### Running scripts

First You need to define script definition that implements hitrix.Script interface:

```go

type TestScript struct {}

func (script *TestScript) Code() string {
    return "test-script"
}

func (script *TestScript) Unique() bool {
    // if true you can't run more than one script at the same time
    return false
}

func (script *TestScript) Description() string {
    return "script description"
}

func (script *TestScript) Run(ctx context.Context, exit hitrix.Exit) {
    // put logic here
	if shouldExitWithCode2 {
        exit.Error()	// you can exit script and specify exit code
    }
}

```

Methods above are required. Optionally you can also implement these interfaces:

```go

// hitrix.ScriptInfinity interface
func (script *TestScript) Infinity() bool {
    // run script and use blocking operation in cases you run all your code in goroutines
    return true
}

// hitrix.ScriptInterval interface
func (script *TestScript) Interval() time.Duration {                                                    
    // run script every minute
    return time.Minute 
}

// hitrix.ScriptIntervalOptional interface
func (script *TestScript) IntervalActive() bool {                                                    
    // only run first day of month
    return time.Now().Day() == 1
}

// hitrix.ScriptIntermediate interface
func (script *TestScript) IsIntermediate() bool {                                                    
    // script is intermediate, for example is listening for data in chain
    return true
}

// hitrix.ScriptOptional interface
func (script *TestScript) Active() bool {                                                    
    // this script is visible only in local mode
    return DIC().App().IsInLocalMode()
}

```

Once you defined script you can run it using RunScript method:

```go
package main
import "github.com/coretrix/hitrix"

func main() {
	h := hitrix.New("app_name", "your secret").Build()
	h.RunBackgroundProcess(func(b *hitrix.BackgroundProcessor) {
		b.RunScript(&TestScript)
	})
}
``` 


You can also register script as dynamic script and run it using program flag:

```go
package main
import "github.com/coretrix/hitrix"

func main() {
	
    hitrix.New("app_name", "your secret").RegisterDIService(
        &registry.ServiceProvider{
            Name:   "my-script",
            
            Script: true, // you need to set true here
            Build: func(ctn di.Container) (interface{}, error) {
                return &TestScript{}, nil
            },
        },
    ).Build()
}
``` 

You can see all available script by using special flag **-list-scripts**:

```shell script
./app -list-scripts
```

To run script:

```shell script
./app -run-script my-script
```

### Built-in services

#### App 
This service contains information about the application like MODE and so on

#### Config
This service provides you access to your config file. We support only YAML file
When you register the service `registry.ServiceProviderConfigDirectory("../config")`
you should provide the folder where are your config files
The folder structure should looks like that
```
config
 - app-name
    - config.yaml
 - hitrix.yaml #optional config where you can define some settings related to built-in services like slack service
```

#### ORM Engine  
Used to access ORM in background scripts. It is one instance for the whole script

You can register it in that way:
`registry.ServiceProviderOrmEngine()`

#### ORM Engine Context
Used to access ORM in foreground scripts like API. It is one instance per every request

You can register it in that way:
`registry.ServiceProviderOrmEngineForContext()`

#### Error Logger
Used to save unhandled errors in error log. It can be used to save custom errors as well.
If you have setup Slack service you also gonna receive notifications in your slack

You can register it in that way:
`registry.ServiceProviderErrorLogger()`

#### SlackAPI
Gives you ability to send slack messages using slack bot. Also it's used to send messages if you use our ErrorLogger service.
The config that needs to be set in hitrix.yaml is:

```yaml
slack:
    token: "your token"
    error_channel: "test" #optional, used by ErrorLogger
    dev_panel_url: "test" #optional, used by ErrorLogger

```

You can register it in that way:
`registry.ServiceProviderSlackAPI()`

#### JWT
You can use that service to encode and decode JWT tokens

You can register it in that way:
`registry.ServiceProviderJWT()`

#### Password
This service it can be used to hash and verify hashed passwords. It's use the secret provided when you make new Hitrix instance

You can register it in that way:
`registry.ServiceProviderPassword()`

#### Amazon S3

This service is used for storing files into amazon s3

You can register amazon s3 service this way:

```go
registry.ServiceProviderAmazonS3(map[string]uint64{"products": 1}) // 1 is the bucket ID for database counter
```
and you should register the entity `S3BucketCounterEntity` into the ORM
Also, you should put your credentials and other configs in `config/hitrix.yml`

```yml
amazon_s3:
  endpoint: "https://somestorage.com" # set to "" if you're using https://s3.amazonaws.com
  access_key_id: ENV[S3_ACCESS_KEY_ID]
  secret_access_key: ENV[S3_SECRET_ACCESS_KEY_ID]
  disable_ssl: false
  region: us-east-1
  url_prefix: prefix
  domain: domain.com
  buckets: # Register your buckets here for each app mode
    products: # bucket name
      prod: bucket-name
      local: bucket-name-local
  public_urls: # Register your public urls for the GetObjectCachedURL method
    product: # bucket name
      prod: "https://somesite.com/{{.StorageKey}}/" # Available variables are: .Environment, .BucketName, .CounterID, and, .StorageKey
      local: "http://127.0.0.1/{{.Environment}}/{{.BucketName}}/{{.StorageKey}}/{{.CounterID}}" # Will output "http://127.0.0.1/local/product/1.jpeg/1"
```

#### Uploader

This service uses TUS protocol to enable fast resumable and multi-part upload of big files.
It provides an easy interface for plug-in whatever data store and locker you want to implement.
Currently, Amazon S3 data store and Redis locker are implemented. For Amazon data store to work, 
you need to register Amazon S3 service before this one, also for Redis locker to work, you need
to register orm service background before this one.

You can register Uploader service this way:

```go
registry.ServiceProviderUploader(tusd.Config{...}, datastore.GetAmazonS3Store, locker.GetRedisLocker)
```

Hitrix also provides REST uploader controller which you can register all handler methods in your
router:

```go
var uploaderController *hitrixController.UploaderController
uploaderGroup := ginEngine.Group("/files/")
uploaderGroup.Use(middleware.AuthorizeWithHeaderStrict())
{
	uploaderGroup.POST("", uploaderController.PostFileAction)
	uploaderGroup.HEAD(":id", uploaderController.HeadFile)
	uploaderGroup.PATCH(":id", uploaderController.PatchFile)
	uploaderGroup.GET(":id", uploaderController.GetFileAction)
	uploaderGroup.DELETE(":id", uploaderController.DeleteFile)
}
```

Also you need bucket name in config:

````yml
uploader:
  bucket: media
````

#### Stripe

Stripe payment integration 

You can register Stripe service this way:

```go
registry.ServiceProviderStripe(),
```

Config sample:

```yml
stripe:
  key: "api_key"
  webhook_secrets: # map of your webhook secrets
    checkout: "key"
```

#### Dynamic link service
This service is used for generating dynamic links, at this moment only Firebase is supported

You can register Dynamic link service this way:

```go
registry.ServiceProviderDynamicLink(),
```

Config sample:

```yml
 api_key: string # required
 dynamic_link_info: # required
   domain_uri_prefix: string # required
   link: string # required
   android_info: # optional
     package_name: string # optional
     fallback_link: string # optional
     min_package_version_code: string # optional
   ios_info: # optional
     bundle_id: string # optional
     fallback_link: string # optional
     custom_scheme: string # optional
     ipad_fallback_link: string # optional
     ipad_bundle_id: string # optional
     app_store_id: string # optional
   navigation_info: # optional
     enable_forced_redirect: boolean # required
   analytics_info: # optional
     google_play_analytics: # optional
       utm_source: string # optional
       utm_medium: string # optional
       utm_campaign: string # optional
       utm_term: string # optional
       utm_content: string # optional
       gcl_id: string # optional
     itunes_connect_analytics: # optional
       at: string # optional
       ct: string # optional
       mt: string # optional
       pt: string # optional
   social_meta_tag_info: # optional
     social_title: string # optional
     social_description: string # optional
     social_image_link: string # optional
 suffix: # optional
   option: string # required, values: "SHORT" or "UNGUESSABLE"
```


#### Firebase cloud messaging (FCM) service
This service is used for sending different types of push notifications

You can register FCM service this way:

```go
registry.ServiceProviderFCM(),
```

Config sample:

expose FIREBASE_CONFIG="path/to/service-account-file.json"

#### OSS Google
This service is used for storage files into google storage

You can register it in that way:
`registry.OSSGoogle(map[string]uint64{"my-bucket-name": 1})`

and you should register the entity `OSSBucketCounterEntity` into the ORM
You should pass parameter as a map that contains all buckets you need as a key and as a value you should pass id. This id should be unique

In your config folder you should put the .oss.json config file that you have from google
Your config file should looks like that:
```json
{
  "type": "...",
  "project_id": "...",
  "private_key_id": "...",
  "private_key": "...",
  "client_email": "...",
  "client_id": "...",
  "auth_uri": "...",
  "token_uri": "...",
  "auth_provider_x509_cert_url": "...",
  "client_x509_cert_url": "..."
}
```

The last thing you need to set in domain that gonna be used for the static files.
You can setup the domain in hitrix.yaml config file like this:
```yaml
oss: 
  domain: myapp.com
```

and the url to access your static files will looks like
`https://static-%s.myapp.com/%s/%s`
where first %s is app mode

second %s is bucket name concatenated with app mode 

and last %s is the id of the file

#### DDOS Protection
This service contains DDOS protection features

You can register it in that way:
`registry.ServiceProviderDDOS()`

You can protect for example login endpoint from many attempts  by using method `ProtectManyAttempts`

#### PDF service
PDF service provides a generating pdf function from html code using Chrome headless.

First you need these in your app config:
```yaml
chrome_headless:
  web_socket_url: ENV[CHROME_HEADLESS_WEB_SOCKET_URL]
```
Register the PDF service:

```go
registry.ServiceProviderPDF()
```

Access the registered DI service:
```go
pdfService := service.DI().PDFService()
```
Using `HtmlToPdf()` function to generate PDF from html:
```go
pdfBytes := pdfService.HtmlToPdf("<html><p>Hi!</p></html>")
```

Recommended docker file for Chrome headless:
```
https://hub.docker.com/r/chromedp/headless-shell/
```
#### Localizer service
Localizer provides you a simple translation service that can pull and push translation pairs from local (file) and external sources (online services).

Currently localizer supports only [POEditor](https://poeditor.com) online source.

Localizer using a bucket key to separate and manage translation parts of your app.

First you need these in your app config:
```yaml
translation:
  poeditor:
    api_key: ENV[POEDITOR_API_KEY]
    project_id: ENV[POEDITOR_PROJECT_ID]
    language: ENV[POEDITOR_LANGUAGE]
```

Register the localizer service:

```go
registry.ServiceProviderLocalizer()
```

Access the registered DI service:
```go
localizerService := service.DI().LocalizerService()
```

Loading translation pairs from map:
```go
bucketKey := "greet-service"
append := false // append or replace?
pairs := map[string]string{
  "app_name": "My App Name",
  "loading_text": "Loading ...",
}
localizerService.LoadBucketFromMap(
  bucketKey, 
  pairs, 
  append,
)
```
Using `Localize()` function to translate a key:
```go
appName, err := localizerService.Localize(bucketKey, "app_name")
if err !nil {
  // handle error
}
```
Loading translation pairs from local file:
```go
localizerService.LoadBucketFromFile(
  bucketKey,
  "locales/greet.en.json",
  append,
)
```
Pull the translations from external source:
```go
err := localizerService.PullBucketFromSource(bucketKey, append)
if err != nil {
  log.Fatal(err)
}
```
Push translations to external source:
```go
err := localizerService.PushBucketToSource(bucketKey)
if err != nil {
  // handle error
}
```

#### File extractor service
File extractor provides you a simple function to search in a path recursively and find terms based on a regular expression. 

Register the localizer service:
```go
registry.ServiceProviderExtractor(),
```
Access the registered DI service:

```go
extractService := service.DI().FileExtractorService()
```
Extract phrase (errors in this example):
```go
errorTerms, err := extractService.Extract(fileextractor.ExtractParams{
  SearchPath: "./",
  Excludes:   []string{},
  Expression: `errors.New[(]*\("([^)]*)"\)`,
})
if err != nil {
  // handle error
}
```

#### API logger service
This service us used to track every api request and response.
You can register it in that way:
`registry.APILogger(&entity.APILogEntity{}),`

The methods that this service provide are:
```go
type APILogger interface {
	LogStart(logType string, request interface{})
	LogError(message string, response interface{})
	LogSuccess(response interface{})
}
```
You should call `LogStart` before you send request to the api

You should call `LogError` in case api return you error

You should call `LogSuccess` in case api return you success

#### WebSocket
This service add support of websockets. It manage the connections and provide you easy way to read and write messages

You can register it in that way:
`registry.ServiceSocketRegistry(registerHandler, unregisterHandler func(s *socket.Socket))`

To be able to handle new connections you should create your own route and create a handler for it.
Your handler should looks like that:
```go
type WebsocketController struct {
}

func (controller *WebsocketController) InitConnection(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	socketRegistryService, has := service.DI().SocketRegistry()
	if !has {
		panic("Socket Registry is not registered")
	}

	errorLoggerService, has := service.DI().ErrorLogger()
	if !has {
		panic("Socket Registry is not registered")
	}

	connection := &socket.Connection{Send: make(chan []byte, 256), Ws: ws}
	socketHolder := &socket.Socket{
		ErrorLogger: errorLoggerService,
		Connection:  connection,
		ID:          "unique connection hash based on userID, deviceID and timestamp",
		Namespace:   model.DefaultNamespace,
	}

	socketRegistryService.Register <- socketHolder

	go socketHolder.WritePump()
	go socketHolder.ReadPump(socketRegistryService, func(rawData []byte) {
		s, _ := socketRegistryService.Sockets.Load(socketHolder.ID)
		
        dto := &DTOMessage{}
        err = json.Unmarshal(rawData, dto)
        if err != nil {
            errorLoggerService.LogError(err)
            retrun
        }
        //handle business logic here
        s.(*socket.Socket).Emit(dto)
	})
}

```
This handler initializes the new incoming connections and have 2 goroutines - one for writing messages, and the second one for reading messages
If you want to send a message you should use ```socketRegistryService.Emit```

If you want to read incoming messages you should do it in the function we are passing as second parameter of ```ReadPump``` method

If you want to select certain connection you can do it by the ID and this method ```s, err := socketRegistryService.Sockets.Load(ID)```

Also websocket service provide you hooks for registering new connections and for unregistering already existing connections.
You can define those handlers when you register the service based on namespace of socket.

#### Clock service
This service is used for `time` operations. It is better to use it everywhere instead of `time.Now()` because it can be mocked and you can set whatever time you want in your tests

You can register it in that way:
`registry.ServiceClock(),`

The methods that this service provide are:
```Now() and NowPointer()```

#### Mail Mandrill service

This service is used for sending transactional emails using Mandrill

You can register the service this way:

```go
registry.MailMandrill()
```
and you should register the entity `MailTrackerEntity` into the ORM
Also, you should put your credentials and other configs in your config file

```yml
mandrill:
  api_key: ...
  default_from_email: test@coretrix.tv
  from_namme: coretrix.com
```

Some of the functions this service provide are:
```go
	SendTemplate(ormService *beeorm.Engine, message *TemplateMessage) error
	SendTemplateAsync(ormService *beeorm.Engine, message *TemplateMessage) error
	SendTemplateWithAttachments(ormService *beeorm.Engine, message *TemplateAttachmentMessage) error
	SendTemplateWithAttachmentsAsync(ormService *beeorm.Engine, message *TemplateAttachmentMessage) error
```

#### Authentication Service
This service is used to making the life easy by doing the whole authentication life cycle using JWT token. the methods that this service provides are as follows:

##### dependencies : 
`JWTService`

`PasswordService`

`ClockService`

`GeneratorService`

`SMSService` # optional , when you need to support for otp

`GoogleService` # optional , when you need to support google login

`FacebookService` # optional , when you need to support facebook login

```go
func Authenticate(ormService *beeorm.Engine, uniqueValue string, password string, entity AuthProviderEntity) (accessToken string, refreshToken string, err error) {}
func VerifyAccessToken(ormService *beeorm.Engine, accessToken string, entity beeorm.Entity) error {}
func VerifySocialLogin(source, token string)
func RefreshToken(ormService *beeorm.Engine, refreshToken string) (newAccessToken string, newRefreshToken string, err error) {}
func LogoutCurrentSession(ormService *beeorm.Engine, accessKey string){}
func LogoutAllSessions(ormService *beeorm.Engine, id uint64)
func GenerateAndSendOTP(ormService *beeorm.Engine, mobile string, country string){}
func VerifyOTP(ormService *beeorm.Engine, code string, input *GenerateOTP) error{}
func AuthenticateOTP(ormService *beeorm.Engine, phone string, entity OTPProviderEntity) (accessToken string, refreshToken string, err error){}
```
1. The `Authenticate` function will take an uniqueValue such as Email or Mobile, a plain password, and generates accessToken and refreshToken. 
   You will also need to pass your entity as third argument, and it will give you the specific user entity related to provided access token
   The entity should implement the `AuthProviderEntity` interface : 
    ```go
       type AuthProviderEntity interface {
        beeorm.Entity
        GetUniqueFieldName() string
        GetPassword() string
       }
    ```
    The example of such entity is as follows:
    ```go
    type UserEntity struct {
	    beeorm.ORM  `orm:"table=users;redisCache;redisSearch=search_pool"`
	    ID       uint64 `orm:"searchable;sortable"`
	    Email    string `orm:"required;unique=Email;searchable"`
	    Password string `orm:"required"`
    }
   
    func (user *UserEntity) GetUniqueFieldName() string {
	    return "Email"
    }
    
    func (user *UserEntity) GetPassword() string {
    return user.Password
    }
    ```
2. The `VerifyAccessToken` will get the AccessToken, process the validation and expiration, and fill the entity param with the authenticated user entity in case of successful authentication.
3. The `RefreshToken` method will generate a new token pair for given user
4. The `LogoutCurrentSession` you can logout the user current session , you need to pass it the `accessKey`  that is the jwt identifier `jti` the exists in both access and refresh token.
5. The `LogoutAllSessions` you can logout the user from all sessions , you need to pass it the `id` (user id).
6. The `GenerateAndSendOTP` only in otp flow, it will generate code and send it to the specified number `Mobile` and also returns `GenerateOTP` inside it we have `Token` that it is the hashed otp credentials that needs to be sent by client when verifying.
    ```go
    type GenerateOTP struct {
    	Mobile         string
    	ExpirationTime time.Time
    	Token          string
    }
    ```
7. The `VerifyOTP` only in otp flow , will compare the `code`(otp code) with the `input`(otp credentials)  provided by client.
8. The `AuthenticateOTP` only in otp flow , will get the `phone` and `entity` that should implement `OTPProviderEntity` and query to find the user and will login the user.careful just call this after you verified the otp code using the previous method `VerifyOTP`
   the response is asa same as the `Authenticate`.
   ```go
   type OTPProviderEntity interface {
	    beeorm.Entity
	    GetPhoneFieldName() string
    }
   ```
9. You need to have a `authentication` key in your config file for this service to work. `secret` key under `authentication` is mandatory but other options are optional:
10. The service can also support `OTP` if you want your service to support otp you should have `support_otp` key set to true under `authentication`
11. The service also needs redis to store its sessions so you need to identify the redis storage name in config , the key is `auth_redis` under `authentication`
```yaml
authentication:
  secret: "a-deep-dark-secret" #mandatory, secret to be used for JWT
  access_token_ttl: 86400 # optional, in seconds, default to 1day
  refresh_token_ttl: 31536000 #optional, in seconds, default to 1year
  auth_redis: default #optional , default is the default redis
  support_otp: true # if you want to support otp flow in your app
  otp_ttl: 120 #optional ,set it when you want to use otp, It is the ttl of otp code , default is 60 seconds
```

#### SMS Service
This service is capable of sending simple message and otp message and also calling by different sms providers .
for now we support 3 sms providers : `twilio` `sinch` `kavenegar`

##### dependencies : 
`ClockService`

and also when registering the service you need to pass it the `LogEntity` that is responsible to log every action made by sms service : 
```go
type LogEntity interface {
    beeorm.Entity
    SetStatus(string)
    SetTo(string)
    SetText(string)
    SetFromPrimaryGateway(string)
    SetFromSecondaryGateway(string)
    SetPrimaryGatewayError(string)
    SetSecondaryGatewayError(string)
    SetType(string)
    SetSentAt(time time.Time)
}
```
for example : 
```go
const (
	SMSTrackerTypeSMS     = "sms"
	SMSTrackerTypeCallout = "callout"
)

type smsTrackerTypeAll struct {
	SMSTrackerTypeSMS     string
	SMSTrackerTypeCallout string
}

var SMSTrackerTypeAll = smsTrackerTypeAll{
	SMSTrackerTypeSMS:     SMSTrackerTypeSMS,
	SMSTrackerTypeCallout: SMSTrackerTypeCallout,
}

type SmsTrackerEntity struct {
	beeorm.ORM               `orm:"table=sms_tracker"`
	ID                    uint64
	Status                string
	To                    string `orm:"varchar=15"`
	Text                  string
	FromPrimaryGateway    string
	FromSecondaryGateway  string
	PrimaryGatewayError   string
	SecondaryGatewayError string
	Type                  string    `orm:"enum=entity.SMSTrackerTypeAll;required"`
	SentAt                time.Time `orm:"time"`
}
```
we have 2 providers active at the same time `primary` `secondary` and when send via primary fails we try to send with the secondary provider.
```go
func SendOTPSMS(*OTP) error{}
func SendOTPCallout(*OTP) error{}
func SendMessage(*Message) error{}
```
1. The `SendOTPSMS` send otp sms by providing the otp data
```go
type OTP struct {
	OTP      string
	Number   string
	CC       string
	Provider *Provider
	Template string
}
```
2. The `SendOTPCallout` used to call and tell the otp code
3. The `SendMessage` used to send simple message
```go
type Message struct {
	Text     string
	Number   string
	Provider *Provider
}
```
##### configs
```yaml
sms:
  twilio:
    sid: ENV[SMS_TWILIO_SID]
    token: ENV[SMS_TWILIO_TOKEN]
    from_number: ENV[SMS_TWILIO_FROM_NUMBER]
    authy_url: ENV[SMS_TWILIO_AUTHY_URL]
    authy_api_key: ENV[SMS_TWILIO_AUTHY_API_KEY]
    verify_url: ENV[SMS_TWILIO_VERIFY_URL]
    verify_sid: ENV[SMS_TWILIO_VERIFY_SID]
  kavenegar:
    api_key: ENV[SMS_KAVENEGAR_API_KEY]
    sender: ENV[SMS_KAVENEGAR_SENDER]
  sinch:
    app_id: ENV[SMS_SINCH_APP_ID]
    app_secret: ENV[SMS_SINCH_APP_SECRET]
    msg_url: ENV[SMS_SINCH_MSG_URL]
    from_number: ENV[SMS_SINCH_FROM_NUMBER]
    call_url: ENV[SMS_SINCH_CALL_URL]
    caller_number: ENV[SMS_SINCH_CALLER_NUMBER]
```

#### CRUD

You can register CRUD service this way:

```go
registry.Crud(),
```
This service it gives you ability to build a query and apply different query parameters to the query that should be used in listing pages

### Validator
We support 2 types of validators. One of them is related to graphql, the other one is related to rest.

#### Graphql validator
There are 2 steps that needs to be executed if you want to use this kind of validator

1. Add `directive @validate(rules: String!) on INPUT_FIELD_DEFINITION` into your `schema.graphqls` file

2. Call `ValidateDirective` into your main.go file
```go
config := generated.Config{Resolvers: &graph.Resolver{}, Directives: generated.DirectiveRoot{Validate: hitrix.ValidateDirective()} }

s.RunServer(4001, generated.NewExecutableSchema(config), func(ginEngine *gin.Engine) {
    commonMiddleware.Cors(ginEngine)
    middleware.Router(ginEngine)
})
```

After that you can define the validation rules in that way:
```graphql
input ApplePurchaseRequest {
  ForceEmail: Boolean!
  Name: String
  Email: String @validate(rules: "email") #for rules param you can use everything supported by https://github.com/go-playground/validator validate.Var(value, rules)
  AppleReceipt: String!
}
```

To handle the errors you need to call function `hitrix.Validate(ctx, nil)` in your resolver
```go
func (r *mutationResolver) RegisterTransactions(ctx context.Context, applePurchaseRequest model.ApplePurchaseRequest) (*model.RegisterTransactionsResponse, error) {
    if !hitrix.Validate(ctx, nil) {
        return nil, nil
    }
    // your logic here...
}
```

The function `hitrix.Validate(ctx, nil)` as second param accept callback where you can define your custom validation related to business logic

### Pre deploy
If you run your binary with argument `-pre-deploy` the program will check for alters and if there is no alters it will exit with code 0 but if there is an alters it will exit with code 1.

### Force alters
If you run your binary with argument `-force-alters` the program will check for DB and RediSearch alters and it will execute them(only in local mode).

You can use this feature during the deployment process check if you need to execute the alters before you deploy it

### Pagination
You can use:
```go
package helper

type URLQueryPager struct {
	// example = ?current_page=1&page_size=25
	CurrentPage int `binding:"min=1" form:"current_page"`
	PageSize    int `binding:"min=1" form:"page_size"`
}
```
in your code that needs pagination like:

```go
package mypackage

import "github.com/coretrix/hitrix/pkg/helper"

type SomeURLQuery struct {
	helper.URLQueryPager
	OtherField1 string `form:"other_field_1"`
	OtherField2 int `form:"other_field_2"`
}
```

### Tests
Hitrix provide you test helper functions which can be used to make requests to your graphql api

In your code you can create similar function that makes new instance of your app

```go
func createContextMyApp(t *testing.T, projectName string, resolvers graphql.ExecutableSchema) *test.Ctx {
	defaultServices := []*service.Definition{
		registry.ServiceProviderConfigDirectory("../example/config"),
		registry.ServiceProviderOrmRegistry(entity.Init),
		registry.ServiceProviderOrmEngine(),
	}

	return test.CreateContext(t,
		projectName,
		resolvers,
		func(ginEngine *gin.Engine) { middleware.Router(ginEngine) },
		defaultServices,
	)
}

```

After that you can call queries or mutations

```go
func TestProcessApplePurchaseWithEmail(t *testing.T) {
	type queryRegisterTransactions struct {
		RegisterTransactionsResponse *model.RegisterTransactionsResponse `graphql:"RegisterTransactions(applePurchaseRequest: $applePurchaseRequest)"`
	}

	variables := map[string]interface{}{
		"applePurchaseRequest": model.ApplePurchaseRequest{
			ForceEmail:   false,
		},
	}

	fakeMail := &mailMock.Sender{}
	fakeMail.On("SendTemplate", "hymn@abv.bg").Return(nil)

	got := &queryRegisterTransactions{}
	projectName, resolver := tests.GetWebAPIResolver()
	ctx := tests.CreateContextWebAPI(t, projectName, resolver, &tests.IoCMocks{MailService: fakeMail})

	err := ctx.HandleMutation(got, variables)
	assert.Nil(t, err)

	//...
	fakeMail.AssertExpectations(t)
}
```

Hitrix supports `parallel` tests
In case you want to execute parallel tests you need to set
`PARALLEL_TESTS=true` env var in your IDE config and be sure you don't have set `-p 1` in `Go tool arguments` 
In case you want to disable `parallel` tests remove `PARALLEL_TESTS` or set it to `false` and set in `Go tool arguments` value `-p 1`


## Other Features

### Database Seeding
Hitrix supports multi-versioned multiple seeds. 

To Seed your database ([wiki](https://en.wikipedia.org/wiki/Database_seeding)), 
you can use the seeding feature that is implemented as script (`DBSeedScript`). 
This `DBSeedScript` needs to be provided a `Seeds map[string]Seed` field, Where the string key is the identifier of the seed that is used to detect seed versioning.
The Script can be implementd in your app by making a type that satisfies the `Seed` interface: 
```go
type Seed interface {
	Execute(*beeorm.Engine)
	Version() int
}
```
Example:

`users_seed.go`
```go
type UsersSeed struct{}

func (seed *UsersSeed) Version() int {
	return 1
}

func (seed *UsersSeed) Execute(ormService *beeorm.Engine) {
	// TODO insert a new user entity to the db
}
```
And after your server definition (such as the one in [server.go](example/server.go)), 
you could use it to run the seeds just like you would run any script (explained earlier in the [scripts seciton](#running-scripts)) as follows:
```go
// Server Definition
s,_ := hitrix.New()

...

s.RunBackgroundProcess(func(b *hitrix.BackgroundProcessor) {
  // seed database
  go b.RunScript(&hitrixScripts.DBSeedScript{
    Seeds: map[string]hitrixScripts.Seed{
      "users-seed": UsersSeed{},
    },
  })

  // ... other scripts
})
```

This seed will only run in the following cases:
* No seeds has ever been ran (in db table settings, no row with `'key'='seeds'` )
* This seed has never been ran
* The seed IS ran before, but with an older version

