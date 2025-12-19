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

	// Route to appropriate help handler
	switch topic {
	case "overview":
		helpText = s.getOverviewHelp(input.Verbose)
	case "tools":
		helpText = s.getToolsHelp(input.Verbose)
	case "resources":
		helpText = s.getResourcesHelp(input.Verbose)
	case "prompts":
		helpText = s.getPromptsHelp(input.Verbose)
	case "workflows":
		helpText = s.getWorkflowsHelp(input.Verbose)
	case "troubleshooting":
		helpText = s.getTroubleshootingHelp(input.Verbose)
	default:
		// Try to find help for a specific tool
		helpText, err = s.getToolHelp(topic, input.Verbose)
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

// getOverviewHelp returns high-level overview of agentic-obs
func (s *Server) getOverviewHelp(verbose bool) string {
	help := `# agentic-obs - OBS Studio Control via MCP

**What is agentic-obs?**
An Model Context Protocol (MCP) server that gives AI assistants programmatic control over OBS Studio through the OBS WebSocket API.

## Quick Start

1. **Get help on available tools**: Use topic="tools" to see all 45 tools grouped by category
2. **Explore resources**: Use topic="resources" to learn about scenes, screenshots, and presets
3. **Try workflows**: Use topic="prompts" to see pre-built workflows for common tasks
4. **Get specific help**: Use topic="<tool_name>" for detailed help on any tool

## Key Features

- **45 Tools** across 6 categories (Core, Sources, Audio, Layout, Visual, Design) + Help
- **4 Resource Types** (scenes, screenshots, screenshot URLs, presets)
- **13 Workflow Prompts** for common streaming/recording tasks
- **Real-time Monitoring** via screenshot sources for AI visual inspection
- **Scene Presets** to save and restore source visibility states
- **Persistent Storage** for configuration and presets via SQLite
`

	if verbose {
		help += `
## Categories

**Core Tools** (13 tools): Scene management, recording, streaming, status
**Sources Tools** (3 tools): List, visibility toggle, settings inspection
**Audio Tools** (4 tools): Mute control, volume adjustment
**Layout Tools** (6 tools): Scene preset save/restore/manage
**Visual Tools** (4 tools): Screenshot source creation and monitoring
**Design Tools** (14 tools): Source creation, transforms, positioning

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
`
	}

	return help
}

// getToolsHelp returns comprehensive list of all tools grouped by category
func (s *Server) getToolsHelp(verbose bool) string {
	help := `# All Available Tools (45 total)

## Core Tools (13 tools) - Scene Management, Recording, Streaming

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

## Help Tool (1 tool) - Always Enabled

- help - Get detailed help on tools, resources, prompts, workflows, or troubleshooting

## Sources Tools (3 tools) - Source Management

- list_sources - List all input sources (audio/video)
- toggle_source_visibility - Show/hide source in scene
- get_source_settings - Retrieve source configuration

## Audio Tools (4 tools) - Audio Control

- get_input_mute - Check if audio input is muted
- toggle_input_mute - Toggle mute state
- set_input_volume - Set volume (dB or multiplier)
- get_input_volume - Get current volume levels

## Layout Tools (6 tools) - Scene Presets

- save_scene_preset - Save current source visibility states
- apply_scene_preset - Restore saved preset
- list_scene_presets - List all saved presets
- get_preset_details - Get preset source states
- rename_scene_preset - Rename a preset
- delete_scene_preset - Delete a preset

## Visual Tools (4 tools) - Screenshot Monitoring

- create_screenshot_source - Create periodic screenshot capture
- remove_screenshot_source - Stop and remove screenshot source
- list_screenshot_sources - List all screenshot sources with URLs
- configure_screenshot_cadence - Update capture interval

## Design Tools (14 tools) - Source Creation & Layout

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
`

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
`
	}

	return help
}

// getResourcesHelp returns information about MCP resources
func (s *Server) getResourcesHelp(verbose bool) string {
	help := `# MCP Resources (4 types)

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
`

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

// getPromptsHelp returns information about workflow prompts
func (s *Server) getPromptsHelp(verbose bool) string {
	help := `# MCP Prompts (13 workflows)

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
`

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

// getWorkflowsHelp returns multi-tool workflows for common tasks
func (s *Server) getWorkflowsHelp(verbose bool) string {
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

// getTroubleshootingHelp returns common issues and solutions
func (s *Server) getTroubleshootingHelp(verbose bool) string {
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

// getToolHelp returns detailed help for a specific tool
func (s *Server) getToolHelp(toolName string, verbose bool) (string, error) {
	// Map of all tools with their detailed help information
	toolHelp := map[string]string{
		// Core - Scenes
		"list_scenes": `# list_scenes

**Category**: Core - Scene Management

**Description**: List all available scenes in OBS and identify the current active scene.

**Input**: None

**Output**:
- scenes: Array of scene names
- current_scene: Name of currently active scene

**Example**:
{
  "scenes": ["Scene 1", "Gaming", "BRB"],
  "current_scene": "Gaming"
}`,

		"set_current_scene": `# set_current_scene

**Category**: Core - Scene Management

**Description**: Switch to a different scene in OBS.

**Input**:
- scene_name (string, required): Name of scene to switch to

**Output**:
- message: Success confirmation

**Example Input**:
{
  "scene_name": "Gaming"
}`,

		"create_scene": `# create_scene

**Category**: Core - Scene Management

**Description**: Create a new empty scene in OBS.

**Input**:
- scene_name (string, required): Name for the new scene

**Output**:
- message: Success confirmation

**Example Input**:
{
  "scene_name": "New Tutorial Scene"
}

**Note**: Scene names must be unique. Use list_scenes to check existing scenes.`,

		"remove_scene": `# remove_scene

**Category**: Core - Scene Management

**Description**: Remove a scene from OBS. Cannot remove currently active scene.

**Input**:
- scene_name (string, required): Name of scene to remove

**Output**:
- message: Success confirmation

**Example Input**:
{
  "scene_name": "Old Scene"
}

**Warning**: This permanently deletes the scene and all its source configurations.`,

		// Core - Recording
		"start_recording": `# start_recording

**Category**: Core - Recording

**Description**: Start recording in OBS using configured output settings.

**Input**: None

**Output**:
- message: Success confirmation

**Requirements**: OBS recording output must be configured (Settings > Output > Recording)`,

		"stop_recording": `# stop_recording

**Category**: Core - Recording

**Description**: Stop the current recording and finalize the output file.

**Input**: None

**Output**:
- message: Success confirmation with output file path

**Example Output**:
{
  "message": "Successfully stopped recording. Output saved to: C:/Videos/2024-01-15_12-30-45.mp4"
}`,

		"get_recording_status": `# get_recording_status

**Category**: Core - Recording

**Description**: Check current recording state and statistics.

**Input**: None

**Output**:
- output_active: Whether recording is active (bool)
- output_paused: Whether recording is paused (bool)
- output_duration: Recording duration in milliseconds
- output_bytes: Total bytes recorded

**Example Output**:
{
  "output_active": true,
  "output_paused": false,
  "output_duration": 125000,
  "output_bytes": 52428800
}`,

		"pause_recording": `# pause_recording

**Category**: Core - Recording

**Description**: Pause an active recording. Recording must be in progress.

**Input**: None

**Output**:
- message: Success confirmation

**Note**: Use resume_recording to continue. Not all output formats support pausing.`,

		"resume_recording": `# resume_recording

**Category**: Core - Recording

**Description**: Resume a paused recording. Recording must be paused.

**Input**: None

**Output**:
- message: Success confirmation`,

		// Core - Streaming
		"start_streaming": `# start_streaming

**Category**: Core - Streaming

**Description**: Start streaming using configured stream settings.

**Input**: None

**Output**:
- message: Success confirmation

**Requirements**: OBS stream settings must be configured (Settings > Stream)`,

		"stop_streaming": `# stop_streaming

**Category**: Core - Streaming

**Description**: Stop the current stream.

**Input**: None

**Output**:
- message: Success confirmation`,

		"get_streaming_status": `# get_streaming_status

**Category**: Core - Streaming

**Description**: Check current streaming state and statistics.

**Input**: None

**Output**:
- output_active: Whether streaming is active (bool)
- output_reconnecting: Whether attempting to reconnect (bool)
- output_duration: Stream duration in milliseconds
- output_bytes: Total bytes sent

**Example Output**:
{
  "output_active": true,
  "output_reconnecting": false,
  "output_duration": 3600000,
  "output_bytes": 1073741824
}`,

		// Core - Status
		"get_obs_status": `# get_obs_status

**Category**: Core - Status

**Description**: Get comprehensive OBS status including version, scenes, and output states.

**Input**: None

**Output**:
- version: OBS Studio version
- websocket_version: OBS WebSocket version
- current_scene: Active scene name
- recording_active: Recording state
- streaming_active: Streaming state

**Example Output**:
{
  "version": "30.0.0",
  "websocket_version": "5.5.6",
  "current_scene": "Gaming",
  "recording_active": false,
  "streaming_active": true
}

**Use Case**: Verify OBS connection and overall state before starting operations.`,

		// Sources
		"list_sources": `# list_sources

**Category**: Sources

**Description**: List all input sources (audio and video) available in OBS.

**Input**: None

**Output**: Array of source objects with:
- name: Source name
- type: Source type/kind
- type_id: OBS internal type ID

**Example Output**:
[
  {"name": "Webcam", "type": "Video Capture Device", "type_id": "dshow_input"},
  {"name": "Microphone", "type": "Audio Input Capture", "type_id": "wasapi_input_capture"}
]`,

		"toggle_source_visibility": `# toggle_source_visibility

**Category**: Sources

**Description**: Toggle the visibility of a source in a specific scene (show/hide).

**Input**:
- scene_name (string, required): Name of scene containing source
- source_id (int, required): Scene item ID of the source

**Output**:
- scene_name: Scene name
- source_id: Source ID
- visible: New visibility state (bool)

**Example Input**:
{
  "scene_name": "Gaming",
  "source_id": 1
}

**Note**: Use list_sources or read obs://scene/{name} resource to get source IDs.`,

		"get_source_settings": `# get_source_settings

**Category**: Sources

**Description**: Retrieve configuration settings for a specific source.

**Input**:
- source_name (string, required): Name of source

**Output**: Source settings object (varies by source type)

**Example Input**:
{
  "source_name": "Webcam"
}

**Use Case**: Inspect source configuration, useful for debugging or verification.`,

		// Audio
		"get_input_mute": `# get_input_mute

**Category**: Audio

**Description**: Check whether an audio input is currently muted.

**Input**:
- input_name (string, required): Name of audio input

**Output**:
- input_name: Audio input name
- is_muted: Mute state (bool)

**Example Input**:
{
  "input_name": "Microphone"
}`,

		"toggle_input_mute": `# toggle_input_mute

**Category**: Audio

**Description**: Toggle the mute state of an audio input (muted <-> unmuted).

**Input**:
- input_name (string, required): Name of audio input

**Output**:
- message: Success confirmation

**Example Input**:
{
  "input_name": "Microphone"
}`,

		"set_input_volume": `# set_input_volume

**Category**: Audio

**Description**: Set the volume level of an audio input. Supports dB or multiplier format.

**Input**:
- input_name (string, required): Name of audio input
- volume_db (float, optional): Volume in decibels (e.g., -6.0)
- volume_mul (float, optional): Volume as linear multiplier (e.g., 0.5)

**Output**:
- message: Success confirmation

**Example Input (dB)**:
{
  "input_name": "Microphone",
  "volume_db": -6.0
}

**Example Input (multiplier)**:
{
  "input_name": "Desktop Audio",
  "volume_mul": 0.75
}

**Note**: Provide either volume_db OR volume_mul, not both. 0 dB = no change, negative = quieter.`,

		"get_input_volume": `# get_input_volume

**Category**: Audio

**Description**: Get the current volume level of an audio input in both dB and multiplier formats.

**Input**:
- input_name (string, required): Name of audio input

**Output**:
- input_name: Audio input name
- volume_db: Volume in decibels
- volume_mul: Volume as linear multiplier

**Example Input**:
{
  "input_name": "Microphone"
}

**Example Output**:
{
  "input_name": "Microphone",
  "volume_db": -6.0,
  "volume_mul": 0.5011872336272722
}`,

		// Layout - Scene Presets
		"save_scene_preset": `# save_scene_preset

**Category**: Layout - Scene Presets

**Description**: Save the current source visibility states from an OBS scene as a named preset.

**Input**:
- preset_name (string, required): Name for the new preset
- scene_name (string, required): Name of OBS scene to capture state from

**Output**:
- id: Preset ID
- preset_name: Preset name
- scene_name: Scene name
- source_count: Number of sources saved
- message: Success confirmation

**Example Input**:
{
  "preset_name": "gaming_webcam_only",
  "scene_name": "Gaming"
}

**Use Case**: Save complex scene layouts so you can restore them later with apply_scene_preset.`,

		"apply_scene_preset": `# apply_scene_preset

**Category**: Layout - Scene Presets

**Description**: Load a saved preset and apply its source visibility states to the target scene.

**Input**:
- preset_name (string, required): Name of preset to apply

**Output**:
- preset_name: Preset name
- scene_name: Scene name
- applied_count: Number of sources updated
- message: Success confirmation

**Example Input**:
{
  "preset_name": "gaming_webcam_only"
}

**Note**: Sources that no longer exist in the scene are skipped automatically.`,

		"list_scene_presets": `# list_scene_presets

**Category**: Layout - Scene Presets

**Description**: List all saved scene presets, optionally filtered by scene name.

**Input**:
- scene_name (string, optional): Filter presets for specific scene

**Output**:
- presets: Array of preset summaries (id, name, scene_name, created_at)
- count: Total number of presets

**Example Input** (all presets):
{}

**Example Input** (filtered):
{
  "scene_name": "Gaming"
}`,

		"get_preset_details": `# get_preset_details

**Category**: Layout - Scene Presets

**Description**: Get full details of a scene preset including all source states.

**Input**:
- preset_name (string, required): Name of preset

**Output**:
- id: Preset ID
- name: Preset name
- scene_name: Scene name
- sources: Array of source states (name, visible)
- created_at: Creation timestamp

**Example Input**:
{
  "preset_name": "gaming_webcam_only"
}`,

		"rename_scene_preset": `# rename_scene_preset

**Category**: Layout - Scene Presets

**Description**: Change the name of an existing scene preset.

**Input**:
- old_name (string, required): Current preset name
- new_name (string, required): New preset name

**Output**:
- message: Success confirmation

**Example Input**:
{
  "old_name": "gaming_preset_1",
  "new_name": "gaming_webcam_only"
}`,

		"delete_scene_preset": `# delete_scene_preset

**Category**: Layout - Scene Presets

**Description**: Permanently remove a scene preset from storage.

**Input**:
- preset_name (string, required): Name of preset to delete

**Output**:
- message: Success confirmation

**Example Input**:
{
  "preset_name": "old_preset"
}

**Warning**: This action cannot be undone.`,

		// Visual - Screenshot Sources
		"create_screenshot_source": `# create_screenshot_source

**Category**: Visual - Screenshot Monitoring

**Description**: Create a periodic screenshot capture source for AI visual monitoring of OBS scenes.

**Input**:
- name (string, required): Unique name for this screenshot source
- source_name (string, required): OBS scene or source name to capture
- cadence_ms (int, optional): Capture interval in milliseconds (default: 5000)
- image_format (string, optional): "png" or "jpg" (default: "png")
- image_width (int, optional): Resize width, 0 = original (default: 0)
- image_height (int, optional): Resize height, 0 = original (default: 0)
- quality (int, optional): Compression quality 0-100 (default: 80)

**Output**:
- id: Screenshot source ID
- name: Source name
- url: HTTP URL for accessing screenshots
- message: Success confirmation

**Example Input**:
{
  "name": "stream_monitor",
  "source_name": "Gaming",
  "cadence_ms": 5000,
  "image_format": "jpg",
  "quality": 85
}

**Use Case**: Enable AI to "see" your stream output for visual verification and issue detection.`,

		"remove_screenshot_source": `# remove_screenshot_source

**Category**: Visual - Screenshot Monitoring

**Description**: Stop and remove a screenshot capture source.

**Input**:
- name (string, required): Name of screenshot source to remove

**Output**:
- message: Success confirmation

**Example Input**:
{
  "name": "stream_monitor"
}

**Note**: This stops capture and deletes all stored screenshots.`,

		"list_screenshot_sources": `# list_screenshot_sources

**Category**: Visual - Screenshot Monitoring

**Description**: List all configured screenshot sources with their status and HTTP URLs.

**Input**: None

**Output**:
- sources: Array of screenshot source objects
- count: Total number of sources

**Example Output**:
{
  "sources": [
    {
      "id": 1,
      "name": "stream_monitor",
      "source_name": "Gaming",
      "cadence_ms": 5000,
      "enabled": true,
      "url": "http://localhost:8765/screenshot/stream_monitor",
      "screenshot_count": 42
    }
  ],
  "count": 1
}`,

		"configure_screenshot_cadence": `# configure_screenshot_cadence

**Category**: Visual - Screenshot Monitoring

**Description**: Update the capture interval for a screenshot source.

**Input**:
- name (string, required): Name of screenshot source
- cadence_ms (int, required): New capture interval in milliseconds

**Output**:
- name: Screenshot source name
- cadence_ms: New cadence value
- message: Success confirmation

**Example Input**:
{
  "name": "stream_monitor",
  "cadence_ms": 10000
}

**Use Case**: Adjust monitoring frequency based on performance or monitoring needs.`,

		// Design - Source Creation
		"create_text_source": `# create_text_source

**Category**: Design - Source Creation

**Description**: Create a text/label source in a scene with customizable font and color.

**Input**:
- scene_name (string, required): Name of scene to add source to
- source_name (string, required): Name for the new text source
- text (string, required): Text content to display
- font_name (string, optional): Font face name (default: "Arial")
- font_size (int, optional): Font size in points (default: 36)
- color (int, optional): Text color as ABGR integer (default: white)

**Output**:
- scene_name: Scene name
- source_name: Source name
- scene_item_id: Scene item ID for positioning
- message: Success confirmation

**Example Input**:
{
  "scene_name": "Gaming",
  "source_name": "Stream Title",
  "text": "Welcome to my stream!",
  "font_size": 48,
  "color": 4294967295
}

**Color Format**: ABGR format where 0xFFFFFFFF = white, 0xFF0000FF = red, 0xFF00FF00 = green, 0xFFFF0000 = blue`,

		"create_image_source": `# create_image_source

**Category**: Design - Source Creation

**Description**: Create an image source in a scene from a file path.

**Input**:
- scene_name (string, required): Name of scene to add source to
- source_name (string, required): Name for the new image source
- file_path (string, required): Path to the image file

**Output**:
- scene_name: Scene name
- source_name: Source name
- scene_item_id: Scene item ID for positioning
- file_path: Image file path
- message: Success confirmation

**Example Input**:
{
  "scene_name": "Gaming",
  "source_name": "Logo",
  "file_path": "C:/Images/logo.png"
}

**Supported Formats**: PNG, JPG, BMP, GIF, etc.`,

		"create_color_source": `# create_color_source

**Category**: Design - Source Creation

**Description**: Create a solid color source in a scene (useful for backgrounds).

**Input**:
- scene_name (string, required): Name of scene to add source to
- source_name (string, required): Name for the new color source
- color (int, required): Color as ABGR integer
- width (int, optional): Width in pixels (default: 1920)
- height (int, optional): Height in pixels (default: 1080)

**Output**:
- scene_name: Scene name
- source_name: Source name
- scene_item_id: Scene item ID for positioning
- width: Width in pixels
- height: Height in pixels
- message: Success confirmation

**Example Input**:
{
  "scene_name": "Gaming",
  "source_name": "Background",
  "color": 4278190080,
  "width": 1920,
  "height": 1080
}

**Color Examples**: 0xFF000000 = black, 0xFFFFFFFF = white, 0xFF1a1a1a = dark gray`,

		"create_browser_source": `# create_browser_source

**Category**: Design - Source Creation

**Description**: Create a browser source in a scene to display web content.

**Input**:
- scene_name (string, required): Name of scene to add source to
- source_name (string, required): Name for the new browser source
- url (string, required): URL to load in the browser source
- width (int, optional): Browser width in pixels (default: 800)
- height (int, optional): Browser height in pixels (default: 600)
- fps (int, optional): Frame rate (default: 30)

**Output**:
- scene_name: Scene name
- source_name: Source name
- scene_item_id: Scene item ID for positioning
- url: Browser URL
- message: Success confirmation

**Example Input**:
{
  "scene_name": "Gaming",
  "source_name": "Chat Overlay",
  "url": "https://example.com/chat-overlay",
  "width": 400,
  "height": 600,
  "fps": 30
}

**Use Cases**: Stream overlays, alerts, chat widgets, web dashboards`,

		"create_media_source": `# create_media_source

**Category**: Design - Source Creation

**Description**: Create a media/video source in a scene from a file path.

**Input**:
- scene_name (string, required): Name of scene to add source to
- source_name (string, required): Name for the new media source
- file_path (string, required): Path to the media file
- loop (bool, optional): Whether to loop the media (default: false)

**Output**:
- scene_name: Scene name
- source_name: Source name
- scene_item_id: Scene item ID for positioning
- file_path: Media file path
- loop: Loop setting
- message: Success confirmation

**Example Input**:
{
  "scene_name": "Intro",
  "source_name": "Intro Video",
  "file_path": "C:/Videos/intro.mp4",
  "loop": false
}

**Supported Formats**: MP4, MOV, AVI, MKV, WebM, etc.`,

		// Design - Layout Control
		"set_source_transform": `# set_source_transform

**Category**: Design - Layout Control

**Description**: Set position, scale, and rotation of a source in a scene.

**Input**:
- scene_name (string, required): Name of scene containing source
- scene_item_id (int, required): Scene item ID of the source
- x (float, optional): X position in pixels
- y (float, optional): Y position in pixels
- scale_x (float, optional): X scale factor (1.0 = 100%)
- scale_y (float, optional): Y scale factor (1.0 = 100%)
- rotation (float, optional): Rotation in degrees

**Output**:
- scene_name: Scene name
- scene_item_id: Scene item ID
- x, y, scale_x, scale_y, rotation: Applied values
- message: Success confirmation

**Example Input**:
{
  "scene_name": "Gaming",
  "scene_item_id": 1,
  "x": 1600,
  "y": 900,
  "scale_x": 0.25,
  "scale_y": 0.25,
  "rotation": 0
}

**Note**: Omit parameters you don't want to change. Only provided values are updated.`,

		"get_source_transform": `# get_source_transform

**Category**: Design - Layout Control

**Description**: Get the current transform properties of a source in a scene.

**Input**:
- scene_name (string, required): Name of scene containing source
- scene_item_id (int, required): Scene item ID of the source

**Output**: Transform object with position, scale, rotation, bounds, crop, size

**Example Input**:
{
  "scene_name": "Gaming",
  "scene_item_id": 1
}

**Use Case**: Get current position/scale before adjusting, or to verify layout.`,

		"set_source_crop": `# set_source_crop

**Category**: Design - Layout Control

**Description**: Set crop values for a source in a scene (trim edges).

**Input**:
- scene_name (string, required): Name of scene containing source
- scene_item_id (int, required): Scene item ID of the source
- crop_top (int, optional): Pixels to crop from top (default: 0)
- crop_bottom (int, optional): Pixels to crop from bottom (default: 0)
- crop_left (int, optional): Pixels to crop from left (default: 0)
- crop_right (int, optional): Pixels to crop from right (default: 0)

**Output**:
- scene_name: Scene name
- scene_item_id: Scene item ID
- crop values: Applied crop values
- message: Success confirmation

**Example Input**:
{
  "scene_name": "Gaming",
  "scene_item_id": 1,
  "crop_top": 50,
  "crop_bottom": 50
}`,

		"set_source_bounds": `# set_source_bounds

**Category**: Design - Layout Control

**Description**: Set bounds type and size for a source (controls scaling behavior).

**Input**:
- scene_name (string, required): Name of scene containing source
- scene_item_id (int, required): Scene item ID of the source
- bounds_type (string, required): Bounds type (see below)
- bounds_width (float, optional): Bounds width in pixels
- bounds_height (float, optional): Bounds height in pixels

**Bounds Types**:
- OBS_BOUNDS_NONE: No bounds
- OBS_BOUNDS_STRETCH: Stretch to bounds
- OBS_BOUNDS_SCALE_INNER: Scale to fit inside bounds
- OBS_BOUNDS_SCALE_OUTER: Scale to cover bounds
- OBS_BOUNDS_SCALE_TO_WIDTH: Scale to width
- OBS_BOUNDS_SCALE_TO_HEIGHT: Scale to height
- OBS_BOUNDS_MAX_ONLY: Scale down only

**Example Input**:
{
  "scene_name": "Gaming",
  "scene_item_id": 1,
  "bounds_type": "OBS_BOUNDS_SCALE_INNER",
  "bounds_width": 1920,
  "bounds_height": 1080
}`,

		"set_source_order": `# set_source_order

**Category**: Design - Layout Control

**Description**: Set the z-order index of a source (layering, 0 = back, higher = front).

**Input**:
- scene_name (string, required): Name of scene containing source
- scene_item_id (int, required): Scene item ID of the source
- index (int, required): New index position (0 = bottom layer)

**Output**:
- scene_name: Scene name
- scene_item_id: Scene item ID
- index: New index
- message: Success confirmation

**Example Input**:
{
  "scene_name": "Gaming",
  "scene_item_id": 1,
  "index": 5
}

**Use Case**: Control which sources appear in front/back (e.g., background at 0, overlays at higher indices).`,

		// Design - Advanced
		"set_source_locked": `# set_source_locked

**Category**: Design - Advanced

**Description**: Lock or unlock a source to prevent/allow accidental changes.

**Input**:
- scene_name (string, required): Name of scene containing source
- scene_item_id (int, required): Scene item ID of the source
- locked (bool, required): Whether source should be locked

**Output**:
- scene_name: Scene name
- scene_item_id: Scene item ID
- locked: Lock state
- message: Success confirmation

**Example Input**:
{
  "scene_name": "Gaming",
  "scene_item_id": 1,
  "locked": true
}

**Use Case**: Lock background elements to prevent accidental movement during stream.`,

		"duplicate_source": `# duplicate_source

**Category**: Design - Advanced

**Description**: Duplicate a source within the same scene or to another scene.

**Input**:
- scene_name (string, required): Name of source scene
- scene_item_id (int, required): Scene item ID to duplicate
- dest_scene_name (string, optional): Destination scene (default: same scene)

**Output**:
- source_scene: Source scene name
- source_item_id: Source scene item ID
- dest_scene: Destination scene name
- new_scene_item_id: New scene item ID
- message: Success confirmation

**Example Input**:
{
  "scene_name": "Gaming",
  "scene_item_id": 1,
  "dest_scene_name": "BRB"
}

**Use Case**: Copy configured sources between scenes without recreating.`,

		"remove_source": `# remove_source

**Category**: Design - Advanced

**Description**: Remove a source from a scene (deletes scene item, not the source itself).

**Input**:
- scene_name (string, required): Name of scene containing source
- scene_item_id (int, required): Scene item ID to remove

**Output**:
- scene_name: Scene name
- scene_item_id: Scene item ID removed
- message: Success confirmation

**Example Input**:
{
  "scene_name": "Gaming",
  "scene_item_id": 1
}

**Note**: This removes the source from the scene only. The source can still exist in other scenes.`,

		"list_input_kinds": `# list_input_kinds

**Category**: Design - Advanced

**Description**: List all available input source types in OBS (useful for knowing what sources can be created).

**Input**: None

**Output**:
- input_kinds: Array of available source type IDs
- count: Total number of types

**Example Output**:
{
  "input_kinds": [
    "dshow_input",
    "wasapi_input_capture",
    "browser_source",
    "image_source",
    "color_source_v3",
    "text_gdiplus_v3"
  ],
  "count": 50
}

**Use Case**: Discover available source types for your OBS installation.`,
	}

	help, exists := toolHelp[toolName]
	if !exists {
		return "", fmt.Errorf("tool '%s' not found", toolName)
	}

	if verbose {
		help += `

## Related Resources

- topic="tools" - See all tools grouped by category
- topic="workflows" - Multi-tool workflows using this tool
- topic="troubleshooting" - Common issues and solutions

## Example Workflow

Combine this tool with others to accomplish complex tasks. See topic="workflows" for examples.
`
	}

	return help, nil
}
