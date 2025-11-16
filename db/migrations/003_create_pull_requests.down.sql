DROP INDEX IF EXISTS idx_pull_requests_author_id;
DROP INDEX IF EXISTS idx_pull_requests_status;

DROP TABLE IF EXISTS pull_requests;

DROP TYPE IF EXISTS pull_request_status;