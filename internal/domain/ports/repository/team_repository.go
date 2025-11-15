package repository

import (
	"context"
	"github.com/google/uuid"
	"pr-review-service/internal/domain/model"
)

type TeamRepository interface {
	Create(ctx context.Context, team *model.Team) error
	Update(ctx context.Context, team *model.Team) error
	GetByID(ctx context.Context, ID uuid.UUID) error
	Exists(ctx context.Context, ID uuid.UUID) (bool, error)
}