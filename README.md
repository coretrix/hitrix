![Check & test](https://github.com/coretrix/hitrix/workflows/Check%20&%20test/badge.svg)
[![codecov](https://codecov.io/gh/coretrix/hitrix/branch/main/graph/badge.svg)](https://codecov.io/gh/coretrix/hitrix)
[![Go Report Card](https://goreportcard.com/badge/github.com/coretrix/hitrix)](https://goreportcard.com/report/github.com/coretrix/hitrix)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)



# hitrix

### Simple Framework designed to build scalable GraphQL services

### Main features:

 * Build on top of [Gqlgen](https://gqlgen.com/]) and [Gin Framework](https://github.com/gin-gonic/gin)
 * Easy to integrate with [hitrix ORM](https://github.com/summer-solutions/orm)
 * Follows [Dependency injection](https://en.wikipedia.org/wiki/Dependency_injection) pattern
 
### Create hitrix instance

```go
package main
import "github.com/coretrix/hitrix"

func main() {
    registry := hitrix.New("app_name").Build()
    //Starting from now you have access to global DI container (DIC)
    container := DIC()
}

``` 
 
 
### Starting GraphQL Server

```go
package main
import "github.com/coretrix/hitrix"

func main() {
	
    graphQLExecutableSchema := ... // setup graphql.ExecutableSchema 
    ginHandler := // setup gin routes and middlewares
    // run http server
    hitrix.New("app_name").Build().RunServer(8080, graphQLExecutableSchema, ginHandler)
}

``` 

#### Setting server port

By default, hitrix server is using port defined in environment variable "**PORT**". If this variable is not
set hitrix will use port number passed as fist argument.

#### Application name

When you setup server using **New** method yo must provide unique application name that can be
checked in code like this:

```go
    hitrix.New("app_name").Build()
    DIC().App().Name()
```

#### Setting mode

By default, hitrix is running in "**hitrix.ModeLocal**" mode. Mode is a string that is available in: 
```go
    hitrix.New("app_name").Build()
    // now you can access current hitrix mode
    DIC().App().Mode()
```

You can define hitrix mode using special environment variable "**hitrix_MODE**".

hitrix provides by default two modes:

 * **hitrix.ModeLocal**
   * should be used on local development machine (developer laptop)
   * errors and stack trace is printed directly to system console
   * log level is set to Debug level
   * log is formatted using human friendly console text formatter
   * Gin Framework is running in GinDebug mode
  * **hitrix.ModeProd**
    * errors and stack trace is printed only using Log
    * log level is set to Warn level
    * log is formatted using json formatter   
    
Mode is just a string. You can define any name you want. Remember that every mode that you create
follows **hitrix.ModeProd** rules explained above.
    
    
In code you can easly check current mode using one of these methods:    

```go
    DIC().App().Mode()
    DIC().App().IsInLocalMode()
    DIC().App().IsInProdMode()
    DIC().App().IsInMode("my_mode")
```

#### Defining DI services

hitrix builds global shared Dependency Injection container. You can register new services using this method:

```go
package main
import "github.com/coretrix/hitrix"

func main() {
    hitrix.New("my_app").RegisterDIService(
      // put service definitions here
    )
}

``` 

Example of DI service definition:

```go
package main
import (
    "github.com/coretrix/hitrix"
)
    
func main() {
    myService := &hitrix.ServiceDefinition{
        Name:   "my_service", // unique service key
        Global: true, // false if this service should be created as separate instance for each http request
        Build: func() (interface{}, error) {
            return &SomeService{}, nil // you can return any data you want
        },
        Close: func(obj interface{}) error { //optional
        },
        Flags: func(registry *hitrix.FlagsRegistry) { //optional
            registry.Bool("my-service-flag", false, "my flag description")
            registry.String("my-other-flag", "default value", "my flag description")
        },
    }
    
    // register it and run server
    hitrix.New("my_app").RegisterDIService(
      myService,
    )

    // you can access flags:
    val := hitrix.DIC().App().Flags().Bool("my-service-flag")
}

```

Now you can access this service in your code using:

```go
import (
    "github.com/coretrix/hitrix"
)

func SomeResolver(ctx context.Context) {

    hitrix.HasService("my_service") // return true
    
    // return error if Build function returned error
    myService, has, err := hitrix.GetServiceSafe("my_service") 
    // will panic if Build function returns error
    myService, has := hitrix.GetServiceOptional("my_service") 
    // will panic if service is not registered or Build function returned errors
    myService := hitrix.GetServiceRequired("my_service") 

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
    "context"
    "github.com/coretrix/hitrix"
)

type dic struct {
}

var dicInstance = &dic{}

type DICInterface interface {
    MyService() *MyService
    MyOtherServiceForContext(ctx context.Context) *MyOtherService
}

func DIC() DICInterface {
    return dicInstance
}

func (d *dic) MyService() MyService {
    return hitrix.GetServiceRequired("service_key").(*MyService)
}

func (d *dic) MyOtherServiceForContext(ctx context.Context) MyOtherService {
    return hitrix.GetServiceForRequestRequired(ctx, "other_service_key").(*MyOtherService)
}

```


#### Running scripts


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
	hitrix.New("app_name").Build().RunScript(&TestScript{})
}
``` 


You can also register script as dynamic script and run it using program flag:

```go
package main
import "github.com/coretrix/hitrix"

func main() {
	
    hitrix.New("app_name").RegisterDIService(
        &hitrix.ServiceDefinition{
            Name:   "my-script",
            Global: true,
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