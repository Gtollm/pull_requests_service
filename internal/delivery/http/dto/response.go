package dto

type TeamResponse struct {
	Team TeamDTO `json:"team"`
}

type UserResponse struct {
	User UserDTO `json:"user"`
}

type PullRequestResponse struct {
	PullRequest PullRequestDTO `json:"pr"`
}

type ReassignResponse struct {
	PullRequest PullRequestDTO `json:"pr"`
	ReplacedBy  string         `json:"replaced_by"`
}

type UserReviewsResponse struct {
	UserID       string                `json:"user_id"`
	PullRequests []PullRequestShortDTO `json:"pull_requests"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}