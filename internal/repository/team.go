package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
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

type teamRepository struct {
	db *sqlx.DB
}

func NewTeamRepository(db *sqlx.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (t *teamRepository) CreateTeam(ctx context.Context, team domain.Team) error {
	query := `INSERT INTO teams (name) VALUES ($1)`
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, team.Name)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (t *teamRepository) GetTeam(ctx context.Context, name domain.TeamName) (*domain.Team, error) {
	query := `SELECT name FROM teams WHERE name = $1`
	var team domain.Team

	err := t.db.GetContext(ctx, &team, query, name)
	if err != nil {
		return nil, err
	}

	return &team, nil
}

func (t *teamRepository) TeamExists(ctx context.Context, name domain.TeamName) (bool, error) {
	query := `SELECT name FROM teams WHERE name = $1`
	var team domain.Team

	err := t.db.GetContext(ctx, &team, query, name)
	if err != nil {
		return false, err
	}
	return team.Name == name, nil
}
