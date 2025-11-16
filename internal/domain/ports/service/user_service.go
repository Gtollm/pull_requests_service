package service

import (
	"context"
	"pr-review-service/internal/domain/model"
)

type UserService interface {
	GetUser(ctx context.Context, ID model.UserID) (*model.User, error)
	SetActive(ctx context.Context, ID model.UserID, active bool) error
	BulkDeactivateUsers(ctx context.Context, userIDs []model.UserID) error
}