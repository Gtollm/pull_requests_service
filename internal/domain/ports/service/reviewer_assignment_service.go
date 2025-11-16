package service

import (
	"context"
	"pull-request-review/internal/domain/model"
)

type ReviewerAssignmentService interface {
	AssignReviewer(ctx context.Context, pullRequestID model.PullRequestID, reviewerID model.UserID) error
	ReassignReviewer(ctx context.Context, pullRequestID model.PullRequestID, oldReviewerID model.UserID) error
}