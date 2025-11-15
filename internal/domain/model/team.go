package model

import (
	"github.com/google/uuid"
	"time"
)

type Team struct {
	TeamID    uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}