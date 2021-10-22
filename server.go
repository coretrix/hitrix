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
	done chan bool
	exit chan int
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
	h.forceAlters()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
		h.done <- true
	}()
	h.await()

	app := service.DI().App()
	defer app.CancelContext()
	if err := srv.Shutdown(app.GlobalContext); err != nil {
		log.Println("Server forced to shutdown")
	}
}

func (h *Hitrix) RunBackgroundProcess(callback func(b *BackgroundProcessor)) {
	h.preDeploy()
	h.forceAlters()
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

	ormService := service.DI().OrmEngine()

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

func (h *Hitrix) forceAlters() {
	app := service.DI().App()

	if !app.IsInLocalMode() {
		return
	}

	forceAltersFlag := app.Flags.Bool("force-alters")

	if !forceAltersFlag {
		return
	}

	ormService := service.DI().OrmEngine()

	dbService := ormService.GetMysql()

	alters := ormService.GetAlters()
	var queries string

	for _, alter := range alters {
		queries += alter.SQL
	}

	if queries != "" {
		_, def := dbService.Query(queries)
		defer def()
	}

	altersSearch := ormService.GetRedisSearchIndexAlters()
	for _, alter := range altersSearch {
		alter.Execute()
	}

	log.Println("FORCE ALTERS executed")
}

func (h *Hitrix) await() {
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	app := service.DI().App()

	select {
	case code := <-h.exit:
		app.CancelContext()
		os.Exit(code)
	case <-h.done:
		app.CancelContext()
	case <-termChan:
		log.Println("TERMINATING")
		app.CancelContext()
		time.Sleep(time.Millisecond * 300)
		log.Println("TERMINATED")
	}
}
