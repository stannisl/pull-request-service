package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/stannisl/avito-test/internal/domain"
)

type PullRequestRepository interface {
	// Create создает новый PR
	Create(ctx context.Context, pr *domain.PullRequest) error

	// GetByID возвращает PR по ID
	GetByID(ctx context.Context, prID domain.PRID) (*domain.PullRequest, error)

	// Update обновляет PR
	Update(ctx context.Context, pr *domain.PullRequest) error

	// GetByReviewerID возвращает все PR, где пользователь назначен ревьювером
	GetByReviewerID(ctx context.Context, reviewerID domain.UserID) ([]domain.PullRequest, error)

	// Exists проверяет существование PR
	Exists(ctx context.Context, prID domain.PRID) (bool, error)
}

type pullRequestRepository struct {
	db *sqlx.DB
}

func NewPullRequestRepository(db *sqlx.DB) PullRequestRepository {
	return &pullRequestRepository{db: db}
}

func (p *pullRequestRepository) Create(ctx context.Context, pr *domain.PullRequest) error {
	query := `INSERT INTO pull_requests (id, name, author_id, merged_at) VALUES (:id, :name, :author_id, :merged_at)`

	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.NamedExecContext(ctx, query, pr)
	if err != nil {
		return err
	}

	if len(pr.AssignedReviewers) > 0 {
		query := `INSERT INTO pull_request_reviewers (pull_request_id, reviewer_id) VALUES ($1, $2)`

		for _, userID := range pr.AssignedReviewers {
			_, err := tx.ExecContext(ctx, query, pr.ID, userID)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (p *pullRequestRepository) GetByID(ctx context.Context, prID domain.PRID) (*domain.PullRequest, error) {
	query := `SELECT * FROM pull_requests WHERE id = $1`
	var pr domain.PullRequest
	err := p.db.GetContext(ctx, &pr, query, prID)
	if err != nil {
		return nil, err
	}

	return &pr, nil
}

func (p *pullRequestRepository) Update(ctx context.Context, pr *domain.PullRequest) error {
	query := `
		UPDATE pull_requests 
		SET status = :status,  merged_at = :merged_at, need_more_reviewers = :need_more_reviewers 
		WHERE id = :id
	`

	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.NamedExecContext(ctx, query, pr)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (p *pullRequestRepository) GetByReviewerID(
	ctx context.Context,
	reviewerID domain.UserID,
) ([]domain.PullRequest, error) {
	panic("implement me")
}

func (p *pullRequestRepository) Exists(ctx context.Context, prID domain.PRID) (bool, error) {
	query := `SELECT * FROM pull_requests WHERE id = $1`
	var pr domain.PullRequest

	err := p.db.GetContext(ctx, &pr, query, prID)
	if err != nil {
		return false, err
	}

	return pr.ID == prID, nil
}
