package handlers

import (
	"github.com/go-chi/chi/v5"

	"github.com/mostanin/avito-test/internal/service"
)

type Handlers struct {
	teams       *TeamHandler
	users       *UserHandler
	pullRequest *PullRequestHandler
	health      *HealthHandler
}

func New(deps service.Dependencies) *Handlers {
	return &Handlers{
		teams:       NewTeamHandler(deps.TeamService),
		users:       NewUserHandler(deps.UserService, deps.PullRequestService),
		pullRequest: NewPullRequestHandler(deps.PullRequestService),
		health:      NewHealthHandler(),
	}
}

func (h *Handlers) Mount(r chi.Router) {
	r.Route("/team", func(r chi.Router) {
		r.Post("/add", h.teams.CreateTeam)
		r.Get("/get", h.teams.GetTeam)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", h.users.SetIsActive)
		r.Get("/getReview", h.users.GetReviewAssignments)
	})

	r.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", h.pullRequest.Create)
		r.Post("/merge", h.pullRequest.Merge)
		r.Post("/reassign", h.pullRequest.Reassign)
	})

	r.Get("/health", h.health.Check)
}
