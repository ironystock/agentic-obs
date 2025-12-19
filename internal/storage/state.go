package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// OBSConfig represents the OBS WebSocket connection configuration.
type OBSConfig struct {
	Host     string // OBS WebSocket host (e.g., "localhost")
	Port     int    // OBS WebSocket port (e.g., 4455)
	Password string // OBS WebSocket password (may be empty)
}

// State keys used in the state table
const (
	StateKeyFirstRun      = "first_run"      // Whether this is the first run
	StateKeyLastConnected = "last_connected" // Timestamp of last successful OBS connection
	StateKeyAppVersion    = "app_version"    // Application version
	StateKeyAutoReconnect = "auto_reconnect" // Auto-reconnect preference
)

// Tool group state keys - control which tool categories are enabled
const (
	StateKeyToolsCore    = "tools_enabled_core"    // Core OBS tools (scenes, recording, streaming, status)
	StateKeyToolsVisual  = "tools_enabled_visual"  // Visual monitoring tools (screenshots)
	StateKeyToolsLayout  = "tools_enabled_layout"  // Layout management tools (scene presets)
	StateKeyToolsAudio   = "tools_enabled_audio"   // Audio control tools
	StateKeyToolsSources = "tools_enabled_sources" // Source management tools
	StateKeyToolsDesign  = "tools_enabled_design"  // Scene design tools (source creation, transforms)
)

// Webserver configuration keys
const (
	ConfigKeyWebServerEnabled = "web_server_enabled" // Whether HTTP server is enabled
	ConfigKeyWebServerHost    = "web_server_host"    // HTTP server host
	ConfigKeyWebServerPort    = "web_server_port"    // HTTP server port
)

// Config keys used in the config table
const (
	ConfigKeyOBSHost     = "obs_host"     // OBS WebSocket host
	ConfigKeyOBSPort     = "obs_port"     // OBS WebSocket port
	ConfigKeyOBSPassword = "obs_password" // OBS WebSocket password
)

// SaveOBSConfig persists OBS connection configuration to the database.
// This updates existing values or inserts new ones atomically.
func (db *DB) SaveOBSConfig(ctx context.Context, cfg OBSConfig) error {
	return db.Transaction(ctx, func(tx *sql.Tx) error {
		// Prepare upsert statement for config table
		stmt, err := tx.PrepareContext(ctx, `
			INSERT INTO config (key, value, updated_at)
			VALUES (?, ?, CURRENT_TIMESTAMP)
			ON CONFLICT(key) DO UPDATE SET
				value = excluded.value,
				updated_at = CURRENT_TIMESTAMP
		`)
		if err != nil {
			return fmt.Errorf("failed to prepare config upsert statement: %w", err)
		}
		defer stmt.Close()

		// Save host
		if _, err := stmt.ExecContext(ctx, ConfigKeyOBSHost, cfg.Host); err != nil {
			return fmt.Errorf("failed to save OBS host: %w", err)
		}

		// Save port
		if _, err := stmt.ExecContext(ctx, ConfigKeyOBSPort, fmt.Sprintf("%d", cfg.Port)); err != nil {
			return fmt.Errorf("failed to save OBS port: %w", err)
		}

		// Save password (may be empty)
		if _, err := stmt.ExecContext(ctx, ConfigKeyOBSPassword, cfg.Password); err != nil {
			return fmt.Errorf("failed to save OBS password: %w", err)
		}

		return nil
	})
}

