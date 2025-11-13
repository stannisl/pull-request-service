package service

import (
	"context"

	"github.com/mostanin/avito-test/internal/domain"
)

type UserService interface {
	SetActiveFlag(ctx context.Context, userID domain.UserID, isActive bool) (domain.User, error)
}

func NewUserServiceStub() UserService {
	return &userServiceStub{}
}

type userServiceStub struct{}

func (s *userServiceStub) SetActiveFlag(ctx context.Context, userID domain.UserID, isActive bool) (domain.User, error) {
	return domain.User{
		ID:       userID,
		IsActive: isActive,
	}, nil
}
