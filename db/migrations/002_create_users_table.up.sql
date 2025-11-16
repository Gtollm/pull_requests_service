CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    username TEXT NOT NULL,
    team_id UUID NOT NULL REFERENCES teams(team_id) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_users_team_id ON users(team_id);
CREATE INDEX idx_users_is_active ON users(is_active);