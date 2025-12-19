// Copyright 2025 The MCP-UI Go SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package mcpui

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUIActionResult_ToUIResponse(t *testing.T) {
	t.Run("success result", func(t *testing.T) {
		result := &UIActionResult{Response: "success data"}
		resp := result.ToUIResponse("msg-123")

		assert.Equal(t, ResponseTypeResponse, resp.Type)
		assert.Equal(t, "msg-123", resp.MessageID)
		assert.True(t, resp.IsSuccess())
		assert.Equal(t, "success data", resp.GetResponse())
	})

	t.Run("error result", func(t *testing.T) {
		result := &UIActionResult{Error: errors.New("something failed")}
		resp := result.ToUIResponse("msg-456")

		assert.Equal(t, ResponseTypeResponse, resp.Type)
		assert.Equal(t, "msg-456", resp.MessageID)
		assert.True(t, resp.IsError())
		assert.Equal(t, "something failed", resp.GetError().Message)
	})
}

func TestRouter_HandleType(t *testing.T) {
	router := NewRouter()

	var called bool
	router.HandleType(ActionTypeTool, func(ctx context.Context, req *UIActionRequest) (*UIActionResult, error) {
		called = true
		return &UIActionResult{Response: "tool handled"}, nil
	})

	action, _ := NewToolAction("msg-1", "test_tool", nil)
	req := &UIActionRequest{Action: action}

	result, err := router.Dispatch(context.Background(), req)
	require.NoError(t, err)
	assert.True(t, called)
	assert.Equal(t, "tool handled", result.Response)
}

func TestRouter_HandleResource(t *testing.T) {
	router := NewRouter()

	var called bool
	router.HandleResource("ui://dashboard/main", func(ctx context.Context, req *UIActionRequest) (*UIActionResult, error) {
		called = true
		return &UIActionResult{Response: "resource handled"}, nil
	})

	action, _ := NewToolAction("msg-1", "test_tool", nil)
	req := &UIActionRequest{
		Action:      action,
		ResourceURI: "ui://dashboard/main",
	}

	result, err := router.Dispatch(context.Background(), req)
	require.NoError(t, err)
	assert.True(t, called)
	assert.Equal(t, "resource handled", result.Response)
}

func TestRouter_Priority(t *testing.T) {
	router := NewRouter()

	// Register handlers in reverse priority order
	router.SetDefault(func(ctx context.Context, req *UIActionRequest) (*UIActionResult, error) {
		return &UIActionResult{Response: "default"}, nil
	})
	router.HandleType(ActionTypeTool, func(ctx context.Context, req *UIActionRequest) (*UIActionResult, error) {
		return &UIActionResult{Response: "type"}, nil
	})
	router.HandleResource("ui://test", func(ctx context.Context, req *UIActionRequest) (*UIActionResult, error) {
		return &UIActionResult{Response: "resource"}, nil
	})

	action, _ := NewToolAction("msg-1", "test", nil)

	t.Run("resource takes priority", func(t *testing.T) {
		req := &UIActionRequest{
			Action:      action,
			ResourceURI: "ui://test",
		}
		result, err := router.Dispatch(context.Background(), req)
		require.NoError(t, err)
		assert.Equal(t, "resource", result.Response)
	})

	t.Run("type handler when no resource match", func(t *testing.T) {
		req := &UIActionRequest{
			Action:      action,
			ResourceURI: "ui://other",
		}
		result, err := router.Dispatch(context.Background(), req)
		require.NoError(t, err)
		assert.Equal(t, "type", result.Response)
	})

	t.Run("default when no match", func(t *testing.T) {
		promptAction, _ := NewPromptAction("msg-1", "test")
		req := &UIActionRequest{
			Action:      promptAction,
			ResourceURI: "ui://other",
		}
		result, err := router.Dispatch(context.Background(), req)
		require.NoError(t, err)
		assert.Equal(t, "default", result.Response)
	})
}

func TestRouter_NoHandler(t *testing.T) {
	router := NewRouter()

	action, _ := NewToolAction("msg-1", "test", nil)
	req := &UIActionRequest{Action: action}

	_, err := router.Dispatch(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no handler")
}

func TestRouter_Handle(t *testing.T) {
	router := NewRouter()
	router.HandleType(ActionTypeTool, func(ctx context.Context, req *UIActionRequest) (*UIActionResult, error) {
		return &UIActionResult{Response: "handled"}, nil
	})

	action, _ := NewToolAction("msg-1", "test", nil)
	req := &UIActionRequest{Action: action}

	// Test that Handle delegates to Dispatch
	result, err := router.Handle(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "handled", result.Response)
}

func TestWrapToolHandler(t *testing.T) {
	handler := WrapToolHandler(func(ctx context.Context, toolName string, params map[string]any) (any, error) {
		return map[string]string{"tool": toolName}, nil
	})

	action := &UIAction{
		Type:    ActionTypeTool,
		Payload: json.RawMessage(`{"toolName":"my_tool","params":{}}`),
	}
	req := &UIActionRequest{Action: action}

	result, err := handler(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"tool": "my_tool"}, result.Response)
}

