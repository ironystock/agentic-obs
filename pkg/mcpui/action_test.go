// Copyright 2025 The MCP-UI Go SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package mcpui

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUIAction_ParsePayload(t *testing.T) {
	tests := []struct {
		name    string
		action  *UIAction
		wantErr bool
		check   func(t *testing.T, payload any)
	}{
		{
			name: "tool action",
			action: &UIAction{
				Type:    ActionTypeTool,
				Payload: json.RawMessage(`{"toolName":"get_status","params":{"verbose":true}}`),
			},
			check: func(t *testing.T, payload any) {
				p, ok := payload.(*ToolActionPayload)
				require.True(t, ok)
				assert.Equal(t, "get_status", p.ToolName)
				assert.Equal(t, true, p.Params["verbose"])
			},
		},
		{
			name: "intent action",
			action: &UIAction{
				Type:    ActionTypeIntent,
				Payload: json.RawMessage(`{"intent":"switch_scene","params":{"scene":"Gaming"}}`),
			},
			check: func(t *testing.T, payload any) {
				p, ok := payload.(*IntentActionPayload)
				require.True(t, ok)
				assert.Equal(t, "switch_scene", p.Intent)
				assert.Equal(t, "Gaming", p.Params["scene"])
			},
		},
		{
			name: "prompt action",
			action: &UIAction{
				Type:    ActionTypePrompt,
				Payload: json.RawMessage(`{"prompt":"How do I start streaming?"}`),
			},
			check: func(t *testing.T, payload any) {
				p, ok := payload.(*PromptActionPayload)
				require.True(t, ok)
				assert.Equal(t, "How do I start streaming?", p.Prompt)
			},
		},
		{
			name: "notify action",
			action: &UIAction{
				Type:    ActionTypeNotify,
				Payload: json.RawMessage(`{"message":"Recording started","level":"info"}`),
			},
			check: func(t *testing.T, payload any) {
				p, ok := payload.(*NotifyActionPayload)
				require.True(t, ok)
				assert.Equal(t, "Recording started", p.Message)
				assert.Equal(t, "info", p.Level)
			},
		},
		{
			name: "link action",
			action: &UIAction{
				Type:    ActionTypeLink,
				Payload: json.RawMessage(`{"url":"https://example.com/docs"}`),
			},
			check: func(t *testing.T, payload any) {
				p, ok := payload.(*LinkActionPayload)
				require.True(t, ok)
				assert.Equal(t, "https://example.com/docs", p.URL)
			},
		},
		{
			name: "ui-size-change action",
			action: &UIAction{
				Type:    ActionTypeUISize,
				Payload: json.RawMessage(`{"height":600,"width":800}`),
			},
			check: func(t *testing.T, payload any) {
				p, ok := payload.(*UISizeActionPayload)
				require.True(t, ok)
				assert.Equal(t, 600, p.Height)
				assert.Equal(t, 800, p.Width)
			},
		},
		{
			name: "unknown action type",
			action: &UIAction{
				Type:    "unknown",
				Payload: json.RawMessage(`{}`),
			},
			wantErr: true,
		},
		{
			name: "invalid payload",
			action: &UIAction{
				Type:    ActionTypeTool,
				Payload: json.RawMessage(`not valid json`),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := tt.action.ParsePayload()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			tt.check(t, payload)
		})
	}
}

func TestUIAction_TypedPayloadAccessors(t *testing.T) {
	t.Run("ToolPayload", func(t *testing.T) {
		action := &UIAction{
			Type:    ActionTypeTool,
			Payload: json.RawMessage(`{"toolName":"test_tool"}`),
		}
		p, err := action.ToolPayload()
		require.NoError(t, err)
		assert.Equal(t, "test_tool", p.ToolName)

		// Wrong type
		action.Type = ActionTypePrompt
		_, err = action.ToolPayload()
		assert.Error(t, err)
	})

	t.Run("IntentPayload", func(t *testing.T) {
		action := &UIAction{
			Type:    ActionTypeIntent,
			Payload: json.RawMessage(`{"intent":"test_intent"}`),
		}
		p, err := action.IntentPayload()
		require.NoError(t, err)
		assert.Equal(t, "test_intent", p.Intent)
	})

	t.Run("PromptPayload", func(t *testing.T) {
		action := &UIAction{
			Type:    ActionTypePrompt,
			Payload: json.RawMessage(`{"prompt":"test prompt"}`),
		}
		p, err := action.PromptPayload()
		require.NoError(t, err)
		assert.Equal(t, "test prompt", p.Prompt)
	})

	t.Run("NotifyPayload", func(t *testing.T) {
		action := &UIAction{
			Type:    ActionTypeNotify,
			Payload: json.RawMessage(`{"message":"test"}`),
		}
		p, err := action.NotifyPayload()
		require.NoError(t, err)
		assert.Equal(t, "test", p.Message)
	})

	t.Run("LinkPayload", func(t *testing.T) {
		action := &UIAction{
			Type:    ActionTypeLink,
			Payload: json.RawMessage(`{"url":"https://test.com"}`),
		}
		p, err := action.LinkPayload()
		require.NoError(t, err)
		assert.Equal(t, "https://test.com", p.URL)
	})

	t.Run("UISizePayload", func(t *testing.T) {
		action := &UIAction{
			Type:    ActionTypeUISize,
			Payload: json.RawMessage(`{"height":100}`),
		}
		p, err := action.UISizePayload()
		require.NoError(t, err)
		assert.Equal(t, 100, p.Height)
	})
}

