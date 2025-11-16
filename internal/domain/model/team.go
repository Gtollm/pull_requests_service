package model

import (
	"github.com/google/uuid"
	"time"
)

type TeamID uuid.UUID

type Team struct {
	TeamID    TeamID    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}