package mcp

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHandleHelp tests the main help handler
func TestHandleHelp(t *testing.T) {
	t.Run("returns overview help by default", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)
		defer db.Close()

		input := HelpInput{}
		_, result, err := server.handleHelp(context.Background(), nil, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "overview", resultMap["topic"])
		assert.Equal(t, false, resultMap["verbose"])

		helpText, ok := resultMap["help"].(string)
		require.True(t, ok)
		assert.Contains(t, helpText, "agentic-obs")
		assert.Contains(t, helpText, "OBS Studio")
	})

	t.Run("returns overview help explicitly", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)
		defer db.Close()

		input := HelpInput{Topic: "overview"}
		_, result, err := server.handleHelp(context.Background(), nil, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "overview", resultMap["topic"])

		helpText, ok := resultMap["help"].(string)
		require.True(t, ok)
		assert.Contains(t, helpText, "What is agentic-obs")
		assert.Contains(t, helpText, "Quick Start")
		assert.Contains(t, helpText, "Key Features")
	})

	t.Run("returns verbose overview help", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)
		defer db.Close()

		input := HelpInput{Topic: "overview", Verbose: true}
		_, result, err := server.handleHelp(context.Background(), nil, input)

		assert.NoError(t, err)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, true, resultMap["verbose"])

		helpText, ok := resultMap["help"].(string)
		require.True(t, ok)
		// Verbose should include additional sections
		assert.Contains(t, helpText, "Categories")
		assert.Contains(t, helpText, "Common Workflows")
		assert.Contains(t, helpText, "Next Steps")
	})

	t.Run("returns tools help", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)
		defer db.Close()

		input := HelpInput{Topic: "tools"}
		_, result, err := server.handleHelp(context.Background(), nil, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "tools", resultMap["topic"])

		helpText, ok := resultMap["help"].(string)
		require.True(t, ok)
		assert.Contains(t, helpText, "All Available Tools")
		assert.Contains(t, helpText, "Core Tools")
		assert.Contains(t, helpText, "Sources Tools")
		assert.Contains(t, helpText, "Audio Tools")
		assert.Contains(t, helpText, "Layout Tools")
		assert.Contains(t, helpText, "Visual Tools")
		assert.Contains(t, helpText, "Design Tools")
		assert.Contains(t, helpText, "list_scenes")
		assert.Contains(t, helpText, "start_recording")
		assert.Contains(t, helpText, "create_screenshot_source")
	})

	t.Run("returns verbose tools help", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)
		defer db.Close()

		input := HelpInput{Topic: "tools", Verbose: true}
		_, result, err := server.handleHelp(context.Background(), nil, input)

		assert.NoError(t, err)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)

		helpText, ok := resultMap["help"].(string)
		require.True(t, ok)
		// Verbose should include additional guidance
		assert.Contains(t, helpText, "Getting Help on Specific Tools")
	})

	t.Run("returns resources help", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)
		defer db.Close()

		input := HelpInput{Topic: "resources"}
		_, result, err := server.handleHelp(context.Background(), nil, input)

		assert.NoError(t, err)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "resources", resultMap["topic"])

		helpText, ok := resultMap["help"].(string)
		require.True(t, ok)
		assert.Contains(t, helpText, "MCP Resources")
		assert.Contains(t, helpText, "OBS Scenes")
		assert.Contains(t, helpText, "Screenshot Images")
		assert.Contains(t, helpText, "Scene Presets")
		assert.Contains(t, helpText, "obs://scene/")
		assert.Contains(t, helpText, "obs://screenshot/")
		assert.Contains(t, helpText, "obs://preset/")
	})

	t.Run("returns verbose resources help", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)
		defer db.Close()

		input := HelpInput{Topic: "resources", Verbose: true}
		_, result, err := server.handleHelp(context.Background(), nil, input)

		assert.NoError(t, err)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)

		helpText, ok := resultMap["help"].(string)
		require.True(t, ok)
		// Verbose should include additional sections
		assert.Contains(t, helpText, "Resource Operations")
		assert.Contains(t, helpText, "Resource Notifications")
		assert.Contains(t, helpText, "Using Resources")
		assert.Contains(t, helpText, "Example Use Cases")
	})

	t.Run("returns prompts help", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)
		defer db.Close()

		input := HelpInput{Topic: "prompts"}
		_, result, err := server.handleHelp(context.Background(), nil, input)

		assert.NoError(t, err)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "prompts", resultMap["topic"])

		helpText, ok := resultMap["help"].(string)
		require.True(t, ok)
		assert.Contains(t, helpText, "MCP Prompts")
		assert.Contains(t, helpText, "stream-launch")
		assert.Contains(t, helpText, "stream-teardown")
		assert.Contains(t, helpText, "health-check")
		assert.Contains(t, helpText, "audio-check")
		assert.Contains(t, helpText, "visual-check")
		assert.Contains(t, helpText, "problem-detection")
		assert.Contains(t, helpText, "preset-switcher")
		assert.Contains(t, helpText, "scene-organizer")
		assert.Contains(t, helpText, "recording-workflow")
		assert.Contains(t, helpText, "scene-designer")
		assert.Contains(t, helpText, "source-management")
		assert.Contains(t, helpText, "visual-setup")
	})

	t.Run("returns verbose prompts help", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)
		defer db.Close()

		input := HelpInput{Topic: "prompts", Verbose: true}
		_, result, err := server.handleHelp(context.Background(), nil, input)

		assert.NoError(t, err)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)

		helpText, ok := resultMap["help"].(string)
		require.True(t, ok)
		// Verbose should include additional sections
		assert.Contains(t, helpText, "Using Prompts")
		assert.Contains(t, helpText, "Prompt Arguments")
		assert.Contains(t, helpText, "Creating Custom Workflows")
		assert.Contains(t, helpText, "Workflow Examples")
	})

	t.Run("returns workflows help", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)
		defer db.Close()

		input := HelpInput{Topic: "workflows"}
		_, result, err := server.handleHelp(context.Background(), nil, input)

		assert.NoError(t, err)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "workflows", resultMap["topic"])

		helpText, ok := resultMap["help"].(string)
		require.True(t, ok)
		assert.Contains(t, helpText, "Common Workflows")
		assert.Contains(t, helpText, "Start Streaming")
		assert.Contains(t, helpText, "Visual Stream Monitoring")
		assert.Contains(t, helpText, "Scene Design from Scratch")
		assert.Contains(t, helpText, "Scene Preset Management")
		assert.Contains(t, helpText, "Audio Configuration")
		assert.Contains(t, helpText, "Multi-Scene Setup")
	})

	t.Run("returns verbose workflows help", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)
		defer db.Close()

		input := HelpInput{Topic: "workflows", Verbose: true}
		_, result, err := server.handleHelp(context.Background(), nil, input)

		assert.NoError(t, err)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)

		helpText, ok := resultMap["help"].(string)
		require.True(t, ok)
		// Verbose should include additional sections
		assert.Contains(t, helpText, "Automated Stream Production")
		assert.Contains(t, helpText, "Workflow Best Practices")
	})

	t.Run("returns troubleshooting help", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)
		defer db.Close()

		input := HelpInput{Topic: "troubleshooting"}
		_, result, err := server.handleHelp(context.Background(), nil, input)

		assert.NoError(t, err)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "troubleshooting", resultMap["topic"])

		helpText, ok := resultMap["help"].(string)
		require.True(t, ok)
		assert.Contains(t, helpText, "Troubleshooting Guide")
		assert.Contains(t, helpText, "Connection Issues")
		assert.Contains(t, helpText, "Tool Errors")
		assert.Contains(t, helpText, "Resource Issues")
		assert.Contains(t, helpText, "Preset Issues")
		assert.Contains(t, helpText, "Screenshot Issues")
	})

	t.Run("returns verbose troubleshooting help", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)
		defer db.Close()

		input := HelpInput{Topic: "troubleshooting", Verbose: true}
		_, result, err := server.handleHelp(context.Background(), nil, input)

		assert.NoError(t, err)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)

		helpText, ok := resultMap["help"].(string)
		require.True(t, ok)
		// Verbose should include additional sections
		assert.Contains(t, helpText, "Diagnostic Steps")
		assert.Contains(t, helpText, "Getting More Help")
		assert.Contains(t, helpText, "Common Error Messages")
	})

	t.Run("returns error for unknown topic", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)
		defer db.Close()

		input := HelpInput{Topic: "unknown_topic_xyz"}
		_, result, err := server.handleHelp(context.Background(), nil, input)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "unknown help topic")
		assert.Contains(t, err.Error(), "unknown_topic_xyz")
	})

	t.Run("handles case-insensitive topics", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)
		defer db.Close()

		input := HelpInput{Topic: "TOOLS"}
		_, result, err := server.handleHelp(context.Background(), nil, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "tools", resultMap["topic"])
	})

	t.Run("handles whitespace in topics", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)
		defer db.Close()

		input := HelpInput{Topic: "  overview  "}
		_, result, err := server.handleHelp(context.Background(), nil, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "overview", resultMap["topic"])
	})
}

