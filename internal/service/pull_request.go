package service

import (
	"context"

	"github.com/stannisl/avito-test/internal/domain"
)

type PullRequestService interface {
	Create(ctx context.Context, prID domain.PRID, prName string, authorID domain.UserID) (domain.PullRequest, error)
	Merge(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID domain.PRID, oldReviewerId domain.UserID) (domain.PullRequest, domain.UserID, error)

	ListByReviewer(ctx context.Context, pr domain.PullRequest) ([]domain.PullRequest, error)
}

type pullRequestService struct{}

func (p *pullRequestService) Create(ctx context.Context, prID domain.PRID, prName string, authorID domain.UserID) (domain.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (p *pullRequestService) Merge(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (p *pullRequestService) ReassignReviewer(ctx context.Context, prID domain.PRID, oldReviewerId domain.UserID) (domain.PullRequest, domain.UserID, error) {
	//TODO implement me
	panic("implement me")
}

func (p *pullRequestService) ListByReviewer(ctx context.Context, pr domain.PullRequest) ([]domain.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func NewPullRequestService() PullRequestService {
	return &pullRequestService{}
}
