package model

import (
	"github.com/google/uuid"
	"time"
)

type PullRequestStatus string

const (
	PRStatusOpen   PullRequestStatus = "OPEN"
	PRStatusMerged PullRequestStatus = "MERGED"
)

type PullRequest struct {
	PullRequestID uuid.UUID         `db:"pull_request_id"`
	Name          string            `db:"name"`
	AuthorID      uuid.UUID         `db:"author_id"`
	Status        PullRequestStatus `db:"status"`
	CreatedAt     time.Time         `db:"created_at"`
	MergedAt      time.Time         `db:"merged_at"`
}
