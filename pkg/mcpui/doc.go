// Copyright 2025 The MCP-UI Go SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package mcpui provides a Go SDK for the MCP-UI protocol (SEP-1865).
//
// MCP-UI enables MCP servers to return interactive UI resources that clients
// can render in sandboxed iframes. This package implements the server-side
// portion of the protocol, allowing Go-based MCP servers to generate UI
// resources.
//
// # Overview
//
// The MCP-UI protocol defines three types of UI content:
//
//   - [HTMLContent]: Inline HTML rendered via iframe srcdoc
//   - [URLContent]: External URL rendered via iframe src
//   - [RemoteDOMContent]: Script-based UI using remote DOM rendering
//
// Each content type implements the [UIContent] interface and can be serialized
// to JSON for transmission to clients.
//
// # Creating UI Resources
//
// To create a simple HTML UI resource:
//
//	resource := &mcpui.UIResource{
//		URI:  "ui://hello/greeting",
//		Name: "Greeting Card",
//		Content: &mcpui.HTMLContent{
//			HTML: `<div style="padding: 20px;">Hello, World!</div>`,
//		},
//	}
//
// To create an external URL resource:
//
//	resource := &mcpui.UIResource{
//		URI:  "ui://dashboard/main",
//		Name: "Dashboard",
//		Content: &mcpui.URLContent{
//			URL: "https://example.com/dashboard",
//		},
//	}
//
// # Handling UI Actions
//
// When users interact with UI resources, the client sends [UIAction] messages.
// Use [UIActionHandler] to process these actions:
//
//	handler := func(ctx context.Context, req *mcpui.UIActionRequest) (*mcpui.UIActionResult, error) {
//		switch req.Action.Type {
//		case mcpui.ActionTypeTool:
//			var payload mcpui.ToolActionPayload
//			if err := json.Unmarshal(req.Action.Payload, &payload); err != nil {
//				return nil, err
//			}
//			// Handle tool call...
//		}
//		return &mcpui.UIActionResult{Response: "OK"}, nil
//	}
//
// # Response Messages
//
// Send responses back to the UI using the response builders:
//
//	// Acknowledge receipt
//	ack := mcpui.NewReceivedResponse(action.MessageID)
//
//	// Send success response
//	resp := mcpui.NewSuccessResponse(action.MessageID, result)
//
//	// Send error response
//	errResp := mcpui.NewErrorResponse(action.MessageID, err)
//
// # Protocol Details
//
// The MCP-UI protocol is documented at https://mcpui.dev/. This SDK implements
// the server-side portion, enabling Go-based MCP servers to generate UI
// resources that compliant clients can render.
//
// MIME types used by the protocol:
//
//   - text/html: Inline HTML content
//   - text/uri-list: External URL content
//   - application/vnd.mcp-ui.remote-dom: Remote DOM script content
//
// # Integration with MCP
//
// This package is designed to work alongside the official MCP Go SDK
// (github.com/modelcontextprotocol/go-sdk). UI resources can be returned
// from MCP tool handlers as additional content alongside regular tool results.
package mcpui
