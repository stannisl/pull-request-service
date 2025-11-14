CREATE TABLE IF NOT EXISTS teams (
   name VARCHAR(255) PRIMARY KEY,
   created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
   id VARCHAR(255) PRIMARY KEY,
   username VARCHAR(255) NOT NULL,
   team_name VARCHAR(255) NOT NULL REFERENCES teams(name) ON DELETE CASCADE,
   is_active BOOLEAN NOT NULL DEFAULT true,
   created_at TIMESTAMP NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pull_requests (
   id VARCHAR(255) PRIMARY KEY,
   name VARCHAR(255) NOT NULL,
   author_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
   status VARCHAR(20) NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'MERGED')),
   need_more_reviewers BOOLEAN NOT NULL DEFAULT false,
   created_at TIMESTAMP NOT NULL DEFAULT NOW(),
   merged_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pull_request_reviewers (
    pull_request_id VARCHAR(255) NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    reviewer_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    PRIMARY KEY (pull_request_id, reviewer_id)
);

CREATE INDEX IF NOT EXISTS idx_users_team_name ON users(team_name);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
CREATE INDEX IF NOT EXISTS idx_users_team_active ON users(team_name, is_active);
CREATE INDEX IF NOT EXISTS idx_pull_requests_author_id ON pull_requests(author_id);
CREATE INDEX IF NOT EXISTS idx_pull_requests_status ON pull_requests(status);
CREATE INDEX IF NOT EXISTS idx_pull_request_reviewers_reviewer_id ON pull_request_reviewers(reviewer_id);
CREATE INDEX IF NOT EXISTS idx_pull_request_reviewers_pr_id ON pull_request_reviewers(pull_request_id);