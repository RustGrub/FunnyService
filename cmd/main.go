package main

import (
	"github.com/RustGrub/FunnyGoService/application"
	"github.com/RustGrub/FunnyGoService/config"
	"github.com/RustGrub/FunnyGoService/logger"
	"github.com/RustGrub/FunnyGoService/provider"
	"github.com/RustGrub/FunnyGoService/services"
	"github.com/RustGrub/FunnyGoService/services/FunnyService/api"
	"github.com/RustGrub/FunnyGoService/services/ListeningService/Listener"
)

func main() {

	cfg := config.LoadConfig()
	l := logger.NewLogger(cfg)
	p := provider.New(cfg, l)

	funnyService := api.New(p, l)
	listeningService := Listener.New(p, l)
	s := []services.Service{funnyService, listeningService}
	app := application.New(cfg, l, s)
	app.Start()
}
