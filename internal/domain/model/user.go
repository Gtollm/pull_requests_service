package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `db:"user_id"`
	Username  string    `db:"username"`
	TeamID    uuid.UUID `db:"team_id"`
	IsActive  bool      `db:"is_active"`
	createdAt time.Time `db:"created_at"`
	updatedAt time.Time `db:"updated_at"`
}