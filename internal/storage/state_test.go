package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveOBSConfig(t *testing.T) {
	t.Run("saves config successfully", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		cfg := OBSConfig{
			Host:     "localhost",
			Port:     4455,
			Password: "secret123",
		}

		err := db.SaveOBSConfig(context.Background(), cfg)
		assert.NoError(t, err)
	})

	t.Run("updates existing config", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		// Save initial config
		err := db.SaveOBSConfig(context.Background(), OBSConfig{
			Host:     "localhost",
			Port:     4455,
			Password: "old-password",
		})
		require.NoError(t, err)

		// Update config
		err = db.SaveOBSConfig(context.Background(), OBSConfig{
			Host:     "192.168.1.100",
			Port:     4456,
			Password: "new-password",
		})
		require.NoError(t, err)

		// Verify update
		loaded, err := db.LoadOBSConfig(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, "192.168.1.100", loaded.Host)
		assert.Equal(t, 4456, loaded.Port)
		assert.Equal(t, "new-password", loaded.Password)
	})

	t.Run("saves config without password", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		cfg := OBSConfig{
			Host:     "localhost",
			Port:     4455,
			Password: "",
		}

		err := db.SaveOBSConfig(context.Background(), cfg)
		assert.NoError(t, err)

		loaded, err := db.LoadOBSConfig(context.Background())
		assert.NoError(t, err)
		assert.Empty(t, loaded.Password)
	})
}

func TestLoadOBSConfig(t *testing.T) {
	t.Run("loads saved config", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		expected := OBSConfig{
			Host:     "192.168.1.50",
			Port:     4455,
			Password: "mypassword",
		}

		err := db.SaveOBSConfig(context.Background(), expected)
		require.NoError(t, err)

		loaded, err := db.LoadOBSConfig(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expected.Host, loaded.Host)
		assert.Equal(t, expected.Port, loaded.Port)
		assert.Equal(t, expected.Password, loaded.Password)
	})

	t.Run("returns error when config not found", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		_, err := db.LoadOBSConfig(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestSetState(t *testing.T) {
	t.Run("sets state value", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		err := db.SetState(context.Background(), "test_key", "test_value")
		assert.NoError(t, err)

		val, err := db.GetState(context.Background(), "test_key")
		assert.NoError(t, err)
		assert.Equal(t, "test_value", val)
	})

	t.Run("updates existing state", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		err := db.SetState(context.Background(), "key", "value1")
		require.NoError(t, err)

		err = db.SetState(context.Background(), "key", "value2")
		require.NoError(t, err)

		val, err := db.GetState(context.Background(), "key")
		assert.NoError(t, err)
		assert.Equal(t, "value2", val)
	})

	t.Run("handles empty value", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		err := db.SetState(context.Background(), "empty", "")
		assert.NoError(t, err)

		val, err := db.GetState(context.Background(), "empty")
		assert.NoError(t, err)
		assert.Empty(t, val)
	})
}

func TestGetState(t *testing.T) {
	t.Run("retrieves state value", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		err := db.SetState(context.Background(), "my_key", "my_value")
		require.NoError(t, err)

		val, err := db.GetState(context.Background(), "my_key")
		assert.NoError(t, err)
		assert.Equal(t, "my_value", val)
	})

	t.Run("returns error for non-existent key", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		_, err := db.GetState(context.Background(), "nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestDeleteState(t *testing.T) {
	t.Run("deletes state value", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		err := db.SetState(context.Background(), "to_delete", "value")
		require.NoError(t, err)

		err = db.DeleteState(context.Background(), "to_delete")
		assert.NoError(t, err)

		_, err = db.GetState(context.Background(), "to_delete")
		assert.Error(t, err)
	})

	t.Run("succeeds even if key does not exist", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		err := db.DeleteState(context.Background(), "nonexistent")
		assert.NoError(t, err) // DELETE is idempotent
	})
}

func TestListState(t *testing.T) {
	t.Run("lists all state entries", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		// Set multiple state entries
		err := db.SetState(context.Background(), "key1", "value1")
		require.NoError(t, err)
		err = db.SetState(context.Background(), "key2", "value2")
		require.NoError(t, err)
		err = db.SetState(context.Background(), "key3", "value3")
		require.NoError(t, err)

		state, err := db.ListState(context.Background())
		assert.NoError(t, err)
		assert.Len(t, state, 3)
		assert.Equal(t, "value1", state["key1"])
		assert.Equal(t, "value2", state["key2"])
		assert.Equal(t, "value3", state["key3"])
	})

	t.Run("returns empty map when no state", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		state, err := db.ListState(context.Background())
		assert.NoError(t, err)
		assert.Empty(t, state)
	})
}

