package service

import (
	"context"
	"pull-request-review/internal/domain/model"
)

type PullRequestService interface {
	CreatePullRequest(ctx context.Context, pullRequest *model.PullRequest) error
	GetPullRequest(ctx context.Context, ID model.PullRequestID) (*model.PullRequest, error)
	GetUserReviews(ctx context.Context, userID model.UserID) ([]model.PullRequest, error)
	ReassignPullRequest(ctx context.Context, ID model.PullRequestID, oldReviewerID model.UserID) error
}