// TestHandleHelpSpecificTools tests help for specific tool names
func TestHandleHelpSpecificTools(t *testing.T) {
	server, _, db := testServerWithStorage(t)
	defer db.Close()

	// Test a sample of tools across different categories
	toolTests := []struct {
		toolName     string
		shouldFind   []string
		categoryHint string
	}{
		{
			toolName:     "list_scenes",
			shouldFind:   []string{"list_scenes", "Category", "Core", "Scene Management", "Description", "Input", "Output"},
			categoryHint: "Core scene management",
		},
		{
			toolName:     "start_recording",
			shouldFind:   []string{"start_recording", "Category", "Core", "Recording", "Description"},
			categoryHint: "Core recording",
		},
		{
			toolName:     "create_screenshot_source",
			shouldFind:   []string{"create_screenshot_source", "Category", "Visual", "Screenshot Monitoring", "Description", "cadence_ms"},
			categoryHint: "Visual monitoring",
		},
		{
			toolName:     "save_scene_preset",
			shouldFind:   []string{"save_scene_preset", "Category", "Layout", "Scene Presets", "preset_name", "scene_name"},
			categoryHint: "Layout presets",
		},
		{
			toolName:     "toggle_input_mute",
			shouldFind:   []string{"toggle_input_mute", "Category", "Audio", "input_name"},
			categoryHint: "Audio control",
		},
		{
			toolName:     "create_text_source",
			shouldFind:   []string{"create_text_source", "Category", "Design", "Source Creation", "font_size", "color"},
			categoryHint: "Design source creation",
		},
		{
			toolName:     "set_source_transform",
			shouldFind:   []string{"set_source_transform", "Category", "Design", "Layout Control", "scale_x", "rotation"},
			categoryHint: "Design layout",
		},
	}

	for _, tt := range toolTests {
		t.Run("returns help for "+tt.toolName, func(t *testing.T) {
			input := HelpInput{Topic: tt.toolName, Verbose: false}
			_, result, err := server.handleHelp(context.Background(), nil, input)

			assert.NoError(t, err, "Should find help for %s (%s)", tt.toolName, tt.categoryHint)
			assert.NotNil(t, result)

			resultMap, ok := result.(map[string]interface{})
			require.True(t, ok)
			assert.Equal(t, tt.toolName, resultMap["topic"])

			helpText, ok := resultMap["help"].(string)
			require.True(t, ok)

			for _, expected := range tt.shouldFind {
				assert.Contains(t, helpText, expected, "Help for %s should contain '%s'", tt.toolName, expected)
			}
		})

		t.Run("returns verbose help for "+tt.toolName, func(t *testing.T) {
			input := HelpInput{Topic: tt.toolName, Verbose: true}
			_, result, err := server.handleHelp(context.Background(), nil, input)

			assert.NoError(t, err)
			assert.NotNil(t, result)

			resultMap, ok := result.(map[string]interface{})
			require.True(t, ok)
			assert.Equal(t, true, resultMap["verbose"])

			helpText, ok := resultMap["help"].(string)
			require.True(t, ok)

			// Verbose tool help should include "Related Resources"
			assert.Contains(t, helpText, "Related Resources")
			assert.Contains(t, helpText, "Example Workflow")
		})
	}
}