func TestNewToolAction(t *testing.T) {
	action, err := NewToolAction("msg-123", "get_status", map[string]any{
		"verbose": true,
	})
	require.NoError(t, err)

	assert.Equal(t, ActionTypeTool, action.Type)
	assert.Equal(t, "msg-123", action.MessageID)

	p, err := action.ToolPayload()
	require.NoError(t, err)
	assert.Equal(t, "get_status", p.ToolName)
	assert.Equal(t, true, p.Params["verbose"])
}

func TestNewIntentAction(t *testing.T) {
	action, err := NewIntentAction("msg-456", "change_scene", map[string]any{
		"scene": "Gaming",
	})
	require.NoError(t, err)

	assert.Equal(t, ActionTypeIntent, action.Type)
	assert.Equal(t, "msg-456", action.MessageID)

	p, err := action.IntentPayload()
	require.NoError(t, err)
	assert.Equal(t, "change_scene", p.Intent)
	assert.Equal(t, "Gaming", p.Params["scene"])
}

func TestNewPromptAction(t *testing.T) {
	action, err := NewPromptAction("msg-789", "How do I configure audio?")
	require.NoError(t, err)

	assert.Equal(t, ActionTypePrompt, action.Type)
	assert.Equal(t, "msg-789", action.MessageID)

	p, err := action.PromptPayload()
	require.NoError(t, err)
	assert.Equal(t, "How do I configure audio?", p.Prompt)
}

func TestNewNotifyAction(t *testing.T) {
	action, err := NewNotifyAction("Stream started!", "info")
	require.NoError(t, err)

	assert.Equal(t, ActionTypeNotify, action.Type)
	assert.Empty(t, action.MessageID) // notify doesn't need messageId

	p, err := action.NotifyPayload()
	require.NoError(t, err)
	assert.Equal(t, "Stream started!", p.Message)
	assert.Equal(t, "info", p.Level)
}

func TestNewLinkAction(t *testing.T) {
	action, err := NewLinkAction("https://docs.example.com")
	require.NoError(t, err)

	assert.Equal(t, ActionTypeLink, action.Type)

	p, err := action.LinkPayload()
	require.NoError(t, err)
	assert.Equal(t, "https://docs.example.com", p.URL)
}

func TestNewUISizeAction(t *testing.T) {
	action, err := NewUISizeAction(600, 800)
	require.NoError(t, err)

	assert.Equal(t, ActionTypeUISize, action.Type)

	p, err := action.UISizePayload()
	require.NoError(t, err)
	assert.Equal(t, 600, p.Height)
	assert.Equal(t, 800, p.Width)
}

func TestUIAction_JSONRoundTrip(t *testing.T) {
	original := &UIAction{
		Type:      ActionTypeTool,
		MessageID: "test-msg-id",
		Payload:   json.RawMessage(`{"toolName":"my_tool","params":{"key":"value"}}`),
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded UIAction
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, original.Type, decoded.Type)
	assert.Equal(t, original.MessageID, decoded.MessageID)

	// Verify payload can be parsed
	p, err := decoded.ToolPayload()
	require.NoError(t, err)
	assert.Equal(t, "my_tool", p.ToolName)
}

func TestActionTypeConstants(t *testing.T) {
	// Verify constants match protocol specification
	assert.Equal(t, "tool", ActionTypeTool)
	assert.Equal(t, "intent", ActionTypeIntent)
	assert.Equal(t, "prompt", ActionTypePrompt)
	assert.Equal(t, "notify", ActionTypeNotify)
	assert.Equal(t, "link", ActionTypeLink)
	assert.Equal(t, "ui-size-change", ActionTypeUISize)
}
