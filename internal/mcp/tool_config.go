// Package mcp provides MCP server implementation for OBS control.
package mcp

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ironystock/agentic-obs/internal/storage"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// ToolGroupMetadata contains static information about each tool group.
type ToolGroupMetadata struct {
	Name        string   // Group name (e.g., "Core", "Audio")
	Description string   // Human-readable description
	ToolCount   int      // Number of tools in this group
	ToolNames   []string // Tool names in this group
}

// ToolGroupOrder defines the canonical ordering of tool groups.
// Used for consistent iteration and validation across the codebase.
var ToolGroupOrder = []string{"Core", "Sources", "Audio", "Layout", "Visual", "Design", "Filters", "Transitions"}

// toolGroupMetadata defines metadata for all tool groups.
var toolGroupMetadata = map[string]*ToolGroupMetadata{
	"Core": {
		Name:        "Core",
		Description: "Core OBS tools: scenes, recording, streaming, status, virtual camera, replay buffer, studio mode, and hotkeys",
		ToolCount:   25,
		ToolNames: []string{
			"list_scenes", "set_current_scene", "create_scene", "remove_scene",
			"start_recording", "stop_recording", "get_recording_status", "pause_recording", "resume_recording",
			"start_streaming", "stop_streaming", "get_streaming_status",
			"get_obs_status",
			"get_virtual_cam_status", "toggle_virtual_cam",
			"get_replay_buffer_status", "toggle_replay_buffer", "save_replay_buffer", "get_last_replay",
			"get_studio_mode_enabled", "toggle_studio_mode", "get_preview_scene", "set_preview_scene",
			"list_hotkeys", "trigger_hotkey_by_name",
		},
	},
	"Sources": {
		Name:        "Sources",
		Description: "Source management: listing sources, visibility control, and settings",
		ToolCount:   3,
		ToolNames:   []string{"list_sources", "toggle_source_visibility", "get_source_settings"},
	},
	"Audio": {
		Name:        "Audio",
		Description: "Audio input control: mute state and volume levels",
		ToolCount:   4,
		ToolNames:   []string{"get_input_mute", "toggle_input_mute", "set_input_volume", "get_input_volume"},
	},
	"Layout": {
		Name:        "Layout",
		Description: "Scene preset management: save, apply, and organize source visibility presets",
		ToolCount:   6,
		ToolNames:   []string{"save_scene_preset", "list_scene_presets", "get_preset_details", "apply_scene_preset", "rename_scene_preset", "delete_scene_preset"},
	},
	"Visual": {
		Name:        "Visual",
		Description: "Visual monitoring: screenshot capture sources for AI visual analysis",
		ToolCount:   4,
		ToolNames:   []string{"create_screenshot_source", "remove_screenshot_source", "list_screenshot_sources", "configure_screenshot_cadence"},
	},
	"Design": {
		Name:        "Design",
		Description: "Scene design: create sources (text, image, browser, media) and control transforms",
		ToolCount:   14,
		ToolNames: []string{
			"create_text_source", "create_image_source", "create_color_source", "create_browser_source", "create_media_source",
			"set_source_transform", "get_source_transform", "set_source_crop", "set_source_bounds", "set_source_order",
			"set_source_locked", "duplicate_source", "remove_source", "list_input_kinds",
		},
	},
	"Filters": {
		Name:        "Filters",
		Description: "Source filter management: create, configure, and toggle filters on sources",
		ToolCount:   7,
		ToolNames:   []string{"list_source_filters", "get_source_filter", "create_source_filter", "remove_source_filter", "toggle_source_filter", "set_source_filter_settings", "list_filter_kinds"},
	},
	"Transitions": {
		Name:        "Transitions",
		Description: "Scene transition control: list, set, and trigger transitions",
		ToolCount:   5,
		ToolNames:   []string{"list_transitions", "get_current_transition", "set_current_transition", "set_transition_duration", "trigger_transition"},
	},
}

// MetaToolNames are tools that are always enabled and cannot be disabled.
var MetaToolNames = []string{"help", "get_tool_config", "set_tool_config", "list_tool_groups"}

