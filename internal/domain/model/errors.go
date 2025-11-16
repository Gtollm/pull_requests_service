package model

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrUserExists          = errors.New("user already exists")
	ErrTeamExists          = errors.New("team already exists")
	ErrPullRequestExists   = errors.New("pull request already exists")
	ErrPullRequestMerged   = errors.New("pull request merged")
	ErrNotAssigned         = errors.New("review not assigned to pull request")
	ErrNoCandidates        = errors.New("no active users to replace reviewer")
	ErrTeamNotFound        = errors.New("team not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrPullRequestNotFound = errors.New("pull request not found")
)