package repository

import (
	"context"

	"github.com/stannisl/avito-test/internal/domain"
)

type TeamRepository interface {
	// CreateTeam создает новую команду
	CreateTeam(ctx context.Context, team domain.Team) error
	
	// GetTeam возвращает команду по имени
	GetTeam(ctx context.Context, name domain.TeamName) (*domain.Team, error)
	
	// TeamExists проверяет существование команды
	TeamExists(ctx context.Context, name domain.TeamName) (bool, error)
}
