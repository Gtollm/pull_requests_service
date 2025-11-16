package service

import (
	"context"
	"pull-request-review/internal/domain/model"
)

type TeamService interface {
	CreateTeam(ctx context.Context, team *model.Team) error
	GetTeam(ctx context.Context, ID model.TeamID) (*model.Team, []model.User, error)
	CreateTeamWithMembers(ctx context.Context, teamName string, members []model.User) error
	GetTeamWithMembers(ctx context.Context, teamName string) (*model.Team, []model.User, error)
	BulkDeactivateTeam(ctx context.Context, ID model.TeamID) error
}