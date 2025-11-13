package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	_ "github.com/stannisl/avito-test/pkg/db/migrations" // Импортируем для регистрации Go-миграций
)

//go:embed migrations
var migrationsFS embed.FS

// RunMigrations выполняет все миграции из директории migrations используя goose
func RunMigrations(ctx context.Context, conn *pgx.Conn) error {

	config := conn.Config()
	sqlDB := stdlib.OpenDB(*config)
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			log.Fatalf("error closing DB conn: %v", err)
		}
	}(sqlDB)

	goose.SetBaseFS(migrationsFS)

	migrationsDir := "migrations"

	if err := goose.Up(sqlDB, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// GetMigrationVersion возвращает текущую версию миграции
func GetMigrationVersion(ctx context.Context, conn *pgx.Conn) (int64, error) {
	config := conn.Config()
	sqlDB := stdlib.OpenDB(*config)
	defer sqlDB.Close()

	goose.SetBaseFS(migrationsFS)

	version, err := goose.GetDBVersion(sqlDB)
	if err != nil {
		return 0, fmt.Errorf("failed to get db version: %w", err)
	}

	return version, nil
}

func Down(ctx context.Context, conn *pgx.Conn) error {
	config := conn.Config()
	sqlDB := stdlib.OpenDB(*config)
	defer sqlDB.Close()

	goose.SetBaseFS(migrationsFS)
	migrationsDir := "migrations"

	if err := goose.Down(sqlDB, migrationsDir); err != nil {
		return fmt.Errorf("failed to rollback db: %w", err)
	}

	return nil
}

func Status(ctx context.Context, conn *pgx.Conn) error {
	config := conn.Config()
	sqlDB := stdlib.OpenDB(*config)
	defer sqlDB.Close()

	goose.SetBaseFS(migrationsFS)
	migrationsDir := "migrations"

	if err := goose.Status(sqlDB, migrationsDir); err != nil {
		return fmt.Errorf("failed to get db status: %w", err)
	}

	return nil
}
