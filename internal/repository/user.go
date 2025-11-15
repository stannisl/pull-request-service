// internal/repository/user.go
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/stannisl/avito-test/internal/domain"
	"github.com/stannisl/avito-test/pkg/db"
)

type UserRepository interface {
	CreateOrUpdateUser(ctx context.Context, user *domain.User) error
	GetUser(ctx context.Context, userID domain.UserID) (*domain.User, error)
	GetActiveUsersByTeamWithLimit(
		ctx context.Context,
		teamName domain.TeamName,
		excludeUserIDs []domain.UserID,
		limit int,
	) ([]domain.User, error)
	SetIsActive(ctx context.Context, userID domain.UserID, isActive bool) error
	GetUsersByTeam(ctx context.Context, teamName domain.TeamName) ([]domain.User, error)
}

type userRepository struct {
	*BaseRepository
}

func NewUserRepository(db *sqlx.DB, txManager db.TransactionManager) UserRepository {
	return &userRepository{
		BaseRepository: NewBaseRepository(db, txManager),
	}
}

func (u *userRepository) CreateOrUpdateUser(ctx context.Context, user *domain.User) error {
	return u.WithTransaction(ctx, func(ctx context.Context) error {
		query := `INSERT INTO users (id, username, team_name, is_active) VALUES ($1, $2, $3, $4) 
                 ON CONFLICT (id) DO UPDATE SET username = $2, team_name = $3, is_active = $4`

		executor := u.GetExecutor(ctx)
		_, err := executor.ExecContext(ctx, query, user.Id, user.Username, user.TeamName, user.IsActive)
		return err
	})
}

func (u *userRepository) GetUser(ctx context.Context, userID domain.UserID) (*domain.User, error) {
	query := `SELECT id, username, team_name, is_active FROM users WHERE id = $1`
	var user domain.User

	executor := u.GetExecutor(ctx)
	err := sqlx.GetContext(ctx, executor, &user, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrEntityNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) GetActiveUsersByTeamWithLimit(
	ctx context.Context,
	teamName domain.TeamName,
	excludeUserIDs []domain.UserID,
	limit int,
) ([]domain.User, error) {
	baseQuery := `SELECT id, username, team_name, is_active FROM users WHERE team_name = $1 AND is_active = true`
	args := []any{teamName}

	if len(excludeUserIDs) > 0 {
		placeholders := make([]string, len(excludeUserIDs))
		for i, id := range excludeUserIDs {
			placeholders[i] = fmt.Sprintf("$%d", len(args)+1)
			args = append(args, id)
		}
		baseQuery += fmt.Sprintf(" AND id NOT IN (%s)", strings.Join(placeholders, ", "))
	}

	baseQuery += " ORDER BY random()"

	if limit > 0 {
		baseQuery += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, limit)
	}

	var users []domain.User
	executor := u.GetExecutor(ctx)
	err := sqlx.SelectContext(ctx, executor, &users, baseQuery, args...)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *userRepository) SetIsActive(ctx context.Context, userID domain.UserID, isActive bool) error {
	return u.WithTransaction(ctx, func(ctx context.Context) error {
		query := `UPDATE users SET is_active = $1 WHERE id = $2`
		executor := u.GetExecutor(ctx)

		_, err := executor.ExecContext(ctx, query, isActive, userID)
		return err
	})
}

func (u *userRepository) GetUsersByTeam(ctx context.Context, teamName domain.TeamName) ([]domain.User, error) {
	query := `SELECT id, username, team_name, is_active FROM users WHERE team_name = $1`
	var users []domain.User

	executor := u.GetExecutor(ctx)
	err := sqlx.SelectContext(ctx, executor, &users, query, teamName)
	if err != nil {
		return nil, err
	}
	return users, nil
}