// TestGetOverviewHelp tests overview help content
func TestGetOverviewHelp(t *testing.T) {
	t.Run("basic overview includes key sections", func(t *testing.T) {
		help := GetOverviewHelp(false)
		assert.Contains(t, help, "agentic-obs")
		assert.Contains(t, help, "What is agentic-obs")
		assert.Contains(t, help, "Quick Start")
		assert.Contains(t, help, "Key Features")
		assert.Contains(t, help, "57 Tools")
		assert.Contains(t, help, "4 Resource Types")
	})

	t.Run("verbose overview includes additional sections", func(t *testing.T) {
		help := GetOverviewHelp(true)
		assert.Contains(t, help, "Categories")
		assert.Contains(t, help, "Common Workflows")
		assert.Contains(t, help, "Next Steps")
		assert.Contains(t, help, "Core Tools")
		assert.Contains(t, help, "Design Tools")
		assert.Contains(t, help, "Filters Tools")
		assert.Contains(t, help, "Transitions Tools")
	})
}

// TestGetToolsHelp tests tools help content
func TestGetToolsHelp(t *testing.T) {
	t.Run("basic tools help lists all categories", func(t *testing.T) {
		help := GetToolsHelp(false)
		assert.Contains(t, help, "All Available Tools")
		assert.Contains(t, help, "Core Tools")
		assert.Contains(t, help, "Help Tool")
		assert.Contains(t, help, "Sources Tools")
		assert.Contains(t, help, "Audio Tools")
		assert.Contains(t, help, "Layout Tools")
		assert.Contains(t, help, "Visual Tools")
		assert.Contains(t, help, "Design Tools")
		assert.Contains(t, help, "Filters Tools")
		assert.Contains(t, help, "Transitions Tools")

		// Check a few tool names
		assert.Contains(t, help, "list_scenes")
		assert.Contains(t, help, "start_recording")
		assert.Contains(t, help, "create_screenshot_source")
		assert.Contains(t, help, "save_scene_preset")
		assert.Contains(t, help, "list_source_filters")
		assert.Contains(t, help, "list_transitions")
	})

	t.Run("verbose tools help includes guidance", func(t *testing.T) {
		help := GetToolsHelp(true)
		assert.Contains(t, help, "Getting Help on Specific Tools")
		assert.Contains(t, help, "Tool Groups")
	})
}

