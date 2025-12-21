package mcp

import "fmt"

// ============================================================================
// SOURCE OF TRUTH: Help Metrics
// ============================================================================
//
// These constants are the SOURCE OF TRUTH for tool, resource, and prompt counts.
// All documentation should reference these values.
//
// VERIFICATION:
//   - Run: ./scripts/verify-docs.sh
//   - The script checks that documentation matches these constants
//
// UPDATE PROCEDURE (when adding tools/resources/prompts):
//  1. Update the appropriate constant below
//  2. Update HelpToolCount if total changes
//  3. Run ./scripts/verify-docs.sh to find docs needing updates
//  4. See scripts/DOC_UPDATE_CHECKLIST.md for full checklist
//
// ============================================================================
const (
	HelpToolCount     = 69 // Total MCP tools (including help)
	HelpResourceCount = 4  // Resource types: scenes, screenshots, screenshot-url, presets
	HelpPromptCount   = 13 // Workflow prompts

	// Tool counts by category (should sum to HelpToolCount)
	HelpCoreToolCount        = 25 // Scene management, recording, streaming, status, virtual cam, replay buffer, studio mode, hotkeys
	HelpHelpToolCount        = 1  // The help tool itself
	HelpSourcesToolCount     = 3  // Source management
	HelpAudioToolCount       = 4  // Audio control
	HelpLayoutToolCount      = 6  // Scene presets
	HelpVisualToolCount      = 4  // Screenshot monitoring
	HelpDesignToolCount      = 14 // Source creation and layout
	HelpFiltersToolCount     = 7  // Filter management (FB-23)
	HelpTransitionsToolCount = 5  // Transition control (FB-24)
)

// GetOverviewHelp returns high-level overview of agentic-obs
func GetOverviewHelp(verbose bool) string {
	help := fmt.Sprintf(`# agentic-obs - OBS Studio Control via MCP

**What is agentic-obs?**
A Model Context Protocol (MCP) server that gives AI assistants programmatic control over OBS Studio through the OBS WebSocket API.

## Quick Start

1. **Get help on available tools**: Use topic="tools" to see all %d tools grouped by category
2. **Explore resources**: Use topic="resources" to learn about scenes, screenshots, and presets
3. **Try workflows**: Use topic="prompts" to see pre-built workflows for common tasks
4. **Get specific help**: Use topic="<tool_name>" for detailed help on any tool

## Key Features

- **%d Tools** across 8 categories (Core, Sources, Audio, Layout, Visual, Design, Filters, Transitions) + Help
- **%d Resource Types** (scenes, screenshots, screenshot URLs, presets)
- **%d Workflow Prompts** for common streaming/recording tasks
- **Real-time Monitoring** via screenshot sources for AI visual inspection
- **Scene Presets** to save and restore source visibility states
- **Persistent Storage** for configuration and presets via SQLite
`, HelpToolCount, HelpToolCount, HelpResourceCount, HelpPromptCount)

	if verbose {
		help += fmt.Sprintf(`
## Categories

**Core Tools** (%d tools): Scene management, recording, streaming, status
**Sources Tools** (%d tools): List, visibility toggle, settings inspection
**Audio Tools** (%d tools): Mute control, volume adjustment
**Layout Tools** (%d tools): Scene preset save/restore/manage
**Visual Tools** (%d tools): Screenshot source creation and monitoring
**Design Tools** (%d tools): Source creation, transforms, positioning
**Filters Tools** (%d tools): Filter creation, toggle, settings
**Transitions Tools** (%d tools): Transition selection, duration, trigger

## Common Workflows

- **Start streaming**: Use 'stream-launch' prompt for pre-flight checklist
- **Stop streaming**: Use 'stream-teardown' prompt for clean shutdown
- **Monitor visually**: Create screenshot sources for AI to "see" your stream
- **Save layouts**: Use scene presets to capture and restore source visibility
- **Health check**: Use 'health-check' prompt for comprehensive diagnostics

## Next Steps

- topic="tools" - See all available tools organized by category
- topic="workflows" - Learn multi-step workflows for common tasks
- topic="prompts" - Discover pre-built workflow prompts
- topic="troubleshooting" - Common issues and solutions
`, HelpCoreToolCount, HelpSourcesToolCount, HelpAudioToolCount,
			HelpLayoutToolCount, HelpVisualToolCount, HelpDesignToolCount,
			HelpFiltersToolCount, HelpTransitionsToolCount)
	}

	return help
}

