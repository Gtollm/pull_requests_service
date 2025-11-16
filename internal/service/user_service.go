package service

import (
	"context"

	"pull-request-review/internal/domain/model"
	"pull-request-review/internal/domain/ports/repository"
	"pull-request-review/internal/domain/rules"
	"pull-request-review/internal/infrastructure/adapters/logger"
)

type UserService struct {
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
	logger   logger.Logger
}

func NewUserService(
	userRepo repository.UserRepository, teamRepo repository.TeamRepository, logger logger.Logger,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		teamRepo: teamRepo,
		logger:   logger,
	}
}

func (s *UserService) GetUser(ctx context.Context, ID model.UserID) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, ID)
	if err != nil {
		s.logger.Error(err, "cannot get user")
		return nil, err
	}
	return user, nil
}

func (s *UserService) SetActive(ctx context.Context, ID model.UserID, active bool) (*model.User, error) {
	err := s.userRepo.UpdateActivity(ctx, ID, active)
	if err != nil {
		s.logger.Error(err, "cannot update user activity")
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, ID)
	if err != nil {
		s.logger.Error(err, "cannot get updated user")
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserWithTeamName(ctx context.Context, ID model.UserID) (*model.User, string, error) {
	user, err := s.userRepo.GetByID(ctx, ID)
	if err != nil {
		s.logger.Error(err, "cannot get user")
		return nil, "", rules.ErrUserNotFound
	}

	team, err := s.teamRepo.GetByID(ctx, model.TeamID(user.TeamID))
	if err != nil {
		s.logger.Error(err, "cannot get team for user")
		return nil, "", err
	}

	return user, team.Name, nil
}