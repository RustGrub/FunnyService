package router

import (
	"github.com/RustGrub/FunnyGoService/http/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

type API interface {
	Router(r *mux.Router, c *middleware.Middleware)
}

func New(mw *middleware.Middleware, apis ...API) http.Handler {
	root := mux.NewRouter().StrictSlash(true).PathPrefix("/").Subrouter()
	api := root.PathPrefix("/api").Subrouter()

	for i := range apis {
		// for auth management
		v := api.PathPrefix("/").Subrouter()
		apis[i].Router(v, mw)
	}

	return mw.CorsMiddleware(root)
}
