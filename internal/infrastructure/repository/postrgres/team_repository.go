package postrgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"pr-review-service/internal/domain/model"
	"pr-review-service/internal/domain/ports/repository"
	"pr-review-service/internal/domain/rules"
)

type TeamRepositoryPgx struct {
	pool *pgxpool.Pool
}

func NewTeamRepository(pool *pgxpool.Pool) repository.TeamRepository {
	return &TeamRepositoryPgx{pool: pool}
}

func (r *TeamRepositoryPgx) Create(ctx context.Context, team *model.Team) error {
	query := `
INSERT INTO teams (team_id, name, created_at)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING
`

	_, err := r.pool.Exec(ctx, query, team.TeamID, team.Name, team.CreatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (r *TeamRepositoryPgx) Update(ctx context.Context, team *model.Team) error {
	query := `
UPDATE teams
SET name = $1 
WHERE team_id = $2
`

	_, err := r.pool.Exec(ctx, query, team.Name, team.TeamID)

	if err != nil {
		return err
	}
	return nil
}

func (r *TeamRepositoryPgx) GetByID(ctx context.Context, ID model.TeamID) (*model.Team, error) {
	query := `
SELECT team_id, name, created_at FROM teams
WHERE team_id = $1
`
	var team model.Team
	err := r.pool.QueryRow(ctx, query, ID).Scan(&team.TeamID, &team.Name, &team.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, rules.ErrTeamNotFound
		}

		return nil, err
	}

	return &team, nil
}

func (r *TeamRepositoryPgx) Exists(ctx context.Context, ID model.TeamID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM teams WHERE team_id = $1)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, ID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}