// GetToolsHelp returns comprehensive list of all tools grouped by category
func GetToolsHelp(verbose bool) string {
	help := fmt.Sprintf(`# All Available Tools (%d total)

## Core Tools (%d tools) - Scene Management, Recording, Streaming

**Scene Management:**
- list_scenes - List all scenes and identify current scene
- set_current_scene - Switch to a different scene
- create_scene - Create a new scene
- remove_scene - Delete a scene

**Recording:**
- start_recording - Begin recording
- stop_recording - End recording (returns output path)
- get_recording_status - Check recording state
- pause_recording - Pause active recording
- resume_recording - Resume paused recording

**Streaming:**
- start_streaming - Begin streaming
- stop_streaming - End stream
- get_streaming_status - Check streaming state

**Status:**
- get_obs_status - Overall OBS connection and state

## Help Tool (%d tool) - Always Enabled

- help - Get detailed help on tools, resources, prompts, workflows, or troubleshooting

## Sources Tools (%d tools) - Source Management

- list_sources - List all input sources (audio/video)
- toggle_source_visibility - Show/hide source in scene
- get_source_settings - Retrieve source configuration

## Audio Tools (%d tools) - Audio Control

- get_input_mute - Check if audio input is muted
- toggle_input_mute - Toggle mute state
- set_input_volume - Set volume (dB or multiplier)
- get_input_volume - Get current volume levels

## Layout Tools (%d tools) - Scene Presets

- save_scene_preset - Save current source visibility states
- apply_scene_preset - Restore saved preset
- list_scene_presets - List all saved presets
- get_preset_details - Get preset source states
- rename_scene_preset - Rename a preset
- delete_scene_preset - Delete a preset

## Visual Tools (%d tools) - Screenshot Monitoring

- create_screenshot_source - Create periodic screenshot capture
- remove_screenshot_source - Stop and remove screenshot source
- list_screenshot_sources - List all screenshot sources with URLs
- configure_screenshot_cadence - Update capture interval

## Design Tools (%d tools) - Source Creation & Layout

**Source Creation:**
- create_text_source - Add text/label with font customization
- create_image_source - Add image from file path
- create_color_source - Add solid color background
- create_browser_source - Add web content display
- create_media_source - Add video/media file

**Layout Control:**
- set_source_transform - Position, scale, rotation
- get_source_transform - Get current transform properties
- set_source_crop - Crop edges of source
- set_source_bounds - Set bounds type and size
- set_source_order - Z-order (layering)

**Advanced:**
- set_source_locked - Lock/unlock to prevent changes
- duplicate_source - Copy source within/between scenes
- remove_source - Delete source from scene
- list_input_kinds - List all available OBS source types

## Filters Tools (%d tools) - Filter Management

- list_source_filters - List all filters on a source
- get_source_filter - Get filter details and settings
- create_source_filter - Add a new filter (color correction, noise suppression, etc.)
- remove_source_filter - Remove a filter from source
- toggle_source_filter - Enable/disable a filter
- set_source_filter_settings - Modify filter configuration
- list_filter_kinds - List all available filter types

## Transitions Tools (%d tools) - Scene Transition Control

- list_transitions - List available transitions and current one
- get_current_transition - Get current transition details
- set_current_transition - Change active transition (Cut, Fade, Swipe, etc.)
- set_transition_duration - Set transition duration in milliseconds
- trigger_transition - Trigger studio mode transition (preview to program)
`, HelpToolCount, HelpCoreToolCount, HelpHelpToolCount, HelpSourcesToolCount,
		HelpAudioToolCount, HelpLayoutToolCount, HelpVisualToolCount, HelpDesignToolCount,
		HelpFiltersToolCount, HelpTransitionsToolCount)

	if verbose {
		help += `
## Getting Help on Specific Tools

Use topic="<tool_name>" for detailed help on any tool. Examples:
- topic="create_text_source" - Full parameters and examples for text sources
- topic="save_scene_preset" - How to save and restore scene layouts
- topic="create_screenshot_source" - Enable AI visual monitoring

## Tool Groups

Tools are organized by capability. You can:
- Use Core tools for basic OBS operations
- Use Visual tools to let AI "see" your stream
- Use Layout tools to save/restore complex setups
- Use Design tools to build scenes programmatically
- Use Filters tools to manage source effects (color correction, noise suppression)
- Use Transitions tools to control scene change animations
`
	}

	return help
}

