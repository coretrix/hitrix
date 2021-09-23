package goroutine

type IGoroutine interface {
	Goroutine(fn func())
	GoroutineWithRestart(fn func())
}
