package storage

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

// DB represents the application's database connection pool.
type DB struct {
	conn *sql.DB
	mu   sync.RWMutex // Protects connection operations
}

// Config holds database configuration options.
type Config struct {
	// Path to the SQLite database file. If empty, defaults to ~/.agentic-obs/db.sqlite
	Path string
}

var (
	// Default database path in user's home directory
	defaultDBPath = filepath.Join(getHomeDir(), ".agentic-obs", "db.sqlite")
)

// New creates a new database connection with the given configuration.
// It initializes the database schema if this is the first run.
//
// Context is used to support cancellation during initialization.
func New(ctx context.Context, cfg Config) (*DB, error) {
	dbPath := cfg.Path
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// Ensure the parent directory exists
	if err := ensureDir(filepath.Dir(dbPath)); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection with appropriate settings
	// modernc.org/sqlite uses the same connection string format as mattn/go-sqlite3
	conn, err := sql.Open("sqlite", dbPath+"?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, fmt.Errorf("failed to open database at %s: %w", dbPath, err)
	}

	// Configure connection pool for single-user, single-tenant usage
	// SQLite performs best with a single writer, multiple readers
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)

	// Verify connection is working
	if err := conn.PingContext(ctx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{conn: conn}

	// Run migrations to initialize or update schema
	if err := db.migrate(ctx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	return db, nil
}

// migrate runs database schema migrations.
// This is idempotent - it can be safely run multiple times.
func (db *DB) migrate(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Create migrations table to track schema version
	migrations := []string{
		// Migration 0: Create schema version tracking
		`CREATE TABLE IF NOT EXISTS schema_version (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Migration 1: Create config table for connection settings
		`CREATE TABLE IF NOT EXISTS config (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Migration 2: Create state table for application state
		`CREATE TABLE IF NOT EXISTS state (
			key TEXT PRIMARY KEY,
			value TEXT,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Migration 3: Create scene_presets table for user-defined presets
		`CREATE TABLE IF NOT EXISTS scene_presets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			scene_name TEXT NOT NULL,
			sources TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Migration 4: Create index on config.updated_at for performance
		`CREATE INDEX IF NOT EXISTS idx_config_updated_at ON config(updated_at)`,

		// Migration 5: Create index on state.updated_at for performance
		`CREATE INDEX IF NOT EXISTS idx_state_updated_at ON state(updated_at)`,

		// Migration 6: Create index on scene_presets.name for lookups
		`CREATE INDEX IF NOT EXISTS idx_scene_presets_name ON scene_presets(name)`,

		// Migration 7: Create index on scene_presets.scene_name for filtering
		`CREATE INDEX IF NOT EXISTS idx_scene_presets_scene_name ON scene_presets(scene_name)`,
	}

	// Execute each migration in a transaction
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin migration transaction: %w", err)
	}
	defer tx.Rollback() // Safe to call even if committed

	for i, migration := range migrations {
		if _, err := tx.ExecContext(ctx, migration); err != nil {
			return fmt.Errorf("failed to execute migration %d: %w", i, err)
		}
	}

	// Record successful migration
	if _, err := tx.ExecContext(ctx,
		"INSERT OR REPLACE INTO schema_version (version) VALUES (?)",
		len(migrations),
	); err != nil {
		return fmt.Errorf("failed to record schema version: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration transaction: %w", err)
	}

	return nil
}

// Close closes the database connection pool.
// It's safe to call multiple times.
func (db *DB) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.conn != nil {
		if err := db.conn.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %w", err)
		}
		db.conn = nil
	}
	return nil
}

// DB returns the underlying sql.DB connection for advanced operations.
// Use with caution - prefer the higher-level methods when possible.
func (db *DB) DB() *sql.DB {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.conn
}

// Ping verifies the database connection is alive.
func (db *DB) Ping(ctx context.Context) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if err := db.conn.PingContext(ctx); err != nil {
		return fmt.Errorf("database connection check failed: is the database accessible? %w", err)
	}
	return nil
}

// Transaction executes a function within a database transaction.
// If the function returns an error, the transaction is rolled back.
// Otherwise, the transaction is committed.
//
// Example:
//
//	err := db.Transaction(ctx, func(tx *sql.Tx) error {
//	    // Perform multiple operations
//	    _, err := tx.ExecContext(ctx, "INSERT INTO ...")
//	    return err
//	})
func (db *DB) Transaction(ctx context.Context, fn func(*sql.Tx) error) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Safe to call even if committed

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ensureDir creates a directory if it doesn't exist.
func ensureDir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}
	return nil
}

// getHomeDir returns the user's home directory.
func getHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home directory is unavailable
		return "."
	}
	return home
}

// IsFirstRun checks if this is the first time the application is being run
// by checking if the database file exists and has data.
func IsFirstRun(cfg Config) bool {
	dbPath := cfg.Path
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// Check if database file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return true
	}

	// Database exists, check if it has any config data
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return true
	}
	defer conn.Close()

	var count int
	err = conn.QueryRow("SELECT COUNT(*) FROM config WHERE key = 'obs_host'").Scan(&count)
	if err != nil || count == 0 {
		return true
	}

	return false
}
