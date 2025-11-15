package repository

import (
	"context"
	"github.com/google/uuid"
	"pr-review-service/internal/domain/model"
)

type UserRepository interface {
	Insert(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	UpdateActivity(ctx context.Context, ID uuid.UUID, isActive bool) error
	GetByID(ctx context.Context, ID uuid.UUID) (*model.User, error)
	GetByTeam(ctx context.Context, teamID uuid.UUID) ([]model.User, error)
	Exists(ctx context.Context, ID uuid.UUID) (bool, error)
	GetActiveByTeamExcluding(ctx context.Context, teamID uuid.UUID) ([]model.User, error)
}