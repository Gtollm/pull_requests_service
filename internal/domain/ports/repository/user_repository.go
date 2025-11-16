package repository

import (
	"context"
	"pr-review-service/internal/domain/model"
)

type UserRepository interface {
	Insert(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	UpdateActivity(ctx context.Context, ID model.UserID, isActive bool) error
	GetByID(ctx context.Context, ID model.UserID) (*model.User, error)
	GetByTeam(ctx context.Context, teamID model.TeamID) ([]model.User, error)
	Exists(ctx context.Context, ID model.UserID) (bool, error)
	GetActiveByTeamExcluding(ctx context.Context, teamID model.TeamID) ([]model.User, error)
}