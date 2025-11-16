package service

import (
	"context"

	"github.com/stannisl/pull-request-service/internal/domain"
	"github.com/stannisl/pull-request-service/internal/repository"
)

type StatsService interface {
	GetStats(ctx context.Context) ([]domain.UserAssignments, error)
}

type statsService struct {
	statsRepository repository.StatsRepository
}

func NewStatsService(statsRepository repository.StatsRepository) StatsService {
	return &statsService{statsRepository: statsRepository}
}

func (s *statsService) GetStats(ctx context.Context) ([]domain.UserAssignments, error) {
	return s.statsRepository.CountAssignments(ctx)
}
