# Introduction

Hitrix is a web framework written in Go (Golang) and support Graphql and REST api.
It is based on top of [Gqlgen](https://gqlgen.com/]) and [Gin Framework](https://github.com/gin-gonic/gin)

#### Why to choose Hitrix?
Hitrix is combination between high performance and speed of development.
There are many build-in features and tools that save development time.
Also this framework helps you from the day zero to start creating new features for your project.
You don't need to spend time thinking about error log, db layer, caching, DI, structure, background jobs and so on.
The only thing you need to do is to use Hitrix and deliver fast to the business

Built-in features:

* It supports all features of [Gqlgen](https://gqlgen.com/]) and [Gin Framework](https://github.com/gin-gonic/gin)
* Integrated with [ORM](https://github.com/latolukasz/beeorm)
* Follows [Dependency injection](https://en.wikipedia.org/wiki/Dependency_injection) pattern
* Provides many DI services that makes your life easier. You can read more about them in our documentation
* Provides [Dev panel](https://github.com/coretrix/dev-frontend) where you can monitor and manage your application(monitoring, error log, db alters redis status and so on)
* Other Features
    * Database Seeding
    * Helpers
    * Background scripts
    * Validators
    * Integration tests
    * and so on...

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
	"github.com/coretrix/hitrix/service/component/app"
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
		registry.ServiceProviderPassword(password.NewSimpleManager), //register pasword DI service
	).RegisterDIRequestService(
		registry.ServiceProviderOrmEngineForContext(false), //register our ORM engine per context used in foreground processes 
	).RegisterRedisPools(&app.RedisPools{Persistent: "your pool here"}).
    Build()
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
   If you register Slack error logger also it will send messages to slack channel

   2.2. Global DI service that loads config file

   2.3. Global DI service that initialize our ORM registry

   2.4. Global DI ORM engine used in background processes

   2.6. Global DI JWT service used by dev panel

   2.7. Global DI Password service used by dev-panel

   2.8. Request DI ORM engine used in foreground processes
3. We register redis pools. Those pools are used by different services as `authentication` service, `dev panel` and so on
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
    RegisterDevPanel(&entity.DevPanelUserEntity{}, middleware.Router). //register our dev-panel and pass the entity where we save admin users, the router and the third param is used for the redis stream pool if its used
    Build()
```



### DI services
We have two types of DI services - Global and Request services
Global services are singletons created once for the whole application
Request services are singletons created once per request

#### Calling DI services
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

### Environment variables

#### APP_FOLDER environment variable
There are another important environment variable called `environment`
You can set path to your app folder for your demo, prod or any other environment
