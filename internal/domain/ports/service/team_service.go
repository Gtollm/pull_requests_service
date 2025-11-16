package service

import (
	"context"
	"pull-request-review/internal/domain/model"
)

type TeamService interface {
	CreateTeam(ctx context.Context, team *model.Team) error
	GetTeam(ctx context.Context, ID model.TeamID) (*model.Team, error)
}