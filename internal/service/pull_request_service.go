package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"pull-request-review/internal/domain/model"
	"pull-request-review/internal/domain/ports/repository"
	"pull-request-review/internal/domain/rules"
	"pull-request-review/internal/infrastructure/adapters/logger"
)

type PullRequestService struct {
	pullRequestRepo      repository.PullRequestRepository
	userRepo             repository.UserRepository
	reviewAssignmentRepo repository.ReviewAssignmentRepository
	logger               logger.Logger
	maxReviewersCount    int
}

func NewPullRequestService(
	pullRequestRepo repository.PullRequestRepository,
	userRepo repository.UserRepository,
	reviewAssignmentRepo repository.ReviewAssignmentRepository,
	logger logger.Logger,
	maxReviewersCount int,
) *PullRequestService {
	return &PullRequestService{
		pullRequestRepo:      pullRequestRepo,
		userRepo:             userRepo,
		reviewAssignmentRepo: reviewAssignmentRepo,
		logger:               logger,
		maxReviewersCount:    maxReviewersCount,
	}
}

func (s *PullRequestService) CreatePullRequest(
	ctx context.Context,
	pullRequest *model.PullRequest,
) (*model.PullRequest, []model.UserID, error) {
	exists, err := s.pullRequestRepo.Exists(ctx, pullRequest.PullRequestID)
	if err != nil {
		s.logger.Error(err, "failed to check PR existence")
		return nil, nil, err
	}
	if exists {
		return nil, nil, rules.ErrPullRequestExists
	}

	author, err := s.userRepo.GetByID(ctx, pullRequest.AuthorID)
	if err != nil {
		s.logger.Error(err, "failed to get author")
		return nil, nil, err
	}

	pullRequest.Status = model.PRStatusOpen
	pullRequest.CreatedAt = time.Now()
	err = s.pullRequestRepo.Create(ctx, pullRequest)
	if err != nil {
		s.logger.Error(err, "failed to create pull request")
		return nil, nil, err
	}

	reviewerIDs, err := s.assignInitialReviewers(ctx, pullRequest, uuid.UUID(author.TeamID))
	if err != nil {
		s.logger.Error(err, "failed to assign reviewers")
		return nil, nil, err
	}

	return pullRequest, reviewerIDs, nil
}

func (s *PullRequestService) GetPullRequest(ctx context.Context, ID model.PullRequestID) (*model.PullRequest, error) {
	pullRequest, err := s.pullRequestRepo.GetByID(ctx, ID)
	if err != nil {
		s.logger.Error(err, "cannot get pull request")
		return nil, err
	}
	return pullRequest, nil
}

func (s *PullRequestService) GetPullRequestReviewers(ctx context.Context, ID model.PullRequestID) (
	[]model.UserID, error,
) {
	reviewers, err := s.reviewAssignmentRepo.GetReviewers(ctx, ID)
	if err != nil {
		s.logger.Error(err, "failed to get reviewers for pull request")
		return nil, err
	}

	reviewerIDs := make([]model.UserID, len(reviewers))
	for i, reviewer := range reviewers {
		reviewerIDs[i] = reviewer.ID
	}

	return reviewerIDs, nil
}

func (s *PullRequestService) GetUserReviews(ctx context.Context, userID model.UserID) ([]model.PullRequest, error) {
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Error(err, "failed to get user")
		return nil, err
	}

	pullRequests, err := s.pullRequestRepo.GetByReviewer(ctx, userID)
	if err != nil {
		s.logger.Error(err, "cannot get pull requests for reviewer")
		return nil, err
	}
	return pullRequests, nil
}

func (s *PullRequestService) ReassignPullRequest(
	ctx context.Context,
	ID model.PullRequestID,
	oldReviewerID model.UserID,
) (*model.PullRequest, model.UserID, error) {
	pullRequest, err := s.pullRequestRepo.GetByID(ctx, ID)
	if err != nil {
		s.logger.Error(err, "failed to get pull request")
		return nil, model.UserID(uuid.Nil), err
	}

	if pullRequest.Status == model.PRStatusMerged {
		return nil, model.UserID(uuid.Nil), rules.ErrPullRequestMerged
	}

	isAssigned, err := s.reviewAssignmentRepo.Exists(ctx, ID, oldReviewerID)
	if err != nil {
		s.logger.Error(err, "failed to check reviewer assignment")
		return nil, model.UserID(uuid.Nil), err
	}
	if !isAssigned {
		return nil, model.UserID(uuid.Nil), rules.ErrNotAssigned
	}

	newReviewerID, err := s.reassignReviewer(ctx, pullRequest, oldReviewerID)
	if err != nil {
		s.logger.Error(err, "failed to reassign reviewer")
		return nil, model.UserID(uuid.Nil), err
	}

	updatedPR, err := s.pullRequestRepo.GetByID(ctx, ID)
	if err != nil {
		s.logger.Error(err, "failed to get updated pull request")
		return nil, model.UserID(uuid.Nil), err
	}

	return updatedPR, newReviewerID, nil
}

