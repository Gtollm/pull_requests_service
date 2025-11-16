DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'pr_status') THEN
            CREATE TYPE pull_request_status AS ENUM ('OPEN', 'MERGED');
        END IF;
    END
$$;


CREATE TABLE pull_requests
(
    pull_request_id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    author_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    status pull_request_status NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
)