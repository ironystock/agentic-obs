package obs

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/events"
	"github.com/andreykaipov/goobs/api/events/subscriptions"
)

// Client wraps the OBS WebSocket client with connection management and state tracking.
type Client struct {
	// Connection configuration
	host     string
	port     string
	password string

	// OBS client instance
	client *goobs.Client

	// Connection state
	mu        sync.RWMutex
	connected bool
	reconnect bool // Whether to attempt auto-reconnect

	// Event handlers
	eventCallback EventCallback

	// Context for managing lifecycle
	ctx    context.Context
	cancel context.CancelFunc
}

// ConnectionConfig holds the parameters needed to connect to OBS.
type ConnectionConfig struct {
	Host     string
	Port     string
	Password string
}

// EventCallback is the interface for handling OBS events and triggering MCP notifications.
type EventCallback interface {
	OnSceneCreated(sceneName string)
	OnSceneRemoved(sceneName string)
	OnCurrentProgramSceneChanged(sceneName string)
}

// NewClient creates a new OBS client with the specified connection configuration.
// The client is not connected until Connect() is called.
func NewClient(config ConnectionConfig) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		host:      config.Host,
		port:      config.Port,
		password:  config.Password,
		connected: false,
		reconnect: true,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// SetEventCallback registers a callback handler for OBS events.
// This should be called before Connect() to ensure no events are missed.
func (c *Client) SetEventCallback(callback EventCallback) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.eventCallback = callback
}

// Connect establishes a connection to OBS WebSocket server.
// Returns an error if the connection fails.
func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return nil // Already connected
	}

	address := fmt.Sprintf("%s:%s", c.host, c.port)

	var client *goobs.Client
	var err error

	// Build connection options
	opts := []goobs.Option{
		goobs.WithEventSubscriptions(subscriptions.Scenes),
	}

	if c.password != "" {
		opts = append(opts, goobs.WithPassword(c.password))
	}

	client, err = goobs.New(address, opts...)
	if err != nil {
		return fmt.Errorf("OBS connection failed. Is OBS Studio running with WebSocket server enabled? Error: %w", err)
	}

	c.client = client
	c.connected = true

	// Set up event subscriptions
	if err := c.setupEventHandlers(); err != nil {
		c.client.Disconnect()
		c.client = nil
		c.connected = false
		return fmt.Errorf("failed to set up OBS event handlers: %w", err)
	}

	// Start auto-reconnect monitor if enabled
	if c.reconnect {
		go c.monitorConnection()
	}

	return nil
}

// Disconnect closes the connection to OBS WebSocket server.
func (c *Client) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return nil // Already disconnected
	}

	c.reconnect = false // Disable auto-reconnect

	if c.client != nil {
		if err := c.client.Disconnect(); err != nil {
			return fmt.Errorf("failed to disconnect from OBS: %w", err)
		}
		c.client = nil
	}

	c.connected = false
	return nil
}

// Close performs a graceful shutdown of the client, disconnecting and cleaning up resources.
func (c *Client) Close() error {
	c.cancel() // Cancel context to stop background goroutines
	return c.Disconnect()
}

// IsConnected returns true if the client is currently connected to OBS.
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// GetConnectionStatus returns detailed connection status information.
func (c *Client) GetConnectionStatus() (ConnectionStatus, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	status := ConnectionStatus{
		Connected: c.connected,
		Host:      c.host,
		Port:      c.port,
	}

	if !c.connected {
		return status, nil
	}

	// Get OBS version and additional status info
	versionResp, err := c.client.General.GetVersion()
	if err != nil {
		// Connection might be broken
		return status, fmt.Errorf("failed to get OBS version (connection may be broken): %w", err)
	}

	status.OBSVersion = versionResp.ObsVersion
	status.WebSocketVersion = versionResp.ObsWebSocketVersion
	status.Platform = versionResp.Platform

	return status, nil
}

// HealthCheck performs a lightweight health check by pinging OBS.
// Returns nil if healthy, error if unhealthy.
func (c *Client) HealthCheck() error {
	c.mu.RLock()
	client := c.client
	connected := c.connected
	c.mu.RUnlock()

	if !connected {
		return fmt.Errorf("OBS client is not connected")
	}

	// Use GetVersion as a lightweight health check
	_, err := client.General.GetVersion()
	if err != nil {
		return fmt.Errorf("OBS health check failed: %w", err)
	}

	return nil
}

// monitorConnection runs in the background and attempts to reconnect if the connection is lost.
func (c *Client) monitorConnection() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return // Client is shutting down
		case <-ticker.C:
			// Check if we need to reconnect
			c.mu.RLock()
			shouldReconnect := c.reconnect && !c.connected
			c.mu.RUnlock()

			if shouldReconnect {
				// Attempt reconnection
				if err := c.Connect(); err != nil {
					// Log error but continue trying
					// In production, you might want to use a proper logger here
					fmt.Printf("Auto-reconnect failed: %v\n", err)
				} else {
					fmt.Println("Successfully reconnected to OBS")
				}
			} else {
				// Perform health check on connected client
				c.mu.RLock()
				isConnected := c.connected
				c.mu.RUnlock()

				if isConnected {
					if err := c.HealthCheck(); err != nil {
						// Mark as disconnected so we can try to reconnect
						c.mu.Lock()
						c.connected = false
						if c.client != nil {
							c.client.Disconnect()
							c.client = nil
						}
						c.mu.Unlock()
						fmt.Printf("OBS connection lost: %v\n", err)
					}
				}
			}
		}
	}
}

// setupEventHandlers subscribes to relevant OBS events and sets up handlers.
func (c *Client) setupEventHandlers() error {
	// Events are subscribed to via WithEventSubscriptions option during client creation
	// We just need to start listening for events
	go c.handleEvents()

	return nil
}

// handleEvents processes incoming OBS events and dispatches to the callback.
func (c *Client) handleEvents() {
	for event := range c.client.IncomingEvents {
		c.mu.RLock()
		callback := c.eventCallback
		c.mu.RUnlock()

		if callback == nil {
			continue // No callback registered
		}

		// Dispatch based on event type
		switch e := event.(type) {
		case *events.SceneCreated:
			callback.OnSceneCreated(e.SceneName)

		case *events.SceneRemoved:
			callback.OnSceneRemoved(e.SceneName)

		case *events.CurrentProgramSceneChanged:
			callback.OnCurrentProgramSceneChanged(e.SceneName)

		// Additional events can be handled here as needed
		default:
			// Ignore other events for now
		}
	}
}

// getClient is a helper that returns the client if connected, or an error if not.
func (c *Client) getClient() (*goobs.Client, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected || c.client == nil {
		return nil, fmt.Errorf("not connected to OBS. Call Connect() first")
	}

	return c.client, nil
}

// ConnectionStatus holds detailed information about the OBS connection.
type ConnectionStatus struct {
	Connected        bool   `json:"connected"`
	Host             string `json:"host"`
	Port             string `json:"port"`
	OBSVersion       string `json:"obs_version,omitempty"`
	WebSocketVersion string `json:"websocket_version,omitempty"`
	Platform         string `json:"platform,omitempty"`
}
