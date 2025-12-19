package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	agenthttp "github.com/ironystock/agentic-obs/internal/http"
	"github.com/ironystock/agentic-obs/internal/obs"
	"github.com/ironystock/agentic-obs/internal/screenshot"
	"github.com/ironystock/agentic-obs/internal/storage"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server represents the MCP server instance for OBS control
type Server struct {
	mcpServer     *mcpsdk.Server
	obsClient     OBSClient
	storage       *storage.DB
	screenshotMgr *screenshot.Manager
	httpServer    *agenthttp.Server
	toolGroups    ToolGroupConfig
	ctx           context.Context
	cancel        context.CancelFunc
}

// ServerConfig holds configuration for server initialization
type ServerConfig struct {
	ServerName    string
	ServerVersion string
	OBSHost       string
	OBSPort       string
	OBSPassword   string
	DBPath        string
	HTTPHost      string // HTTP server host for screenshot serving (default: localhost)
	HTTPPort      int    // HTTP server port for screenshot serving (default: 8765)
	HTTPEnabled   bool   // Whether to enable HTTP server (default: true)
	ToolGroups    ToolGroupConfig
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

// DefaultToolGroupConfig returns config with all tool groups enabled
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

// NewServer creates a new MCP server instance
func NewServer(config ServerConfig) (*Server, error) {
	ctx, cancel := context.WithCancel(context.Background())

	s := &Server{
		ctx:        ctx,
		cancel:     cancel,
		toolGroups: config.ToolGroups,
	}

	// Initialize storage
	db, err := storage.New(ctx, storage.Config{Path: config.DBPath})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}
	s.storage = db

	// Initialize OBS client
	obsClient := obs.NewClient(obs.ConnectionConfig{
		Host:     config.OBSHost,
		Port:     config.OBSPort,
		Password: config.OBSPassword,
	})

	// Set up event callback to dispatch MCP notifications
	eventHandler := obs.NewEventHandler(s.handleOBSEventNotification)
	obsClient.SetEventCallback(eventHandler)

	s.obsClient = obsClient

	// Initialize HTTP server for screenshot serving (if enabled)
	if config.HTTPEnabled {
		httpCfg := agenthttp.DefaultConfig()
		if config.HTTPHost != "" {
			httpCfg.Host = config.HTTPHost
		}
		if config.HTTPPort != 0 {
			httpCfg.Port = config.HTTPPort
		}
		s.httpServer = agenthttp.NewServer(db, httpCfg)
	}

	// Initialize screenshot manager (works without HTTP server for MCP resource access)
	screenshotCfg := screenshot.DefaultConfig()
	s.screenshotMgr = screenshot.NewManager(obsClient, db, screenshotCfg)

	// Create MCP server with completion handler
	mcpServer := mcpsdk.NewServer(
		&mcpsdk.Implementation{
			Name:    config.ServerName,
			Version: config.ServerVersion,
		},
		&mcpsdk.ServerOptions{
			CompletionHandler: s.handleCompletion,
		},
	)
	s.mcpServer = mcpServer

	// Register resource handlers
	s.registerResourceHandlers()

	// Register tool handlers (conditional based on tool groups)
	s.registerToolHandlers()

	// Register prompt handlers
	s.registerPrompts()

	log.Printf("MCP Server initialized: %s v%s", config.ServerName, config.ServerVersion)
	return s, nil
}

// Start establishes OBS connection and starts background services
func (s *Server) Start() error {
	// Connect to OBS
	if err := s.obsClient.Connect(); err != nil {
		return fmt.Errorf("failed to connect to OBS: %w", err)
	}
	log.Println("Connected to OBS successfully")

	// Start HTTP server for screenshot serving (if enabled)
	if s.httpServer != nil {
		if err := s.httpServer.Start(); err != nil {
			return fmt.Errorf("failed to start HTTP server: %w", err)
		}
		log.Printf("HTTP server started at %s", s.httpServer.GetAddr())
	} else {
		log.Println("HTTP server disabled")
	}

	// Start screenshot manager (if visual tools enabled)
	if s.toolGroups.Visual {
		if err := s.screenshotMgr.Start(s.ctx); err != nil {
			return fmt.Errorf("failed to start screenshot manager: %w", err)
		}
		log.Printf("Screenshot manager started with %d workers", s.screenshotMgr.GetWorkerCount())
	}

	return nil
}

