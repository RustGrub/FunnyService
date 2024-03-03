package api

import (
	"fmt"
	"github.com/RustGrub/FunnyGoService/http/middleware"
	"github.com/RustGrub/FunnyGoService/logger"
	"github.com/RustGrub/FunnyGoService/provider"
	"github.com/RustGrub/FunnyGoService/services/FunnyService/usecases"
	"github.com/gorilla/mux"
)

type FunnyService struct {
	provider provider.Provider
	logger   logger.Logger
	usecases *usecases.FunnyService
}

func New(p provider.Provider, l logger.Logger) *FunnyService {
	uc := usecases.New(p.GetTxManager(), p.GetGoodsRepository(), p.GetGoodsCache(), p.GetBroker())
	return &FunnyService{
		provider: p,
		logger:   l,
		usecases: uc,
	}
}

func (s *FunnyService) Router(r *mux.Router, c *middleware.Middleware) {
	v1 := r.PathPrefix("/v1/").Subrouter()
	v1.HandleFunc("/good", s.getGood).Methods("GET")
	v1.HandleFunc("/good/create", s.createGood).Methods("POST")
	v1.HandleFunc("/good/update", s.updateGood).Methods("PATCH")
	v1.HandleFunc("/good/remove", s.removeGood).Methods("DELETE")
	v1.HandleFunc("/goods/list", s.getList).Methods("GET")
	// Не знаю, ошибка или нет, в задании ../reprioritiIze
	v1.HandleFunc("/good/reprioritize", s.reprioritizeGoods).Methods("PATCH")
}

func (s *FunnyService) Start() error {
	s.logger.Info(fmt.Sprintf("Service %v has started...", s.GetName()))

	return nil
}

func (s *FunnyService) Stop() error {
	s.logger.Info(fmt.Printf("Service %v has stopped...", s.GetName()))
	return nil
}

func (s *FunnyService) GetName() string {
	return "FunnyService"
}
