package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/stannisl/avito-test/internal/domain"
	"github.com/stannisl/avito-test/internal/repository"
)

type PullRequestService interface {
	Create(ctx context.Context, prID domain.PRID, prName string, authorID domain.UserID) (*domain.PullRequest, error)
	Merge(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error)
	ReassignReviewer(
		ctx context.Context,
		prID domain.PRID,
		oldReviewerId domain.UserID,
	) (*domain.PullRequest, *domain.UserID, error)
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

	availableActiveUsers, err := p.userRepo.GetActiveUsersByTeamWithLimit(ctx, user.TeamName, []domain.UserID{user.Id}, 2)
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
		Status:            domain.PullRequestStatusOpen,
	}

	if err := p.prRepo.Create(ctx, pr); err != nil {
		log.Printf("error creating pull request: %v", err)
		return nil, err
	}
	return pr, nil
}

func (p *pullRequestService) Merge(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error) {
	pullRequest, err := p.prRepo.GetByID(ctx, pr.ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEntityNotFound
		}
		log.Printf("error getting pull request: %v", err)
		return nil, err
	}

	t := time.Now()
	pullRequest.MergedAt = &t
	pullRequest.Status = domain.PullRequestStatusMerged

	err = p.prRepo.Update(ctx, pullRequest)
	if err != nil {
		log.Printf("error updating pull request status: %v", err)
		return nil, err
	}

	return pullRequest, nil
}

func (p *pullRequestService) ReassignReviewer(
	ctx context.Context,
	prID domain.PRID,
	oldReviewerId domain.UserID,
) (*domain.PullRequest, *domain.UserID, error) {
	pullRequest, err := p.prRepo.GetByID(ctx, prID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, domain.ErrEntityNotFound
		}
		log.Printf("error getting pull request: %v", err)
		return nil, nil, err
	}

	if pullRequest.Status == domain.PullRequestStatusMerged {
		return nil, nil, domain.ErrPRMerged
	}

	isAssigned := false
	for _, reviewer := range pullRequest.AssignedReviewers {
		if reviewer == oldReviewerId {
			isAssigned = true
			break
		}
	}
	if !isAssigned {
		return nil, nil, domain.ErrNotAssigned
	}

	oldUser, err := p.userRepo.GetUser(ctx, oldReviewerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, domain.ErrEntityNotFound
		}
		log.Printf("error getting old user: %v", err)
		return nil, nil, err
	}

	excludeUsers := []domain.UserID{oldReviewerId, pullRequest.AuthorID}
	excludeUsers = append(excludeUsers, pullRequest.AssignedReviewers...)
	availableReviewers, err := p.userRepo.GetActiveUsersByTeamWithLimit(ctx, oldUser.TeamName, excludeUsers, -1)
	if err != nil {
		log.Printf("error getting active users: %v", err)
		return nil, nil, err
	}

	if len(availableReviewers) == 0 {
		return nil, nil, domain.ErrNoCandidate
	}

	newReviewer := availableReviewers[rand.Intn(len(availableReviewers))]

	newReviewers := make([]domain.UserID, 0, len(pullRequest.AssignedReviewers))
	for _, reviewer := range pullRequest.AssignedReviewers {
		if reviewer == oldReviewerId {
			newReviewers = append(newReviewers, newReviewer.Id)
		} else {
			newReviewers = append(newReviewers, reviewer)
		}
	}
	pullRequest.AssignedReviewers = newReviewers

	if err := p.prRepo.Update(ctx, pullRequest); err != nil {
		log.Printf("error updating pull request: %v", err)
		return nil, nil, err
	}

	return pullRequest, &newReviewer.Id, nil
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
