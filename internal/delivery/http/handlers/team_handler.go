package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"pull-request-review/internal/delivery/http/dto"
	"pull-request-review/internal/domain/ports/service"
)

type TeamHandler struct {
	teamService service.TeamService
}

func NewTeamHandler(teamService service.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

// AddTeam handles POST /team/add
func (h *TeamHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
	var req dto.TeamRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, err)
		return
	}

	if strings.TrimSpace(req.TeamName) == "" {
		WriteError(w, &ValidationError{Message: "team_name is required"})
		return
	}

	if len(req.Members) == 0 {
		WriteError(w, &ValidationError{Message: "members array cannot be empty"})
		return
	}

	for i, member := range req.Members {
		if strings.TrimSpace(member.UserID) == "" {
			WriteError(w, &ValidationError{Message: fmt.Sprintf("member user_id is required at index %d", i)})
			return
		}
		if strings.TrimSpace(member.Username) == "" {
			WriteError(w, &ValidationError{Message: fmt.Sprintf("member username is required at index %d", i)})
			return
		}
	}

	members, err := dto.TeamMembersToUsers(req.Members)
	if err != nil {
		WriteError(w, &ValidationError{Message: "invalid user_id format in members"})
		return
	}

	err = h.teamService.CreateTeamWithMembers(r.Context(), req.TeamName, members)
	if err != nil {
		WriteError(w, err)
		return
	}

	team, users, err := h.teamService.GetTeamWithMembers(r.Context(), req.TeamName)
	if err != nil {
		WriteError(w, err)
		return
	}

	response := dto.TeamResponse{
		Team: dto.TeamToDTO(team, users),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

// GetTeam handles GET /team/get
func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")

	if strings.TrimSpace(teamName) == "" {
		WriteError(w, &ValidationError{Message: "team_name query parameter is required"})
		return
	}

	team, users, err := h.teamService.GetTeamWithMembers(r.Context(), teamName)
	if err != nil {
		WriteError(w, err)
		return
	}

	response := dto.TeamResponse{
		Team: dto.TeamToDTO(team, users),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}