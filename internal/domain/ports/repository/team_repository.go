package repository

import (
	"context"
	"pull-request-review/internal/domain/model"
)

type TeamRepository interface {
	Create(ctx context.Context, team *model.Team) error
	Update(ctx context.Context, team *model.Team) error
	GetByID(ctx context.Context, ID model.TeamID) (*model.Team, error)
	Exists(ctx context.Context, ID model.TeamID) (bool, error)
}