package dto

import (
	"github.com/google/uuid"
	"pull-request-review/internal/domain/model"
)

type TeamDTO struct {
	TeamName string          `json:"team_name"`
	Members  []TeamMemberDTO `json:"members"`
}

type TeamMemberDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

func TeamToDTO(team *model.Team, users []model.User) TeamDTO {
	return TeamDTO{
		TeamName: team.Name,
		Members:  UsersToTeamMemberDTOs(users),
	}
}

func UsersToTeamMemberDTOs(users []model.User) []TeamMemberDTO {
	members := make([]TeamMemberDTO, len(users))
	for i, user := range users {
		members[i] = UserToTeamMemberDTO(&user)
	}
	return members
}

func UserToTeamMemberDTO(user *model.User) TeamMemberDTO {
	return TeamMemberDTO{
		UserID:   uuid.UUID(user.ID).String(),
		Username: user.Username,
		IsActive: user.IsActive,
	}
}

func TeamMemberToUser(member TeamMember) (model.User, error) {
	userUUID, err := uuid.Parse(member.UserID)
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		ID:       model.UserID(userUUID),
		Username: member.Username,
		IsActive: member.IsActive,
	}, nil
}

func TeamMembersToUsers(members []TeamMember) ([]model.User, error) {
	users := make([]model.User, len(members))
	for i, member := range members {
		user, err := TeamMemberToUser(member)
		if err != nil {
			return nil, err
		}
		users[i] = user
	}
	return users, nil
}