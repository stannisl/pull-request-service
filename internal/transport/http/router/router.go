package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/mostanin/avito-test/internal/config"
	"github.com/mostanin/avito-test/internal/service"
	handlerspkg "github.com/mostanin/avito-test/internal/transport/http/handlers"
)

type Router interface {
	http.Handler
}

type chiRouter struct {
	mux *chi.Mux
}

func New(cfg config.Config, deps service.Dependencies) Router {
	mux := chi.NewRouter()
	registerMiddlewares(mux)

	handlers := handlerspkg.New(deps)
	handlers.Mount(mux)

	return &chiRouter{
		mux: mux,
	}
}

func registerMiddlewares(mux *chi.Mux) {
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.AllowContentType("application/json"))
}

func (r *chiRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