// Run starts the MCP server and blocks until context is cancelled
func (s *Server) Run() error {
	// Create stdio transport
	transport := &mcpsdk.StdioTransport{}

	// Run server (blocks until context is cancelled)
	if err := s.mcpServer.Run(s.ctx, transport); err != nil {
		return fmt.Errorf("MCP server error: %w", err)
	}

	return nil
}

// Stop performs graceful shutdown of the server
func (s *Server) Stop() error {
	log.Println("Shutting down MCP server...")

	// Cancel context to stop all operations
	s.cancel()

	// Stop screenshot manager first (depends on OBS connection)
	if s.screenshotMgr != nil {
		s.screenshotMgr.Stop()
		log.Println("Screenshot manager stopped")
	}

	// Stop HTTP server
	if s.httpServer != nil {
		if err := s.httpServer.Stop(context.Background()); err != nil {
			log.Printf("Warning: error stopping HTTP server: %v", err)
		} else {
			log.Println("HTTP server stopped")
		}
	}

	// Disconnect from OBS
	if err := s.obsClient.Disconnect(); err != nil {
		log.Printf("Warning: error disconnecting from OBS: %v", err)
	}

	// Close storage
	if err := s.storage.Close(); err != nil {
		log.Printf("Warning: error closing storage: %v", err)
	}

	log.Println("MCP server stopped")
	return nil
}

// recordAction logs a tool action to the action history database.
// This should be called at the end of each tool handler.
func (s *Server) recordAction(toolName, action string, input interface{}, output interface{}, success bool, duration time.Duration) {
	// Skip if storage is not initialized (e.g., in tests)
	if s.storage == nil {
		return
	}

	// Convert input/output to JSON strings
	inputStr := ""
	if input != nil {
		if b, err := json.Marshal(input); err == nil {
			inputStr = string(b)
		}
	}

	outputStr := ""
	if output != nil {
		if b, err := json.Marshal(output); err == nil {
			outputStr = string(b)
		}
	}

	record := storage.ActionRecord{
		Action:     action,
		ToolName:   toolName,
		Input:      inputStr,
		Output:     outputStr,
		Success:    success,
		DurationMs: duration.Milliseconds(),
	}

	if _, err := s.storage.RecordAction(s.ctx, record); err != nil {
		log.Printf("Warning: failed to record action history: %v", err)
	}
}

// SendResourceUpdated notifies clients that a specific resource has been updated
func (s *Server) SendResourceUpdated(ctx context.Context, uri string) error {
	return s.mcpServer.ResourceUpdated(ctx, &mcpsdk.ResourceUpdatedNotificationParams{
		URI: uri,
	})
}

// handleOBSEventNotification processes OBS event notifications and dispatches MCP resource notifications
func (s *Server) handleOBSEventNotification(eventType obs.EventType, data map[string]interface{}) {
	ctx := s.ctx

	// Check if list changed (scene created or removed)
	if obs.ShouldTriggerListChanged(eventType) {
		// Resource list changed - clients should re-list resources
		// The MCP SDK handles this automatically when resources are added/removed dynamically
		log.Printf("Scene list changed for event: %s", eventType)
	}

	// Check if specific resource updated (scene changed)
	if obs.ShouldTriggerResourceUpdated(eventType) {
		if sceneName, ok := data["scene_name"].(string); ok {
			uri := obs.GetResourceURIForScene(sceneName)
			if err := s.SendResourceUpdated(ctx, uri); err != nil {
				log.Printf("Error sending resource updated notification: %v", err)
			} else {
				log.Printf("Sent resource updated notification for scene: %s", sceneName)
			}
		}
	}
}

// GetOBSClient returns the OBS client instance (for internal use)
func (s *Server) GetOBSClient() OBSClient {
	return s.obsClient
}

// SetOBSClient sets the OBS client (primarily for testing with mocks)
func (s *Server) SetOBSClient(client OBSClient) {
	s.obsClient = client
}

// GetStorage returns the storage instance (for internal use)
func (s *Server) GetStorage() *storage.DB {
	return s.storage
}

// GetScreenshotManager returns the screenshot manager instance (for internal use)
func (s *Server) GetScreenshotManager() *screenshot.Manager {
	return s.screenshotMgr
}

// GetHTTPServer returns the HTTP server instance (for internal use)
func (s *Server) GetHTTPServer() *agenthttp.Server {
	return s.httpServer
}
