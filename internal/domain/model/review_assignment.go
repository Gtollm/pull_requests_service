package model

import (
	"github.com/google/uuid"
	"time"
)

type ReviewAssignment struct {
	ID         uuid.UUID `db:"id"`
	ReviewerID uuid.UUID `db:"reviewer_id"`
	AssignedAt time.Time `db:"assigned_at"`
}