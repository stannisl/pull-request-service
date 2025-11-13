package service

import (
	"context"

	"github.com/stannisl/avito-test/internal/domain"
)

type UserService interface {
	SetIsActive(ctx context.Context, userID domain.UserID, isActive bool) error
	GetReview(ctx context.Context, userID domain.UserID) (*domain.User, error)
}

type userService struct{}

func (u *userService) SetIsActive(ctx context.Context, userID domain.UserID, isActive bool) error {
	//TODO implement me
	panic("implement me")
}

func (u *userService) GetReview(ctx context.Context, userID domain.UserID) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func NewUserService() UserService {
	return &userService{}
}
