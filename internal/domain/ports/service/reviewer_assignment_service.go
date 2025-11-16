package service

import (
	"context"
	"github.com/google/uuid"
	"pull-request-review/internal/domain/model"
)

type ReviewerAssignmentService interface {
	AssignReviewer(ctx context.Context, pullRequestID model.PullRequestID, reviewerID model.UserID) error
	AssignInitialReviewers(ctx context.Context, pr *model.PullRequest, authorTeamID uuid.UUID) ([]model.UserID, error)
	ReassignReviewer(ctx context.Context, pr *model.PullRequest, oldReviewerID model.UserID) (model.UserID, error)
}