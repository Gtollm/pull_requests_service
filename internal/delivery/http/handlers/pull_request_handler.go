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

type PullRequestHandler struct {
	pullRequestService service.PullRequestService
}

func NewPullRequestHandler(
	pullRequestService service.PullRequestService,
) *PullRequestHandler {
	return &PullRequestHandler{
		pullRequestService: pullRequestService,
	}
}

// CreatePullRequest handles POST /pullRequest/create
func (h *PullRequestHandler) CreatePullRequest(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePRRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, err)
		return
	}

	if strings.TrimSpace(req.PullRequestID) == "" {
		WriteError(w, &ValidationError{Message: "pull_request_id is required"})
		return
	}
	if strings.TrimSpace(req.PullRequestName) == "" {
		WriteError(w, &ValidationError{Message: "pull_request_name is required"})
		return
	}
	if strings.TrimSpace(req.AuthorID) == "" {
		WriteError(w, &ValidationError{Message: "author_id is required"})
		return
	}

	prUUID, err := uuid.Parse(req.PullRequestID)
	if err != nil {
		WriteError(w, &ValidationError{Message: "invalid pull_request_id format"})
		return
	}
	prID := model.PullRequestID(prUUID)

	authorUUID, err := uuid.Parse(req.AuthorID)
	if err != nil {
		WriteError(w, &ValidationError{Message: "invalid author_id format"})
		return
	}
	authorID := model.UserID(authorUUID)

	pr := &model.PullRequest{
		PullRequestID: prID,
		Name:          req.PullRequestName,
		AuthorID:      authorID,
		Status:        model.PRStatusOpen,
	}

	createdPR, reviewerIDs, err := h.pullRequestService.CreatePullRequest(r.Context(), pr)
	if err != nil {
		WriteError(w, err)
		return
	}

	reviewerIDStrings := make([]string, len(reviewerIDs))
	for i, id := range reviewerIDs {
		reviewerIDStrings[i] = uuid.UUID(id).String()
	}

	response := dto.PullRequestResponse{
		PullRequest: dto.PullRequestToDTO(createdPR, reviewerIDStrings),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

// MergePullRequest handles POST /pullRequest/merge
func (h *PullRequestHandler) MergePullRequest(w http.ResponseWriter, r *http.Request) {
	var req dto.MergePRRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, err)
		return
	}

	if strings.TrimSpace(req.PullRequestID) == "" {
		WriteError(w, &ValidationError{Message: "pull_request_id is required"})
		return
	}

	prUUID, err := uuid.Parse(req.PullRequestID)
	if err != nil {
		WriteError(w, &ValidationError{Message: "invalid pull_request_id format"})
		return
	}
	prID := model.PullRequestID(prUUID)

	mergedPR, err := h.pullRequestService.MergePullRequest(r.Context(), prID)
	if err != nil {
		WriteError(w, err)
		return
	}

	reviewerIDs, err := h.pullRequestService.GetPullRequestReviewers(r.Context(), prID)
	if err != nil {
		WriteError(w, err)
		return
	}

	reviewerIDStrings := make([]string, len(reviewerIDs))
	for i, id := range reviewerIDs {
		reviewerIDStrings[i] = uuid.UUID(id).String()
	}

	response := dto.PullRequestResponse{
		PullRequest: dto.PullRequestToDTO(mergedPR, reviewerIDStrings),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

// ReassignReviewer handles POST /pullRequest/reassign
func (h *PullRequestHandler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	var req dto.ReassignRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, err)
		return
	}

	if strings.TrimSpace(req.PullRequestID) == "" {
		WriteError(w, &ValidationError{Message: "pull_request_id is required"})
		return
	}
	if strings.TrimSpace(req.OldUserID) == "" {
		WriteError(w, &ValidationError{Message: "old_user_id is required"})
		return
	}

	prUUID, err := uuid.Parse(req.PullRequestID)
	if err != nil {
		WriteError(w, &ValidationError{Message: "invalid pull_request_id format"})
		return
	}
	prID := model.PullRequestID(prUUID)

	oldUserUUID, err := uuid.Parse(req.OldUserID)
	if err != nil {
		WriteError(w, &ValidationError{Message: "invalid old_user_id format"})
		return
	}
	oldUserID := model.UserID(oldUserUUID)

	updatedPR, newReviewerID, err := h.pullRequestService.ReassignPullRequest(r.Context(), prID, oldUserID)
	if err != nil {
		WriteError(w, err)
		return
	}

	reviewerIDs, err := h.pullRequestService.GetPullRequestReviewers(r.Context(), prID)
	if err != nil {
		WriteError(w, err)
		return
	}

	reviewerIDStrings := make([]string, len(reviewerIDs))
	for i, id := range reviewerIDs {
		reviewerIDStrings[i] = uuid.UUID(id).String()
	}

	response := dto.ReassignResponse{
		PullRequest: dto.PullRequestToDTO(updatedPR, reviewerIDStrings),
		ReplacedBy:  uuid.UUID(newReviewerID).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}