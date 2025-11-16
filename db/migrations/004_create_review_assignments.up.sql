CREATE TABLE review_assignments(
    pull_request_id UUID NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (pull_request_id, user_id)
);

CREATE INDEX idx_reviewer_assignments_user_id ON review_assignments(user_id);