package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"pull-request-review/internal/domain/model"
	"pull-request-review/internal/domain/ports/repository"
	"pull-request-review/internal/domain/rules"
	"pull-request-review/internal/infrastructure/adapters/logger"
)

type TeamService struct {
	teamRepo repository.TeamRepository
	userRepo repository.UserRepository
	logger   logger.Logger
}

func NewTeamService(
	teamRepo repository.TeamRepository,
	userRepo repository.UserRepository,
	logger logger.Logger,
) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, team *model.Team) error {
	err := s.teamRepo.Create(ctx, team)
	if err != nil {
		s.logger.Error(err, "cannot create team")
		return err
	}
	return nil
}

func (s *TeamService) GetTeam(ctx context.Context, ID model.TeamID) (*model.Team, []model.User, error) {
	team, err := s.teamRepo.GetByID(ctx, ID)
	if err != nil {
		s.logger.Error(err, "cannot get team")
		return nil, nil, err
	}

	users, err := s.userRepo.GetByTeam(ctx, ID)
	if err != nil {
		s.logger.Error(err, "cannot get team members")
		return nil, nil, err
	}

	return team, users, nil
}

func (s *TeamService) CreateTeamWithMembers(ctx context.Context, teamName string, members []model.User) error {
	exists, err := s.teamRepo.ExistsByName(ctx, teamName)
	if err != nil {
		s.logger.Error(err, "cannot check team existence")
		return err
	}
	if exists {
		return rules.ErrTeamExists
	}

	teamID := model.TeamID(uuid.New())
	now := time.Now()
	team := &model.Team{
		TeamID:    teamID,
		Name:      teamName,
		CreatedAt: now,
	}

	for i := range members {
		members[i].TeamID = uuid.UUID(teamID)
		if members[i].CreatedAt.IsZero() {
			members[i].CreatedAt = now
		}
		if members[i].UpdatedAt.IsZero() {
			members[i].UpdatedAt = now
		}
	}

	err = s.teamRepo.CreateWithMembers(ctx, team, members)
	if err != nil {
		s.logger.Error(err, "cannot create team with members")
		return err
	}

	return nil
}

func (s *TeamService) GetTeamWithMembers(ctx context.Context, teamName string) (*model.Team, []model.User, error) {
	team, err := s.teamRepo.GetByName(ctx, teamName)
	if err != nil {
		s.logger.Error(err, "cannot get team by name")
		return nil, nil, err
	}

	users, err := s.userRepo.GetByTeam(ctx, team.TeamID)
	if err != nil {
		s.logger.Error(err, "cannot get team members")
		return nil, nil, err
	}

	return team, users, nil
}

func (s *TeamService) BulkDeactivateTeam(ctx context.Context, ID model.TeamID) error {
	err := s.teamRepo.BulkDeactivateTeam(ctx, ID)
	if err != nil {
		s.logger.Error(err, "cannot bulk deactivate team")
		return err
	}
	return nil
}