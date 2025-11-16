package postrgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"pr-review-service/internal/domain/model"
	"pr-review-service/internal/domain/ports/repository"
	"pr-review-service/internal/domain/rules"
	"time"
)

type PullRequestRepositoryPgx struct {
	pool *pgxpool.Pool
}

func NewPullRequestRepositoryPgx(pool *pgxpool.Pool) repository.PullRequestRepository {
	return &PullRequestRepositoryPgx{pool: pool}
}

func (r *PullRequestRepositoryPgx) Create(ctx context.Context, pullRequest *model.PullRequest) error {
	query := `
INSERT INTO pull_requests (pull_request_id, name, author_id, status, created_at, merged_at)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT DO NOTHING;
`
	_, err := r.pool.Exec(
		ctx, query,
		pullRequest.PullRequestID,
		pullRequest.Name,
		pullRequest.AuthorID,
		pullRequest.Status,
		pullRequest.CreatedAt,
		pullRequest.MergedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *PullRequestRepositoryPgx) GetByID(ctx context.Context, ID model.PullRequestID) (*model.PullRequest, error) {
	query := `
SELECT pull_request_id, name, author_id, status, created_at, merged_at
FROM pull_requests
WHERE pull_request_id = $1
	`

	var pr model.PullRequest
	err := r.pool.QueryRow(ctx, query, ID).Scan(
		&pr.PullRequestID,
		&pr.Name,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, rules.ErrNotFound
		}
		return nil, err
	}

	return &pr, nil
}

func (r *PullRequestRepositoryPgx) Exists(ctx context.Context, ID model.PullRequestID) (bool, error) {
	query := `
SELECT EXISTS (SELECT 1 FROM pull_requests WHERE pull_request_id = $1)
`
	var exists bool
	err := r.pool.QueryRow(ctx, query, ID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *PullRequestRepositoryPgx) UpdateStatus(
	ctx context.Context, ID model.PullRequestID, status model.PullRequestStatus, mergedAt time.Time,
) error {
	query := `
UPDATE pull_requests
SET status = $1, merged_at = $2
WHERE pull_request_id = $3
	`

	result, err := r.pool.Exec(ctx, query, status, mergedAt, ID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return rules.ErrPullRequestNotFound
	}

	return nil
}

func (r *PullRequestRepositoryPgx) GetByReviewer(ctx context.Context, ID model.UserID) ([]model.PullRequest, error) {
	query := `
SELECT DISTINCT p.pull_request_id, p.name, p.author_id, p.status, p.created_at, p.merged_at
FROM pull_requests p
INNER JOIN review_assignments ra ON p.pull_request_id = ra.pull_request_id
WHERE ra.user_id = $1
`

	rows, err := r.pool.Query(ctx, query, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pullRequests []model.PullRequest

	for rows.Next() {
		var pr model.PullRequest
		err := rows.Scan(
			&pr.PullRequestID,
			&pr.Name,
			&pr.AuthorID,
			&pr.Status,
			&pr.CreatedAt,
			&pr.MergedAt,
		)
		if err != nil {
			return nil, err
		}
		pullRequests = append(pullRequests, pr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pullRequests, nil
}