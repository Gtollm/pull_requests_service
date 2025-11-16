package dto

import (
	"github.com/google/uuid"
	"pull-request-review/internal/domain/model"
	"time"
)

type PullRequestDTO struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	CreatedAt         *string  `json:"createdAt,omitempty"`
	MergedAt          *string  `json:"mergedAt,omitempty"`
}

type PullRequestShortDTO struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

func PullRequestToDTO(pr *model.PullRequest, reviewerIDs []string) PullRequestDTO {
	dto := PullRequestDTO{
		PullRequestID:     uuid.UUID(pr.PullRequestID).String(),
		PullRequestName:   pr.Name,
		AuthorID:          uuid.UUID(pr.AuthorID).String(),
		Status:            string(pr.Status),
		AssignedReviewers: reviewerIDs,
	}

	if !pr.CreatedAt.IsZero() {
		createdAt := pr.CreatedAt.Format(time.RFC3339)
		dto.CreatedAt = &createdAt
	}

	if !pr.MergedAt.IsZero() {
		mergedAt := pr.MergedAt.Format(time.RFC3339)
		dto.MergedAt = &mergedAt
	}

	return dto
}

func PullRequestsToShortDTOs(prs []model.PullRequest) []PullRequestShortDTO {
	dtos := make([]PullRequestShortDTO, len(prs))
	for i, pr := range prs {
		dtos[i] = PullRequestToShortDTO(&pr)
	}
	return dtos
}

func PullRequestToShortDTO(pr *model.PullRequest) PullRequestShortDTO {
	return PullRequestShortDTO{
		PullRequestID:   uuid.UUID(pr.PullRequestID).String(),
		PullRequestName: pr.Name,
		AuthorID:        uuid.UUID(pr.AuthorID).String(),
		Status:          string(pr.Status),
	}
}