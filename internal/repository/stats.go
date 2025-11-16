package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/stannisl/pull-request-service/internal/domain"
	"github.com/stannisl/pull-request-service/pkg/db"
)

type StatsRepository interface {
	CountAssignments(ctx context.Context) ([]domain.UserAssignments, error)
}

type statsRepository struct {
	*BaseRepository
}

func NewStatsRepository(db *sqlx.DB, txManager db.TransactionManager) StatsRepository {
	return &statsRepository{
		BaseRepository: NewBaseRepository(db, txManager),
	}
}

func (r *statsRepository) CountAssignments(ctx context.Context) ([]domain.UserAssignments, error) {
	query := `SELECT reviewer_id, COUNT(*) AS assigned_count FROM pull_request_reviewers group by reviewer_id`

	var rows []domain.UserAssignments

	err := r.db.SelectContext(ctx, &rows, query)

	return rows, err
}