// TestGetResourcesHelp tests resources help content
func TestGetResourcesHelp(t *testing.T) {
	t.Run("basic resources help describes all types", func(t *testing.T) {
		help := GetResourcesHelp(false)
		assert.Contains(t, help, "MCP Resources")
		assert.Contains(t, help, "OBS Scenes")
		assert.Contains(t, help, "Screenshot Images")
		assert.Contains(t, help, "Screenshot URLs")
		assert.Contains(t, help, "Scene Presets")
		assert.Contains(t, help, "obs://scene/")
		assert.Contains(t, help, "obs://screenshot/")
		assert.Contains(t, help, "obs://preset/")
	})

	t.Run("verbose resources help includes operations and examples", func(t *testing.T) {
		help := GetResourcesHelp(true)
		assert.Contains(t, help, "Resource Operations")
		assert.Contains(t, help, "Resource Notifications")
		assert.Contains(t, help, "Using Resources")
		assert.Contains(t, help, "Example Use Cases")
	})
}

// TestGetPromptsHelp tests prompts help content
func TestGetPromptsHelp(t *testing.T) {
	t.Run("basic prompts help lists all prompts", func(t *testing.T) {
		help := GetPromptsHelp(false)
		assert.Contains(t, help, "MCP Prompts")
		assert.Contains(t, help, "stream-launch")
		assert.Contains(t, help, "stream-teardown")
		assert.Contains(t, help, "recording-workflow")
		assert.Contains(t, help, "health-check")
		assert.Contains(t, help, "audio-check")
		assert.Contains(t, help, "visual-check")
		assert.Contains(t, help, "problem-detection")
		assert.Contains(t, help, "preset-switcher")
		assert.Contains(t, help, "scene-organizer")
		assert.Contains(t, help, "quick-status")
		assert.Contains(t, help, "scene-designer")
		assert.Contains(t, help, "source-management")
		assert.Contains(t, help, "visual-setup")
	})

	t.Run("verbose prompts help includes usage and examples", func(t *testing.T) {
		help := GetPromptsHelp(true)
		assert.Contains(t, help, "Using Prompts")
		assert.Contains(t, help, "Prompt Arguments")
		assert.Contains(t, help, "Creating Custom Workflows")
		assert.Contains(t, help, "Workflow Examples")
	})
}

