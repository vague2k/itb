package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	db "itb.ihatedoing.work/server/database/generated"
	_ "modernc.org/sqlite"
)

// set at build/compile time using ldflags to determine which env file to load
//
// downside to using ldflags is that they work **ONLY** if it's a var type string
// so this is a perfect usecase
var buildEnv = "development"

const (
	envDEV  = "development"
	envSTG  = "staging"
	envPROD = "production"

	dbFilename = "itb.db"
)

type Config struct {
	Environment string
	Port        string
	DatabaseDir string

	UnderlyingDB *sql.DB
	Database     *db.Queries
}

func Init() *Config {
	envfile := envFile(buildEnv)
	if envfile != "" {
		if err := godotenv.Load(envfile); err != nil {
			log.Fatalf("Fatal config error: loading env file '%s': %v", envfile, err)
		}
	}

	c := &Config{
		Environment: buildEnv,
		Port:        requireEnv("PORT"),
		DatabaseDir: requireEnv("DB_PATH"),
	}

	if err := os.MkdirAll(c.DatabaseDir, 0o755); err != nil {
		log.Fatalf("Fatal config error: creating database dir: %v", err)
	}

	conn, err := openDatabase(c.DatabaseDir)
	if err != nil {
		log.Fatalf("Fatal config error: opening database: %v", err)
	}

	c.UnderlyingDB = conn
	c.Database = db.New(conn)

	logStartup(c, envfile)

	return c
}

func (c *Config) Production() bool {
	return c.Environment == envPROD
}

// Deployed reports whether the binary is a deployed build (staging or
// production) rather than a local dev run; deployed builds serve embedded assets.
func (c *Config) Deployed() bool {
	return c.Environment != envDEV
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

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment variable '%s' is not set and is required", key)
	}
	return v
}

func envFile(env string) string {
	switch env {
	case envDEV:
		return ".env"
	case envSTG:
		return ".env.staging"
	case envPROD:
		return ".env.production"
	default:
		return ""
	}
}

func logStartup(c *Config, envPath string) {
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  ── server config ──────────────────────")
	fmt.Fprintf(os.Stderr, "  %-18s %s\n", "environment:", c.Environment)
	fmt.Fprintf(os.Stderr, "  %-18s %v\n", "production:", c.Production())
	fmt.Fprintf(os.Stderr, "  %-18s %s\n", "env_file:", envPath)
	fmt.Fprintf(os.Stderr, "  %-18s %s\n", "port:", c.Port)
	fmt.Fprintf(os.Stderr, "  %-18s %s\n", "database_dir:", c.DatabaseDir)
	fmt.Fprintf(os.Stderr, "  %-18s %s\n", "database_path:", filepath.Join(c.DatabaseDir, dbFilename))
	fmt.Fprintf(os.Stderr, "  %-18s %s/%s\n", "goos/goarch:", runtime.GOOS, runtime.GOARCH)
	fmt.Fprintln(os.Stderr, "  ─────────────────────────────────────────")
	fmt.Fprintln(os.Stderr, "")
}
