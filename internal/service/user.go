package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/stannisl/avito-test/internal/domain"
	"github.com/stannisl/avito-test/internal/repository"
)

type UserService interface {
	SetIsActive(ctx context.Context, userID domain.UserID, isActive bool) (*domain.User, error)
	GetReview(ctx context.Context, userID domain.UserID) ([]domain.PullRequest, error)
}

type userService struct {
	userRepo repository.UserRepository
	prRepo   repository.PullRequestRepository
}

func (u *userService) SetIsActive(ctx context.Context, userID domain.UserID, isActive bool) (*domain.User, error) {
	err := u.userRepo.SetIsActive(ctx, userID, isActive)

	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userService) GetReview(ctx context.Context, userID domain.UserID) ([]domain.PullRequest, error) {
	_, err := u.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, err
	}

	prs, err := u.prRepo.GetByReviewerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return prs, nil
}

func NewUserService(userRepo repository.UserRepository, prRepo repository.PullRequestRepository) UserService {
	return &userService{
		userRepo: userRepo,
		prRepo:   prRepo,
	}
}
