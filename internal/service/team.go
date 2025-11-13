package service

import (
	"context"

	"github.com/stannisl/avito-test/internal/domain"
)

type TeamService interface {
	CreateTeam(ctx context.Context, team domain.Team) (*domain.Team, error)
	GetTeam(ctx context.Context, name domain.TeamName) (*domain.Team, error)
}

type teamService struct{}

func (t *teamService) CreateTeam(ctx context.Context, team domain.Team) (*domain.Team, error) {
	//TODO implement me
	panic("implement me")
}

func (t *teamService) GetTeam(ctx context.Context, name domain.TeamName) (*domain.Team, error) {
	//TODO implement me
	panic("implement me")
}

func NewTeamService() TeamService {
	return &teamService{}
}
