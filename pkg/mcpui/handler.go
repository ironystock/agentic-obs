// Copyright 2025 The MCP-UI Go SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package mcpui

import (
	"context"
	"fmt"
	"sync"
)

// UIActionHandler handles UI actions from embedded resources.
// This follows the pattern of mcp.ResourceHandler and mcp.ToolHandler.
type UIActionHandler func(context.Context, *UIActionRequest) (*UIActionResult, error)

// UIActionRequest contains the action and session context.
type UIActionRequest struct {
	// Action is the UI action to process.
	Action *UIAction
	// ResourceURI is the URI of the resource that triggered the action.
	ResourceURI string
	// Session can hold session-specific data (e.g., mcp.ServerSession).
	Session any
}

// UIActionResult is the result of handling a UI action.
type UIActionResult struct {
	// Response contains the successful result data.
	Response any
	// Error contains error information if the action failed.
	Error error
}

// ToUIResponse converts the result to a UIResponse.
func (r *UIActionResult) ToUIResponse(messageID string) *UIResponse {
	if r.Error != nil {
		return NewErrorResponse(messageID, r.Error)
	}
	return NewSuccessResponse(messageID, r.Response)
}

// Router dispatches UI actions to appropriate handlers.
// It provides a way to register handlers for different action types and resources.
type Router struct {
	mu sync.RWMutex
	// handlers by action type
	typeHandlers map[string]UIActionHandler
	// handlers by resource URI pattern
	resourceHandlers map[string]UIActionHandler
	// default handler for unmatched actions
	defaultHandler UIActionHandler
}

// NewRouter creates a new Router.
func NewRouter() *Router {
	return &Router{
		typeHandlers:     make(map[string]UIActionHandler),
		resourceHandlers: make(map[string]UIActionHandler),
	}
}

// HandleType registers a handler for a specific action type.
func (r *Router) HandleType(actionType string, handler UIActionHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.typeHandlers[actionType] = handler
}

// HandleResource registers a handler for a specific resource URI.
func (r *Router) HandleResource(resourceURI string, handler UIActionHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.resourceHandlers[resourceURI] = handler
}

// SetDefault sets the default handler for unmatched actions.
func (r *Router) SetDefault(handler UIActionHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.defaultHandler = handler
}

// Dispatch routes an action to the appropriate handler.
// Priority order:
// 1. Resource-specific handler (exact URI match)
// 2. Action type handler
// 3. Default handler
func (r *Router) Dispatch(ctx context.Context, req *UIActionRequest) (*UIActionResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Check for resource-specific handler first
	if req.ResourceURI != "" {
		if handler, ok := r.resourceHandlers[req.ResourceURI]; ok {
			return handler(ctx, req)
		}
	}

	// Check for action type handler
	if req.Action != nil {
		if handler, ok := r.typeHandlers[req.Action.Type]; ok {
			return handler(ctx, req)
		}
	}

	// Fall back to default handler
	if r.defaultHandler != nil {
		return r.defaultHandler(ctx, req)
	}

	return nil, fmt.Errorf("no handler for action type %q from resource %q", req.Action.Type, req.ResourceURI)
}

// Handle implements the UIActionHandler interface, making Router itself a handler.
func (r *Router) Handle(ctx context.Context, req *UIActionRequest) (*UIActionResult, error) {
	return r.Dispatch(ctx, req)
}

// ToolHandler is a convenience type for handling tool actions.
// It is called when an embedded UI requests a tool execution.
type ToolHandler func(ctx context.Context, toolName string, params map[string]any) (any, error)

// WrapToolHandler wraps a ToolHandler as a UIActionHandler.
func WrapToolHandler(handler ToolHandler) UIActionHandler {
	return func(ctx context.Context, req *UIActionRequest) (*UIActionResult, error) {
		if req.Action.Type != ActionTypeTool {
			return nil, fmt.Errorf("expected tool action, got %s", req.Action.Type)
		}
		payload, err := req.Action.ToolPayload()
		if err != nil {
			return nil, err
		}
		result, err := handler(ctx, payload.ToolName, payload.Params)
		if err != nil {
			return &UIActionResult{Error: err}, nil
		}
		return &UIActionResult{Response: result}, nil
	}
}

