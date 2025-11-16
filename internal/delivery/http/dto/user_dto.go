package dto

import (
	"github.com/google/uuid"
	"pull-request-review/internal/domain/model"
)

type UserDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

func UserToDTO(user *model.User, teamName string) UserDTO {
	return UserDTO{
		UserID:   uuid.UUID(user.ID).String(),
		Username: user.Username,
		TeamName: teamName,
		IsActive: user.IsActive,
	}
}