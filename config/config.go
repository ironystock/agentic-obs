package config

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ironystock/agentic-obs/internal/storage"
)

// Config holds the application configuration
type Config struct {
	// Server configuration
	ServerName    string
	ServerVersion string

	// OBS WebSocket configuration
	OBSHost     string
	OBSPort     string
	OBSPassword string

	// Database configuration
	DBPath string
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	return &Config{
		ServerName:    "agentic-obs",
		ServerVersion: "0.1.0",
		OBSHost:       "localhost",
		OBSPort:       "4455",
		OBSPassword:   "",
		DBPath:        filepath.Join(homeDir, ".agentic-obs", "db.sqlite"),
	}
}

// DetectOrPrompt attempts to auto-detect OBS configuration or prompts the user interactively.
// This is called when no valid configuration exists or the default connection fails.
func (c *Config) DetectOrPrompt() error {
	log.Println("OBS configuration needed. Starting interactive setup...")
	log.Println("Note: Make sure OBS Studio is running with WebSocket server enabled.")
	log.Println("      (Tools > WebSocket Server Settings in OBS Studio)")

	scanner := bufio.NewScanner(os.Stdin)

	// Prompt for host
	fmt.Printf("\nOBS WebSocket Host [%s]: ", c.OBSHost)
	if scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input != "" {
			c.OBSHost = input
		}
	}

	// Prompt for port
	fmt.Printf("OBS WebSocket Port [%s]: ", c.OBSPort)
	if scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input != "" {
			c.OBSPort = input
		}
	}

	// Prompt for password
	fmt.Print("OBS WebSocket Password (leave empty if none): ")
	if scanner.Scan() {
		c.OBSPassword = strings.TrimSpace(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	log.Printf("Configuration set: %s:%s", c.OBSHost, c.OBSPort)
	return nil
}

// LoadFromStorage loads configuration from the database
func LoadFromStorage(ctx context.Context, dbPath string) (*Config, error) {
	cfg := DefaultConfig()
	cfg.DBPath = dbPath

	// Open database
	db, err := storage.New(ctx, storage.Config{Path: dbPath})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Check if this is first run
	isFirstRun, err := db.IsFirstRun(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check first run status: %w", err)
	}

	if isFirstRun {
		log.Println("First run detected - using default configuration")
		return cfg, nil
	}

	// Load OBS configuration from database
	obsConfig, err := db.LoadOBSConfig(ctx)
	if err != nil {
		log.Printf("Warning: failed to load OBS config from database: %v", err)
		log.Println("Using default OBS configuration")
		return cfg, nil
	}

	cfg.OBSHost = obsConfig.Host
	cfg.OBSPort = fmt.Sprintf("%d", obsConfig.Port)
	cfg.OBSPassword = obsConfig.Password

	log.Printf("Loaded configuration from database: OBS at %s:%s", cfg.OBSHost, cfg.OBSPort)
	return cfg, nil
}

// SaveToStorage persists configuration to the database
func SaveToStorage(ctx context.Context, cfg *Config) error {
	// Open database
	db, err := storage.New(ctx, storage.Config{Path: cfg.DBPath})
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Parse port as integer
	var port int
	if _, err := fmt.Sscanf(cfg.OBSPort, "%d", &port); err != nil {
		return fmt.Errorf("invalid OBS port '%s': %w", cfg.OBSPort, err)
	}

	// Save OBS configuration
	obsConfig := storage.OBSConfig{
		Host:     cfg.OBSHost,
		Port:     port,
		Password: cfg.OBSPassword,
	}

	if err := db.SaveOBSConfig(ctx, obsConfig); err != nil {
		return fmt.Errorf("failed to save OBS config: %w", err)
	}

	// Mark first run as complete
	if err := db.MarkFirstRunComplete(ctx); err != nil {
		return fmt.Errorf("failed to mark first run complete: %w", err)
	}

	// Record app version
	if err := db.SetAppVersion(ctx, cfg.ServerVersion); err != nil {
		log.Printf("Warning: failed to save app version: %v", err)
	}

	log.Printf("Configuration saved to database: OBS at %s:%s", cfg.OBSHost, cfg.OBSPort)
	return nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.ServerName == "" {
		return fmt.Errorf("server name cannot be empty")
	}

	if c.ServerVersion == "" {
		return fmt.Errorf("server version cannot be empty")
	}

	if c.OBSHost == "" {
		return fmt.Errorf("OBS host cannot be empty")
	}

	if c.OBSPort == "" {
		return fmt.Errorf("OBS port cannot be empty")
	}

	// Validate port is a number
	var port int
	if _, err := fmt.Sscanf(c.OBSPort, "%d", &port); err != nil {
		return fmt.Errorf("OBS port must be a valid number: %w", err)
	}

	if port < 1 || port > 65535 {
		return fmt.Errorf("OBS port must be between 1 and 65535, got %d", port)
	}

	if c.DBPath == "" {
		return fmt.Errorf("database path cannot be empty")
	}

	return nil
}

// String returns a string representation of the config (masking sensitive data)
func (c *Config) String() string {
	password := "<empty>"
	if c.OBSPassword != "" {
		password = "<set>"
	}

	return fmt.Sprintf(
		"Config{ServerName: %s, ServerVersion: %s, OBS: %s:%s, Password: %s, DBPath: %s}",
		c.ServerName,
		c.ServerVersion,
		c.OBSHost,
		c.OBSPort,
		password,
		c.DBPath,
	)
}
