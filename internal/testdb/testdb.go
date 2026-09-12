package testdb

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"

	"itb.ihatedoing.work/server/database"
	db "itb.ihatedoing.work/server/database/generated"
)

// Open returns a fully migrated SQLite database backed by a file in a
// per-test temp directory, using the same DSN pragmas as the app.
func Open(t *testing.T) (*sql.DB, *db.Queries) {
	t.Helper()
	goose.SetLogger(goose.NopLogger())

	sqlDB, err := sql.Open("sqlite", fmt.Sprintf(
		"file:%s?_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)",
		filepath.Join(t.TempDir(), "itb.db"),
	))
	if err != nil {
		t.Fatalf("testdb: open sqlite: %v", err)
	}
	t.Cleanup(func() { sqlDB.Close() })

	if err := database.MigrateUp(sqlDB); err != nil {
		t.Fatalf("testdb: migrate: %v", err)
	}
	return sqlDB, db.New(sqlDB)
}
