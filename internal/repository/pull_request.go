// internal/repository/pull_request.go
package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stannisl/pull-request-service/internal/domain"
	"github.com/stannisl/pull-request-service/pkg/db"
)

type PullRequestRepository interface {
	Create(ctx context.Context, pr *domain.PullRequest) error
	GetByID(ctx context.Context, prID domain.PRID) (*domain.PullRequest, error)
	Update(ctx context.Context, pr *domain.PullRequest) error
	GetByReviewerID(ctx context.Context, reviewerID domain.UserID) ([]domain.PullRequest, error)
	Exists(ctx context.Context, prID domain.PRID) (bool, error)
}

type pullRequestRepository struct {
	*BaseRepository
}

func NewPullRequestRepository(db *sqlx.DB, txManager db.TransactionManager) PullRequestRepository {
	return &pullRequestRepository{
		BaseRepository: NewBaseRepository(db, txManager),
	}
}

func (p *pullRequestRepository) Create(ctx context.Context, pr *domain.PullRequest) error {
	return p.WithTransaction(ctx, func(ctx context.Context) error {
		executor := p.GetExecutor(ctx)

		// Вставляем основной PR
		query := `INSERT INTO pull_requests (id, name, author_id, status, need_more_reviewers, created_at, merged_at) 
                 VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err := executor.ExecContext(ctx, query,
			pr.ID, pr.Name, pr.AuthorID, pr.Status, pr.NeedMoreReviewers, pr.CreatedAt, pr.MergedAt)
		if err != nil {
			var pgErr *pq.Error
			if errors.As(err, &pgErr) {
				if pgErr.Code == "23505" {
					return domain.ErrPRExists
				}
			}

			return err
		}

		// Вставляем ревьюверов
		if len(pr.AssignedReviewers) > 0 {
			reviewerQuery := `INSERT INTO pull_request_reviewers (pull_request_id, reviewer_id) VALUES ($1, $2)`
			for _, reviewerID := range pr.AssignedReviewers {
				_, err := executor.ExecContext(ctx, reviewerQuery, pr.ID, reviewerID)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (p *pullRequestRepository) GetByID(ctx context.Context, prID domain.PRID) (*domain.PullRequest, error) {
	executor := p.GetExecutor(ctx)

	// Получаем основной PR
	query := `SELECT id, name, author_id, status, need_more_reviewers, created_at, merged_at 
              FROM pull_requests WHERE id = $1`
	var pr domain.PullRequest
	err := sqlx.GetContext(ctx, executor, &pr, query, prID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrEntityNotFound
		}
		return nil, err
	}

	// Получаем ревьюверов
	reviewersQuery := `SELECT reviewer_id FROM pull_request_reviewers WHERE pull_request_id = $1`
	var reviewers []domain.UserID
	err = sqlx.SelectContext(ctx, executor, &reviewers, reviewersQuery, prID)
	if err != nil {
		return nil, err
	}

	pr.AssignedReviewers = reviewers
	return &pr, nil
}

func (p *pullRequestRepository) Update(ctx context.Context, pr *domain.PullRequest) error {
	return p.WithTransaction(ctx, func(ctx context.Context) error {
		executor := p.GetExecutor(ctx)

		// Обновляем основной PR
		query := `UPDATE pull_requests 
                 SET name = $1, status = $2, need_more_reviewers = $3, merged_at = $4 
                 WHERE id = $5`
		_, err := executor.ExecContext(ctx, query,
			pr.Name, pr.Status, pr.NeedMoreReviewers, pr.MergedAt, pr.ID)
		if err != nil {
			return err
		}

		// Удаляем старых ревьюверов и добавляем новых
		deleteQuery := `DELETE FROM pull_request_reviewers WHERE pull_request_id = $1`
		_, err = executor.ExecContext(ctx, deleteQuery, pr.ID)
		if err != nil {
			return err
		}

		if len(pr.AssignedReviewers) > 0 {
			insertQuery := `INSERT INTO pull_request_reviewers (pull_request_id, reviewer_id) VALUES ($1, $2)`
			for _, reviewerID := range pr.AssignedReviewers {
				_, err := executor.ExecContext(ctx, insertQuery, pr.ID, reviewerID)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (p *pullRequestRepository) GetByReviewerID(
	ctx context.Context,
	reviewerID domain.UserID,
) ([]domain.PullRequest, error) {
	query := `
		SELECT pr.id, pr.name, pr.author_id, pr.status, pr.need_more_reviewers, pr.created_at, pr.merged_at
		FROM pull_requests pr
		JOIN pull_request_reviewers prr ON pr.id = prr.pull_request_id
		WHERE prr.reviewer_id = $1
	`
	var prs []domain.PullRequest
	executor := p.GetExecutor(ctx)
	err := sqlx.SelectContext(ctx, executor, &prs, query, reviewerID)
	if err != nil {
		return nil, err
	}

	// Для каждого PR загружаем ревьюверов
	for i := range prs {
		reviewersQuery := `SELECT reviewer_id FROM pull_request_reviewers WHERE pull_request_id = $1`
		var reviewers []domain.UserID
		err := sqlx.SelectContext(ctx, executor, &reviewers, reviewersQuery, prs[i].ID)
		if err != nil {
			return nil, err
		}
		prs[i].AssignedReviewers = reviewers
	}

	return prs, nil
}

func (p *pullRequestRepository) Exists(ctx context.Context, prID domain.PRID) (bool, error) {
	query := `SELECT id FROM pull_requests WHERE id = $1`
	var id string

	executor := p.GetExecutor(ctx)
	err := sqlx.GetContext(ctx, executor, &id, query, prID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
