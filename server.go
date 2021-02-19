package hitrix

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

type Hitrix struct {
	ctx    context.Context
	cancel context.CancelFunc
	done   chan bool
	exit   chan int
}

func (h *Hitrix) RunServer(defaultPort uint, server graphql.ExecutableSchema, ginInitHandler GinInitHandler) {
	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", defaultPort)
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: InitGin(server, ginInitHandler),
	}
	func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
		h.done <- true
	}()
	h.await()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("server forced to shutdown")
	}
}

func (h *Hitrix) await() {
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case code := <-h.exit:
		h.cancel()
		os.Exit(code)
	case <-h.done:
		h.cancel()
	case <-termChan:
		log.Println("TERMINATING")
		h.cancel()
		time.Sleep(time.Millisecond * 300)
		log.Println("TERMINATED")
	}
}
