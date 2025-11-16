package postrgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"pull-request-review/internal/domain/model"
	"pull-request-review/internal/domain/ports/repository"
	"pull-request-review/internal/domain/rules"
)

type UserRepositoryPgx struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) repository.UserRepository {
	return &UserRepositoryPgx{pool: pool}
}

func (r *UserRepositoryPgx) Insert(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (user_id, username, team_id, is_active, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (user_id) DO NOTHING 
`

	_, err := r.pool.Exec(
		ctx, query,
		user.ID,
		user.Username,
		user.TeamID,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryPgx) Update(ctx context.Context, user *model.User) error {
	query := `UPDATE users
SET username = $1, team_id = $2, is_active = $3, updated_at = $4
WHERE user_id = $5
`

	result, err := r.pool.Exec(
		ctx, query,
		user.Username,
		user.TeamID,
		user.IsActive,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return rules.ErrUserNotFound
	}

	return nil
}

func (r *UserRepositoryPgx) UpdateActivity(ctx context.Context, ID model.UserID, isActive bool) error {
	query := `
UPDATE users
SET is_active = $1, updated_at = now()
WHERE user_id = $2`

	result, err := r.pool.Exec(ctx, query, isActive, ID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return rules.ErrUserNotFound
	}

	return nil
}

func (r *UserRepositoryPgx) GetByID(ctx context.Context, ID model.UserID) (*model.User, error) {
	query := `
SELECT user_id, username, team_id, is_active, created_at, updated_at FROM users WHERE user_id = $1
`
	var user model.User

	err := r.pool.QueryRow(ctx, query, ID).Scan(
		&user.ID,
		&user.Username,
		&user.TeamID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, rules.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryPgx) GetByTeam(ctx context.Context, teamID model.TeamID) ([]model.User, error) {
	query := `
SELECT user_id, username, team_id, is_active, created_at, updated_at FROM users WHERE team_id = $1
`

	rows, err := r.pool.Query(ctx, query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []model.User

	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.TeamID,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
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

func (r *UserRepositoryPgx) Exists(ctx context.Context, ID model.UserID) (bool, error) {
	query := `
SELECT EXISTS(SELECT 1
FROM users
WHERE user_id = $1)
`

	var exists bool
	err := r.pool.QueryRow(ctx, query, ID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UserRepositoryPgx) GetActiveByTeamExcluding(
	ctx context.Context, teamID model.TeamID, excludedUserIDs []model.UserID,
) ([]model.User, error) {
	query := `
SELECT user_id, username, team_id, is_active, created_at, updated_at FROM users WHERE team_id = $1 AND is_active = true AND user_id != ALL($2)
`

	rows, err := r.pool.Query(ctx, query, teamID, excludedUserIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []model.User

	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.TeamID,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
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