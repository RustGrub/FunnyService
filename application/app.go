package application

import (
	"context"
	"fmt"
	"github.com/RustGrub/FunnyGoService/config"
	"github.com/RustGrub/FunnyGoService/http/middleware"
	"github.com/RustGrub/FunnyGoService/http/router"
	"github.com/RustGrub/FunnyGoService/logger"
	"github.com/RustGrub/FunnyGoService/services"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Application struct {
	services []services.Service
	server   *http.Server
	logger   logger.Logger
}

func New(cfg *config.Config, l logger.Logger, serv []services.Service) *Application {
	handler := setApis(serv, l)
	return &Application{
		services: serv,
		server: &http.Server{
			Addr:    ":" + cfg.AppPort,
			Handler: handler,
		},
		logger: l,
	}
}

func (app *Application) Start() {
	listenErr := make(chan error, 1)
	go func() {
		listenErr <- app.server.ListenAndServe()
	}()
	app.logger.Info("http server started at port", app.server.Addr)

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	app.startServices()

	select {
	case err := <-listenErr:
		if err != nil {
			app.logger.Fatal(err)
		}
	case s := <-osSignals:
		app.logger.Info("SIGNAL:", s.String())
		app.server.SetKeepAlivesEnabled(false)
		app.stopServices()
		timeout := time.Second * 5
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer func() {
			cancel()
		}()
		if err := app.server.Shutdown(ctx); err != nil {
			app.logger.Fatal(err)
		}
	}
	app.logger.Info("Service stopped")
	app.logger.Close()
}

func (app *Application) startServices() {
	app.logger.Info("Starting listeners")
	for i := range app.services {
		if err := app.services[i].Start(); err != nil {
			app.logger.Fatal(fmt.Sprintf("Couldn't start listeners %s: %v", app.services[i].GetName(), err))
		}
	}
}

func (app *Application) stopServices() {
	for i := range app.services {
		if err := app.services[i].Stop(); err != nil {
			app.logger.Error(fmt.Sprintf("error while stopping listeners %s: %v", app.services[i].GetName(), err))
		}
	}
	app.logger.Info("Stopping listeners...")
}

func setApis(s []services.Service, l logger.Logger) http.Handler {
	mw := middleware.New().WithStartTime().WithLogger(l)
	return router.New(mw, getApis(s)...)
}

func getApis(services []services.Service) (apis []router.API) {
	for i := range services {
		if v, ok := services[i].(router.API); ok {
			apis = append(apis, v)
		}
	}

	return apis
}
