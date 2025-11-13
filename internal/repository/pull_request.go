package repository

import (
	"context"

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
