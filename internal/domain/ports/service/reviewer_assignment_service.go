package service

import (
	"context"
	"pr-review-service/internal/domain/model"
)

type ReviewerAssignmentService interface {
	AssignReviewer(ctx context.Context, pullRequestID model.PullRequestID, reviewerID model.UserID) error
	ReassignReviewer(ctx context.Context, pullRequestID model.PullRequestID, oldReviewerID model.UserID) error
}