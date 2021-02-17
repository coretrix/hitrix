package hitrix

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

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