func TestWrapToolHandler_WrongType(t *testing.T) {
	handler := WrapToolHandler(func(ctx context.Context, toolName string, params map[string]any) (any, error) {
		return nil, nil
	})

	action, _ := NewPromptAction("msg-1", "test")
	req := &UIActionRequest{Action: action}

	_, err := handler(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected tool action")
}

func TestWrapToolHandler_Error(t *testing.T) {
	handler := WrapToolHandler(func(ctx context.Context, toolName string, params map[string]any) (any, error) {
		return nil, errors.New("tool failed")
	})

	action, _ := NewToolAction("msg-1", "test", nil)
	req := &UIActionRequest{Action: action}

	result, err := handler(context.Background(), req)
	require.NoError(t, err) // Handler error goes in result, not return error
	assert.NotNil(t, result.Error)
	assert.Equal(t, "tool failed", result.Error.Error())
}

func TestWrapIntentHandler(t *testing.T) {
	handler := WrapIntentHandler(func(ctx context.Context, intent string, params map[string]any) (any, error) {
		return map[string]string{"intent": intent}, nil
	})

	action, _ := NewIntentAction("msg-1", "switch_scene", nil)
	req := &UIActionRequest{Action: action}

	result, err := handler(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"intent": "switch_scene"}, result.Response)
}

func TestWrapPromptHandler(t *testing.T) {
	handler := WrapPromptHandler(func(ctx context.Context, prompt string) (any, error) {
		return "Got: " + prompt, nil
	})

	action, _ := NewPromptAction("msg-1", "Hello")
	req := &UIActionRequest{Action: action}

	result, err := handler(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "Got: Hello", result.Response)
}

func TestWrapNotifyHandler(t *testing.T) {
	var receivedMessage, receivedLevel string
	handler := WrapNotifyHandler(func(ctx context.Context, message string, level string) error {
		receivedMessage = message
		receivedLevel = level
		return nil
	})

	action, _ := NewNotifyAction("Test notification", "info")
	req := &UIActionRequest{Action: action}

	result, err := handler(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "acknowledged", result.Response)
	assert.Equal(t, "Test notification", receivedMessage)
	assert.Equal(t, "info", receivedLevel)
}

func TestWrapLinkHandler(t *testing.T) {
	var receivedURL string
	handler := WrapLinkHandler(func(ctx context.Context, url string) error {
		receivedURL = url
		return nil
	})

	action, _ := NewLinkAction("https://example.com")
	req := &UIActionRequest{Action: action}

	result, err := handler(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "opened", result.Response)
	assert.Equal(t, "https://example.com", receivedURL)
}

func TestWrapUISizeHandler(t *testing.T) {
	var receivedHeight, receivedWidth int
	handler := WrapUISizeHandler(func(ctx context.Context, height, width int) error {
		receivedHeight = height
		receivedWidth = width
		return nil
	})

	action, _ := NewUISizeAction(600, 800)
	req := &UIActionRequest{Action: action}

	result, err := handler(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "acknowledged", result.Response)
	assert.Equal(t, 600, receivedHeight)
	assert.Equal(t, 800, receivedWidth)
}

func TestWrappedHandlers_WrongType(t *testing.T) {
	toolAction, _ := NewToolAction("msg", "test", nil)

	tests := []struct {
		name    string
		handler UIActionHandler
		action  *UIAction
	}{
		{"IntentHandler", WrapIntentHandler(func(ctx context.Context, intent string, params map[string]any) (any, error) { return nil, nil }), toolAction},
		{"PromptHandler", WrapPromptHandler(func(ctx context.Context, prompt string) (any, error) { return nil, nil }), toolAction},
		{"NotifyHandler", WrapNotifyHandler(func(ctx context.Context, message string, level string) error { return nil }), toolAction},
		{"LinkHandler", WrapLinkHandler(func(ctx context.Context, url string) error { return nil }), toolAction},
		{"UISizeHandler", WrapUISizeHandler(func(ctx context.Context, height, width int) error { return nil }), toolAction},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &UIActionRequest{Action: tt.action}
			_, err := tt.handler(context.Background(), req)
			assert.Error(t, err)
		})
	}
}
