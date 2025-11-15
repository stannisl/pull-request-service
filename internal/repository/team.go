// internal/repository/team.go
package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/stannisl/avito-test/internal/domain"
	"github.com/stannisl/avito-test/pkg/db"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, team domain.Team) error
	GetTeam(ctx context.Context, name domain.TeamName) (*domain.Team, error)
	TeamExists(ctx context.Context, name domain.TeamName) (bool, error)
}

type teamRepository struct {
	*BaseRepository
}

func NewTeamRepository(db *sqlx.DB, txManager db.TransactionManager) TeamRepository {
	return &teamRepository{
		BaseRepository: NewBaseRepository(db, txManager),
	}
}

func (t *teamRepository) CreateTeam(ctx context.Context, team domain.Team) error {
	return t.WithTransaction(ctx, func(ctx context.Context) error {
		query := `INSERT INTO teams (name) VALUES ($1)`
		executor := t.GetExecutor(ctx)

		_, err := executor.ExecContext(ctx, query, team.Name)
		return err
	})
}

func (t *teamRepository) GetTeam(ctx context.Context, name domain.TeamName) (*domain.Team, error) {
	query := `SELECT name FROM teams WHERE name = $1`
	var team domain.Team

	executor := t.GetExecutor(ctx)
	err := sqlx.GetContext(ctx, executor, &team, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrEntityNotFound
		}
		return nil, err
	}

	return &team, nil
}

func (t *teamRepository) TeamExists(ctx context.Context, name domain.TeamName) (bool, error) {
	query := `SELECT name FROM teams WHERE name = $1`
	var team domain.Team

	executor := t.GetExecutor(ctx)
	err := sqlx.GetContext(ctx, executor, &team, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return team.Name == name, nil
}
