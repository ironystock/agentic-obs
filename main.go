package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ironystock/agentic-obs/config"
	"github.com/ironystock/agentic-obs/internal/mcp"
	"github.com/ironystock/agentic-obs/internal/storage"
	"github.com/ironystock/agentic-obs/internal/tui"
)

const appName = "agentic-obs"

func main() {
	// Parse command-line flags
	tuiMode := flag.Bool("tui", false, "Run in TUI dashboard mode instead of MCP server mode")
	flag.BoolVar(tuiMode, "t", false, "Run in TUI dashboard mode (shorthand)")
	showHelp := flag.Bool("help", false, "Show usage information")
	flag.BoolVar(showHelp, "h", false, "Show usage information (shorthand)")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.BoolVar(showVersion, "v", false, "Show version information (shorthand)")
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s %s\n", appName, Version())
		os.Exit(0)
	}

	if *showHelp {
		printUsage()
		os.Exit(0)
	}

	// Set up logging with timestamps and source location
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("========================================")
	log.Printf("Starting %s v%s", appName, Version())
	if commit != "none" {
		log.Printf("Build: %s (%s)", commit, date)
	}
	log.Printf("========================================")

	// Create context for initialization
	ctx := context.Background()

	// Step 1: Initialize storage layer
	log.Println("[1/5] Initializing storage layer...")
	cfg, err := loadConfig(ctx)
	if err != nil {
		log.Fatalf("FATAL: Failed to load configuration: %v", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("FATAL: Invalid configuration: %v", err)
	}

	log.Printf("Configuration loaded: %s", cfg)

	// Check if TUI mode is requested
	if *tuiMode {
		log.Println("Starting TUI dashboard mode...")
		if err := runTUIMode(ctx, cfg); err != nil {
			log.Fatalf("TUI error: %v", err)
		}
		return
	}

	// Step 2: Load or detect OBS configuration
	log.Println("[2/5] Checking OBS connection configuration...")

	// Step 3: Create MCP server (includes storage and OBS client initialization)
	log.Println("[3/5] Initializing MCP server...")
	serverConfig := mcp.ServerConfig{
		ServerName:        cfg.ServerName,
		ServerVersion:     cfg.ServerVersion,
		OBSHost:           cfg.OBSHost,
		OBSPort:           cfg.OBSPort,
		OBSPassword:       cfg.OBSPassword,
		DBPath:            cfg.DBPath,
		HTTPEnabled:       cfg.WebServer.Enabled,
		HTTPHost:          cfg.WebServer.Host,
		HTTPPort:          cfg.WebServer.Port,
		ThumbnailCacheSec: cfg.WebServer.ThumbnailCacheSec,
		ToolGroups: mcp.ToolGroupConfig{
			Core:        cfg.ToolGroups.Core,
			Visual:      cfg.ToolGroups.Visual,
			Layout:      cfg.ToolGroups.Layout,
			Audio:       cfg.ToolGroups.Audio,
			Sources:     cfg.ToolGroups.Sources,
			Design:      cfg.ToolGroups.Design,
			Filters:     cfg.ToolGroups.Filters,
			Transitions: cfg.ToolGroups.Transitions,
		},
	}

	server, err := mcp.NewServer(serverConfig)
	if err != nil {
		log.Fatalf("FATAL: Failed to create MCP server: %v", err)
	}
	defer func() {
		log.Println("Cleaning up resources...")
		if err := server.Stop(); err != nil {
			log.Printf("Warning during shutdown: %v", err)
		}
	}()

	// Step 4: Connect to OBS with retry logic
	log.Println("[4/5] Connecting to OBS WebSocket...")
	if err := connectToOBSWithRetry(server, cfg, ctx); err != nil {
		log.Fatalf("FATAL: Failed to establish OBS connection: %v", err)
	}

	// Record successful connection in storage
	if err := recordSuccessfulConnection(ctx, cfg); err != nil {
		log.Printf("Warning: Failed to record successful connection: %v", err)
	}

	// Step 5: Start MCP server and setup event loop
	log.Println("[5/5] Starting MCP server event loop...")

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("\nReceived shutdown signal: %v", sig)
		log.Println("Initiating graceful shutdown...")
		server.Stop()
	}()

	// Run the MCP server (blocks until shutdown)
	log.Println("========================================")
	log.Println("MCP server is running on stdio")
	log.Println("Waiting for client connections...")
	log.Println("Press Ctrl+C to stop")
	log.Println("========================================")

	if err := server.Run(); err != nil {
		log.Printf("MCP server error: %v", err)
		os.Exit(1)
	}

	log.Println("========================================")
	log.Println("Server stopped successfully")
	log.Println("========================================")
}

