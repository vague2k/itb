package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	db "itb.ihatedoing.work/server/database/generated"
	_ "modernc.org/sqlite"
)

const dbFilename = "itb.db"

type Config struct {
	Port         string
	DatabaseDir  string
	UnderlyingDB *sql.DB
	Database     *db.Queries
}

func Init() *Config {
	c := &Config{
		Port:        os.Getenv("PORT"),
		DatabaseDir: envOr("DB_PATH", "data"),
	}

	if err := os.MkdirAll(c.DatabaseDir, 0o755); err != nil {
		log.Fatalf("could not create database dir: %v", err)
	}

	conn, err := openDatabase(c.DatabaseDir)
	if err != nil {
		log.Fatalf("could not open database: %v", err)
	}

	c.UnderlyingDB = conn
	c.Database = db.New(conn)

	return c
}

func openDatabase(dir string) (*sql.DB, error) {
	sqlDB, err := sql.Open("sqlite", fmt.Sprintf(
		"file:%s?_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)",
		filepath.Join(dir, dbFilename),
	))
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return sqlDB, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
