CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    username TEXT NOT NULL,
    team_id UUID NOT NULL REFERENCES teams(team_id) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
)