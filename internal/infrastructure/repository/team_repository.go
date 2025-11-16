package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"pull-request-review/internal/domain/model"
	"pull-request-review/internal/domain/ports/repository"
	"pull-request-review/internal/domain/rules"
	"pull-request-review/internal/infrastructure/database"
)

type TeamRepositoryPgx struct {
	database *database.Database
}

func NewTeamRepository(database *database.Database) repository.TeamRepository {
	return &TeamRepositoryPgx{database: database}
}

func (r *TeamRepositoryPgx) Create(ctx context.Context, team *model.Team) error {
	query := `
INSERT INTO teams (team_id, name, created_at)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING
`

	_, err := r.database.GetPool().Exec(ctx, query, team.TeamID, team.Name, team.CreatedAt)

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

	_, err := r.database.GetPool().Exec(ctx, query, team.Name, team.TeamID)

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
	err := r.database.GetPool().QueryRow(ctx, query, ID).Scan(&team.TeamID, &team.Name, &team.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, rules.ErrTeamNotFound
		}

		return nil, err
	}

	return &team, nil
}

func (r *TeamRepositoryPgx) GetByName(ctx context.Context, name string) (*model.Team, error) {
	query := `
SELECT team_id, name, created_at FROM teams
WHERE name = $1
`
	var team model.Team
	err := r.database.GetPool().QueryRow(ctx, query, name).Scan(&team.TeamID, &team.Name, &team.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, rules.ErrTeamNotFound
		}

		return nil, err
	}

	return &team, nil
}

func (r *TeamRepositoryPgx) ExistsByName(ctx context.Context, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM teams WHERE name = $1)`

	var exists bool
	err := r.database.GetPool().QueryRow(ctx, query, name).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *TeamRepositoryPgx) Exists(ctx context.Context, ID model.TeamID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM teams WHERE team_id = $1)`

	var exists bool
	err := r.database.GetPool().QueryRow(ctx, query, ID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *TeamRepositoryPgx) GetMembers(ctx context.Context, ID model.TeamID) ([]*model.User, error) {
	query := `
SELECT user_id, username, team_id, is_active, created_at, updated_at 
FROM users
WHERE team_id = $1
`
	rows, err := r.database.GetPool().Query(ctx, query, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Username, &user.TeamID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *TeamRepositoryPgx) BulkDeactivateTeam(ctx context.Context, ID model.TeamID) error {
	query := `
UPDATE users
SET is_active = false, updated_at = now()
WHERE team_id = $1 AND is_active = true
`
	_, err := r.database.GetPool().Exec(ctx, query, ID)
	return err
}

func (r *TeamRepositoryPgx) CreateWithMembers(ctx context.Context, team *model.Team, members []model.User) error {
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

	teamQuery := `INSERT INTO teams (team_id, name, created_at) VALUES ($1, $2, $3)`
	_, err = tx.Exec(ctx, teamQuery, team.TeamID, team.Name, team.CreatedAt)
	if err != nil {
		return err
	}

	for _, member := range members {
		upsertQuery := `INSERT INTO users (user_id, username, team_id, is_active, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (user_id) DO UPDATE 
SET username = EXCLUDED.username, 
    team_id = EXCLUDED.team_id, 
    is_active = EXCLUDED.is_active, 
    updated_at = EXCLUDED.updated_at`

		_, err = tx.Exec(
			ctx, upsertQuery, member.ID, member.Username, member.TeamID, member.IsActive, member.CreatedAt,
			member.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	return nil
}