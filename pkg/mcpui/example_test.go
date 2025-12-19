// Copyright 2025 The MCP-UI Go SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package mcpui_test

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ironystock/agentic-obs/pkg/mcpui"
)

// ExampleHTMLContent demonstrates creating HTML content for a UI resource.
func ExampleHTMLContent() {
	content := &mcpui.HTMLContent{
		HTML: "<div>Hello, World!</div>",
	}

	data, _ := content.MarshalJSON()

	// Parse to verify structure
	var m map[string]any
	json.Unmarshal(data, &m)
	fmt.Printf("mimeType: %s\n", m["mimeType"])
	fmt.Printf("has text: %v\n", m["text"] != nil)
	// Output:
	// mimeType: text/html
	// has text: true
}

// ExampleURLContent demonstrates creating URL content for an external resource.
func ExampleURLContent() {
	content := &mcpui.URLContent{
		URL: "https://example.com/dashboard",
	}

	data, _ := content.MarshalJSON()
	fmt.Println(string(data))
	// Output: {"mimeType":"text/uri-list","text":"https://example.com/dashboard"}
}

// ExampleRemoteDOMContent demonstrates creating Remote DOM content with React.
func ExampleRemoteDOMContent() {
	content := &mcpui.RemoteDOMContent{
		Script:    "React.createElement('div', null, 'Hello from React!');",
		Framework: mcpui.FrameworkReact,
	}

	data, _ := content.MarshalJSON()
	fmt.Println(string(data))
	// Output: {"mimeType":"application/vnd.mcp-ui.remote-dom+javascript; framework=react","text":"React.createElement('div', null, 'Hello from React!');"}
}

// ExampleUIResource demonstrates creating a UI resource definition.
func ExampleUIResource() {
	resource := &mcpui.UIResource{
		URI:         "ui://dashboard/main",
		Name:        "main-dashboard",
		Title:       "Main Dashboard",
		Description: "The primary dashboard view for monitoring",
		MIMEType:    mcpui.MIMETypeHTML,
	}

	if err := resource.Validate(); err != nil {
		fmt.Println("Invalid:", err)
		return
	}

	data, _ := json.Marshal(resource)
	fmt.Println(string(data))
	// Output: {"uri":"ui://dashboard/main","name":"main-dashboard","title":"Main Dashboard","description":"The primary dashboard view for monitoring","mimeType":"text/html"}
}

// ExampleNewUIResourceContents demonstrates creating resource contents from content.
func ExampleNewUIResourceContents() {
	content := &mcpui.HTMLContent{
		HTML: "<p>Hello</p>",
	}

	rc, err := mcpui.NewUIResourceContents("ui://greeting/hello", content)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("URI: %s, MIME: %s\n", rc.URI, rc.MIMEType)
	// Output: URI: ui://greeting/hello, MIME: text/html
}

// ExampleUIAction_ParsePayload demonstrates parsing action payloads.
func ExampleUIAction_ParsePayload() {
	action := &mcpui.UIAction{
		Type:      mcpui.ActionTypeTool,
		MessageID: "msg-123",
		Payload:   json.RawMessage(`{"toolName":"get_status","params":{"verbose":true}}`),
	}

	payload, err := action.ParsePayload()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	toolPayload := payload.(*mcpui.ToolActionPayload)
	fmt.Printf("Tool: %s, Verbose: %v\n", toolPayload.ToolName, toolPayload.Params["verbose"])
	// Output: Tool: get_status, Verbose: true
}

