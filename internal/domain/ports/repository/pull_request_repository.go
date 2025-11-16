package repository

import (
	"context"
	"time"

	"pull-request-review/internal/domain/model"
)

type PullRequestRepository interface {
	Create(ctx context.Context, pullRequest *model.PullRequest) error
	GetByID(ctx context.Context, ID model.PullRequestID) (*model.PullRequest, error)
	Exists(ctx context.Context, ID model.PullRequestID) (bool, error)
	UpdateStatus(ctx context.Context, ID model.PullRequestID, status model.PullRequestStatus, mergedAt time.Time) error
	GetByReviewer(ctx context.Context, ID model.UserID) ([]model.PullRequest, error)
	GetPullRequestCountsByStatus(ctx context.Context) (map[string]int, error)
}