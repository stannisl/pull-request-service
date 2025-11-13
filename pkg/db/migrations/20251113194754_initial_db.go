package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInitialDb, downInitialDb)
}

func upInitialDb(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, `
		CREATE TABLE teams (
			name VARCHAR(255) PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		CREATE TABLE users (
			id VARCHAR(255) PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			team_name VARCHAR(255) NOT NULL REFERENCES teams(name) ON DELETE CASCADE,
			is_active BOOLEAN NOT NULL DEFAULT true,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		CREATE TABLE pull_requests (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			author_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
			status VARCHAR(20) NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'MERGED')),
			need_more_reviewers BOOLEAN NOT NULL DEFAULT false,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			merged_at TIMESTAMP
		)
	`); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		CREATE TABLE pull_request_reviewers (
			pull_request_id VARCHAR(255) NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
			reviewer_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
			PRIMARY KEY (pull_request_id, reviewer_id)
		)
	`); err != nil {
		return err
	}

	indexes := []string{
		"CREATE INDEX idx_users_team_name ON users(team_name)",
		"CREATE INDEX idx_users_is_active ON users(is_active)",
		"CREATE INDEX idx_users_team_active ON users(team_name, is_active)",
		"CREATE INDEX idx_pull_requests_author_id ON pull_requests(author_id)",
		"CREATE INDEX idx_pull_requests_status ON pull_requests(status)",
		"CREATE INDEX idx_pull_request_reviewers_reviewer_id ON pull_request_reviewers(reviewer_id)",
		"CREATE INDEX idx_pull_request_reviewers_pr_id ON pull_request_reviewers(pull_request_id)",
	}

	for _, indexSQL := range indexes {
		if _, err := tx.ExecContext(ctx, indexSQL); err != nil {
			return err
		}
	}

	return nil
}

func downInitialDb(ctx context.Context, tx *sql.Tx) error {
	dropIndexes := []string{
		"DROP INDEX IF EXISTS idx_pull_request_reviewers_pr_id",
		"DROP INDEX IF EXISTS idx_pull_request_reviewers_reviewer_id",
		"DROP INDEX IF EXISTS idx_pull_requests_status",
		"DROP INDEX IF EXISTS idx_pull_requests_author_id",
		"DROP INDEX IF EXISTS idx_users_team_active",
		"DROP INDEX IF EXISTS idx_users_is_active",
		"DROP INDEX IF EXISTS idx_users_team_name",
	}

	for _, dropSQL := range dropIndexes {
		if _, err := tx.ExecContext(ctx, dropSQL); err != nil {
			return err
		}
	}

	dropTables := []string{
		"DROP TABLE IF EXISTS pull_request_reviewers",
		"DROP TABLE IF EXISTS pull_requests",
		"DROP TABLE IF EXISTS users",
		"DROP TABLE IF EXISTS teams",
	}

	for _, dropSQL := range dropTables {
		if _, err := tx.ExecContext(ctx, dropSQL); err != nil {
			return err
		}
	}

	return nil
}