// TestGetWorkflowsHelp tests workflows help content
func TestGetWorkflowsHelp(t *testing.T) {
	t.Run("basic workflows help describes common workflows", func(t *testing.T) {
		help := GetWorkflowsHelp(false)
		assert.Contains(t, help, "Common Workflows")
		assert.Contains(t, help, "Start Streaming")
		assert.Contains(t, help, "Visual Stream Monitoring")
		assert.Contains(t, help, "Scene Design from Scratch")
		assert.Contains(t, help, "Scene Preset Management")
		assert.Contains(t, help, "Audio Configuration")
		assert.Contains(t, help, "Multi-Scene Setup")
	})

	t.Run("verbose workflows help includes advanced patterns", func(t *testing.T) {
		help := GetWorkflowsHelp(true)
		assert.Contains(t, help, "Automated Stream Production")
		assert.Contains(t, help, "Workflow Best Practices")
		assert.Contains(t, help, "Pre-stream")
		assert.Contains(t, help, "During stream")
		assert.Contains(t, help, "Post-stream")
	})
}

// TestGetTroubleshootingHelp tests troubleshooting help content
func TestGetTroubleshootingHelp(t *testing.T) {
	t.Run("basic troubleshooting help covers common issues", func(t *testing.T) {
		help := GetTroubleshootingHelp(false)
		assert.Contains(t, help, "Troubleshooting Guide")
		assert.Contains(t, help, "Connection Issues")
		assert.Contains(t, help, "Tool Errors")
		assert.Contains(t, help, "Resource Issues")
		assert.Contains(t, help, "Preset Issues")
		assert.Contains(t, help, "Screenshot Issues")
		assert.Contains(t, help, "Not connected to OBS")
		assert.Contains(t, help, "Scene not found")
	})

	t.Run("verbose troubleshooting help includes diagnostic steps", func(t *testing.T) {
		help := GetTroubleshootingHelp(true)
		assert.Contains(t, help, "Diagnostic Steps")
		assert.Contains(t, help, "Getting More Help")
		assert.Contains(t, help, "Common Error Messages")
		assert.Contains(t, help, "For connection issues")
		assert.Contains(t, help, "For tool errors")
	})
}

