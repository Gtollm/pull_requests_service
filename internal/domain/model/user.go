package model

import (
	"github.com/google/uuid"
	"time"
)

type UserID uuid.UUID

type User struct {
	ID        UserID    `db:"user_id"`
	Username  string    `db:"username"`
	TeamID    uuid.UUID `db:"team_id"`
	IsActive  bool      `db:"is_active"`
	createdAt time.Time `db:"created_at"`
	updatedAt time.Time `db:"updated_at"`
}