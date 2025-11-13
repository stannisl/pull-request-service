package app

import (
	"context"
	"fmt"

	"github.com/mostanin/avito-test/internal/config"
	"github.com/mostanin/avito-test/internal/service"
	"github.com/mostanin/avito-test/internal/transport/http/router"
)

type App struct {
	router router.Router
}

func New(ctx context.Context, cfg config.Config) (*App, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	services := service.Dependencies{
		TeamService:        service.NewTeamServiceStub(),
		UserService:        service.NewUserServiceStub(),
		PullRequestService: service.NewPullRequestServiceStub(),
	}

	httpRouter := router.New(cfg, services)

	return &App{
		router: httpRouter,
	}, nil
}

func (a *App) Router() router.Router {
	return a.router
}