// ExampleNewToolAction demonstrates creating a tool action.
func ExampleNewToolAction() {
	action, err := mcpui.NewToolAction("msg-456", "set_volume", map[string]any{
		"source": "Microphone",
		"level":  0.8,
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	data, _ := json.Marshal(action)
	fmt.Println(string(data))
	// Output: {"type":"tool","messageId":"msg-456","payload":{"toolName":"set_volume","params":{"level":0.8,"source":"Microphone"}}}
}

// ExampleNewReceivedResponse demonstrates creating an acknowledgment response.
func ExampleNewReceivedResponse() {
	resp := mcpui.NewReceivedResponse("msg-123")

	data, _ := json.Marshal(resp)
	fmt.Println(string(data))
	// Output: {"type":"ui-message-received","messageId":"msg-123"}
}

// ExampleNewSuccessResponse demonstrates creating a success response.
func ExampleNewSuccessResponse() {
	result := map[string]any{
		"status": "ok",
		"volume": 0.8,
	}
	resp := mcpui.NewSuccessResponse("msg-456", result)

	fmt.Printf("Success: %v, Response: %v\n", resp.IsSuccess(), resp.GetResponse() != nil)
	// Output: Success: true, Response: true
}

// ExampleNewErrorResponse demonstrates creating an error response.
func ExampleNewErrorResponse() {
	resp := mcpui.NewErrorResponse("msg-789", fmt.Errorf("source not found"))

	fmt.Printf("Error: %v, Message: %s\n", resp.IsError(), resp.GetError().Message)
	// Output: Error: true, Message: source not found
}

// ExampleRouter demonstrates using the action router.
func ExampleRouter() {
	router := mcpui.NewRouter()

	// Register handler for tool actions
	router.HandleType(mcpui.ActionTypeTool, mcpui.WrapToolHandler(
		func(ctx context.Context, toolName string, params map[string]any) (any, error) {
			return map[string]string{"executed": toolName}, nil
		},
	))

	// Register handler for a specific resource
	router.HandleResource("ui://dashboard/main", func(ctx context.Context, req *mcpui.UIActionRequest) (*mcpui.UIActionResult, error) {
		return &mcpui.UIActionResult{Response: "dashboard action handled"}, nil
	})

	// Create and dispatch an action
	action, _ := mcpui.NewToolAction("msg-1", "get_status", nil)
	req := &mcpui.UIActionRequest{Action: action}

	result, err := router.Dispatch(context.Background(), req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Result:", result.Response)
	// Output: Result: map[executed:get_status]
}

// ExampleWrapToolHandler demonstrates wrapping a tool handler.
func ExampleWrapToolHandler() {
	handler := mcpui.WrapToolHandler(func(ctx context.Context, toolName string, params map[string]any) (any, error) {
		switch toolName {
		case "greet":
			name := params["name"].(string)
			return fmt.Sprintf("Hello, %s!", name), nil
		default:
			return nil, fmt.Errorf("unknown tool: %s", toolName)
		}
	})

	action, _ := mcpui.NewToolAction("msg-1", "greet", map[string]any{"name": "World"})
	req := &mcpui.UIActionRequest{Action: action}

	result, _ := handler(context.Background(), req)
	fmt.Println(result.Response)
	// Output: Hello, World!
}

// Example demonstrates a complete workflow of creating and handling UI resources.
func Example() {
	// 1. Create a UI resource
	resource := &mcpui.UIResource{
		URI:      "ui://greeting/welcome",
		Name:     "welcome",
		Title:    "Welcome Card",
		MIMEType: mcpui.MIMETypeHTML,
	}

	// 2. Create content for the resource
	content := &mcpui.HTMLContent{
		HTML: "<h1>Welcome!</h1><button onclick=\"sendAction()\">Click Me</button>",
	}

	// 3. Create resource contents
	rc, _ := mcpui.NewUIResourceContents(resource.URI, content)

	// 4. Set up a router to handle actions
	router := mcpui.NewRouter()
	router.HandleType(mcpui.ActionTypeTool, mcpui.WrapToolHandler(
		func(ctx context.Context, toolName string, params map[string]any) (any, error) {
			return "Tool executed: " + toolName, nil
		},
	))

	// 5. Simulate receiving an action from the UI
	action, _ := mcpui.NewToolAction("msg-1", "button_clicked", nil)
	req := &mcpui.UIActionRequest{
		Action:      action,
		ResourceURI: resource.URI,
	}

	// 6. Handle the action and create response
	result, _ := router.Dispatch(context.Background(), req)
	resp := result.ToUIResponse(action.MessageID)

	fmt.Printf("Resource: %s\n", rc.URI)
	fmt.Printf("Response type: %s\n", resp.Type)
	fmt.Printf("Success: %v\n", resp.IsSuccess())
	// Output:
	// Resource: ui://greeting/welcome
	// Response type: ui-message-response
	// Success: true
}
