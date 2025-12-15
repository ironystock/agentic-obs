package mcp

import (
	"context"
	"fmt"
	"log"

	"github.com/ironystock/agentic-obs/internal/obs"
	"github.com/ironystock/agentic-obs/internal/storage"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server represents the MCP server instance for OBS control
type Server struct {
	mcpServer *mcpsdk.Server
	obsClient OBSClient
	storage   *storage.DB
	ctx       context.Context
	cancel    context.CancelFunc
}

// ServerConfig holds configuration for server initialization
type ServerConfig struct {
	ServerName    string
	ServerVersion string
	OBSHost       string
	OBSPort       string
	OBSPassword   string
	DBPath        string
}

// NewServer creates a new MCP server instance
func NewServer(config ServerConfig) (*Server, error) {
	ctx, cancel := context.WithCancel(context.Background())

	s := &Server{
		ctx:    ctx,
		cancel: cancel,
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

	// Create MCP server
	mcpServer := mcpsdk.NewServer(
		&mcpsdk.Implementation{
			Name:    config.ServerName,
			Version: config.ServerVersion,
		},
		nil, // No additional server options for now
	)
	s.mcpServer = mcpServer

	// Register resource handlers
	s.registerResourceHandlers()

	// Register tool handlers
	s.registerToolHandlers()

	log.Printf("MCP Server initialized: %s v%s", config.ServerName, config.ServerVersion)
	return s, nil
}

// Start establishes OBS connection
func (s *Server) Start() error {
	// Connect to OBS
	if err := s.obsClient.Connect(); err != nil {
		return fmt.Errorf("failed to connect to OBS: %w", err)
	}

	log.Println("Connected to OBS successfully")
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
