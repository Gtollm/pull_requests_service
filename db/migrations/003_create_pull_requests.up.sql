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
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    merged_at TIMESTAMPTZ
);

CREATE INDEX idx_pull_requests_author_id ON pull_requests(author_id);
CREATE INDEX idx_pull_requests_status ON pull_requests(status);