// LoadOBSConfig retrieves OBS connection configuration from the database.
// Returns an error if the configuration is not found or incomplete.
func (db *DB) LoadOBSConfig(ctx context.Context) (*OBSConfig, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	cfg := &OBSConfig{}

	// Query host
	var host string
	err := db.conn.QueryRowContext(ctx,
		"SELECT value FROM config WHERE key = ?",
		ConfigKeyOBSHost,
	).Scan(&host)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("OBS configuration not found. Please run the setup process")
	} else if err != nil {
		return nil, fmt.Errorf("failed to load OBS host: %w", err)
	}
	cfg.Host = host

	// Query port
	var portStr string
	err = db.conn.QueryRowContext(ctx,
		"SELECT value FROM config WHERE key = ?",
		ConfigKeyOBSPort,
	).Scan(&portStr)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("OBS port not configured. Please run the setup process")
	} else if err != nil {
		return nil, fmt.Errorf("failed to load OBS port: %w", err)
	}

	var port int
	if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
		return nil, fmt.Errorf("invalid OBS port value '%s': %w", portStr, err)
	}
	cfg.Port = port

	// Query password (may be empty)
	var password string
	err = db.conn.QueryRowContext(ctx,
		"SELECT value FROM config WHERE key = ?",
		ConfigKeyOBSPassword,
	).Scan(&password)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to load OBS password: %w", err)
	}
	cfg.Password = password

	return cfg, nil
}

// SetState saves a key-value pair to the application state table.
// This is useful for storing application-level settings and flags.
func (db *DB) SetState(ctx context.Context, key, value string) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	_, err := db.conn.ExecContext(ctx, `
		INSERT INTO state (key, value, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(key) DO UPDATE SET
			value = excluded.value,
			updated_at = CURRENT_TIMESTAMP
	`, key, value)

	if err != nil {
		return fmt.Errorf("failed to set state key '%s': %w", key, err)
	}

	return nil
}

// GetState retrieves a value from the application state table.
// Returns sql.ErrNoRows if the key doesn't exist.
func (db *DB) GetState(ctx context.Context, key string) (string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var value string
	err := db.conn.QueryRowContext(ctx,
		"SELECT value FROM state WHERE key = ?",
		key,
	).Scan(&value)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("state key '%s' not found", key)
	} else if err != nil {
		return "", fmt.Errorf("failed to get state key '%s': %w", key, err)
	}

	return value, nil
}

// DeleteState removes a key-value pair from the application state table.
func (db *DB) DeleteState(ctx context.Context, key string) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	_, err := db.conn.ExecContext(ctx,
		"DELETE FROM state WHERE key = ?",
		key,
	)

	if err != nil {
		return fmt.Errorf("failed to delete state key '%s': %w", key, err)
	}

	return nil
}

// ListState returns all key-value pairs from the application state table.
// This is useful for debugging and administrative operations.
func (db *DB) ListState(ctx context.Context) (map[string]string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	rows, err := db.conn.QueryContext(ctx,
		"SELECT key, value FROM state ORDER BY key",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list state: %w", err)
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, fmt.Errorf("failed to scan state row: %w", err)
		}
		result[key] = value
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating state rows: %w", err)
	}

	return result, nil
}

// MarkFirstRunComplete marks that the first-run setup has been completed.
// This is typically called after successful OBS connection during initial setup.
func (db *DB) MarkFirstRunComplete(ctx context.Context) error {
	return db.SetState(ctx, StateKeyFirstRun, "false")
}

// IsFirstRun checks if this is the first run of the application.
// Returns true if the first_run state is not set or is "true".
func (db *DB) IsFirstRun(ctx context.Context) (bool, error) {
	value, err := db.GetState(ctx, StateKeyFirstRun)
	if err != nil {
		// If key doesn't exist, this is the first run
		if err.Error() == fmt.Sprintf("state key '%s' not found", StateKeyFirstRun) {
			return true, nil
		}
		return false, err
	}

	return value == "true", nil
}

// RecordSuccessfulConnection records the timestamp of the last successful
// OBS connection. This is useful for diagnostics and connection health monitoring.
func (db *DB) RecordSuccessfulConnection(ctx context.Context) error {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	return db.SetState(ctx, StateKeyLastConnected, timestamp)
}

