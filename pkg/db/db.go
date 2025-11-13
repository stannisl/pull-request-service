package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectPoolWithRetry пытается подключиться к БД с повторными попытками
func ConnectPoolWithRetry(ctx context.Context, connStr string, maxRetries int, retryInterval time.Duration) (*pgxpool.Pool, error) {

	var pool *pgxpool.Pool
	var err error

	for i := 0; i < maxRetries; i++ {
		pool, err = pgxpool.New(ctx, connStr)
		if err == nil {
			return pool, nil
		}

		if i < maxRetries-1 {
			log.Printf("Failed to connect to database (attempt %d/%d): %v. Retrying in %v...", i+1, maxRetries, err, retryInterval)
			time.Sleep(retryInterval)
		}
	}

	return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, err)
}

type ReleaseFunc func()

func GetConnFromPool(ctx context.Context, pool *pgxpool.Pool) (*pgxpool.Conn, ReleaseFunc, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, nil, err
	}

	return conn, func() { conn.Release() }, nil
}