// ToolGroupInfo represents information about a tool group for API responses.
type ToolGroupInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Enabled     bool     `json:"enabled"`
	ToolCount   int      `json:"tool_count"`
	Tools       []string `json:"tools,omitempty"` // Only included with verbose=true
}

// GetToolConfigInput is the input for querying tool configuration.
type GetToolConfigInput struct {
	Group   string `json:"group,omitempty" jsonschema:"Filter by group name (Core, Visual, Audio, Layout, Sources, Design, Filters, Transitions)"`
	Verbose bool   `json:"verbose,omitempty" jsonschema:"Include list of tool names per group"`
}

// SetToolConfigInput is the input for modifying tool configuration.
type SetToolConfigInput struct {
	Group   string `json:"group" jsonschema:"Tool group name to enable/disable"`
	Enabled bool   `json:"enabled" jsonschema:"True to enable the group, false to disable"`
	Persist bool   `json:"persist,omitempty" jsonschema:"Save to database for future sessions (default: session-only)"`
}

// ListToolGroupsInput is the input for listing tool groups.
type ListToolGroupsInput struct {
	IncludeDisabled bool `json:"include_disabled,omitempty" jsonschema:"Include disabled groups in listing (default: true)"`
}

// handleGetToolConfig returns the current tool configuration.
func (s *Server) handleGetToolConfig(ctx context.Context, request *mcpsdk.CallToolRequest, input GetToolConfigInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Getting tool configuration (group=%s, verbose=%v)", input.Group, input.Verbose)

	s.toolGroupMutex.RLock()
	defer s.toolGroupMutex.RUnlock()

	var groups []ToolGroupInfo

	// Build group info based on current state
	for _, groupName := range ToolGroupOrder {
		// If filtering by group, skip non-matching groups
		if input.Group != "" && input.Group != groupName {
			continue
		}

		meta := toolGroupMetadata[groupName]
		if meta == nil {
			continue
		}

		info := ToolGroupInfo{
			Name:        meta.Name,
			Description: meta.Description,
			Enabled:     s.getGroupEnabled(groupName),
			ToolCount:   meta.ToolCount,
		}

		if input.Verbose {
			info.Tools = meta.ToolNames
		}

		groups = append(groups, info)
	}

	// Calculate totals
	totalTools := 0
	enabledTools := 0
	for _, g := range groups {
		totalTools += g.ToolCount
		if g.Enabled {
			enabledTools += g.ToolCount
		}
	}

	// Add meta-tools to count (always enabled)
	totalTools += len(MetaToolNames)
	enabledTools += len(MetaToolNames)

	result := map[string]interface{}{
		"groups":        groups,
		"total_tools":   totalTools,
		"enabled_tools": enabledTools,
		"meta_tools":    MetaToolNames,
		"message":       fmt.Sprintf("%d of %d tools enabled across %d groups", enabledTools, totalTools, len(groups)),
	}

	s.recordAction("get_tool_config", "Get tool configuration", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleSetToolConfig enables or disables a tool group.
func (s *Server) handleSetToolConfig(ctx context.Context, request *mcpsdk.CallToolRequest, input SetToolConfigInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Setting tool config: group=%s, enabled=%v, persist=%v", input.Group, input.Enabled, input.Persist)

	// Validate group name (before acquiring lock to reduce contention on invalid input)
	meta := toolGroupMetadata[input.Group]
	if meta == nil {
		return nil, nil, fmt.Errorf("invalid group name '%s'. Valid groups: %v", input.Group, ToolGroupOrder)
	}

	s.toolGroupMutex.Lock()
	previousState := s.getGroupEnabled(input.Group)
	s.setGroupEnabled(input.Group, input.Enabled)
	s.toolGroupMutex.Unlock()

	// Persist if requested
	persisted := false
	var persistError string
	if input.Persist && s.storage != nil {
		if err := s.storage.SaveToolGroupConfig(ctx, s.convertToStorageConfig()); err != nil {
			log.Printf("Warning: failed to persist tool config: %v", err)
			persistError = err.Error()
		} else {
			persisted = true
		}
	}

	action := "enabled"
	if !input.Enabled {
		action = "disabled"
	}

	result := map[string]interface{}{
		"group":          input.Group,
		"previous_state": previousState,
		"new_state":      input.Enabled,
		"tools_affected": meta.ToolCount,
		"persisted":      persisted,
		"message":        fmt.Sprintf("Tool group '%s' (%d tools) %s", input.Group, meta.ToolCount, action),
	}

	// Include persistence error if it occurred
	if persistError != "" {
		result["persist_error"] = persistError
		result["message"] = fmt.Sprintf("Tool group '%s' (%d tools) %s (persistence failed: %s)", input.Group, meta.ToolCount, action, persistError)
	}

	// Note: In Phase 2, we would dynamically register/unregister tools here
	// For now, the configuration is tracked but tools remain registered
	// (handlers check group enabled state before executing)

	s.recordAction("set_tool_config", "Set tool configuration", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleListToolGroups lists all available tool groups.
func (s *Server) handleListToolGroups(ctx context.Context, request *mcpsdk.CallToolRequest, input ListToolGroupsInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Listing tool groups")

	s.toolGroupMutex.RLock()
	defer s.toolGroupMutex.RUnlock()

	var groups []ToolGroupInfo

	for _, groupName := range ToolGroupOrder {
		meta := toolGroupMetadata[groupName]
		if meta == nil {
			continue
		}

		enabled := s.getGroupEnabled(groupName)

		// Skip disabled groups if not including them
		if !input.IncludeDisabled && !enabled {
			continue
		}

		groups = append(groups, ToolGroupInfo{
			Name:        meta.Name,
			Description: meta.Description,
			Enabled:     enabled,
			ToolCount:   meta.ToolCount,
		})
	}

	result := map[string]interface{}{
		"groups":     groups,
		"count":      len(groups),
		"meta_tools": MetaToolNames,
		"message":    fmt.Sprintf("Found %d tool groups", len(groups)),
	}

	s.recordAction("list_tool_groups", "List tool groups", input, result, true, time.Since(start))
	return nil, result, nil
}

// getGroupEnabled returns whether a tool group is enabled.
// Must be called with toolGroupMutex held.
func (s *Server) getGroupEnabled(group string) bool {
	switch group {
	case "Core":
		return s.toolGroups.Core
	case "Sources":
		return s.toolGroups.Sources
	case "Audio":
		return s.toolGroups.Audio
	case "Layout":
		return s.toolGroups.Layout
	case "Visual":
		return s.toolGroups.Visual
	case "Design":
		return s.toolGroups.Design
	case "Filters":
		return s.toolGroups.Filters
	case "Transitions":
		return s.toolGroups.Transitions
	default:
		return false
	}
}

// setGroupEnabled sets whether a tool group is enabled.
// Must be called with toolGroupMutex held.
func (s *Server) setGroupEnabled(group string, enabled bool) {
	switch group {
	case "Core":
		s.toolGroups.Core = enabled
	case "Sources":
		s.toolGroups.Sources = enabled
	case "Audio":
		s.toolGroups.Audio = enabled
	case "Layout":
		s.toolGroups.Layout = enabled
	case "Visual":
		s.toolGroups.Visual = enabled
	case "Design":
		s.toolGroups.Design = enabled
	case "Filters":
		s.toolGroups.Filters = enabled
	case "Transitions":
		s.toolGroups.Transitions = enabled
	}
}

// convertToStorageConfig converts the server's tool group config to storage format.
func (s *Server) convertToStorageConfig() storage.ToolGroupConfig {
	return storage.ToolGroupConfig{
		Core:        s.toolGroups.Core,
		Visual:      s.toolGroups.Visual,
		Layout:      s.toolGroups.Layout,
		Audio:       s.toolGroups.Audio,
		Sources:     s.toolGroups.Sources,
		Design:      s.toolGroups.Design,
		Filters:     s.toolGroups.Filters,
		Transitions: s.toolGroups.Transitions,
	}
}
