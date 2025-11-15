package repository

import (
	"context"
	"github.com/google/uuid"
	"pr-review-service/internal/domain/model"
)

type ReviewAssignmentRepository interface {
	AssignReviewer(ctx context.Context, pullRequestID uuid.UUID, reviewerID uuid.UUID) error
	AssignReviewers(ctx context.Context, pullRequestID uuid.UUID, reviewerIDs []uuid.UUID) error
	GetByReviewer(ctx context.Context, pullRequestID uuid.UUID) ([]model.PullRequest, error)
	Exists(ctx context.Context, pullRequestID uuid.UUID, reviewerID uuid.UUID) (bool, error)
	GetReviewers(ctx context.Context, pullRequestID uuid.UUID) ([]model.User, error)
	ReplaceReviewer(ctx context.Context, pullRequestID uuid.UUID, reviewerID uuid.UUID) error
}