// GetResourcesHelp returns information about MCP resources
func GetResourcesHelp(verbose bool) string {
	help := fmt.Sprintf(`# MCP Resources (%d types)

Resources provide direct access to OBS data through MCP resource URIs.

## 1. OBS Scenes
- **URI**: obs://scene/{sceneName}
- **Type**: application/json
- **Description**: Scene configuration and source layouts
- **Notifications**: Updates when scene is modified or becomes active

## 2. Screenshot Images
- **URI**: obs://screenshot/{sourceName}
- **Type**: image/png or image/jpeg
- **Description**: Binary screenshot data from screenshot sources
- **Usage**: AI can visually inspect stream output

## 3. Screenshot URLs
- **URI**: obs://screenshot-url/{sourceName}
- **Type**: application/json
- **Description**: HTTP URL for screenshot access (lightweight alternative)
- **Example**: {"url": "http://localhost:8765/screenshot/main"}

## 4. Scene Presets
- **URI**: obs://preset/{presetName}
- **Type**: application/json
- **Description**: Saved source visibility configurations
- **Usage**: Save and restore scene layouts
`, HelpResourceCount)

	if verbose {
		help += `
## Resource Operations

- **resources/list**: List all available resources
- **resources/read**: Get detailed resource content (JSON or binary)
- **resources/subscribe**: Subscribe to resource change notifications (future)

## Resource Notifications

The server sends notifications when resources change:
- **notifications/resources/updated**: Specific resource was modified
- **notifications/resources/list_changed**: Resources were created/deleted

## Using Resources

Resources complement tools by providing read-only access to OBS state:
- Tools: Modify OBS state (set scene, create source, etc.)
- Resources: Inspect OBS state (read scene config, view screenshots)

## Example Use Cases

1. **Visual Monitoring**: Create screenshot source, read via obs://screenshot/{name}
2. **Scene Inspection**: Read obs://scene/{name} to see all sources and settings
3. **Preset Management**: Save preset via tool, read via obs://preset/{name}
4. **Change Detection**: Subscribe to resource notifications for real-time updates
`
	}

	return help
}

// GetPromptsHelp returns information about workflow prompts
func GetPromptsHelp(verbose bool) string {
	help := fmt.Sprintf(`# MCP Prompts (%d workflows)

Pre-built prompts combine multiple tools into guided workflows.

## Streaming & Recording

**stream-launch** - Pre-stream checklist and setup guidance
- Verifies OBS connection, scene configuration, audio setup
- Ensures everything is ready before going live

**stream-teardown** - Post-stream cleanup and shutdown
- Stops streaming/recording, switches to offline scene
- Clean graceful shutdown workflow

**recording-workflow** - Complete recording session management
- Guide through start/stop, verify scene and audio
- Status checks and confirmation

## Diagnostics & Monitoring

**health-check** - Comprehensive OBS diagnostic
- Connection state, scenes, sources, streaming status
- Complete system health overview

**audio-check** - Audio verification and diagnostics
- Check all audio inputs, mute states, volume levels
- Identify audio configuration issues

**visual-check** - Visual layout analysis (requires screenshot_source)
- AI analyzes screenshot to verify layout
- Detects visual issues, alignment problems

**problem-detection** - Automated issue detection (requires screenshot_source)
- Detects black screens, frozen frames, wrong scenes
- Proactive monitoring for common problems

## Management

**preset-switcher** - Scene preset management and switching
- List available presets, optionally apply one
- Quick scene layout switching

**scene-organizer** - Scene organization and cleanup
- Helps organize, rename, delete unused scenes
- Keeps OBS workspace tidy

**quick-status** - Brief status summary for rapid checks
- Fast overview of current state
- Ideal for quick verification

## Scene Design

**scene-designer** - Visual layout creation (requires scene_name)
- Guide creating layouts using 14 Design tools
- Source creation, positioning, transforms
- Optional action argument for specific operations

**source-management** - Manage source visibility and properties (requires scene_name)
- Toggle source visibility in specified scene
- Configure source settings and properties
- Inventory and organize scene sources

**visual-setup** - Configure screenshot monitoring (optional monitor_scene)
- Set up screenshot sources for AI visual monitoring
- Configure capture cadence and image settings
- Enable real-time visual inspection
`, HelpPromptCount)

	if verbose {
		help += `
## Using Prompts

Prompts are invoked via the MCP prompts interface. In Claude Desktop or other MCP clients, use:
- prompts/list - See all available prompts
- prompts/get - Get a specific prompt with arguments

## Prompt Arguments

Some prompts require arguments:
- **visual-check**: screenshot_source (required) - Name of screenshot source to analyze
- **problem-detection**: screenshot_source (required) - Name of screenshot source to check
- **preset-switcher**: preset_name (optional) - Name of preset to apply
- **scene-designer**: scene_name (required), action (optional) - Scene to design and optional operation
- **source-management**: scene_name (required) - Scene to manage sources in
- **visual-setup**: monitor_scene (optional) - Scene to configure for visual monitoring

## Creating Custom Workflows

Combine multiple tools to create your own workflows:
1. Use tools for individual operations
2. Chain operations together logically
3. Add error checking and validation
4. Build reusable sequences for your use case

## Workflow Examples

**Pre-stream setup:**
1. health-check - Verify everything works
2. audio-check - Confirm audio configuration
3. visual-check - Verify layout via screenshot
4. set_current_scene - Switch to starting scene
5. start_streaming - Go live!

**Scene design:**
1. create_scene - Make new scene
2. create_color_source - Add background
3. create_text_source - Add title
4. set_source_transform - Position elements
5. save_scene_preset - Save the layout
`
	}

	return help
}