func (s *PullRequestService) MergePullRequest(ctx context.Context, ID model.PullRequestID) (*model.PullRequest, error) {
	pullRequest, err := s.pullRequestRepo.GetByID(ctx, ID)
	if err != nil {
		s.logger.Error(err, "failed to get pull request")
		return nil, err
	}

	if pullRequest.Status == model.PRStatusMerged {
		return pullRequest, nil
	}

	now := time.Now()
	err = s.pullRequestRepo.UpdateStatus(ctx, ID, model.PRStatusMerged, now)
	if err != nil {
		s.logger.Error(err, "failed to update pull request status")
		return nil, err
	}

	updatedPR, err := s.pullRequestRepo.GetByID(ctx, ID)
	if err != nil {
		s.logger.Error(err, "failed to get updated pull request")
		return nil, err
	}

	return updatedPR, nil
}

func (s *PullRequestService) assignInitialReviewers(
	ctx context.Context,
	pr *model.PullRequest,
	authorTeamID uuid.UUID,
) ([]model.UserID, error) {
	teamID := model.TeamID(authorTeamID)
	excludedUserIDs := []model.UserID{pr.AuthorID}

	candidates, err := s.userRepo.GetActiveByTeamExcluding(ctx, teamID, excludedUserIDs)
	if err != nil {
		s.logger.Error(err, "failed to get active team members for reviewer assignment")
		return nil, err
	}

	if len(candidates) == 0 {
		return []model.UserID{}, nil
	}

	selectedReviewers := s.randomSelectReviewers(candidates, s.maxReviewersCount)

	reviewerIDs := make([]model.UserID, len(selectedReviewers))
	for i, reviewer := range selectedReviewers {
		reviewerIDs[i] = reviewer.ID
	}

	if len(reviewerIDs) > 0 {
		err = s.reviewAssignmentRepo.AssignReviewers(ctx, pr.PullRequestID, reviewerIDs)
		if err != nil {
			s.logger.Error(err, "failed to assign reviewers to pull request")
			return nil, err
		}
	}

	return reviewerIDs, nil
}

func (s *PullRequestService) reassignReviewer(
	ctx context.Context,
	pr *model.PullRequest,
	oldReviewerID model.UserID,
) (model.UserID, error) {
	oldReviewer, err := s.userRepo.GetByID(ctx, oldReviewerID)
	if err != nil {
		s.logger.Error(err, "failed to get old reviewer")
		return model.UserID(uuid.Nil), err
	}

	currentReviewers, err := s.reviewAssignmentRepo.GetReviewers(ctx, pr.PullRequestID)
	if err != nil {
		s.logger.Error(err, "failed to get current reviewers")
		return model.UserID(uuid.Nil), err
	}

	excludedUserIDs := []model.UserID{pr.AuthorID}
	for _, reviewer := range currentReviewers {
		excludedUserIDs = append(excludedUserIDs, reviewer.ID)
	}

	candidates, err := s.userRepo.GetActiveByTeamExcluding(ctx, model.TeamID(oldReviewer.TeamID), excludedUserIDs)
	if err != nil {
		s.logger.Error(err, "failed to get candidate reviewers for reassignment")
		return model.UserID(uuid.Nil), err
	}

	if len(candidates) == 0 {
		return model.UserID(uuid.Nil), rules.ErrNoCandidates
	}

	selectedReviewers := s.randomSelectReviewers(candidates, 1)
	if len(selectedReviewers) == 0 {
		return model.UserID(uuid.Nil), rules.ErrNoCandidates
	}

	newReviewerID := selectedReviewers[0].ID

	err = s.reviewAssignmentRepo.ReplaceReviewer(ctx, pr.PullRequestID, oldReviewerID, newReviewerID)
	if err != nil {
		s.logger.Error(err, "failed to replace reviewer")
		return model.UserID(uuid.Nil), err
	}

	return newReviewerID, nil
}

func (s *PullRequestService) randomSelectReviewers(candidates []model.User, maxCount int) []model.User {
	if len(candidates) == 0 {
		return []model.User{}
	}

	selectCount := min(maxCount, len(candidates))

	shuffled := make([]model.User, len(candidates))
	copy(shuffled, candidates)

	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled[:selectCount]
}