package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/stannisl/avito-test/internal/domain"
)

type PullRequestRepository interface {
	// Create создает новый PR
	Create(ctx context.Context, pr domain.PullRequest) error

	// GetByID возвращает PR по ID
	GetByID(ctx context.Context, prID domain.PRID) (*domain.PullRequest, error)

	// Update обновляет PR
	Update(ctx context.Context, pr domain.PullRequest) error

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

func (p *pullRequestRepository) Create(ctx context.Context, pr domain.PullRequest) error {

	//TODO implement me
	panic("implement me")
}

func (p *pullRequestRepository) GetByID(ctx context.Context, prID domain.PRID) (*domain.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (p *pullRequestRepository) Update(ctx context.Context, pr domain.PullRequest) error {
	//TODO implement me
	panic("implement me")
}

func (p *pullRequestRepository) GetByReviewerID(ctx context.Context, reviewerID domain.UserID) ([]domain.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (p *pullRequestRepository) Exists(ctx context.Context, prID domain.PRID) (bool, error) {
	//TODO implement me
	panic("implement me")
}