// GetWorkflowsHelp returns multi-tool workflows for common tasks
func GetWorkflowsHelp(verbose bool) string {
	help := `# Common Workflows

Multi-tool sequences for typical tasks.

## Workflow: Start Streaming

1. get_obs_status - Verify OBS is connected
2. list_scenes - Check available scenes
3. set_current_scene - Switch to starting scene
4. get_input_mute - Verify microphone is unmuted
5. create_screenshot_source - Enable visual monitoring (optional)
6. start_streaming - Go live!

## Workflow: Visual Stream Monitoring

1. create_screenshot_source - Create screenshot of output
   - name: "stream_monitor"
   - source_name: "Program" (or scene name)
   - cadence_ms: 5000 (capture every 5 seconds)
2. list_screenshot_sources - Get HTTP URL
3. Read obs://screenshot/stream_monitor - AI inspects visually
4. Use problem-detection prompt - Automated issue detection

## Workflow: Scene Design from Scratch

1. create_scene - Make new scene
   - scene_name: "Tutorial"
2. create_color_source - Add background
   - color: 0xFF1a1a1a (dark gray)
3. create_browser_source - Add web overlay
   - url: "https://example.com/overlay"
4. create_text_source - Add title
   - text: "Tutorial Stream"
   - font_size: 48
5. get_source_transform - Get item IDs for positioning
6. set_source_transform - Position each element
7. set_source_order - Set layering (background at bottom)
8. save_scene_preset - Save layout as "tutorial_default"

## Workflow: Scene Preset Management

1. list_scenes - Find scene to save
2. set_current_scene - Switch to target scene
3. Configure sources visibility via toggle_source_visibility
4. save_scene_preset - Save current state
   - preset_name: "gaming_webcam_only"
   - scene_name: "Gaming"
5. Later: apply_scene_preset to restore
   - preset_name: "gaming_webcam_only"

## Workflow: Audio Configuration

1. list_sources - Find audio inputs
2. get_input_mute - Check mute state
3. get_input_volume - Check volume level
4. set_input_volume - Adjust if needed
   - volume_db: -6.0 (reduce by 6dB)
5. audio-check prompt - Verify all audio

## Workflow: Multi-Scene Setup

1. Create multiple scenes:
   - create_scene "Intro"
   - create_scene "Main Content"
   - create_scene "BRB"
   - create_scene "Outro"
2. Design each scene with Design tools
3. Save presets for each scene layout
4. Switch between scenes during stream:
   - set_current_scene for transitions
`

	if verbose {
		help += `
## Advanced Workflow: Automated Stream Production

1. **Pre-stream**:
   - stream-launch prompt - Verify readiness
   - create_screenshot_source - Enable monitoring
   - save_scene_preset for each scene - Save all layouts

2. **During stream**:
   - set_current_scene - Switch scenes
   - apply_scene_preset - Quick layout changes
   - problem-detection prompt - Monitor for issues
   - get_recording_status - Verify recording

3. **Post-stream**:
   - stream-teardown prompt - Clean shutdown
   - get_recording_status - Get output path
   - remove_screenshot_source - Clean up monitoring

## Workflow Best Practices

- **Always verify status** before starting operations
- **Use presets** for repeatable scene configurations
- **Enable monitoring** for long streams (screenshot sources)
- **Check audio** before going live (audio-check prompt)
- **Save your layouts** so you can restore them later
- **Use prompts** for complex multi-step operations
`
	}

	return help
}

