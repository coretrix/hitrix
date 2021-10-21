# Goroutine
Go provides simple solution for concurrency using go keyword, but in any case that goroutine panics you have to handle that
yourself. Hitrix provides `Goroutine` function that does this automatically for you and you don't have to be worry about
your goroutine panic. Hitrix Goroutine function would recover the panic for you and log the error.
You can use that like this:
```go
package mypackage

import "github.com/coretrix/hitrix"

func Myfunc(){
  hitrix.Goroutine(someFunc())
}
```

Hitrix also provides another function named `GoroutineWithRestart`. If you have a function that must be run all the time
you can use this function. In any case of panics it would log the error and automatically start the function again.
You can use this function like this:
```go
package mypackage

import "github.com/coretrix/hitrix"

func Myfunc(){
  hitrix.GoroutineWithRestart(someFunc())
}
```

