package service

import "context"

type StatisticsService interface {
	GetStatistics(ctx context.Context) (map[string]any, error)
}