// GetTroubleshootingHelp returns common issues and solutions
func GetTroubleshootingHelp(verbose bool) string {
	help := `# Troubleshooting Guide

Common issues and how to resolve them.

## Connection Issues

**Problem: "Not connected to OBS"**
- Verify OBS Studio is running
- Check OBS WebSocket server is enabled (Tools > WebSocket Server Settings)
- Confirm connection details (default: localhost:4455)
- Check OBS WebSocket password matches configuration

**Problem: Connection drops during use**
- OBS may have crashed or restarted
- Check OBS application is still running
- Server will auto-reconnect when OBS is available

## Tool Errors

**Problem: "Scene not found"**
- Use list_scenes to see available scenes
- Check scene name spelling (case-sensitive)
- Scene may have been deleted

**Problem: "Source not found"**
- Use list_sources to see available sources
- Check source name spelling (case-sensitive)
- Source may have been removed

**Problem: "Cannot start recording/streaming"**
- Check OBS output settings are configured
- Verify not already recording/streaming
- Use get_recording_status or get_streaming_status first

## Resource Issues

**Problem: Screenshot resource returns no data**
- Screenshot source may not have captured yet (wait for first capture)
- Check screenshot source exists via list_screenshot_sources
- Verify cadence_ms is reasonable (not too long)
- Check OBS scene/source exists

**Problem: Scene resource shows empty sources**
- Scene may actually be empty
- Sources may be in different scene
- Use list_sources to verify sources exist

## Preset Issues

**Problem: "Preset not found"**
- Use list_scene_presets to see available presets
- Check preset name spelling (case-sensitive)
- Preset may have been deleted

**Problem: Apply preset doesn't restore all sources**
- Sources may have been deleted since preset was saved
- Scene configuration may have changed
- Save a new preset with current sources

## Screenshot Issues

**Problem: Screenshot source not capturing**
- Verify OBS source exists and is rendering
- Check cadence_ms is set (default: 5000ms)
- Use list_screenshot_sources to verify source status
- Source scene may not be active (screenshots still work on inactive scenes)

**Problem: HTTP screenshot URL not accessible**
- Verify HTTP server is enabled (default: localhost:8765)
- Check firewall settings
- Use obs://screenshot/{name} resource for direct binary access
`

	if verbose {
		help += `
## Diagnostic Steps

**For connection issues:**
1. get_obs_status - Check overall status
2. Restart OBS Studio
3. Verify WebSocket settings in OBS
4. Check server logs for connection errors

**For tool errors:**
1. Use list_* tools to verify names
2. Check tool input parameters carefully
3. Use get_obs_status to verify OBS state
4. Try operation manually in OBS to confirm it works

**For resource issues:**
1. Use resources/list to see available resources
2. Check resource URI format matches obs://{type}/{name}
3. Verify resource exists (scene, screenshot source, preset)
4. Use corresponding list tool to find correct names

**For performance issues:**
1. Reduce screenshot source cadence_ms (longer interval)
2. Reduce screenshot image quality or size
3. Check OBS CPU usage
4. Limit number of active screenshot sources

## Getting More Help

- topic="overview" - High-level overview
- topic="<tool_name>" - Detailed help on specific tool
- Use health-check prompt - Comprehensive diagnostics
- Check server logs for error details
- Verify OBS logs for OBS-side issues

## Common Error Messages

- "not connected": OBS WebSocket not available
- "not found": Resource/scene/source doesn't exist
- "already exists": Trying to create duplicate
- "already active": Recording/streaming already running
- "not active": Trying to stop when not running
- "invalid settings": Parameter validation failed
`
	}

	return help
}

// VerboseToolSuffix is appended to all tool help when verbose is true
const VerboseToolSuffix = `

## Related Resources

- topic="tools" - See all tools grouped by category
- topic="workflows" - Multi-tool workflows using this tool
- topic="troubleshooting" - Common issues and solutions

## Example Workflow

Combine this tool with others to accomplish complex tasks. See topic="workflows" for examples.
`
