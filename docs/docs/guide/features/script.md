# Running scripts
This feature is used to run background scripts(cron jobs).

First you need to define script definition that implements hitrix.Script interface:

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