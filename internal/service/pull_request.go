package service

import (
	"context"
	"log"

	"github.com/stannisl/avito-test/internal/domain"
	"github.com/stannisl/avito-test/internal/repository"
)

type PullRequestService interface {
	Create(ctx context.Context, prID domain.PRID, prName string, authorID domain.UserID) (*domain.PullRequest, error)
	Merge(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error)
	ReassignReviewer(
		ctx context.Context,
		prID domain.PRID,
		oldReviewerId domain.UserID,
	) (domain.PullRequest, domain.UserID, error)
	ListByReviewer(ctx context.Context, pr domain.PullRequest) ([]domain.PullRequest, error)
}

type pullRequestService struct {
	prRepo   repository.PullRequestRepository
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func (p *pullRequestService) Create(
	ctx context.Context,
	prID domain.PRID,
	prName string,
	authorID domain.UserID,
) (*domain.PullRequest, error) {
	user, err := p.userRepo.GetUser(ctx, authorID)
	if err != nil {
		log.Printf("error getting user: %v", err)
		return nil, err
	}

	availableActiveUsers, err := p.userRepo.GetActiveUsersByTeamWithLimit(ctx, user.TeamName, user.Id, 2)
	if err != nil {
		log.Printf("error getting active users: %v", err)
	}

	needMoreReviewers := len(availableActiveUsers) < 2

	availableActiveUsersIds := make([]domain.UserID, len(availableActiveUsers))
	for i, user := range availableActiveUsers {
		if i > 2 { // Мы имеем <= 2 всегда в этом списке
			break
		}
		availableActiveUsersIds[i] = user.Id
	}

	pr := &domain.PullRequest{
		ID:                prID,
		Name:              prName,
		AuthorID:          authorID,
		NeedMoreReviewers: needMoreReviewers,
		AssignedReviewers: availableActiveUsersIds,
	}

	if err := p.prRepo.Create(ctx, pr); err != nil {
		log.Printf("error creating pull request: %v", err)
		return nil, err
	}
	return pr, nil
}

func (p *pullRequestService) Merge(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (p *pullRequestService) ReassignReviewer(
	ctx context.Context,
	prID domain.PRID,
	oldReviewerId domain.UserID,
) (domain.PullRequest, domain.UserID, error) {
	//TODO implement me
	panic("implement me")
}

func (p *pullRequestService) ListByReviewer(ctx context.Context, pr domain.PullRequest) ([]domain.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func NewPullRequestService(
	prRepo repository.PullRequestRepository,
	userRepo repository.UserRepository,
	teamRepo repository.TeamRepository,
) PullRequestService {
	return &pullRequestService{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}
