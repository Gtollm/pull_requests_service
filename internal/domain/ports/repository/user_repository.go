package repository

import (
	"context"
	"pull-request-review/internal/domain/model"
)

type UserRepository interface {
	Insert(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Upsert(ctx context.Context, user *model.User) error
	UpdateActivity(ctx context.Context, ID model.UserID, isActive bool) error
	GetByID(ctx context.Context, ID model.UserID) (*model.User, error)
	GetByTeam(ctx context.Context, teamID model.TeamID) ([]model.User, error)
	Exists(ctx context.Context, ID model.UserID) (bool, error)
	GetActiveByTeamExcluding(ctx context.Context, teamID model.TeamID, excludedUserIDs []model.UserID) (
		[]model.User, error,
	)
}