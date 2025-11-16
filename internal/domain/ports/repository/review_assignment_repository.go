package repository

import (
	"context"

	"pull-request-review/internal/domain/model"
)

type ReviewAssignmentRepository interface {
	AssignReviewer(ctx context.Context, pullRequestID model.PullRequestID, reviewerID model.UserID) error
	AssignReviewers(ctx context.Context, pullRequestID model.PullRequestID, reviewerIDs []model.UserID) error
	GetByReviewer(ctx context.Context, pullRequestID model.PullRequestID) ([]model.PullRequest, error)
	Exists(ctx context.Context, pullRequestID model.PullRequestID, reviewerID model.UserID) (bool, error)
	GetReviewers(ctx context.Context, pullRequestID model.PullRequestID) ([]model.User, error)
	ReplaceReviewer(
		ctx context.Context, pullRequestID model.PullRequestID, oldReviewerID model.UserID, newReviewerID model.UserID,
	) error
	GetAssignmentCounts(ctx context.Context) (map[string]int, error)
}