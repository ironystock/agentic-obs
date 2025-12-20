package mcp

import (
	"context"
	"fmt"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// getSession safely extracts the session from a CallToolRequest.
// Returns nil if the request is nil.
func getSession(request *mcpsdk.CallToolRequest) *mcpsdk.ServerSession {
	if request == nil {
		return nil
	}
	return request.Session
}

// ElicitConfirmation requests user confirmation for a potentially dangerous action.
// Returns true if the user confirmed, false if declined/cancelled.
// If session is nil (e.g., in tests), returns true to allow the action.
func ElicitConfirmation(ctx context.Context, session *mcpsdk.ServerSession, message string) (bool, error) {
	if session == nil {
		// In tests or when session is unavailable, skip confirmation
		return true, nil
	}

	result, err := session.Elicit(ctx, &mcpsdk.ElicitParams{
		Message: message,
		RequestedSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"confirmed": map[string]any{
					"type":        "boolean",
					"description": "Confirm this action",
				},
			},
			"required": []string{"confirmed"},
		},
	})
	if err != nil {
		return false, fmt.Errorf("elicitation failed: %w", err)
	}

	// Check the action - accept, decline, or cancel
	if result.Action != "accept" {
		return false, nil
	}

	// Check if user confirmed
	confirmed, ok := result.Content["confirmed"].(bool)
	return ok && confirmed, nil
}

// ElicitStreamingConfirmation requests confirmation before starting a stream.
func ElicitStreamingConfirmation(ctx context.Context, session *mcpsdk.ServerSession) (bool, error) {
	return ElicitConfirmation(ctx, session,
		"You are about to start a live stream. This will begin broadcasting immediately. Continue?")
}

// ElicitStopStreamingConfirmation requests confirmation before stopping a stream.
func ElicitStopStreamingConfirmation(ctx context.Context, session *mcpsdk.ServerSession) (bool, error) {
	return ElicitConfirmation(ctx, session,
		"You are about to stop the live stream. This will end the broadcast. Continue?")
}

// ElicitDeleteConfirmation requests confirmation before deleting something.
func ElicitDeleteConfirmation(ctx context.Context, session *mcpsdk.ServerSession, itemType, itemName string) (bool, error) {
	return ElicitConfirmation(ctx, session,
		fmt.Sprintf("You are about to permanently delete %s '%s'. This cannot be undone. Continue?", itemType, itemName))
}

// CancelledResult returns a result indicating the action was cancelled by the user.
func CancelledResult(action string) SimpleResult {
	return SimpleResult{Message: fmt.Sprintf("%s cancelled by user", action)}
}
