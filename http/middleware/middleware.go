package middleware

import (
	"github.com/RustGrub/FunnyGoService/logger"
	"github.com/rs/cors"
	"net/http"
	"time"
)

type Middleware struct {
	start  time.Time
	logger logger.Logger
}

func New() *Middleware {
	return &Middleware{}
}

func (m *Middleware) CorsMiddleware(h http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		Debug:            false,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PATCH"},
	})

	return c.Handler(h)
}

func (m *Middleware) AuthMiddleware(h http.Handler) http.Handler {
	// no auth yet
	return nil
}

func (m *Middleware) WithLogger(l logger.Logger) *Middleware {
	m.logger = l
	return m
}

func (m *Middleware) WithStartTime() *Middleware {
	m.start = time.Now()
	return m
}

func (m *Middleware) GetStartTime() time.Time {
	return m.start
}
