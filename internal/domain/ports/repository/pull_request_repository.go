package repository

import (
	"context"
	"github.com/google/uuid"
	"pr-review-service/internal/domain/model"
	"time"
)

type PullRequestRepository interface {
	Create(ctx context.Context, pullRequest *model.PullRequest) error
	GetByID(ctx context.Context, ID uuid.UUID) (*model.PullRequest, error)
	Exists(ctx context.Context, ID uuid.UUID) (bool, error)
	UpdateStatus(ctx context.Context, ID uuid.UUID, status model.PullRequestStatus, mergedAt time.Time) error
	GetByReviewer(ctx context.Context, ID uuid.UUID) ([]model.PullRequest, error)
}