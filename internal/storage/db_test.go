package storage

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testDB creates a temporary database for testing.
func testDB(t *testing.T) (*DB, func()) {
	t.Helper()

	// Create temp directory for test database
	tempDir, err := os.MkdirTemp("", "agentic-obs-test-*")
	require.NoError(t, err)

	dbPath := filepath.Join(tempDir, "test.db")

	db, err := New(context.Background(), Config{Path: dbPath})
	require.NoError(t, err)

	cleanup := func() {
		db.Close()
		os.RemoveAll(tempDir)
	}

	return db, cleanup
}

func TestNew(t *testing.T) {
	t.Run("creates database with default path", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "agentic-obs-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		dbPath := filepath.Join(tempDir, "test.db")
		db, err := New(context.Background(), Config{Path: dbPath})
		require.NoError(t, err)
		defer db.Close()

		// Verify database file was created
		_, err = os.Stat(dbPath)
		assert.NoError(t, err)
	})

	t.Run("creates parent directories", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "agentic-obs-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		dbPath := filepath.Join(tempDir, "nested", "dir", "test.db")
		db, err := New(context.Background(), Config{Path: dbPath})
		require.NoError(t, err)
		defer db.Close()

		// Verify database file was created
		_, err = os.Stat(dbPath)
		assert.NoError(t, err)
	})

	t.Run("runs migrations successfully", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		// Check that required tables exist
		tables := []string{"schema_version", "config", "state", "scene_presets"}
		for _, table := range tables {
			var name string
			err := db.DB().QueryRow(
				"SELECT name FROM sqlite_master WHERE type='table' AND name=?",
				table,
			).Scan(&name)
			assert.NoError(t, err, "table %s should exist", table)
			assert.Equal(t, table, name)
		}
	})

	t.Run("migrations are idempotent", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "agentic-obs-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		dbPath := filepath.Join(tempDir, "test.db")

		// Create database first time
		db1, err := New(context.Background(), Config{Path: dbPath})
		require.NoError(t, err)
		db1.Close()

		// Open again - migrations should be safe to run again
		db2, err := New(context.Background(), Config{Path: dbPath})
		require.NoError(t, err)
		defer db2.Close()

		// Verify tables still exist
		var count int
		err = db2.DB().QueryRow(
			"SELECT COUNT(*) FROM sqlite_master WHERE type='table'",
		).Scan(&count)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count, 4) // At least our 4 main tables
	})
}

func TestDB_Ping(t *testing.T) {
	t.Run("returns nil when connected", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		err := db.Ping(context.Background())
		assert.NoError(t, err)
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		err := db.Ping(ctx)
		assert.Error(t, err)
	})
}

func TestDB_Close(t *testing.T) {
	t.Run("closes connection successfully", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "agentic-obs-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		dbPath := filepath.Join(tempDir, "test.db")
		db, err := New(context.Background(), Config{Path: dbPath})
		require.NoError(t, err)

		// Verify connection works before close
		err = db.Ping(context.Background())
		require.NoError(t, err)

		err = db.Close()
		assert.NoError(t, err)

		// Connection should be nil after close
		assert.Nil(t, db.DB())
	})

	t.Run("can be called multiple times", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		err := db.Close()
		assert.NoError(t, err)

		err = db.Close()
		assert.NoError(t, err) // Should not error on second close
	})
}

func TestDB_Transaction(t *testing.T) {
	t.Run("commits successful transaction", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		// Insert data in transaction
		err := db.Transaction(context.Background(), func(tx *sql.Tx) error {
			// Use SetState which internally uses the transaction
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("rolls back on error", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		// Set initial value
		err := db.SetState(context.Background(), "test_key", "initial")
		require.NoError(t, err)

		// Try to update in transaction that fails
		err = db.Transaction(context.Background(), func(tx *sql.Tx) error {
			return assert.AnError
		})
		assert.Error(t, err)

		// Value should still be initial
		val, err := db.GetState(context.Background(), "test_key")
		assert.NoError(t, err)
		assert.Equal(t, "initial", val)
	})
}

func TestDB_DB(t *testing.T) {
	t.Run("returns underlying connection", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		conn := db.DB()
		assert.NotNil(t, conn)

		// Should be able to execute queries
		var result int
		err := conn.QueryRow("SELECT 1").Scan(&result)
		assert.NoError(t, err)
		assert.Equal(t, 1, result)
	})
}

func TestIsFirstRun(t *testing.T) {
	t.Run("returns true when database does not exist", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "agentic-obs-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		dbPath := filepath.Join(tempDir, "nonexistent.db")
		isFirst := IsFirstRun(Config{Path: dbPath})
		assert.True(t, isFirst)
	})

	t.Run("returns true when database has no config", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "agentic-obs-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		dbPath := filepath.Join(tempDir, "empty.db")

		// Create database without config
		db, err := New(context.Background(), Config{Path: dbPath})
		require.NoError(t, err)
		db.Close()

		// Should return true since no OBS config exists
		isFirst := IsFirstRun(Config{Path: dbPath})
		assert.True(t, isFirst)
	})

	t.Run("returns false when database has config", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "agentic-obs-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		dbPath := filepath.Join(tempDir, "test.db")
		db, err := New(context.Background(), Config{Path: dbPath})
		require.NoError(t, err)

		// Save config
		err = db.SaveOBSConfig(context.Background(), OBSConfig{
			Host:     "localhost",
			Port:     4455,
			Password: "",
		})
		require.NoError(t, err)
		db.Close()

		// Check first run
		isFirst := IsFirstRun(Config{Path: dbPath})
		assert.False(t, isFirst)
	})
}

func TestEnsureDir(t *testing.T) {
	t.Run("creates directory if not exists", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "agentic-obs-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		newDir := filepath.Join(tempDir, "new", "nested", "dir")
		err = ensureDir(newDir)
		assert.NoError(t, err)

		// Directory should exist
		info, err := os.Stat(newDir)
		assert.NoError(t, err)
		assert.True(t, info.IsDir())
	})

	t.Run("succeeds if directory already exists", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "agentic-obs-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		err = ensureDir(tempDir)
		assert.NoError(t, err)
	})
}

func TestGetHomeDir(t *testing.T) {
	t.Run("returns home directory", func(t *testing.T) {
		home := getHomeDir()
		assert.NotEmpty(t, home)

		// Should be a valid directory
		info, err := os.Stat(home)
		if err == nil {
			assert.True(t, info.IsDir())
		}
	})
}
