package service

import (
	"context"
	"pull-request-review/internal/domain/model"
)

type UserService interface {
	GetUser(ctx context.Context, ID model.UserID) (*model.User, error)
	SetActive(ctx context.Context, ID model.UserID, active bool) error
	BulkDeactivateUsers(ctx context.Context, userIDs []model.UserID) error
}