// GetLastConnectedTime retrieves the timestamp of the last successful OBS connection.
// Returns nil if never connected.
func (db *DB) GetLastConnectedTime(ctx context.Context) (*time.Time, error) {
	value, err := db.GetState(ctx, StateKeyLastConnected)
	if err != nil {
		// Never connected
		if err.Error() == fmt.Sprintf("state key '%s' not found", StateKeyLastConnected) {
			return nil, nil
		}
		return nil, err
	}

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, fmt.Errorf("failed to parse last connected timestamp '%s': %w", value, err)
	}

	return &t, nil
}

// SetAutoReconnect saves the auto-reconnect preference.
func (db *DB) SetAutoReconnect(ctx context.Context, enabled bool) error {
	value := "false"
	if enabled {
		value = "true"
	}
	return db.SetState(ctx, StateKeyAutoReconnect, value)
}

// GetAutoReconnect retrieves the auto-reconnect preference.
// Defaults to true if not set.
func (db *DB) GetAutoReconnect(ctx context.Context) (bool, error) {
	value, err := db.GetState(ctx, StateKeyAutoReconnect)
	if err != nil {
		// Default to true if not set
		if err.Error() == fmt.Sprintf("state key '%s' not found", StateKeyAutoReconnect) {
			return true, nil
		}
		return false, err
	}

	return value == "true", nil
}

// SetAppVersion records the application version.
// This is useful for tracking which version of the app created/modified the database.
func (db *DB) SetAppVersion(ctx context.Context, version string) error {
	return db.SetState(ctx, StateKeyAppVersion, version)
}

// GetAppVersion retrieves the stored application version.
func (db *DB) GetAppVersion(ctx context.Context) (string, error) {
	return db.GetState(ctx, StateKeyAppVersion)
}

// ToolGroupConfig represents the enabled/disabled state of each tool group.
type ToolGroupConfig struct {
	Core    bool // Core OBS tools (scenes, recording, streaming, status)
	Visual  bool // Visual monitoring tools (screenshots)
	Layout  bool // Layout management tools (scene presets)
	Audio   bool // Audio control tools
	Sources bool // Source management tools
	Design  bool // Scene design tools (source creation, transforms)
}

// DefaultToolGroupConfig returns tool group config with all groups enabled.
func DefaultToolGroupConfig() ToolGroupConfig {
	return ToolGroupConfig{
		Core:    true,
		Visual:  true,
		Layout:  true,
		Audio:   true,
		Sources: true,
		Design:  true,
	}
}

// SaveToolGroupConfig persists tool group preferences to the database.
func (db *DB) SaveToolGroupConfig(ctx context.Context, cfg ToolGroupConfig) error {
	boolToStr := func(b bool) string {
		if b {
			return "true"
		}
		return "false"
	}

	if err := db.SetState(ctx, StateKeyToolsCore, boolToStr(cfg.Core)); err != nil {
		return fmt.Errorf("failed to save core tools preference: %w", err)
	}
	if err := db.SetState(ctx, StateKeyToolsVisual, boolToStr(cfg.Visual)); err != nil {
		return fmt.Errorf("failed to save visual tools preference: %w", err)
	}
	if err := db.SetState(ctx, StateKeyToolsLayout, boolToStr(cfg.Layout)); err != nil {
		return fmt.Errorf("failed to save layout tools preference: %w", err)
	}
	if err := db.SetState(ctx, StateKeyToolsAudio, boolToStr(cfg.Audio)); err != nil {
		return fmt.Errorf("failed to save audio tools preference: %w", err)
	}
	if err := db.SetState(ctx, StateKeyToolsSources, boolToStr(cfg.Sources)); err != nil {
		return fmt.Errorf("failed to save sources tools preference: %w", err)
	}
	if err := db.SetState(ctx, StateKeyToolsDesign, boolToStr(cfg.Design)); err != nil {
		return fmt.Errorf("failed to save design tools preference: %w", err)
	}

	return nil
}

