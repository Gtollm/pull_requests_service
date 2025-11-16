package postrgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"pull-request-review/internal/domain/model"
	"pull-request-review/internal/domain/ports/repository"
	"pull-request-review/internal/domain/rules"
	"time"
)

type ReviewAssignmentRepository struct {
	pool *pgxpool.Pool
}

func NewReviewAssignmentRepository(pool *pgxpool.Pool) repository.ReviewAssignmentRepository {
	return &ReviewAssignmentRepository{pool: pool}
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

	tx, err := r.pool.Begin(ctx)
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

	rows, err := r.pool.Query(ctx, query, pullRequestID)
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
	err := r.pool.QueryRow(ctx, query, pullRequestID, reviewerID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *ReviewAssignmentRepository) GetReviewers(ctx context.Context, pullRequestID model.PullRequestID) (
	[]model.User, error,
) {
	query := `
SELECT pull_request_id, user_id, assigned_at FROM review_assignments
WHERE pull_request_id = $1`
	rows, err := r.pool.Query(ctx, query, pullRequestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user)
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
	tx, err := r.pool.Begin(ctx)
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

	return tx.Commit(ctx)
}