func TestMarkFirstRunComplete(t *testing.T) {
	t.Run("marks first run as complete", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		err := db.MarkFirstRunComplete(context.Background())
		assert.NoError(t, err)

		isFirst, err := db.IsFirstRun(context.Background())
		assert.NoError(t, err)
		assert.False(t, isFirst)
	})
}

func TestDB_IsFirstRun(t *testing.T) {
	t.Run("returns true when state not set", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		isFirst, err := db.IsFirstRun(context.Background())
		assert.NoError(t, err)
		assert.True(t, isFirst)
	})

	t.Run("returns false after MarkFirstRunComplete", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		err := db.MarkFirstRunComplete(context.Background())
		require.NoError(t, err)

		isFirst, err := db.IsFirstRun(context.Background())
		assert.NoError(t, err)
		assert.False(t, isFirst)
	})

	t.Run("returns true when state is 'true'", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		err := db.SetState(context.Background(), StateKeyFirstRun, "true")
		require.NoError(t, err)

		isFirst, err := db.IsFirstRun(context.Background())
		assert.NoError(t, err)
		assert.True(t, isFirst)
	})
}

func TestRecordSuccessfulConnection(t *testing.T) {
	t.Run("records connection timestamp", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		// Truncate to second precision since RFC3339 stores only seconds
		before := time.Now().UTC().Truncate(time.Second)
		err := db.RecordSuccessfulConnection(context.Background())
		assert.NoError(t, err)
		after := time.Now().UTC().Add(time.Second).Truncate(time.Second)

		lastConnected, err := db.GetLastConnectedTime(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, lastConnected)

		// Timestamp should be between before and after (with second precision)
		truncatedConnected := lastConnected.Truncate(time.Second)
		assert.True(t, !truncatedConnected.Before(before), "lastConnected %v should not be before %v", truncatedConnected, before)
		assert.True(t, !truncatedConnected.After(after), "lastConnected %v should not be after %v", truncatedConnected, after)
	})
}

func TestGetLastConnectedTime(t *testing.T) {
	t.Run("returns nil when never connected", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		lastConnected, err := db.GetLastConnectedTime(context.Background())
		assert.NoError(t, err)
		assert.Nil(t, lastConnected)
	})

	t.Run("returns timestamp after connection", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		err := db.RecordSuccessfulConnection(context.Background())
		require.NoError(t, err)

		lastConnected, err := db.GetLastConnectedTime(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, lastConnected)
	})
}

func TestSetAutoReconnect(t *testing.T) {
	t.Run("sets auto-reconnect enabled", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		err := db.SetAutoReconnect(context.Background(), true)
		assert.NoError(t, err)

		enabled, err := db.GetAutoReconnect(context.Background())
		assert.NoError(t, err)
		assert.True(t, enabled)
	})

	t.Run("sets auto-reconnect disabled", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		err := db.SetAutoReconnect(context.Background(), false)
		assert.NoError(t, err)

		enabled, err := db.GetAutoReconnect(context.Background())
		assert.NoError(t, err)
		assert.False(t, enabled)
	})
}

func TestGetAutoReconnect(t *testing.T) {
	t.Run("defaults to true when not set", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		enabled, err := db.GetAutoReconnect(context.Background())
		assert.NoError(t, err)
		assert.True(t, enabled)
	})
}

func TestSetAppVersion(t *testing.T) {
	t.Run("sets app version", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		err := db.SetAppVersion(context.Background(), "1.2.3")
		assert.NoError(t, err)

		version, err := db.GetAppVersion(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, "1.2.3", version)
	})
}

func TestGetAppVersion(t *testing.T) {
	t.Run("returns error when not set", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		_, err := db.GetAppVersion(context.Background())
		assert.Error(t, err)
	})

	t.Run("returns version when set", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		err := db.SetAppVersion(context.Background(), "0.1.0")
		require.NoError(t, err)

		version, err := db.GetAppVersion(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, "0.1.0", version)
	})
}

func TestStateConstants(t *testing.T) {
	t.Run("state keys are defined", func(t *testing.T) {
		assert.Equal(t, "first_run", StateKeyFirstRun)
		assert.Equal(t, "last_connected", StateKeyLastConnected)
		assert.Equal(t, "app_version", StateKeyAppVersion)
		assert.Equal(t, "auto_reconnect", StateKeyAutoReconnect)
	})

	t.Run("config keys are defined", func(t *testing.T) {
		assert.Equal(t, "obs_host", ConfigKeyOBSHost)
		assert.Equal(t, "obs_port", ConfigKeyOBSPort)
		assert.Equal(t, "obs_password", ConfigKeyOBSPassword)
	})
}
