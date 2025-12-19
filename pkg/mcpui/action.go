// Copyright 2025 The MCP-UI Go SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package mcpui

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// ActionType constants define the types of UI actions.
const (
	// ActionTypeTool triggers a tool call on the MCP server.
	ActionTypeTool = "tool"
	// ActionTypeIntent signals a user intent for the AI to interpret.
	ActionTypeIntent = "intent"
	// ActionTypePrompt sends a prompt message to the AI.
	ActionTypePrompt = "prompt"
	// ActionTypeNotify sends a notification message to the host.
	ActionTypeNotify = "notify"
	// ActionTypeLink requests to open an external link.
	ActionTypeLink = "link"
	// ActionTypeUISize indicates a UI size change.
	ActionTypeUISize = "ui-size-change"
)

// UIAction represents a user interaction from embedded UI.
// Actions are sent from the iframe to the host via postMessage.
type UIAction struct {
	// Type is the action type (tool, intent, prompt, notify, link, ui-size-change).
	Type string `json:"type"`
	// MessageID is an optional identifier for correlating async responses.
	MessageID string `json:"messageId,omitempty"`
	// Payload contains the action-specific data.
	Payload json.RawMessage `json:"payload"`
}

// ParsePayload parses the action payload into the appropriate type.
func (a *UIAction) ParsePayload() (any, error) {
	switch a.Type {
	case ActionTypeTool:
		var p ToolActionPayload
		if err := json.Unmarshal(a.Payload, &p); err != nil {
			return nil, fmt.Errorf("invalid tool payload: %w", err)
		}
		return &p, nil
	case ActionTypeIntent:
		var p IntentActionPayload
		if err := json.Unmarshal(a.Payload, &p); err != nil {
			return nil, fmt.Errorf("invalid intent payload: %w", err)
		}
		return &p, nil
	case ActionTypePrompt:
		var p PromptActionPayload
		if err := json.Unmarshal(a.Payload, &p); err != nil {
			return nil, fmt.Errorf("invalid prompt payload: %w", err)
		}
		return &p, nil
	case ActionTypeNotify:
		var p NotifyActionPayload
		if err := json.Unmarshal(a.Payload, &p); err != nil {
			return nil, fmt.Errorf("invalid notify payload: %w", err)
		}
		return &p, nil
	case ActionTypeLink:
		var p LinkActionPayload
		if err := json.Unmarshal(a.Payload, &p); err != nil {
			return nil, fmt.Errorf("invalid link payload: %w", err)
		}
		return &p, nil
	case ActionTypeUISize:
		var p UISizeActionPayload
		if err := json.Unmarshal(a.Payload, &p); err != nil {
			return nil, fmt.Errorf("invalid ui-size-change payload: %w", err)
		}
		return &p, nil
	default:
		return nil, fmt.Errorf("unknown action type: %s", a.Type)
	}
}

