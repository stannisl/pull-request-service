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
	router *gin.Engine
}

func New(deps service.Dependencies) Router {
	router := gin.New()
	registerMiddlewares(router)

	healthHandler := handlers.NewHealthHandler()
	teamHandler := handlers.NewTeamHandler(deps.TeamService)
	pullRequestHandler := handlers.NewPullRequestHandler(deps.PullRequestService)

	health := router.Group("/health")
	{
		health.GET("", healthHandler.Check)
	}

	team := router.Group("/team")
	{
		team.GET("/get", teamHandler.GetTeam)
		team.POST("/add", teamHandler.AddTeam)
	}

	pullRequest := router.Group("/pullRequest")
	{
		pullRequest.POST("/create", pullRequestHandler.Create)
	}

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
