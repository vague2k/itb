package database

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

// MigrateUp applies all pending migrations.
func MigrateUp(db *sql.DB) error {
	if err := setupGoose(); err != nil {
		return fmt.Errorf("setup migrations: %w", err)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}

// MigrateDown rolls back the most recently applied migration.
func MigrateDown(db *sql.DB) error {
	if err := setupGoose(); err != nil {
		return fmt.Errorf("setup migrations: %w", err)
	}
	if err := goose.Down(db, "migrations"); err != nil {
		return fmt.Errorf("roll back migration: %w", err)
	}
	return nil
}

// MigrateStatus prints the current migration status for each version.
func MigrateStatus(db *sql.DB) error {
	if err := setupGoose(); err != nil {
		return fmt.Errorf("setup migrations: %w", err)
	}
	if err := goose.Status(db, "migrations"); err != nil {
		return fmt.Errorf("read migration status: %w", err)
	}
	return nil
}

func setupGoose() error {
	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}
	return nil
}
