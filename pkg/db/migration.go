package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"strings"
)

//go:embed migrations/SCHEMA.sql
var schemaFS embed.FS

//go:embed migrations/DROP_SCHEMA.sql
var dropSchemaFS embed.FS

type Migrator struct {
	conn      *sql.Conn
	closeFunc ReleaseFunc
}

func NewMigrator(conn *sql.Conn, closeFunc ReleaseFunc) *Migrator {
	return &Migrator{
		conn:      conn,
		closeFunc: closeFunc,
	}
}

func (m *Migrator) Close() {
	m.closeFunc()
	m.conn = nil
}

// Run создает схему из migrations/SCHEMA.sql
func (m *Migrator) Run(ctx context.Context) error {
	if m.conn == nil {
		return fmt.Errorf("conn is closed")
	}

	schemaSQL, err := schemaFS.ReadFile("migrations/SCHEMA.sql")

	if err != nil {
		return fmt.Errorf("failed to read schema: %w", err)
	}

	queries := strings.Split(string(schemaSQL), ";")

	for _, query := range queries {
		tx, err := m.conn.BeginTx(ctx, nil)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to start transaction: %w", err)
		}

		_, err = tx.ExecContext(ctx, query)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute query: %w", err)
		}

		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
	}
	return nil
}
