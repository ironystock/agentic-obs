package mcp

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// HelpInput is the input for the help tool
type HelpInput struct {
	Topic   string `json:"topic,omitempty" jsonschema:"description=Topic to get help on: 'overview', 'tools', 'resources', 'prompts', 'workflows', 'troubleshooting', or a specific tool name"`
	Verbose bool   `json:"verbose,omitempty" jsonschema:"description=Include examples and detailed explanations"`
}

// handleHelp provides comprehensive help on agentic-obs features
func (s *Server) handleHelp(ctx context.Context, request *mcpsdk.CallToolRequest, input HelpInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	topic := strings.ToLower(strings.TrimSpace(input.Topic))

	// Default to overview if no topic specified
	if topic == "" {
		topic = "overview"
	}

	log.Printf("Help requested for topic: %s (verbose: %v)", topic, input.Verbose)

	var helpText string
	var err error

	// Route to appropriate help handler (using extracted help content from help_content.go)
	switch topic {
	case "overview":
		helpText = GetOverviewHelp(input.Verbose)
	case "tools":
		helpText = GetToolsHelp(input.Verbose)
	case "resources":
		helpText = GetResourcesHelp(input.Verbose)
	case "prompts":
		helpText = GetPromptsHelp(input.Verbose)
	case "workflows":
		helpText = GetWorkflowsHelp(input.Verbose)
	case "troubleshooting":
		helpText = GetTroubleshootingHelp(input.Verbose)
	default:
		// Try to find help for a specific tool (using extracted content from help_tools.go)
		helpText, err = getToolHelp(topic, input.Verbose)
		if err != nil {
			s.recordAction("help", "Get help", input, nil, false, time.Since(start))
			return nil, nil, fmt.Errorf("unknown help topic '%s'. Try 'overview', 'tools', 'resources', 'prompts', 'workflows', 'troubleshooting', or a specific tool name", topic)
		}
	}

	result := map[string]interface{}{
		"topic":   topic,
		"help":    helpText,
		"verbose": input.Verbose,
	}

	s.recordAction("help", "Get help", input, result, true, time.Since(start))
	return nil, result, nil
}

// getToolHelp returns detailed help for a specific tool using extracted content from help_tools.go
func getToolHelp(toolName string, verbose bool) (string, error) {
	help, exists := GetToolHelpContent(toolName)
	if !exists {
		return "", fmt.Errorf("tool '%s' not found", toolName)
	}

	if verbose {
		help += VerboseToolSuffix
	}

	return help, nil
}
