package Listener

import (
	"fmt"
	"github.com/RustGrub/FunnyGoService/broker"
	"github.com/RustGrub/FunnyGoService/http/middleware"
	"github.com/RustGrub/FunnyGoService/logger"
	"github.com/RustGrub/FunnyGoService/provider"
	"github.com/RustGrub/FunnyGoService/services/ListeningService/usecases"
	"github.com/gorilla/mux"
	"net/http"
)

type ListeningService struct {
	provider  provider.Provider
	logger    logger.Logger
	msgBroker broker.MessageBroker
	useCases  *usecases.UseCases
}

func New(p provider.Provider, l logger.Logger) *ListeningService {
	uc := usecases.New(l, p.GetLsRepo())
	return &ListeningService{
		provider:  p,
		logger:    l,
		msgBroker: p.GetBroker(),
		useCases:  uc,
	}
}

func (s *ListeningService) Router(r *mux.Router, c *middleware.Middleware) {
	v1 := r.PathPrefix("/v1/").Subrouter()
	v1.HandleFunc("/broker/alive", func(writer http.ResponseWriter, request *http.Request) {
		// пингануть брокер, типа живой еще
	})
}

func (s *ListeningService) Start() error {
	err := s.SubscribeToNewGoods()
	if err != nil {
		// Наверн не надо fatal
		s.logger.Fatal(fmt.Sprintf("Service %v can't start...", s.GetName()))
		return err
	}

	s.logger.Info(fmt.Sprintf("Service %v has started...", s.GetName()))
	return nil
}

func (s *ListeningService) Stop() error {
	s.logger.Info(fmt.Printf("Service %v has stopped...", s.GetName()))
	return nil
}

func (s *ListeningService) GetName() string {
	return "UseCases"
}
