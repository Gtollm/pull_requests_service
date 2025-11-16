package model

import (
	"github.com/google/uuid"
	"time"
)

type ReviewAssignmentID uuid.UUID

type ReviewAssignment struct {
	ID         ReviewAssignmentID `db:"id"`
	ReviewerID UserID             `db:"reviewer_id"`
	AssignedAt time.Time          `db:"assigned_at"`
}