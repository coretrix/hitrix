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

	"github.com/coretrix/hitrix/service"
	"github.com/fatih/color"

	"github.com/99designs/gqlgen/graphql"
)

type Hitrix struct {
	ctx    context.Context
	cancel context.CancelFunc
	done   chan bool
	exit   chan int
}

func (h *Hitrix) RunServer(defaultPort uint, server graphql.ExecutableSchema, ginInitHandler GinInitHandler, gqlServerInitHandler GQLServerInitHandler) {
	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", defaultPort)
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: InitGin(server, ginInitHandler, gqlServerInitHandler),
	}

	h.preDeploy()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
		h.done <- true
	}()
	h.await()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server forced to shutdown")
	}
}

func (h *Hitrix) RunBackgroundProcess(callback func(b *BackgroundProcessor)) {
	h.preDeploy()
	callback(&BackgroundProcessor{Server: h})
	h.await()
}

func (h *Hitrix) runDynamicScrips(ctx context.Context, code string) {
	scripts := service.DI().App().Scripts
	if len(scripts) == 0 {
		panic(fmt.Sprintf("unknown script %s", code))
	}
	for _, defCode := range scripts {
		if defCode == code {
			def, has := service.GetServiceOptional(defCode)
			if !has {
				panic(fmt.Sprintf("unknown script %s", code))
			}
			defScript := def.(Script)
			defScript.Run(ctx, &exit{s: h})
			return
		}
	}
	panic(fmt.Sprintf("unknown script %s", code))
}

func (h *Hitrix) preDeploy() {
	app := service.DI().App()

	if app.IsInTestMode() {
		return
	}

	preDeployFlag := app.Flags.Bool("pre-deploy")

	if !preDeployFlag {
		return
	}

	ormService, has := service.DI().OrmEngine()
	if !has {
		return
	}

	alters := ormService.GetAlters()

	hasAlters := false
	for _, alter := range alters {
		if alter.Safe {
			color.Green("%s\n\n", alter.SQL)
		} else {
			color.Red("%s\n\n", alter.SQL)
		}
		hasAlters = true
	}

	if hasAlters {
		os.Exit(1)
	}

	os.Exit(0)
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