// LoadToolGroupConfig retrieves tool group preferences from the database.
// Returns default config (all enabled) if preferences are not set.
func (db *DB) LoadToolGroupConfig(ctx context.Context) (ToolGroupConfig, error) {
	cfg := DefaultToolGroupConfig()

	strToBool := func(s string) bool {
		return s == "true"
	}

	// Load each preference, keeping default if not set
	if val, err := db.GetState(ctx, StateKeyToolsCore); err == nil {
		cfg.Core = strToBool(val)
	}
	if val, err := db.GetState(ctx, StateKeyToolsVisual); err == nil {
		cfg.Visual = strToBool(val)
	}
	if val, err := db.GetState(ctx, StateKeyToolsLayout); err == nil {
		cfg.Layout = strToBool(val)
	}
	if val, err := db.GetState(ctx, StateKeyToolsAudio); err == nil {
		cfg.Audio = strToBool(val)
	}
	if val, err := db.GetState(ctx, StateKeyToolsSources); err == nil {
		cfg.Sources = strToBool(val)
	}
	if val, err := db.GetState(ctx, StateKeyToolsDesign); err == nil {
		cfg.Design = strToBool(val)
	}

	return cfg, nil
}

// WebServerConfig represents HTTP server configuration.
type WebServerConfig struct {
	Enabled bool   // Whether HTTP server is enabled
	Host    string // HTTP server host (default: localhost)
	Port    int    // HTTP server port (default: 8765)
}

// DefaultWebServerConfig returns webserver config with defaults.
func DefaultWebServerConfig() WebServerConfig {
	return WebServerConfig{
		Enabled: true,
		Host:    "localhost",
		Port:    8765,
	}
}

// SaveWebServerConfig persists webserver configuration to the database.
func (db *DB) SaveWebServerConfig(ctx context.Context, cfg WebServerConfig) error {
	return db.Transaction(ctx, func(tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx, `
			INSERT INTO config (key, value, updated_at)
			VALUES (?, ?, CURRENT_TIMESTAMP)
			ON CONFLICT(key) DO UPDATE SET
				value = excluded.value,
				updated_at = CURRENT_TIMESTAMP
		`)
		if err != nil {
			return fmt.Errorf("failed to prepare config upsert: %w", err)
		}
		defer stmt.Close()

		enabledStr := "false"
		if cfg.Enabled {
			enabledStr = "true"
		}

		if _, err := stmt.ExecContext(ctx, ConfigKeyWebServerEnabled, enabledStr); err != nil {
			return fmt.Errorf("failed to save web server enabled: %w", err)
		}
		if _, err := stmt.ExecContext(ctx, ConfigKeyWebServerHost, cfg.Host); err != nil {
			return fmt.Errorf("failed to save web server host: %w", err)
		}
		if _, err := stmt.ExecContext(ctx, ConfigKeyWebServerPort, fmt.Sprintf("%d", cfg.Port)); err != nil {
			return fmt.Errorf("failed to save web server port: %w", err)
		}

		return nil
	})
}

// LoadWebServerConfig retrieves webserver configuration from the database.
// Returns default config if not set.
func (db *DB) LoadWebServerConfig(ctx context.Context) (WebServerConfig, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	cfg := DefaultWebServerConfig()

	// Load enabled
	var enabledStr string
	err := db.conn.QueryRowContext(ctx,
		"SELECT value FROM config WHERE key = ?",
		ConfigKeyWebServerEnabled,
	).Scan(&enabledStr)
	if err == nil {
		cfg.Enabled = enabledStr == "true"
	}

	// Load host
	var host string
	err = db.conn.QueryRowContext(ctx,
		"SELECT value FROM config WHERE key = ?",
		ConfigKeyWebServerHost,
	).Scan(&host)
	if err == nil && host != "" {
		cfg.Host = host
	}

	// Load port
	var portStr string
	err = db.conn.QueryRowContext(ctx,
		"SELECT value FROM config WHERE key = ?",
		ConfigKeyWebServerPort,
	).Scan(&portStr)
	if err == nil {
		var port int
		if _, parseErr := fmt.Sscanf(portStr, "%d", &port); parseErr == nil && port > 0 {
			cfg.Port = port
		}
	}

	return cfg, nil
}
