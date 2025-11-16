package model

import (
	"time"
)

type ReviewAssignment struct {
	PullRequestID PullRequestID `db:"pull_request_id"`
	ReviewerID    UserID        `db:"reviewer_id"`
	AssignedAt    time.Time     `db:"assigned_at"`
}