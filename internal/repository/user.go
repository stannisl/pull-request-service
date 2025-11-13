package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stannisl/avito-test/internal/domain"
)

type UserRepository interface {
	// CreateOrUpdateUser создает или обновляет пользователя
	CreateOrUpdateUser(ctx context.Context, user *domain.User) error

	// GetUser возвращает пользователя по ID
	GetUser(ctx context.Context, userID domain.UserID) (*domain.User, error)

	// GetActiveUsersByTeam возвращает активных пользователей команды (исключая указанного)
	GetActiveUsersByTeam(ctx context.Context, teamName domain.TeamName, excludeUserID domain.UserID) ([]domain.User, error)

	// SetIsActive устанавливает флаг активности пользователя
	SetIsActive(ctx context.Context, userID domain.UserID, isActive bool) error

	// GetUsersByTeam возвращает всех пользователей команды
	GetUsersByTeam(ctx context.Context, teamName domain.TeamName) ([]domain.User, error)
}

type userRepository struct {
	conn *pgxpool.Pool
}

func (u *userRepository) CreateOrUpdateUser(ctx context.Context, user *domain.User) error {
	//query := `insert into users () values ($1, $2, $3, $4, $5)`
	panic("implement me")
}

func (u *userRepository) GetUser(ctx context.Context, userID domain.UserID) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) GetActiveUsersByTeam(ctx context.Context, teamName domain.TeamName, excludeUserID domain.UserID) ([]domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) SetIsActive(ctx context.Context, userID domain.UserID, isActive bool) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) GetUsersByTeam(ctx context.Context, teamName domain.TeamName) ([]domain.User, error) {
	//TODO implement me
	panic("implement me")
}
