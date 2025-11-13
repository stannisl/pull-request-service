package service

import (
	"context"

	"github.com/mostanin/avito-test/internal/domain"
)

type TeamService interface {
	CreateTeam(ctx context.Context, team domain.Team) (domain.Team, error)
	GetTeam(ctx context.Context, teamName domain.TeamName) (domain.Team, error)
}

func NewTeamServiceStub() TeamService {
	return &teamServiceStub{}
}

type teamServiceStub struct{}

func (s *teamServiceStub) CreateTeam(ctx context.Context, team domain.Team) (domain.Team, error) {
	return team, nil
}

func (s *teamServiceStub) GetTeam(ctx context.Context, teamName domain.TeamName) (domain.Team, error) {
	return domain.Team{Name: teamName}, nil
}
