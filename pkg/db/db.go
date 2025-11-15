package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Драйвер для postgres
)

// OptionsDB настройки базы данных
type OptionsDB struct {
	ConnStr       string
	MaxRetries    int
	RetryInterval time.Duration
	DriverName    string
}

// ConnectPoolWithRetry пытается подключиться к бд с повторными попытками
func ConnectPoolWithRetry(ctx context.Context, opts *OptionsDB) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	for i := 0; i < opts.MaxRetries; i++ {
		db, err = sqlx.ConnectContext(ctx, opts.DriverName, opts.ConnStr)
		if err == nil {
			return db, nil
		}

		if i < opts.MaxRetries-1 {
			log.Printf(
				"Failed to connect to database (attempt %d/%d): %v. Retrying in %v...",
				i+1,
				opts.MaxRetries,
				err,
				opts.RetryInterval,
			)
			time.Sleep(opts.RetryInterval)
		}
	}

	return nil, fmt.Errorf("failed to connect after %d attempts: %w", opts.MaxRetries, err)
}

type ReleaseFunc func() error

// GetConnFromPool берет 1 соединение из пула соединений
func GetConnFromPool(ctx context.Context, pool *sqlx.DB) (*sql.Conn, ReleaseFunc, error) {
	conn, err := pool.Conn(ctx)
	if err != nil {
		return nil, nil, err
	}

	return conn, func() error { return conn.Close() }, nil
}