// IntentHandler is a convenience type for handling intent actions.
// It is called when an embedded UI signals a user intent.
type IntentHandler func(ctx context.Context, intent string, params map[string]any) (any, error)

// WrapIntentHandler wraps an IntentHandler as a UIActionHandler.
func WrapIntentHandler(handler IntentHandler) UIActionHandler {
	return func(ctx context.Context, req *UIActionRequest) (*UIActionResult, error) {
		if req.Action.Type != ActionTypeIntent {
			return nil, fmt.Errorf("expected intent action, got %s", req.Action.Type)
		}
		payload, err := req.Action.IntentPayload()
		if err != nil {
			return nil, err
		}
		result, err := handler(ctx, payload.Intent, payload.Params)
		if err != nil {
			return &UIActionResult{Error: err}, nil
		}
		return &UIActionResult{Response: result}, nil
	}
}

// PromptHandler is a convenience type for handling prompt actions.
// It is called when an embedded UI sends a prompt message.
type PromptHandler func(ctx context.Context, prompt string) (any, error)

// WrapPromptHandler wraps a PromptHandler as a UIActionHandler.
func WrapPromptHandler(handler PromptHandler) UIActionHandler {
	return func(ctx context.Context, req *UIActionRequest) (*UIActionResult, error) {
		if req.Action.Type != ActionTypePrompt {
			return nil, fmt.Errorf("expected prompt action, got %s", req.Action.Type)
		}
		payload, err := req.Action.PromptPayload()
		if err != nil {
			return nil, err
		}
		result, err := handler(ctx, payload.Prompt)
		if err != nil {
			return &UIActionResult{Error: err}, nil
		}
		return &UIActionResult{Response: result}, nil
	}
}

// NotifyHandler is a convenience type for handling notify actions.
// It is called when an embedded UI sends a notification.
type NotifyHandler func(ctx context.Context, message string, level string) error

// WrapNotifyHandler wraps a NotifyHandler as a UIActionHandler.
func WrapNotifyHandler(handler NotifyHandler) UIActionHandler {
	return func(ctx context.Context, req *UIActionRequest) (*UIActionResult, error) {
		if req.Action.Type != ActionTypeNotify {
			return nil, fmt.Errorf("expected notify action, got %s", req.Action.Type)
		}
		payload, err := req.Action.NotifyPayload()
		if err != nil {
			return nil, err
		}
		if err := handler(ctx, payload.Message, payload.Level); err != nil {
			return &UIActionResult{Error: err}, nil
		}
		return &UIActionResult{Response: "acknowledged"}, nil
	}
}

// LinkHandler is a convenience type for handling link actions.
// It is called when an embedded UI requests to open a link.
type LinkHandler func(ctx context.Context, url string) error

// WrapLinkHandler wraps a LinkHandler as a UIActionHandler.
func WrapLinkHandler(handler LinkHandler) UIActionHandler {
	return func(ctx context.Context, req *UIActionRequest) (*UIActionResult, error) {
		if req.Action.Type != ActionTypeLink {
			return nil, fmt.Errorf("expected link action, got %s", req.Action.Type)
		}
		payload, err := req.Action.LinkPayload()
		if err != nil {
			return nil, err
		}
		if err := handler(ctx, payload.URL); err != nil {
			return &UIActionResult{Error: err}, nil
		}
		return &UIActionResult{Response: "opened"}, nil
	}
}

// UISizeHandler is a convenience type for handling UI size change actions.
// It is called when an embedded UI reports a size change.
type UISizeHandler func(ctx context.Context, height, width int) error

// WrapUISizeHandler wraps a UISizeHandler as a UIActionHandler.
func WrapUISizeHandler(handler UISizeHandler) UIActionHandler {
	return func(ctx context.Context, req *UIActionRequest) (*UIActionResult, error) {
		if req.Action.Type != ActionTypeUISize {
			return nil, fmt.Errorf("expected ui-size-change action, got %s", req.Action.Type)
		}
		payload, err := req.Action.UISizePayload()
		if err != nil {
			return nil, err
		}
		if err := handler(ctx, payload.Height, payload.Width); err != nil {
			return &UIActionResult{Error: err}, nil
		}
		return &UIActionResult{Response: "acknowledged"}, nil
	}
}
