package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"pull-request-review/internal/delivery/http/dto"
	"pull-request-review/internal/domain/model"
	"pull-request-review/internal/domain/ports/service"
)

type UserHandler struct {
	userService        service.UserService
	pullRequestService service.PullRequestService
}

func NewUserHandler(userService service.UserService, pullRequestService service.PullRequestService) *UserHandler {
	return &UserHandler{
		userService:        userService,
		pullRequestService: pullRequestService,
	}
}

// SetIsActive handles POST /users/setIsActive
func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var req dto.SetActiveRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, err)
		return
	}

	if strings.TrimSpace(req.UserID) == "" {
		WriteError(w, &ValidationError{Message: "user_id is required"})
		return
	}

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		WriteError(w, &ValidationError{Message: "invalid user_id format"})
		return
	}
	userID := model.UserID(userUUID)

	_, err = h.userService.SetActive(r.Context(), userID, req.IsActive)
	if err != nil {
		WriteError(w, err)
		return
	}

	user, teamName, err := h.userService.GetUserWithTeamName(r.Context(), userID)
	if err != nil {
		WriteError(w, err)
		return
	}

	response := dto.UserResponse{
		User: dto.UserToDTO(user, teamName),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

// GetReviews handles GET /users/getReview
func (h *UserHandler) GetReviews(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")

	if strings.TrimSpace(userIDStr) == "" {
		WriteError(w, &ValidationError{Message: "user_id query parameter is required"})
		return
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		WriteError(w, &ValidationError{Message: "invalid user_id format"})
		return
	}
	userID := model.UserID(userUUID)

	pullRequests, err := h.pullRequestService.GetUserReviews(r.Context(), userID)
	if err != nil {
		WriteError(w, err)
		return
	}

	response := dto.UserReviewsResponse{
		UserID:       userIDStr,
		PullRequests: dto.PullRequestsToShortDTOs(pullRequests),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}