// TestGetToolHelp tests specific tool help retrieval
func TestGetToolHelp(t *testing.T) {
	t.Run("returns help for valid tool", func(t *testing.T) {
		help, err := getToolHelp("list_scenes", false)
		assert.NoError(t, err)
		assert.NotEmpty(t, help)
		assert.Contains(t, help, "list_scenes")
		assert.Contains(t, help, "Category")
	})

	t.Run("returns verbose help for valid tool", func(t *testing.T) {
		help, err := getToolHelp("start_recording", true)
		assert.NoError(t, err)
		assert.NotEmpty(t, help)
		assert.Contains(t, help, "start_recording")
		assert.Contains(t, help, "Related Resources")
		assert.Contains(t, help, "Example Workflow")
	})

	t.Run("returns error for invalid tool", func(t *testing.T) {
		help, err := getToolHelp("invalid_tool_name", false)
		assert.Error(t, err)
		assert.Empty(t, help)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("covers all major tool categories", func(t *testing.T) {
		// Test at least one tool from each category to ensure coverage
		categories := map[string]string{
			"list_scenes":                  "Core",
			"start_recording":              "Core",
			"list_sources":                 "Sources",
			"toggle_input_mute":            "Audio",
			"save_scene_preset":            "Layout",
			"create_screenshot_source":     "Visual",
			"create_text_source":           "Design",
			"set_source_transform":         "Design",
			"duplicate_source":             "Design",
			"get_input_volume":             "Audio",
			"apply_scene_preset":           "Layout",
			"configure_screenshot_cadence": "Visual",
		}

		for toolName, category := range categories {
			help, err := getToolHelp(toolName, false)
			assert.NoError(t, err, "Should find help for %s", toolName)
			assert.Contains(t, help, toolName)
			assert.Contains(t, help, "Category")
			assert.Contains(t, help, category, "Tool %s should be in %s category", toolName, category)
		}
	})
}

// TestHelpContentCompleteness tests that help content is comprehensive
func TestHelpContentCompleteness(t *testing.T) {
	server, _, db := testServerWithStorage(t)
	defer db.Close()

	t.Run("all help topics are accessible", func(t *testing.T) {
		topics := []string{"overview", "tools", "resources", "prompts", "workflows", "troubleshooting"}

		for _, topic := range topics {
			input := HelpInput{Topic: topic}
			_, result, err := server.handleHelp(context.Background(), nil, input)

			assert.NoError(t, err, "Topic %s should be accessible", topic)
			assert.NotNil(t, result)

			resultMap, ok := result.(map[string]interface{})
			require.True(t, ok)

			helpText, ok := resultMap["help"].(string)
			require.True(t, ok)
			assert.NotEmpty(t, helpText, "Topic %s should have non-empty help text", topic)
			assert.Greater(t, len(helpText), 100, "Topic %s should have substantial help text", topic)
		}
	})

	t.Run("all tool names have help entries", func(t *testing.T) {
		// List of all 57 tools that should have help
		allTools := []string{
			// Core (13 tools)
			"list_scenes", "set_current_scene", "create_scene", "remove_scene",
			"start_recording", "stop_recording", "get_recording_status", "pause_recording", "resume_recording",
			"start_streaming", "stop_streaming", "get_streaming_status",
			"get_obs_status",
			// Sources (3 tools)
			"list_sources", "toggle_source_visibility", "get_source_settings",
			// Audio (4 tools)
			"get_input_mute", "toggle_input_mute", "set_input_volume", "get_input_volume",
			// Layout (6 tools)
			"save_scene_preset", "apply_scene_preset", "list_scene_presets",
			"get_preset_details", "rename_scene_preset", "delete_scene_preset",
			// Visual (4 tools)
			"create_screenshot_source", "remove_screenshot_source",
			"list_screenshot_sources", "configure_screenshot_cadence",
			// Design (14 tools)
			"create_text_source", "create_image_source", "create_color_source",
			"create_browser_source", "create_media_source",
			"set_source_transform", "get_source_transform", "set_source_crop",
			"set_source_bounds", "set_source_order",
			"set_source_locked", "duplicate_source", "remove_source", "list_input_kinds",
			// Filters (7 tools)
			"list_source_filters", "get_source_filter", "create_source_filter",
			"remove_source_filter", "toggle_source_filter", "set_source_filter_settings",
			"list_filter_kinds",
			// Transitions (5 tools)
			"list_transitions", "get_current_transition", "set_current_transition",
			"set_transition_duration", "trigger_transition",
		}

		for _, toolName := range allTools {
			help, err := getToolHelp(toolName, false)
			assert.NoError(t, err, "Tool %s should have help entry", toolName)
			assert.NotEmpty(t, help, "Tool %s should have non-empty help", toolName)
			assert.Contains(t, strings.ToLower(help), strings.ToLower(toolName),
				"Tool %s help should mention the tool name", toolName)
		}
	})
}

// TestHelpIntegration tests help integration scenarios
func TestHelpIntegration(t *testing.T) {
	server, _, db := testServerWithStorage(t)
	defer db.Close()

	t.Run("help can guide users through discovery flow", func(t *testing.T) {
		// Simulate user discovering the system
		// 1. Start with overview
		input := HelpInput{Topic: "overview"}
		_, result, err := server.handleHelp(context.Background(), nil, input)
		assert.NoError(t, err)

		resultMap := result.(map[string]interface{})
		overviewHelp := resultMap["help"].(string)
		assert.Contains(t, overviewHelp, "tools")
		assert.Contains(t, overviewHelp, "resources")

		// 2. Explore tools
		input = HelpInput{Topic: "tools"}
		_, result, err = server.handleHelp(context.Background(), nil, input)
		assert.NoError(t, err)

		resultMap = result.(map[string]interface{})
		toolsHelp := resultMap["help"].(string)
		assert.Contains(t, toolsHelp, "list_scenes")

		// 3. Get specific tool help
		input = HelpInput{Topic: "list_scenes"}
		_, result, err = server.handleHelp(context.Background(), nil, input)
		assert.NoError(t, err)

		resultMap = result.(map[string]interface{})
		toolHelp := resultMap["help"].(string)
		assert.Contains(t, toolHelp, "Description")
		assert.Contains(t, toolHelp, "Input")
		assert.Contains(t, toolHelp, "Output")
	})

	t.Run("help supports troubleshooting workflow", func(t *testing.T) {
		// User encounters issue and seeks help
		input := HelpInput{Topic: "troubleshooting"}
		_, result, err := server.handleHelp(context.Background(), nil, input)
		assert.NoError(t, err)

		resultMap := result.(map[string]interface{})
		troubleshootingHelp := resultMap["help"].(string)

		// Should guide user to relevant solutions
		assert.Contains(t, troubleshootingHelp, "Connection Issues")
		assert.Contains(t, troubleshootingHelp, "Tool Errors")
		assert.Contains(t, strings.ToLower(troubleshootingHelp), "not connected")
	})
}
