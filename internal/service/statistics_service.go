package service

import (
	"context"

	"pull-request-review/internal/domain/ports/repository"
	"pull-request-review/internal/infrastructure/adapters/logger"
)

type StatisticsService struct {
	reviewAssignmentRepo repository.ReviewAssignmentRepository
	prRepo               repository.PullRequestRepository
	logger               logger.Logger
}

func NewStatisticsService(
	reviewAssignmentRepo repository.ReviewAssignmentRepository,
	prRepo repository.PullRequestRepository,
	logger logger.Logger,
) *StatisticsService {
	return &StatisticsService{
		reviewAssignmentRepo: reviewAssignmentRepo,
		prRepo:               prRepo,
		logger:               logger,
	}
}

func (s *StatisticsService) GetStatistics(ctx context.Context) (map[string]any, error) {
	userAssignments, err := s.reviewAssignmentRepo.GetAssignmentCounts(ctx)
	if err != nil {
		s.logger.Warn("Failed to get assignment counts")
		userAssignments = make(map[string]int)
	}

	prCounts, err := s.prRepo.GetPullRequestCountsByStatus(ctx)
	if err != nil {
		s.logger.Warn("Failed to get PR counts by status")
		prCounts = make(map[string]int)
	}

	totalAssignments := 0
	for _, count := range userAssignments {
		totalAssignments += count
	}

	return map[string]any{
		"user_assignments":  userAssignments,
		"pr_counts":         prCounts,
		"total_assignments": totalAssignments,
	}, nil
}