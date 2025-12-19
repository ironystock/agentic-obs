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

	// Tool group preferences
	ToolGroups ToolGroupConfig

	// HTTP server configuration
	WebServer WebServerConfig
}

// ToolGroupConfig controls which tool categories are enabled
type ToolGroupConfig struct {
	Core    bool // Core OBS tools (scenes, recording, streaming, status)
	Visual  bool // Visual monitoring tools (screenshots)
	Layout  bool // Layout management tools (scene presets)
	Audio   bool // Audio control tools
	Sources bool // Source management tools
	Design  bool // Scene design tools (source creation, transforms)
}

// WebServerConfig controls HTTP server settings
type WebServerConfig struct {
	Enabled bool   // Whether HTTP server is enabled
	Host    string // HTTP server host
	Port    int    // HTTP server port
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
		ToolGroups: ToolGroupConfig{
			Core:    true,
			Visual:  true,
			Layout:  true,
			Audio:   true,
			Sources: true,
			Design:  true,
		},
		WebServer: WebServerConfig{
			Enabled: true,
			Host:    "localhost",
			Port:    8765,
		},
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

// PromptFirstRunSetup performs initial setup prompts for tool groups and webserver.
// This is called during first run to let users customize their experience.
func (c *Config) PromptFirstRunSetup() error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("\n=== Feature Configuration ===")
	fmt.Println("Choose which features to enable (press Enter for defaults):")

	// Helper function to prompt yes/no
	promptBool := func(prompt string, defaultVal bool) bool {
		defaultStr := "Y/n"
		if !defaultVal {
			defaultStr = "y/N"
		}
		fmt.Printf("%s [%s]: ", prompt, defaultStr)
		if scanner.Scan() {
			input := strings.TrimSpace(strings.ToLower(scanner.Text()))
			if input == "" {
				return defaultVal
			}
			return input == "y" || input == "yes"
		}
		return defaultVal
	}

	// Tool group prompts
	fmt.Println("\n--- Tool Groups ---")
	c.ToolGroups.Core = promptBool("Core OBS control (scenes, recording, streaming)", c.ToolGroups.Core)
	c.ToolGroups.Visual = promptBool("Visual monitoring (screenshot capture)", c.ToolGroups.Visual)
	c.ToolGroups.Layout = promptBool("Layout management (scene presets)", c.ToolGroups.Layout)
	c.ToolGroups.Audio = promptBool("Audio control (mute, volume)", c.ToolGroups.Audio)
	c.ToolGroups.Sources = promptBool("Source management (visibility, settings)", c.ToolGroups.Sources)
	c.ToolGroups.Design = promptBool("Scene design (create sources, transforms)", c.ToolGroups.Design)

	// Webserver prompt
	fmt.Println("\n--- HTTP Server ---")
	c.WebServer.Enabled = promptBool("Enable HTTP server for screenshot URLs", c.WebServer.Enabled)

	if c.WebServer.Enabled {
		fmt.Printf("HTTP server port [%d]: ", c.WebServer.Port)
		if scanner.Scan() {
			input := strings.TrimSpace(scanner.Text())
			if input != "" {
				var port int
				if _, err := fmt.Sscanf(input, "%d", &port); err == nil && port > 0 && port < 65536 {
					c.WebServer.Port = port
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	// Summary
	fmt.Println("\n--- Configuration Summary ---")
	fmt.Printf("Core tools: %v\n", c.ToolGroups.Core)
	fmt.Printf("Visual tools: %v\n", c.ToolGroups.Visual)
	fmt.Printf("Layout tools: %v\n", c.ToolGroups.Layout)
	fmt.Printf("Audio tools: %v\n", c.ToolGroups.Audio)
	fmt.Printf("Source tools: %v\n", c.ToolGroups.Sources)
	fmt.Printf("Design tools: %v\n", c.ToolGroups.Design)
	fmt.Printf("HTTP server: %v", c.WebServer.Enabled)
	if c.WebServer.Enabled {
		fmt.Printf(" (port %d)", c.WebServer.Port)
	}
	fmt.Println()

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
	} else {
		cfg.OBSHost = obsConfig.Host
		cfg.OBSPort = fmt.Sprintf("%d", obsConfig.Port)
		cfg.OBSPassword = obsConfig.Password
	}

	// Load tool group preferences
	toolGroups, err := db.LoadToolGroupConfig(ctx)
	if err != nil {
		log.Printf("Warning: failed to load tool group config: %v", err)
	} else {
		cfg.ToolGroups = ToolGroupConfig{
			Core:    toolGroups.Core,
			Visual:  toolGroups.Visual,
			Layout:  toolGroups.Layout,
			Audio:   toolGroups.Audio,
			Sources: toolGroups.Sources,
			Design:  toolGroups.Design,
		}
	}

	// Load webserver configuration
	webCfg, err := db.LoadWebServerConfig(ctx)
	if err != nil {
		log.Printf("Warning: failed to load webserver config: %v", err)
	} else {
		cfg.WebServer = WebServerConfig{
			Enabled: webCfg.Enabled,
			Host:    webCfg.Host,
			Port:    webCfg.Port,
		}
	}

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

	// Save tool group preferences
	toolGroups := storage.ToolGroupConfig{
		Core:    cfg.ToolGroups.Core,
		Visual:  cfg.ToolGroups.Visual,
		Layout:  cfg.ToolGroups.Layout,
		Audio:   cfg.ToolGroups.Audio,
		Sources: cfg.ToolGroups.Sources,
		Design:  cfg.ToolGroups.Design,
	}
	if err := db.SaveToolGroupConfig(ctx, toolGroups); err != nil {
		return fmt.Errorf("failed to save tool group config: %w", err)
	}

	// Save webserver configuration
	webCfg := storage.WebServerConfig{
		Enabled: cfg.WebServer.Enabled,
		Host:    cfg.WebServer.Host,
		Port:    cfg.WebServer.Port,
	}
	if err := db.SaveWebServerConfig(ctx, webCfg); err != nil {
		return fmt.Errorf("failed to save webserver config: %w", err)
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
