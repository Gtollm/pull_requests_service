package service

import (
	"context"
	"pull-request-review/internal/domain/model"
)

type UserService interface {
	GetUser(ctx context.Context, ID model.UserID) (*model.User, error)
	GetUserWithTeamName(ctx context.Context, ID model.UserID) (*model.User, string, error)
	SetActive(ctx context.Context, ID model.UserID, active bool) (*model.User, error)
}