// loadConfig loads configuration from storage or returns default config
func loadConfig(ctx context.Context) (*config.Config, error) {
	// Get default config to determine database path
	defaultCfg := config.DefaultConfig()
	defaultCfg.ServerName = appName
	defaultCfg.ServerVersion = version

	// Check if this is first run (before loading from storage)
	isFirstRun, err := checkFirstRun(ctx, defaultCfg.DBPath)
	if err != nil {
		log.Printf("Warning: failed to check first run status: %v", err)
	}

	// Try to load from storage
	cfg, err := config.LoadFromStorage(ctx, defaultCfg.DBPath)
	if err != nil {
		log.Printf("Warning: failed to load config from storage: %v", err)
		log.Println("Using default configuration")
		cfg = defaultCfg
	}

	// Ensure server name and version are set correctly
	cfg.ServerName = appName
	cfg.ServerVersion = version

	// Apply environment variable overrides (takes precedence over stored config)
	cfg.ApplyEnvOverrides()

	// Run first-run setup prompts if this is first run
	if isFirstRun {
		log.Println("First run detected - running initial setup...")
		if err := cfg.PromptFirstRunSetup(); err != nil {
			log.Printf("Warning: first-run setup failed: %v", err)
		}
	}

	// Save configuration if it was loaded from defaults or overridden by environment
	if err := config.SaveToStorage(ctx, cfg); err != nil {
		log.Printf("Warning: failed to save configuration: %v", err)
	}

	return cfg, nil
}

// checkFirstRun checks if this is the first run of the application
func checkFirstRun(ctx context.Context, dbPath string) (bool, error) {
	db, err := storage.New(ctx, storage.Config{Path: dbPath})
	if err != nil {
		return true, err // Assume first run if we can't check
	}
	defer db.Close()

	return db.IsFirstRun(ctx)
}

// connectToOBSWithRetry attempts to connect to OBS with retry logic and user prompting
func connectToOBSWithRetry(server *mcp.Server, cfg *config.Config, ctx context.Context) error {
	const maxRetries = 3
	retryDelay := 2 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("Connection attempt %d/%d to OBS at %s:%s...", attempt, maxRetries, cfg.OBSHost, cfg.OBSPort)

		err := server.Start()
		if err == nil {
			log.Println("Successfully connected to OBS!")
			return nil
		}

		log.Printf("Connection failed: %v", err)

		// If this isn't the last attempt, wait before retrying
		if attempt < maxRetries {
			log.Printf("Retrying in %v...", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	// All automatic retries failed - prompt user for configuration
	log.Println("========================================")
	log.Println("Failed to connect to OBS after multiple attempts.")
	log.Println("This could mean:")
	log.Println("  1. OBS Studio is not running")
	log.Println("  2. OBS WebSocket server is not enabled")
	log.Println("  3. Connection details are incorrect")
	log.Println("========================================")
	fmt.Print("\nWould you like to reconfigure OBS connection settings? (y/n): ")

	var response string
	fmt.Scanln(&response)

	if response == "y" || response == "Y" || response == "yes" {
		// Prompt for new configuration
		if err := cfg.DetectOrPrompt(); err != nil {
			return fmt.Errorf("failed to get new configuration: %w", err)
		}

		// Validate and save new configuration
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("invalid configuration: %w", err)
		}

		if err := config.SaveToStorage(ctx, cfg); err != nil {
			log.Printf("Warning: failed to save new configuration: %v", err)
		}

		// Try connecting with new configuration
		log.Println("Attempting to connect with new configuration...")

		// Need to recreate server with new config
		// For now, just try to start with existing server
		// (In a real implementation, you might want to update the OBS client configuration)
		if err := server.Start(); err != nil {
			return fmt.Errorf("connection failed with new configuration: %w\nPlease verify OBS is running and WebSocket server is enabled", err)
		}

		log.Println("Successfully connected to OBS with new configuration!")
		return nil
	}

	return fmt.Errorf("OBS connection failed and user declined reconfiguration")
}

// recordSuccessfulConnection records the successful OBS connection timestamp in storage
func recordSuccessfulConnection(ctx context.Context, cfg *config.Config) error {
	db, err := storage.New(ctx, storage.Config{Path: cfg.DBPath})
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	if err := db.RecordSuccessfulConnection(ctx); err != nil {
		return fmt.Errorf("failed to record connection timestamp: %w", err)
	}

	log.Println("Recorded successful OBS connection in database")
	return nil
}

// runTUIMode starts the terminal user interface dashboard
func runTUIMode(ctx context.Context, cfg *config.Config) error {
	// Initialize storage
	db, err := storage.New(ctx, storage.Config{Path: cfg.DBPath})
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer db.Close()

	// Create and run TUI
	app := tui.New(db, cfg, appName, version)
	return app.Run()
}

// printUsage prints usage information
func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: %s [OPTIONS]

An MCP server that provides AI assistants with programmatic control over OBS Studio.

Options:
  -t, --tui       Run in TUI dashboard mode instead of MCP server mode
  -v, --version   Show version information
  -h, --help      Show this help message

Environment Variables:
  OBS_HOST                   OBS WebSocket host (default: localhost)
  OBS_PORT                   OBS WebSocket port (default: 4455)
  OBS_PASSWORD               OBS WebSocket password (default: empty)
  AGENTIC_OBS_DB             Database file path (default: ~/.agentic-obs/db.sqlite)
  AGENTIC_OBS_HTTP_PORT      HTTP server port (default: 8765)
  AGENTIC_OBS_HTTP_ENABLED   Enable/disable HTTP server (default: true)

Examples:
  # Run MCP server (default mode)
  %s

  # Run TUI dashboard
  %s --tui

  # Run with custom OBS host and port
  OBS_HOST=192.168.1.100 OBS_PORT=4456 %s

  # Run with password
  OBS_PASSWORD=mysecret %s

For more information, see: https://github.com/ironystock/agentic-obs
`, appName, appName, appName, appName, appName)
}
