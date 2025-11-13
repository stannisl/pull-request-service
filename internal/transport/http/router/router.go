package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stannisl/avito-test/internal/service"
	"github.com/stannisl/avito-test/internal/transport/http/handlers"
)

type Router interface {
	http.Handler
}

type ginRouter struct {
	//healthHandler *handlers.HealthHandler
	router *gin.Engine
}

func New(deps service.Dependencies) Router {
	router := gin.New()

	healthHandler := handlers.NewHealthHandler()

	health := router.Group("/health")
	{
		health.GET("/", healthHandler.Check)
	}

	registerMiddlewares(router)

	return &ginRouter{
		router: router,
	}
}

func registerMiddlewares(engine *gin.Engine) {
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
}

func (r *ginRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