// ToolPayload returns the payload as a ToolActionPayload if the action type is "tool".
func (a *UIAction) ToolPayload() (*ToolActionPayload, error) {
	if a.Type != ActionTypeTool {
		return nil, fmt.Errorf("action type is %s, not tool", a.Type)
	}
	var p ToolActionPayload
	if err := json.Unmarshal(a.Payload, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// IntentPayload returns the payload as an IntentActionPayload if the action type is "intent".
func (a *UIAction) IntentPayload() (*IntentActionPayload, error) {
	if a.Type != ActionTypeIntent {
		return nil, fmt.Errorf("action type is %s, not intent", a.Type)
	}
	var p IntentActionPayload
	if err := json.Unmarshal(a.Payload, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// PromptPayload returns the payload as a PromptActionPayload if the action type is "prompt".
func (a *UIAction) PromptPayload() (*PromptActionPayload, error) {
	if a.Type != ActionTypePrompt {
		return nil, fmt.Errorf("action type is %s, not prompt", a.Type)
	}
	var p PromptActionPayload
	if err := json.Unmarshal(a.Payload, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// NotifyPayload returns the payload as a NotifyActionPayload if the action type is "notify".
func (a *UIAction) NotifyPayload() (*NotifyActionPayload, error) {
	if a.Type != ActionTypeNotify {
		return nil, fmt.Errorf("action type is %s, not notify", a.Type)
	}
	var p NotifyActionPayload
	if err := json.Unmarshal(a.Payload, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// LinkPayload returns the payload as a LinkActionPayload if the action type is "link".
func (a *UIAction) LinkPayload() (*LinkActionPayload, error) {
	if a.Type != ActionTypeLink {
		return nil, fmt.Errorf("action type is %s, not link", a.Type)
	}
	var p LinkActionPayload
	if err := json.Unmarshal(a.Payload, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// UISizePayload returns the payload as a UISizeActionPayload if the action type is "ui-size-change".
func (a *UIAction) UISizePayload() (*UISizeActionPayload, error) {
	if a.Type != ActionTypeUISize {
		return nil, fmt.Errorf("action type is %s, not ui-size-change", a.Type)
	}
	var p UISizeActionPayload
	if err := json.Unmarshal(a.Payload, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// ToolActionPayload is the payload for tool actions.
// It requests the host to execute an MCP tool.
type ToolActionPayload struct {
	// ToolName is the name of the tool to call.
	ToolName string `json:"toolName"`
	// Params are the tool parameters.
	Params map[string]any `json:"params,omitempty"`
}

// IntentActionPayload is the payload for intent actions.
// It signals a user intent for the AI to interpret and act upon.
type IntentActionPayload struct {
	// Intent is the user intent identifier.
	Intent string `json:"intent"`
	// Params are optional parameters for the intent.
	Params map[string]any `json:"params,omitempty"`
}

// PromptActionPayload is the payload for prompt actions.
// It sends a prompt message to the AI conversation.
type PromptActionPayload struct {
	// Prompt is the text to send to the AI.
	Prompt string `json:"prompt"`
}

// NotifyActionPayload is the payload for notify actions.
// It sends a notification message to the host.
type NotifyActionPayload struct {
	// Message is the notification text.
	Message string `json:"message"`
	// Level is an optional severity level (info, warning, error).
	Level string `json:"level,omitempty"`
}

// LinkActionPayload is the payload for link actions.
// It requests to open an external URL.
type LinkActionPayload struct {
	// URL is the external URL to open.
	URL string `json:"url"`
}

// Validate checks that the LinkActionPayload has a valid URL.
func (p *LinkActionPayload) Validate() error {
	if p.URL == "" {
		return fmt.Errorf("link payload URL is required")
	}
	parsed, err := url.Parse(p.URL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("URL must have http or https scheme, got: %s", parsed.Scheme)
	}
	if parsed.Host == "" {
		return fmt.Errorf("URL must have a host")
	}
	return nil
}

// UISizeActionPayload is the payload for ui-size-change actions.
// It notifies the host of UI dimension changes.
type UISizeActionPayload struct {
	// Height is the new height in pixels.
	Height int `json:"height"`
	// Width is the optional new width in pixels.
	Width int `json:"width,omitempty"`
}

// NewToolAction creates a new tool action.
func NewToolAction(messageID, toolName string, params map[string]any) (*UIAction, error) {
	payload := ToolActionPayload{
		ToolName: toolName,
		Params:   params,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &UIAction{
		Type:      ActionTypeTool,
		MessageID: messageID,
		Payload:   data,
	}, nil
}

// NewIntentAction creates a new intent action.
func NewIntentAction(messageID, intent string, params map[string]any) (*UIAction, error) {
	payload := IntentActionPayload{
		Intent: intent,
		Params: params,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &UIAction{
		Type:      ActionTypeIntent,
		MessageID: messageID,
		Payload:   data,
	}, nil
}

// NewPromptAction creates a new prompt action.
func NewPromptAction(messageID, prompt string) (*UIAction, error) {
	payload := PromptActionPayload{
		Prompt: prompt,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &UIAction{
		Type:      ActionTypePrompt,
		MessageID: messageID,
		Payload:   data,
	}, nil
}

// NewNotifyAction creates a new notify action.
func NewNotifyAction(message string, level string) (*UIAction, error) {
	payload := NotifyActionPayload{
		Message: message,
		Level:   level,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &UIAction{
		Type:    ActionTypeNotify,
		Payload: data,
	}, nil
}

// NewLinkAction creates a new link action.
// The URL is validated to ensure it is a valid absolute URL with http or https scheme.
func NewLinkAction(rawURL string) (*UIAction, error) {
	// Validate URL
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("URL must have http or https scheme, got: %s", parsed.Scheme)
	}
	if parsed.Host == "" {
		return nil, fmt.Errorf("URL must have a host")
	}

	payload := LinkActionPayload{
		URL: rawURL,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &UIAction{
		Type:    ActionTypeLink,
		Payload: data,
	}, nil
}

// NewUISizeAction creates a new UI size change action.
func NewUISizeAction(height, width int) (*UIAction, error) {
	payload := UISizeActionPayload{
		Height: height,
		Width:  width,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &UIAction{
		Type:    ActionTypeUISize,
		Payload: data,
	}, nil
}
