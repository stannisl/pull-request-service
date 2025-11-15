package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
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
	conn *sqlx.DB
}

func NewUserRepository(conn *sqlx.DB) UserRepository {
	return &userRepository{
		conn: conn,
	}
}

func (u *userRepository) CreateOrUpdateUser(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, username, team_name, is_active) VALUES (:id, :username, :team_name, :is_active)`

	tx, err := u.conn.BeginTxx(ctx, nil)
	if err != nil {
		errTr := tx.Rollback()
		if errTr != nil {
			return errTr
		}
		return err
	}

	if _, err = tx.NamedExecContext(ctx, query, user); err != nil {
		errTr := tx.Rollback()
		if errTr != nil {
			return errTr
		}
		return err
	}

	return tx.Commit()
}

func (u *userRepository) GetUser(ctx context.Context, userID domain.UserID) (*domain.User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	var user domain.User

	if err := u.conn.Get(&user, query, userID); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) GetActiveUsersByTeam(ctx context.Context, teamName domain.TeamName, excludeUserID domain.UserID) ([]domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) SetIsActive(ctx context.Context, userID domain.UserID, isActive bool) error {
	query := `UPDATE users SET is_active = $1 WHERE id = $2`

	tx, err := u.conn.BeginTxx(ctx, nil)
	if err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.ExecContext(ctx, query, isActive, userID); err != nil {
		errTr := tx.Rollback()
		if errTr != nil {
			return errTr
		}
		return err
	}

	return tx.Commit()
}

func (u *userRepository) GetUsersByTeam(ctx context.Context, teamName domain.TeamName) ([]domain.User, error) {
	query := `SELECT * FROM users WHERE team_name = $1`
	var users []domain.User

	err := u.conn.SelectContext(ctx, &users, query, teamName)
	if err != nil {
		return nil, err
	}
	return users, nil
}
