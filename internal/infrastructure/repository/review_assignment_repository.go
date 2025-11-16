package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"pull-request-review/internal/domain/model"
	"pull-request-review/internal/domain/ports/repository"
	"pull-request-review/internal/domain/rules"
	"pull-request-review/internal/infrastructure/database"
	"time"
)

type ReviewAssignmentRepository struct {
	database *database.Database
}

func NewReviewAssignmentRepository(database *database.Database) repository.ReviewAssignmentRepository {
	return &ReviewAssignmentRepository{database: database}
}

func (r *ReviewAssignmentRepository) AssignReviewer(
	ctx context.Context, pullRequestID model.PullRequestID, reviewerIDs model.UserID,
) error {
	return r.AssignReviewers(ctx, pullRequestID, []model.UserID{reviewerIDs})
}

func (r *ReviewAssignmentRepository) AssignReviewers(
	ctx context.Context, pullRequestID model.PullRequestID, reviewerIDs []model.UserID,
) error {
	if len(reviewerIDs) == 0 {
		return nil
	}

	tx, err := r.database.GetPool().Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		var e error
		if err == nil {
			e = tx.Commit(ctx)
		} else {
			e = tx.Rollback(ctx)
		}

		if err == nil && e != nil {
			err = fmt.Errorf("finishing transaction: %w", e)
		}
	}()

	query := `
INSERT INTO review_assignments (pull_request_id, user_id, assigned_at)
VALUES ($1, $2, $3)
	`

	now := time.Now()
	for _, userID := range reviewerIDs {
		_, err := tx.Exec(ctx, query, pullRequestID, userID, now)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ReviewAssignmentRepository) GetByReviewer(
	ctx context.Context, pullRequestID model.PullRequestID,
) ([]model.PullRequest, error) {
	query := `
SELECT pull_request_id, user_id, assigned_at FROM review_assignments
WHERE user_id = $1`

	rows, err := r.database.GetPool().Query(ctx, query, pullRequestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pullRequests []model.PullRequest
	for rows.Next() {
		var pullRequest model.PullRequest
		err := rows.Scan(&pullRequest)
		if err != nil {
			return nil, err
		}
		pullRequests = append(pullRequests, pullRequest)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pullRequests, nil
}

func (r *ReviewAssignmentRepository) Exists(
	ctx context.Context, pullRequestID model.PullRequestID, reviewerID model.UserID,
) (bool, error) {
	query := `
SELECT EXISTS (SELECT 1 FROM review_assignments WHERE pull_request_id = $1 AND user_id = $2)
`
	var exists bool
	err := r.database.GetPool().QueryRow(ctx, query, pullRequestID, reviewerID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *ReviewAssignmentRepository) GetReviewers(ctx context.Context, pullRequestID model.PullRequestID) (
	[]model.User, error,
) {
	query := `
SELECT u.user_id, u.username, u.team_id, u.is_active, u.created_at, u.updated_at 
FROM users u 
INNER JOIN review_assignments ra ON u.user_id = ra.user_id 
WHERE ra.pull_request_id = $1`
	rows, err := r.database.GetPool().Query(ctx, query, pullRequestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Username, &user.TeamID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *ReviewAssignmentRepository) ReplaceReviewer(
	ctx context.Context, pullRequestID model.PullRequestID, oldReviewerID model.UserID, newReviewerID model.UserID,
) error {
	tx, err := r.database.GetPool().Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		var e error
		if err == nil {
			e = tx.Commit(ctx)
		} else {
			e = tx.Rollback(ctx)
		}

		if err == nil && e != nil {
			err = fmt.Errorf("finishing transaction: %w", e)
		}
	}()
	deleteQuery := `
DELETE FROM review_assignments 
WHERE pull_request_id = $1 AND user_id = $2
	`

	result, err := tx.Exec(ctx, deleteQuery, pullRequestID, oldReviewerID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return rules.ErrNotAssigned
	}

	insertQuery := `
INSERT INTO review_assignments (pull_request_id, user_id, assigned_at)
VALUES ($1, $2, $3)
	`

	_, err = tx.Exec(ctx, insertQuery, pullRequestID, newReviewerID, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (r *ReviewAssignmentRepository) GetAssignmentCounts(ctx context.Context) (map[string]int, error) {
	query := `
SELECT user_id, COUNT(*) as count
FROM review_assignments
GROUP BY user_id
	`

	rows, err := r.database.GetPool().Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var userID model.UserID
		var count int
		err := rows.Scan(&userID, &count)
		if err != nil {
			return nil, err
		}
		counts[uuid.UUID(userID).String()] = count
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return counts, nil
}