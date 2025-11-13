package service

import (
	"context"

	"github.com/mostanin/avito-test/internal/domain"
)

type PullRequestService interface {
	Create(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error)
	Merge(ctx context.Context, prID domain.PRID) (domain.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID domain.PRID, reviewerID domain.UserID) (domain.PullRequest, domain.UserID, error)
	ListByReviewer(ctx context.Context, reviewerID domain.UserID) ([]domain.PullRequest, error)
}

func NewPullRequestServiceStub() PullRequestService {
	return &pullRequestServiceStub{}
}

type pullRequestServiceStub struct{}

func (s *pullRequestServiceStub) Create(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error) {
	pr.Status = domain.PullRequestStatusOpen
	return pr, nil
}

func (s *pullRequestServiceStub) Merge(ctx context.Context, prID domain.PRID) (domain.PullRequest, error) {
	return domain.PullRequest{
		ID:     prID,
		Status: domain.PullRequestStatusMerged,
	}, nil
}

func (s *pullRequestServiceStub) ReassignReviewer(ctx context.Context, prID domain.PRID, reviewerID domain.UserID) (domain.PullRequest, domain.UserID, error) {
	return domain.PullRequest{
		ID:     prID,
		Status: domain.PullRequestStatusOpen,
	}, reviewerID, nil
}

func (s *pullRequestServiceStub) ListByReviewer(ctx context.Context, reviewerID domain.UserID) ([]domain.PullRequest, error) {
	return []domain.PullRequest{
		{
			ID:       "stub-pr",
			Name:     "Stub Pull Request",
			AuthorID: "stub-author",
			Status:   domain.PullRequestStatusOpen,
		},
	}, nil
}
