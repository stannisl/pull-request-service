package service

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/stannisl/avito-test/internal/domain"
	"github.com/stannisl/avito-test/internal/repository"
)

type TeamService interface {
	CreateTeam(ctx context.Context, team domain.Team) (*domain.Team, error)
	GetTeam(ctx context.Context, name domain.TeamName) (*domain.Team, error)
}

type teamService struct {
	userRepository repository.UserRepository
	teamRepository repository.TeamRepository
}

func (t *teamService) CreateTeam(ctx context.Context, team domain.Team) (*domain.Team, error) {
	err := t.teamRepository.CreateTeam(ctx, team)
	if err != nil {
		log.Printf("Error creating team: %v\n", err)
		return nil, err
	}

	for _, member := range team.Members {
		err = t.userRepository.CreateOrUpdateUser(ctx, &member)
		if err != nil {
			log.Printf("Error adding user to team: %v\n", err)
			return nil, err
		}
	}

	newTeam, err := t.GetTeam(ctx, team.Name)
	if err != nil {
		log.Printf("Error getting team: %v\n", err)
		return nil, err
	}

	return newTeam, nil
}

func (t *teamService) GetTeam(ctx context.Context, name domain.TeamName) (*domain.Team, error) {
	team, err := t.teamRepository.GetTeam(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEntityNotFound
		}
		log.Printf("Error getting team: %v\n", err)
		return nil, err
	}

	usersByTeam, err := t.userRepository.GetUsersByTeam(ctx, team.Name)
	if err != nil {
		log.Printf("Error getting users by team: %v\n", err)
		return nil, err
	}

	team.Members = usersByTeam

	return team, nil
}

func NewTeamService(userRepository repository.UserRepository, teamRepository repository.TeamRepository) TeamService {
	return &teamService{
		userRepository: userRepository,
		teamRepository: teamRepository,
	}
}
