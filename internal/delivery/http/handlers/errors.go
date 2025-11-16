package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"pull-request-review/internal/delivery/http/dto"
	"pull-request-review/internal/domain/rules"
)

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func getErrorCode(err error) string {
	switch {
	case errors.Is(err, rules.ErrTeamExists):
		return "TEAM_EXISTS"
	case errors.Is(err, rules.ErrPullRequestExists):
		return "PR_EXISTS"
	case errors.Is(err, rules.ErrPullRequestMerged):
		return "PR_MERGED"
	case errors.Is(err, rules.ErrNotAssigned):
		return "NOT_ASSIGNED"
	case errors.Is(err, rules.ErrNoCandidates):
		return "NO_CANDIDATE"
	case errors.Is(err, rules.ErrNotFound),
		errors.Is(err, rules.ErrTeamNotFound),
		errors.Is(err, rules.ErrUserNotFound),
		errors.Is(err, rules.ErrPullRequestNotFound):
		return "NOT_FOUND"
	default:
		return "INTERNAL_ERROR"
	}
}

func getHTTPStatus(err error) int {
	switch {
	case errors.Is(err, rules.ErrTeamExists):
		return http.StatusBadRequest
	case errors.Is(err, rules.ErrPullRequestExists):
		return http.StatusConflict
	case errors.Is(err, rules.ErrPullRequestMerged),
		errors.Is(err, rules.ErrNotAssigned),
		errors.Is(err, rules.ErrNoCandidates):
		return http.StatusConflict
	case errors.Is(err, rules.ErrNotFound),
		errors.Is(err, rules.ErrTeamNotFound),
		errors.Is(err, rules.ErrUserNotFound),
		errors.Is(err, rules.ErrPullRequestNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func WriteError(w http.ResponseWriter, err error) {
	code := getErrorCode(err)
	status := getHTTPStatus(err)

	response := dto.ErrorResponse{
		Error: dto.ErrorDetail{
			Code:    code,
			Message: err.Error(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}