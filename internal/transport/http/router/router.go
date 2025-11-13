package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stannisl/avito-test/internal/service"
)

type Router interface {
	http.Handler
}

type ginRouter struct {
	router *gin.Engine
}

func New(deps service.Dependencies) Router {
	handler := gin.New()

	registerMiddlewares(handler)

	// register routes
	// handler.Group()

	return &ginRouter{
		router: handler,
	}
}

func registerMiddlewares(engine *gin.Engine) {
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
}

func (r *ginRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
