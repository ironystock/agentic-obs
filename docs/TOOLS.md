# MCP Tool Reference

Comprehensive documentation for all 30 Model Context Protocol (MCP) tools provided by the agentic-obs server.

## Table of Contents

- [Overview](#overview)
- [Scene Management](#scene-management)
  - [list_scenes](#list_scenes)
  - [set_current_scene](#set_current_scene)
  - [create_scene](#create_scene)
  - [remove_scene](#remove_scene)
- [Scene Presets](#scene-presets)
  - [save_scene_preset](#save_scene_preset)
  - [list_scene_presets](#list_scene_presets)
  - [get_preset_details](#get_preset_details)
  - [apply_scene_preset](#apply_scene_preset)
  - [rename_scene_preset](#rename_scene_preset)
  - [delete_scene_preset](#delete_scene_preset)
- [Recording](#recording)
  - [start_recording](#start_recording)
  - [stop_recording](#stop_recording)
  - [pause_recording](#pause_recording)
  - [resume_recording](#resume_recording)
  - [get_recording_status](#get_recording_status)
- [Streaming](#streaming)
  - [start_streaming](#start_streaming)
  - [stop_streaming](#stop_streaming)
  - [get_streaming_status](#get_streaming_status)
- [Sources](#sources)
  - [list_sources](#list_sources)
  - [toggle_source_visibility](#toggle_source_visibility)
  - [get_source_settings](#get_source_settings)
- [Audio](#audio)
  - [get_input_mute](#get_input_mute)
  - [toggle_input_mute](#toggle_input_mute)
  - [set_input_volume](#set_input_volume)
  - [get_input_volume](#get_input_volume)
- [Screenshot Sources](#screenshot-sources)
  - [create_screenshot_source](#create_screenshot_source)
  - [remove_screenshot_source](#remove_screenshot_source)
  - [list_screenshot_sources](#list_screenshot_sources)
  - [configure_screenshot_cadence](#configure_screenshot_cadence)
- [Status](#status)
  - [get_obs_status](#get_obs_status)
- [Common Patterns](#common-patterns)
- [Error Handling](#error-handling)

---

## Overview

The agentic-obs MCP server provides 30 tools organized into 8 categories for comprehensive OBS Studio control. All tools communicate with OBS via WebSocket (default port 4455) and return structured JSON responses.

| Category | Tools | Description |
|----------|-------|-------------|
| Scene Management | 4 | List, switch, create, remove scenes |
| Scene Presets | 6 | Save and restore source visibility configurations |
| Recording | 5 | Start, stop, pause, resume, status |
| Streaming | 3 | Start, stop, status |
| Sources | 3 | List, toggle visibility, get settings |
| Audio | 4 | Mute, volume control |
| Screenshot Sources | 4 | AI visual monitoring of stream output |
| Status | 1 | Overall OBS status |

**General Prerequisites:**
- OBS Studio 28+ running with WebSocket server enabled
- agentic-obs MCP server connected to OBS
- Proper authentication configured (if password-protected)

---

## Scene Management

### list_scenes

**Purpose:** Retrieve all available scenes in OBS and identify which scene is currently active.

**Parameters:** None

**Return Value Schema:**
```json
{
  "scenes": ["Scene 1", "Scene 2", "Gaming", "Chatting"],
  "current_scene": "Gaming"
}
```

**Return Fields:**
- `scenes` (array of strings): List of all scene names in OBS
- `current_scene` (string): Name of the currently active scene

**Use Cases:**
- Display available scenes to users before switching
- Verify a scene exists before attempting to activate it
- Build scene selection menus in AI-driven interfaces
- Audit scene configuration across multiple OBS instances
- Create scene-based workflows and automation

**Example Natural Language Prompts:**
- "What scenes do I have in OBS?"
- "Show me all available scenes and which one is active"
- "List my OBS scenes"
- "Which scene am I currently on?"
- "What's my current scene setup?"

**Error Scenarios:**
- OBS not connected: "failed to list scenes: OBS client not connected"
- WebSocket error: Connection timeout or authentication failure
- No scenes exist: Returns empty array (unlikely, OBS always has at least one scene)

**Prerequisites:**
- Active OBS WebSocket connection

**Related Tools:**
- `set_current_scene` - Switch to a different scene
- `create_scene` - Add new scenes
- `remove_scene` - Delete existing scenes
- `get_obs_status` - Includes current scene in overall status

**Best Practices:**
- Call this before `set_current_scene` to validate scene names
- Cache results for short periods to reduce OBS queries
- Use current_scene to avoid redundant scene switches
- Consider scenes as MCP resources for real-time updates via notifications

---

### set_current_scene

**Purpose:** Switch the active program scene in OBS to a different scene.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| scene_name | string | Yes | Exact name of the scene to activate |

**Return Value Schema:**
```json
{
  "message": "Successfully switched to scene: Gaming"
}
```

**Use Cases:**
- Automate scene transitions during streams (e.g., "Switch to BRB scene")
- Respond to events (chat commands, timers, external triggers)
- Create multi-scene workflows for presentations or tutorials
- Quick scene changes via voice commands
- Coordinate scene changes with other OBS operations

**Example Natural Language Prompts:**
- "Switch to my Gaming scene"
- "Change to the BRB scene"
- "Go to the Chatting layout"
- "Activate my intro scene"
- "Switch OBS to the main camera view"

**Error Scenarios:**
- Scene doesn't exist: "failed to set current scene to 'InvalidName': Scene may not exist"
- Invalid scene name: Same as above
- OBS not connected: "OBS client not connected"
- Empty scene name: OBS WebSocket will reject the request

**Prerequisites:**
- Scene must exist in OBS (use `list_scenes` to verify)
- Active OBS connection

**Related Tools:**
- `list_scenes` - Verify scene exists before switching
- `create_scene` - Create scene if it doesn't exist
- `get_obs_status` - Verify scene change was successful

**Best Practices:**
- Always verify scene exists using `list_scenes` first
- Scene names are case-sensitive - match exactly
- Consider transition duration when coordinating with other actions
- Use with recording/streaming controls for automated workflows
- Combine with source visibility toggles for complex scene states

---

### create_scene

**Purpose:** Create a new empty scene in OBS.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| scene_name | string | Yes | Name for the new scene (must be unique) |

**Return Value Schema:**
```json
{
  "message": "Successfully created scene: NewScene"
}
```

**Use Cases:**
- Dynamically create scenes for different content types
- Build custom scene layouts programmatically
- Automated setup for recurring events or shows
- Template-based scene generation
- User-requested custom scenes via chat commands

**Example Natural Language Prompts:**
- "Create a new scene called 'Tutorial Mode'"
- "Make a scene named 'Guest Interview'"
- "Add a new OBS scene for the Q&A session"
- "Set up a scene called 'Screen Share'"
- "Create a blank scene called 'Backup'"

**Error Scenarios:**
- Scene already exists: "failed to create scene 'Gaming': Scene may already exist"
- Empty scene name: OBS WebSocket validation error
- Special characters: May be rejected depending on OBS version
- Maximum scenes reached: OBS resource limitation (rare)

**Prerequisites:**
- Active OBS connection
- Unique scene name not currently in use

**Related Tools:**
- `list_scenes` - Verify scene was created
- `remove_scene` - Delete scenes no longer needed
- `set_current_scene` - Switch to newly created scene
- MCP resources - New scenes trigger `resources/list_changed` notification

**Best Practices:**
- Check if scene exists first using `list_scenes` to avoid errors
- Use descriptive, unique names
- Newly created scenes are empty - add sources separately
- Consider scene naming conventions for organization
- Clean up unused scenes periodically with `remove_scene`

---

### remove_scene

**Purpose:** Delete an existing scene from OBS.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| scene_name | string | Yes | Exact name of the scene to remove |

**Return Value Schema:**
```json
{
  "message": "Successfully removed scene: OldScene"
}
```

**Use Cases:**
- Clean up temporary or event-specific scenes
- Remove outdated scene configurations
- Automated scene lifecycle management
- Declutter OBS workspace programmatically
- Delete scenes created for one-time events

**Example Natural Language Prompts:**
- "Delete the 'Old Gaming' scene"
- "Remove the BRB scene"
- "Get rid of the test scene"
- "Delete the scene called 'Backup'"
- "Remove all scenes with 'temp' in the name" (requires iteration)

**Error Scenarios:**
- Scene doesn't exist: "failed to remove scene 'InvalidName': Scene may not exist"
- Last remaining scene: "Scene may not exist or may be the only scene" (OBS requires at least one scene)
- Currently active scene: Some OBS versions prevent deletion of active scene
- Scene in use: May fail if scene is referenced elsewhere

**Prerequisites:**
- Scene must exist in OBS
- OBS must have at least one other scene (cannot delete the only scene)
- Scene should not be currently active (switch first if needed)

**Related Tools:**
- `list_scenes` - Verify scene exists before attempting removal
- `set_current_scene` - Switch away from scene before deleting
- `create_scene` - Create scenes to replace deleted ones
- MCP resources - Deletion triggers `resources/list_changed` notification

**Best Practices:**
- Always check if scene exists first
- Switch to a different scene before deleting the current one
- Confirm with user before deleting scenes (destructive operation)
- Cannot undo scene deletion - all sources in scene are lost
- Use scene naming conventions to identify temporary scenes

---

## Scene Presets

Scene presets allow you to save and restore the visibility state of all sources within a scene. This is useful for quickly switching between different "looks" or configurations without creating separate scenes.

### save_scene_preset

**Purpose:** Save the current visibility state of all sources in a scene as a named preset.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| name | string | Yes | Unique name for the preset |
| scene_name | string | Yes | Name of the scene to capture |
| description | string | No | Optional description of the preset |

**Return Value Schema:**
```json
{
  "id": 1,
  "name": "gaming-webcam-on",
  "scene_name": "Gaming",
  "message": "Preset saved successfully"
}
```

**Use Cases:**
- Save different overlay configurations (e.g., "with chat", "without chat")
- Store webcam on/off states for quick toggling
- Create content-specific source configurations
- Save pre-stream and during-stream layouts

**Example Natural Language Prompts:**
- "Save my current Gaming scene layout as 'gaming-full-overlay'"
- "Create a preset called 'minimal-view' from my current scene"
- "Save this configuration as 'webcam-only' with description 'Just webcam, no overlays'"
- "Capture my current source visibility as a preset named 'sponsor-mode'"

---

### list_scene_presets

**Purpose:** List all saved scene presets, optionally filtered by scene name.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| scene_name | string | No | Filter presets by scene name |

**Return Value Schema:**
```json
{
  "presets": [
    {
      "id": 1,
      "name": "gaming-webcam-on",
      "scene_name": "Gaming",
      "description": "Full overlay with webcam",
      "created_at": "2025-01-15T10:30:00Z"
    }
  ],
  "count": 1
}
```

**Example Natural Language Prompts:**
- "What presets do I have saved?"
- "Show me all my scene presets"
- "List presets for my Gaming scene"
- "What configurations have I saved?"

---

### get_preset_details

**Purpose:** Get detailed information about a specific preset, including all source visibility states.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| name | string | Yes | Name of the preset |

**Return Value Schema:**
```json
{
  "id": 1,
  "name": "gaming-webcam-on",
  "scene_name": "Gaming",
  "description": "Full overlay with webcam",
  "sources": [
    {"name": "Webcam", "visible": true},
    {"name": "Chat Overlay", "visible": true},
    {"name": "Alert Box", "visible": true}
  ],
  "created_at": "2025-01-15T10:30:00Z"
}
```

**Example Natural Language Prompts:**
- "Show me the details of my 'gaming-webcam-on' preset"
- "What sources are saved in the 'minimal-view' preset?"
- "Describe the 'sponsor-mode' preset configuration"

---

### apply_scene_preset

**Purpose:** Apply a saved preset, restoring the source visibility states to the target scene.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| name | string | Yes | Name of the preset to apply |

**Return Value Schema:**
```json
{
  "message": "Preset 'gaming-webcam-on' applied successfully",
  "changes": 3
}
```

**Use Cases:**
- Quickly switch between overlay configurations during stream
- Restore a known-good layout
- Toggle between different content modes
- Apply sponsor overlays on command

**Example Natural Language Prompts:**
- "Apply my 'minimal-view' preset"
- "Switch to the 'gaming-webcam-on' configuration"
- "Restore my 'sponsor-mode' layout"
- "Use the 'interview-setup' preset"

---

### rename_scene_preset

**Purpose:** Rename an existing preset.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| old_name | string | Yes | Current name of the preset |
| new_name | string | Yes | New name for the preset |

**Return Value Schema:**
```json
{
  "message": "Preset renamed from 'old-name' to 'new-name'"
}
```

**Example Natural Language Prompts:**
- "Rename my 'test-preset' to 'production-layout'"
- "Change the name of 'gaming1' to 'gaming-with-cam'"

---

### delete_scene_preset

**Purpose:** Delete a saved preset.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| name | string | Yes | Name of the preset to delete |

**Return Value Schema:**
```json
{
  "message": "Preset 'old-preset' deleted successfully"
}
```

**Example Natural Language Prompts:**
- "Delete the 'test-config' preset"
- "Remove my old 'broken-layout' preset"
- "Get rid of the 'temp-setup' configuration"

---

## Recording

### start_recording

**Purpose:** Begin recording the OBS output to a file.

**Parameters:** None

**Return Value Schema:**
```json
{
  "message": "Successfully started recording"
}
```

**Use Cases:**
- Start recording sessions on command or schedule
- Automated recording for meetings, streams, or events
- Voice-activated recording controls
- Event-triggered recording (e.g., game starts, meeting begins)
- Programmatic recording workflows

**Example Natural Language Prompts:**
- "Start recording"
- "Begin recording my stream"
- "Start capturing OBS output"
- "Hit record in OBS"
- "Start saving the video"

**Error Scenarios:**
- Recording already active: "failed to start recording: output already active"
- Invalid output path: "Check OBS recording settings and output path"
- Disk full: OBS may fail to start recording
- Codec/encoder not available: OBS configuration issue
- No recording output configured: OBS settings incomplete

**Prerequisites:**
- OBS recording settings configured (output path, format, encoder)
- Sufficient disk space at output location
- Recording output must not already be active
- Valid encoder available

**Related Tools:**
- `stop_recording` - End the recording session
- `pause_recording` - Temporarily pause recording
- `get_recording_status` - Check if recording is active
- `get_obs_status` - Overall status including recording state

**Best Practices:**
- Check recording status first to avoid errors
- Verify disk space before long recording sessions
- Configure output path in OBS settings beforehand
- Use with scene management for automated recording workflows
- Consider pausing instead of stopping for temporary interruptions

---

### stop_recording

**Purpose:** Stop the current recording session and save the output file.

**Parameters:** None

**Return Value Schema:**
```json
{
  "message": "Successfully stopped recording. Output saved to: C:/Videos/2024-01-15_12-30-45.mp4"
}
```

**Return Fields:**
- `message` (string): Success message including the full path to the saved recording file

**Use Cases:**
- End recording sessions after completion
- Stop recording on command or schedule
- Save recording when specific events occur
- Automated recording session management
- Emergency stop for unwanted recordings

**Example Natural Language Prompts:**
- "Stop recording"
- "End the recording"
- "Finish recording and save the file"
- "Stop capturing"
- "Save and stop the current recording"

**Error Scenarios:**
- Recording not active: "failed to stop recording: Recording may not be active"
- Disk write error: File system issues during save
- OBS not connected: WebSocket connection lost
- Recording already stopped: Redundant stop command

**Prerequisites:**
- Recording must be currently active
- Write permissions to output directory
- Sufficient disk space to finalize file

**Related Tools:**
- `start_recording` - Begin a recording session
- `get_recording_status` - Verify recording is active before stopping
- `pause_recording` - Pause instead of stopping
- `resume_recording` - Resume a paused recording

**Best Practices:**
- Check recording status before attempting to stop
- Save the output path for reference or processing
- Allow a few seconds for file finalization
- Verify file was saved successfully after stopping
- Use pausing for temporary interruptions instead of stop/start

---

### pause_recording

**Purpose:** Temporarily pause an active recording without stopping it.

**Parameters:** None

**Return Value Schema:**
```json
{
  "message": "Successfully paused recording"
}
```

**Use Cases:**
- Pause during breaks or interruptions
- Exclude unwanted content from recordings
- Temporary stop during technical difficulties
- Privacy protection during sensitive moments
- Conserve disk space during idle periods

**Example Natural Language Prompts:**
- "Pause the recording"
- "Hold the recording for a moment"
- "Pause recording temporarily"
- "Stop recording but don't end it"
- "Put the recording on hold"

**Error Scenarios:**
- Recording not active: "failed to pause recording: Recording may not be active"
- Recording already paused: May succeed silently or return error depending on OBS version
- OBS not connected: WebSocket connection issue

**Prerequisites:**
- Recording must be currently active (not paused)
- OBS version with pause support (OBS 27+)

**Related Tools:**
- `resume_recording` - Continue recording after pause
- `get_recording_status` - Check if recording is paused
- `stop_recording` - Completely end recording
- `start_recording` - Begin new recording

**Best Practices:**
- Use pause instead of stop/start to maintain single file
- Check status to verify pause was successful
- Resume promptly to avoid confusion
- Paused recordings don't capture any content
- Use `get_recording_status` to track paused state

---

### resume_recording

**Purpose:** Resume a paused recording session.

**Parameters:** None

**Return Value Schema:**
```json
{
  "message": "Successfully resumed recording"
}
```

**Use Cases:**
- Continue recording after breaks
- Resume after addressing interruptions
- Restart capture after privacy pause
- Automated resume after scheduled pause
- Continue recording after technical fixes

**Example Natural Language Prompts:**
- "Resume recording"
- "Continue the recording"
- "Unpause the recording"
- "Keep recording"
- "Start recording again"

**Error Scenarios:**
- Recording not paused: "failed to resume recording: Recording may not be paused"
- Recording not active: Must be both active and paused
- OBS not connected: WebSocket connection issue

**Prerequisites:**
- Recording must be active and paused (not stopped)
- Same recording session must still be valid

**Related Tools:**
- `pause_recording` - Pause an active recording
- `get_recording_status` - Verify recording is paused before resuming
- `start_recording` - Start new recording if stopped
- `stop_recording` - End recording session

**Best Practices:**
- Verify recording is paused before attempting resume
- Check status after resume to confirm success
- Pause/resume maintains single output file
- Announce resume to viewers if streaming
- Use status monitoring to track recording state

---

### get_recording_status

**Purpose:** Retrieve detailed information about the current recording state.

**Parameters:** None

**Return Value Schema:**
```json
{
  "active": true,
  "paused": false,
  "timecode": "00:15:32",
  "output_bytes": 1048576000
}
```

**Return Fields:**
- `active` (boolean): Whether recording is currently active
- `paused` (boolean): Whether active recording is paused
- `timecode` (string): Current recording duration (HH:MM:SS format)
- `output_bytes` (integer): Current file size in bytes

**Use Cases:**
- Monitor recording duration and file size
- Verify recording state before operations
- Display recording status to users
- Automated monitoring and alerts
- Check if recording is paused vs stopped

**Example Natural Language Prompts:**
- "Is OBS recording?"
- "Check the recording status"
- "How long have I been recording?"
- "What's the current recording duration?"
- "Show me the recording details"

**Error Scenarios:**
- OBS not connected: "failed to get recording status: OBS client not connected"
- WebSocket error: Connection timeout or authentication failure

**Prerequisites:**
- Active OBS connection (recording doesn't need to be active)

**Related Tools:**
- `start_recording` - Start recording if not active
- `stop_recording` - Stop active recording
- `pause_recording` - Pause if active and not paused
- `resume_recording` - Resume if paused
- `get_obs_status` - Overall OBS status including recording

**Best Practices:**
- Check status before start/stop/pause/resume operations
- Monitor timecode for long recording sessions
- Use output_bytes to estimate disk space usage
- Poll periodically for status dashboard
- Combine active and paused flags for complete state

---

## Streaming

### start_streaming

**Purpose:** Begin streaming to the configured streaming service.

**Parameters:** None

**Return Value Schema:**
```json
{
  "message": "Successfully started streaming"
}
```

**Use Cases:**
- Start scheduled streams automatically
- Voice-activated stream start
- Event-triggered streaming (calendar, external API)
- Automated go-live workflows
- Remote stream control

**Example Natural Language Prompts:**
- "Start streaming"
- "Go live"
- "Begin the stream"
- "Start broadcasting"
- "Turn on the stream"

**Error Scenarios:**
- Stream already active: "output already active"
- Invalid stream settings: "Check OBS stream settings and credentials"
- Authentication failure: Stream key or credentials invalid
- Network unavailable: Cannot connect to streaming service
- Stream service unavailable: Platform-side issues

**Prerequisites:**
- OBS stream settings configured (server, stream key)
- Valid streaming service credentials
- Network connectivity to streaming service
- Stream output must not already be active

**Related Tools:**
- `stop_streaming` - End streaming session
- `get_streaming_status` - Check stream state and stats
- `get_obs_status` - Overall status including streaming state
- `set_current_scene` - Change scene before/during stream

**Best Practices:**
- Verify stream settings before first use
- Check network connectivity before starting
- Test stream key validity beforehand
- Use `get_streaming_status` to verify stream started
- Coordinate with scene management for smooth go-live
- Consider audio/video checks before going live

---

### stop_streaming

**Purpose:** Stop the current streaming session and disconnect from the streaming service.

**Parameters:** None

**Return Value Schema:**
```json
{
  "message": "Successfully stopped streaming"
}
```

**Use Cases:**
- End scheduled streams
- Emergency stream stop
- Automated stream ending
- Remote stream shutdown
- Disconnect from streaming service

**Example Natural Language Prompts:**
- "Stop streaming"
- "End the stream"
- "Go offline"
- "Stop broadcasting"
- "Disconnect the stream"

**Error Scenarios:**
- Stream not active: "failed to stop streaming: Stream may not be active"
- OBS not connected: WebSocket connection issue
- Stream already stopped: Redundant stop command

**Prerequisites:**
- Stream must be currently active
- Active OBS connection

**Related Tools:**
- `start_streaming` - Begin streaming session
- `get_streaming_status` - Verify stream is active before stopping
- `get_obs_status` - Check overall streaming state

**Best Practices:**
- Check streaming status before attempting to stop
- Allow graceful disconnection time
- Coordinate with scene management for end screens
- Verify stream stopped successfully
- Monitor chat/alerts after stopping
- Use proper ending scenes before stopping

---

### get_streaming_status

**Purpose:** Retrieve detailed information about the current streaming state and statistics.

**Parameters:** None

**Return Value Schema:**
```json
{
  "active": true,
  "reconnecting": false,
  "timecode": "01:45:22",
  "total_bytes": 5368709120,
  "total_frames": 189000
}
```

**Return Fields:**
- `active` (boolean): Whether streaming is currently active
- `reconnecting` (boolean): Whether stream is attempting to reconnect
- `timecode` (string): Current stream duration (HH:MM:SS format)
- `total_bytes` (integer): Total bytes sent to streaming service
- `total_frames` (integer): Total frames streamed (OBS 28+, may be 0 on older versions)

**Use Cases:**
- Monitor stream health and duration
- Detect connection issues via reconnecting flag
- Display stream stats to viewers
- Automated stream monitoring and alerts
- Calculate bitrate and performance metrics

**Example Natural Language Prompts:**
- "Is OBS streaming?"
- "Check the stream status"
- "How long have I been live?"
- "Show me stream statistics"
- "Are we still connected to the stream?"

**Error Scenarios:**
- OBS not connected: "failed to get streaming status: OBS client not connected"
- WebSocket error: Connection or authentication issue

**Prerequisites:**
- Active OBS connection (streaming doesn't need to be active)

**Related Tools:**
- `start_streaming` - Start stream if not active
- `stop_streaming` - Stop active stream
- `get_obs_status` - Overall OBS status including streaming
- `get_recording_status` - Similar status for recording

**Best Practices:**
- Monitor reconnecting flag for connection issues
- Poll periodically during streams for health monitoring
- Use total_bytes to estimate bandwidth usage
- Calculate average bitrate: (total_bytes * 8) / (timecode_seconds)
- Alert on reconnecting state for stream quality issues
- Combine with recording status for multi-output monitoring

---

## Sources

### list_sources

**Purpose:** Retrieve all input sources (audio and video) available in OBS.

**Parameters:** None

**Return Value Schema:**
```json
[
  {
    "inputName": "Desktop Audio",
    "inputKind": "wasapi_output_capture",
    "unversionedInputKind": "wasapi_output_capture"
  },
  {
    "inputName": "Microphone",
    "inputKind": "wasapi_input_capture",
    "unversionedInputKind": "wasapi_input_capture"
  },
  {
    "inputName": "Webcam",
    "inputKind": "dshow_input",
    "unversionedInputKind": "dshow_input"
  }
]
```

**Return Fields (per source):**
- `inputName` (string): Display name of the source
- `inputKind` (string): OBS source type identifier
- `unversionedInputKind` (string): Base source type without version

**Use Cases:**
- Discover available audio/video sources
- Build source selection interfaces
- Validate source names before operations
- Audit source configuration
- Identify source types for appropriate handling

**Example Natural Language Prompts:**
- "What sources do I have in OBS?"
- "List all my audio and video inputs"
- "Show me available sources"
- "What cameras and microphones are configured?"
- "Display all OBS inputs"

**Error Scenarios:**
- OBS not connected: "failed to list sources from OBS: OBS client not connected"
- WebSocket error: Connection or authentication issue
- No sources configured: Returns empty array

**Prerequisites:**
- Active OBS connection
- Sources configured in OBS (not required, will return empty array)

**Related Tools:**
- `get_source_settings` - Get detailed configuration for a source
- `toggle_source_visibility` - Show/hide source in scene
- `get_input_mute` - Check audio source mute state
- `set_input_volume` - Adjust audio source volume

**Best Practices:**
- Use to validate source names before operations
- Cache source list to reduce OBS queries
- Filter by inputKind to find specific source types
- Sources are global - can be added to multiple scenes
- Note: This lists input sources, not scene items
- Use scene resource queries for per-scene source instances

---

### toggle_source_visibility

**Purpose:** Toggle the visibility/enabled state of a source within a specific scene.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| scene_name | string | Yes | Name of the scene containing the source |
| source_id | integer | Yes | Scene item ID of the source (not source name) |

**Return Value Schema:**
```json
{
  "scene_name": "Gaming",
  "source_id": 5,
  "visible": true
}
```

**Return Fields:**
- `scene_name` (string): Name of the scene
- `source_id` (integer): Scene item ID that was toggled
- `visible` (boolean): New visibility state after toggle

**Use Cases:**
- Show/hide overlays during streams
- Toggle webcam visibility on command
- Reveal/hide alerts or notifications
- Create interactive scene elements
- Automated source visibility workflows

**Example Natural Language Prompts:**
- "Toggle my webcam in the Gaming scene"
- "Hide the chat overlay"
- "Show source ID 5 in the current scene"
- "Toggle visibility of the alert box"
- "Switch the webcam on/off"

**Error Scenarios:**
- Scene doesn't exist: "failed to get visibility state for source X in scene 'InvalidScene'"
- Invalid source ID: "failed to get visibility state for source 999"
- Source not in scene: Source ID doesn't exist in specified scene
- OBS not connected: WebSocket connection issue

**Prerequisites:**
- Scene must exist in OBS
- Source ID must be valid for that scene
- Active OBS connection

**Related Tools:**
- `list_scenes` - Get available scenes
- Scene resources (MCP) - Get source IDs within scenes
- `list_sources` - List available sources (but not scene item IDs)
- `get_source_settings` - Get source configuration

**Best Practices:**
- Use MCP scene resources to get valid source IDs
- Source ID is per-scene, not global source name
- Same source can have different IDs in different scenes
- Returns new state for confirmation
- Consider using scene resources/read to get source details
- Toggle is atomic - gets current state and flips it

---

### get_source_settings

**Purpose:** Retrieve the configuration settings for a specific input source.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| source_name | string | Yes | Exact name of the source |

**Return Value Schema:**
```json
{
  "device_id": "default",
  "use_device_timing": false,
  "sample_rate": 48000
}
```
(Schema varies by source type - example shows audio input)

**Return Fields:**
- Varies by source type
- Common fields: device settings, paths, URLs, credentials
- Returns source-specific configuration as key-value pairs

**Use Cases:**
- Inspect source configuration programmatically
- Verify device settings
- Debug source issues
- Audit source configuration across setups
- Document source settings for backup/restore

**Example Natural Language Prompts:**
- "What are the settings for my Webcam source?"
- "Show me the configuration for Desktop Audio"
- "Get the settings for the 'Game Capture' source"
- "What device is my microphone using?"
- "Display source settings for 'Browser Source'"

**Error Scenarios:**
- Source doesn't exist: "failed to get settings for source 'InvalidName': Source may not exist"
- Source name incorrect: Same as above (case-sensitive)
- OBS not connected: WebSocket connection issue

**Prerequisites:**
- Source must exist in OBS
- Exact source name required
- Active OBS connection

**Related Tools:**
- `list_sources` - Get available source names
- `toggle_source_visibility` - Control source visibility
- `get_input_mute` - Audio source mute state
- `set_input_volume` - Audio source volume

**Best Practices:**
- Use `list_sources` to get exact source names first
- Source names are case-sensitive
- Settings structure varies by source type
- Some settings may contain sensitive data (URLs, credentials)
- Read-only operation, doesn't modify settings
- Future enhancement: set_source_settings for modification

---

## Audio

### get_input_mute

**Purpose:** Check whether a specific audio input source is currently muted.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| input_name | string | Yes | Exact name of the audio input |

**Return Value Schema:**
```json
{
  "input_name": "Microphone",
  "is_muted": false
}
```

**Return Fields:**
- `input_name` (string): Name of the audio input queried
- `is_muted` (boolean): Current mute state (true = muted, false = unmuted)

**Use Cases:**
- Check microphone mute state before speaking
- Display mute status in dashboards
- Verify audio input state before recording/streaming
- Conditional muting workflows
- Audio monitoring and alerts

**Example Natural Language Prompts:**
- "Is my microphone muted?"
- "Check if Desktop Audio is muted"
- "What's the mute status of my mic?"
- "Am I muted?"
- "Is the Microphone input active?"

**Error Scenarios:**
- Input doesn't exist: "failed to get mute state for input 'InvalidName': Input may not exist"
- Not an audio input: May fail for non-audio sources
- Input name incorrect: Case-sensitive name matching
- OBS not connected: WebSocket connection issue

**Prerequisites:**
- Audio input must exist in OBS
- Exact input name required
- Active OBS connection

**Related Tools:**
- `toggle_input_mute` - Toggle mute state
- `list_sources` - Get available audio input names
- `set_input_volume` - Adjust input volume
- `get_source_settings` - Get detailed input configuration

**Best Practices:**
- Use before toggling to know expected outcome
- Audio inputs only (video sources will fail)
- Names are case-sensitive
- Check state before operations to avoid errors
- Useful for conditional logic in workflows
- Combine with volume controls for complete audio management

---

### toggle_input_mute

**Purpose:** Toggle the mute state of an audio input (muted becomes unmuted and vice versa).

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| input_name | string | Yes | Exact name of the audio input |

**Return Value Schema:**
```json
{
  "message": "Successfully toggled mute for input: Microphone"
}
```

**Use Cases:**
- Quick mute/unmute controls for microphones
- Push-to-mute or push-to-talk workflows
- Voice-activated mute controls
- Privacy protection during streams
- Automated audio management

**Example Natural Language Prompts:**
- "Mute my microphone"
- "Toggle mic mute"
- "Unmute Desktop Audio"
- "Switch my microphone mute state"
- "Turn off my mic"

**Error Scenarios:**
- Input doesn't exist: "failed to toggle mute for input 'InvalidName': Input may not exist"
- Not an audio input: May fail for non-audio sources
- Input name incorrect: Case-sensitive name matching
- OBS not connected: WebSocket connection issue

**Prerequisites:**
- Audio input must exist in OBS
- Exact input name required
- Active OBS connection

**Related Tools:**
- `get_input_mute` - Check current mute state
- `list_sources` - Get available audio input names
- `set_input_volume` - Adjust input volume
- `get_source_settings` - Get input configuration

**Best Practices:**
- Use `get_input_mute` first if you need to know final state
- Toggle is atomic - gets current state and flips it
- Works only on audio inputs
- Names are case-sensitive
- Consider user feedback after toggle
- Combine with visual indicators for mute state
- Useful for hotkey-style controls

---

### set_input_volume

**Purpose:** Set the volume level of an audio input using either decibel or multiplier format.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| input_name | string | Yes | Exact name of the audio input |
| volume_db | float | No | Volume in decibels (-100.0 to 26.0) |
| volume_mul | float | No | Volume multiplier (0.0 to 20.0, where 1.0 = 100%) |

**Note:** Provide either `volume_db` OR `volume_mul`, not both. At least one must be provided.

**Return Value Schema:**
```json
{
  "message": "Successfully set volume for input: Microphone"
}
```

**Use Cases:**
- Adjust microphone gain programmatically
- Balance audio levels across inputs
- Normalize audio during streams
- Automated volume adjustments based on content
- Remote audio mixing controls

**Example Natural Language Prompts:**
- "Set microphone volume to -10 dB"
- "Lower Desktop Audio to 50%"
- "Increase mic volume to 1.5x"
- "Set Microphone to 80% volume"
- "Turn down the music by 6 dB"

**Error Scenarios:**
- Input doesn't exist: "failed to set volume for input 'InvalidName': Input may not exist"
- Not an audio input: May fail for non-audio sources
- Neither parameter provided: Validation error
- Invalid range: Values outside allowed ranges
- OBS not connected: WebSocket connection issue

**Prerequisites:**
- Audio input must exist in OBS
- Exact input name required
- Active OBS connection

**Related Tools:**
- `get_input_mute` - Check/toggle mute state
- `toggle_input_mute` - Mute/unmute input
- `list_sources` - Get available audio input names
- `get_source_settings` - Get detailed input configuration

**Best Practices:**
- Use `volume_db` for precise professional audio control
- Use `volume_mul` for percentage-based adjustments (1.0 = 100%)
- 0 dB / 1.0 multiplier = original volume
- Negative dB reduces volume, positive increases
- Volume above 0 dB may cause clipping/distortion
- Test volume changes before live use
- Consider gradual adjustments instead of dramatic changes
- Volume ranges: dB (-100 to 26), multiplier (0 to 20)

---

### get_input_volume

**Purpose:** Retrieve the current volume level of an audio input in both decibel and multiplier formats.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| input_name | string | Yes | Exact name of the audio input |

**Return Value Schema:**
```json
{
  "input_name": "Microphone",
  "volume_db": -6.0,
  "volume_mul": 0.501
}
```

**Return Fields:**
- `input_name` (string): Name of the audio input queried
- `volume_db` (float): Current volume in decibels
- `volume_mul` (float): Current volume as multiplier (1.0 = 100%)

**Use Cases:**
- Check current volume before making adjustments
- Display volume levels in dashboards
- Verify volume changes were applied
- Monitor audio levels programmatically

**Example Natural Language Prompts:**
- "What's the current volume on my microphone?"
- "Check the volume level of Desktop Audio"
- "Show me my mic volume in dB"
- "What's the audio level on my game capture?"

**Related Tools:**
- `set_input_volume` - Adjust volume levels
- `get_input_mute` - Check mute state
- `toggle_input_mute` - Toggle mute

---

## Screenshot Sources

Screenshot sources enable AI assistants to visually observe your OBS output through periodic image capture. This transforms AI from a blind controller into a seeing collaborator that can verify changes, detect problems, and provide layout feedback.

For detailed documentation on use cases and best practices, see [SCREENSHOTS.md](SCREENSHOTS.md).

### create_screenshot_source

**Purpose:** Create a new periodic screenshot capture source for AI visual monitoring.

**Parameters:**
| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| name | string | Yes | - | Unique identifier for this source |
| source_name | string | Yes | - | OBS source or scene name to capture |
| cadence_ms | integer | No | 5000 | Capture interval in milliseconds |
| image_format | string | No | "png" | Image format: "png" or "jpg" |
| image_width | integer | No | 0 | Resize width (0 = original) |
| image_height | integer | No | 0 | Resize height (0 = original) |
| quality | integer | No | 80 | Compression quality (1-100) |

**Return Value Schema:**
```json
{
  "id": 1,
  "name": "stream-monitor",
  "url": "http://localhost:8765/screenshot/stream-monitor",
  "message": "Screenshot source created successfully"
}
```

**Use Cases:**
- Visual verification after scene changes
- Problem detection (black screens, frozen sources, missing overlays)
- Layout and composition feedback
- Stream quality monitoring
- Pre-stream checks

**Example Natural Language Prompts:**
- "Create a screenshot source called 'stream-view' that captures every 5 seconds"
- "Set up visual monitoring of my Gaming scene"
- "I want you to be able to see my stream - set that up"
- "Create a high-quality screenshot capture at 1080p called 'hd-monitor'"
- "Set up fast 2-second JPG captures for debugging called 'quick-check'"

**Best Practices:**
- Use 5-10 second intervals for general monitoring
- Use 1-2 second intervals for active debugging
- Use PNG for quality-critical captures, JPG for frequent monitoring
- Create descriptive names that indicate purpose

---

### remove_screenshot_source

**Purpose:** Stop and delete a screenshot capture source.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| name | string | Yes | Name of the source to remove |

**Return Value Schema:**
```json
{
  "message": "Screenshot source 'stream-monitor' removed successfully"
}
```

**Example Natural Language Prompts:**
- "Remove the screenshot source called 'test-capture'"
- "Stop capturing screenshots for 'old-monitor'"
- "Delete the 'temp-check' screenshot source"
- "Turn off visual monitoring"

---

### list_screenshot_sources

**Purpose:** List all configured screenshot sources with their status and HTTP URLs.

**Parameters:** None

**Return Value Schema:**
```json
{
  "sources": [
    {
      "id": 1,
      "name": "stream-monitor",
      "source_name": "Gaming",
      "cadence_ms": 5000,
      "image_format": "png",
      "enabled": true,
      "url": "http://localhost:8765/screenshot/stream-monitor",
      "created_at": "2025-01-15T10:30:00Z"
    }
  ],
  "count": 1
}
```

**Example Natural Language Prompts:**
- "What screenshot sources do I have set up?"
- "Show me all my visual monitoring configurations"
- "List screenshot URLs for all my sources"
- "What streams am I currently monitoring?"

---

### configure_screenshot_cadence

**Purpose:** Update the capture interval for an existing screenshot source.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| name | string | Yes | Name of the source to update |
| cadence_ms | integer | Yes | New capture interval in milliseconds |

**Return Value Schema:**
```json
{
  "message": "Screenshot cadence updated to 10000ms for source 'stream-monitor'"
}
```

**Example Natural Language Prompts:**
- "Change my stream-monitor to capture every 10 seconds"
- "Speed up the screenshot capture to every 2 seconds"
- "Slow down 'quick-check' to 30-second intervals"
- "Update the capture rate for 'hd-monitor' to 5 seconds"

---

## Status

### get_obs_status

**Purpose:** Retrieve comprehensive OBS status information including version, connection state, performance metrics, and current operational status.

**Parameters:** None

**Return Value Schema:**
```json
{
  "version": "30.0.2",
  "websocket_version": "5.5.6",
  "platform": "windows",
  "current_scene": "Gaming",
  "recording": false,
  "streaming": true,
  "fps": 59.94,
  "frame_time_ms": 16.68,
  "frames": 324000,
  "dropped_frames": 42
}
```

**Return Fields:**
- `version` (string): OBS Studio version
- `websocket_version` (string): OBS WebSocket plugin version
- `platform` (string): Operating system (windows/macos/linux)
- `current_scene` (string): Currently active scene name
- `recording` (boolean): Whether recording is active
- `streaming` (boolean): Whether streaming is active
- `fps` (float): Current frames per second
- `frame_time_ms` (float): Average frame render time in milliseconds
- `frames` (integer): Total output frames
- `dropped_frames` (integer): Total dropped/skipped frames

**Use Cases:**
- Health monitoring dashboards
- Verify OBS connection and version
- Performance monitoring and diagnostics
- Comprehensive status checks before operations
- Detect performance issues (dropped frames)
- System compatibility verification

**Example Natural Language Prompts:**
- "What's the OBS status?"
- "Check if OBS is working properly"
- "Show me OBS performance stats"
- "Are we streaming or recording?"
- "What version of OBS am I running?"
- "How many frames are being dropped?"

**Error Scenarios:**
- OBS not connected: "failed to get OBS status: OBS client not connected"
- Partial failure: Some fields may be empty if specific queries fail (non-fatal)
- WebSocket error: Connection or authentication issue

**Prerequisites:**
- Active OBS connection
- OBS Studio 28+ (for full feature support)

**Related Tools:**
- `get_recording_status` - Detailed recording information
- `get_streaming_status` - Detailed streaming information
- `list_scenes` - Scene list and current scene
- All other tools - Status check before operations

**Best Practices:**
- Use as initial connection verification
- Monitor dropped_frames for performance issues
- Check platform for OS-specific operations
- Verify version for feature compatibility
- Poll periodically for health monitoring
- High dropped_frames indicates performance problems
- Use fps and frame_time_ms for performance analysis
- Good for pre-flight checks before streaming/recording
- Includes both streaming and recording states in one call

---

## Common Patterns

### Pre-Flight Checks
Before starting a stream or recording:
```
1. get_obs_status - Verify OBS connection and performance
2. list_scenes - Confirm scenes exist
3. list_sources - Verify sources are configured
4. get_recording_status / get_streaming_status - Ensure not already active
5. set_current_scene - Switch to starting scene
6. start_streaming / start_recording - Begin output
```

### Scene Workflow
Creating and managing scenes:
```
1. list_scenes - Check if scene exists
2. create_scene - Create if needed
3. [Add sources to scene - future capability]
4. set_current_scene - Switch to new scene
5. toggle_source_visibility - Show/hide elements
```

### Audio Management
Complete audio control:
```
1. list_sources - Find audio inputs
2. get_input_mute - Check current state
3. toggle_input_mute - Mute/unmute
4. set_input_volume - Adjust levels
5. get_source_settings - Verify configuration
```

### Recording Session
Full recording workflow:
```
1. get_recording_status - Verify not already recording
2. set_current_scene - Switch to recording scene
3. start_recording - Begin capture
4. [During recording: pause_recording / resume_recording as needed]
5. get_recording_status - Monitor duration and size
6. stop_recording - End and save file
```

### Streaming Session
Complete streaming workflow:
```
1. get_streaming_status - Verify not already streaming
2. get_obs_status - Check performance
3. set_current_scene - Switch to starting soon scene
4. start_streaming - Go live
5. set_current_scene - Switch to main scene
6. [During stream: monitor with get_streaming_status]
7. set_current_scene - Switch to ending scene
8. stop_streaming - End stream
```

---

## Error Handling

### Common Error Patterns

**OBS Not Connected:**
```
Error: "OBS client not connected"
Solution: Verify OBS is running and WebSocket server is enabled
Check: OBS → Tools → WebSocket Server Settings
```

**Resource Not Found:**
```
Error: "Scene/Source/Input may not exist"
Solution: Use list_* tools to verify resource names
Check: Names are case-sensitive and must match exactly
```

**Already Active:**
```
Error: "output already active" (recording/streaming)
Solution: Check status before starting
Use: get_recording_status or get_streaming_status first
```

**Invalid State:**
```
Error: "Recording may not be paused" (when trying to resume)
Solution: Verify current state matches required precondition
Use: Status tools to confirm state before operations
```

**Authentication Failed:**
```
Error: WebSocket authentication failure
Solution: Verify password in agentic-obs configuration
Check: OBS WebSocket settings and stored credentials
```

### Error Recovery Strategies

1. **Connection Issues:** Check OBS is running, WebSocket enabled, correct port/password
2. **Resource Issues:** Use list tools to verify names before operations
3. **State Issues:** Check status before state-changing operations
4. **Permission Issues:** Verify disk space, write permissions for recording
5. **Configuration Issues:** Validate OBS settings for streaming/recording

### Validation Best Practices

- Always verify resources exist before operating on them
- Check status before start/stop operations to avoid errors
- Use exact, case-sensitive names for scenes/sources/inputs
- Handle errors gracefully with user-friendly messages
- Provide actionable error messages with suggested solutions
- Log errors for debugging and monitoring

---

**Document Version:** 2.0
**Last Updated:** 2025-12-15
**agentic-obs Version:** Phase 3 Complete
**Total Tools:** 30
