package mcp

// toolHelpContent maps tool names to their detailed help text.
// This is separate from help_content.go for maintainability.
var toolHelpContent = map[string]string{
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

	// Filters (FB-23)
	"list_source_filters": `# list_source_filters

**Category**: Filters

**Description**: List all filters applied to a source.

**Input**:
- source_name (string, required): Name of the source to list filters for

**Output**:
- source_name: Source name
- filters: Array of filter objects (name, kind, index, enabled)
- count: Total number of filters

**Example Input**:
{
  "source_name": "Webcam"
}

**Example Output**:
{
  "source_name": "Webcam",
  "filters": [
    {"name": "Color Correction", "kind": "color_filter_v2", "index": 0, "enabled": true},
    {"name": "Sharpen", "kind": "sharpness_filter_v2", "index": 1, "enabled": true}
  ],
  "count": 2
}`,

	"get_source_filter": `# get_source_filter

**Category**: Filters

**Description**: Get detailed information about a specific filter on a source.

**Input**:
- source_name (string, required): Name of the source containing the filter
- filter_name (string, required): Name of the filter to get details for

**Output**:
- source_name: Source name
- filter: Filter details (name, kind, index, enabled, settings)

**Example Input**:
{
  "source_name": "Webcam",
  "filter_name": "Color Correction"
}

**Use Case**: Inspect filter configuration before modifying settings.`,

	"create_source_filter": `# create_source_filter

**Category**: Filters

**Description**: Add a new filter to a source (e.g., color correction, noise suppression).

**Input**:
- source_name (string, required): Name of the source to add the filter to
- filter_name (string, required): Name for the new filter
- filter_kind (string, required): Type of filter (use list_filter_kinds to see available types)
- filter_settings (object, optional): Initial settings for the filter

**Output**:
- source_name: Source name
- filter_name: Filter name
- filter_kind: Filter type
- message: Success confirmation

**Example Input**:
{
  "source_name": "Webcam",
  "filter_name": "My Color Correction",
  "filter_kind": "color_filter_v2",
  "filter_settings": {"brightness": 0.1, "contrast": 0.05}
}

**Common Filter Types**:
- color_filter_v2: Color correction (brightness, contrast, saturation)
- sharpness_filter_v2: Image sharpening
- noise_suppress_filter_v2: Audio noise suppression
- compressor_filter: Audio compressor
- chroma_key_filter_v2: Green screen removal`,

	"remove_source_filter": `# remove_source_filter

**Category**: Filters

**Description**: Remove a filter from a source.

**Input**:
- source_name (string, required): Name of the source containing the filter
- filter_name (string, required): Name of the filter to remove

**Output**:
- source_name: Source name
- filter_name: Filter name
- message: Success confirmation

**Example Input**:
{
  "source_name": "Webcam",
  "filter_name": "Old Filter"
}

**Warning**: This action cannot be undone. The filter configuration will be lost.`,

	"toggle_source_filter": `# toggle_source_filter

**Category**: Filters

**Description**: Enable or disable a filter on a source.

**Input**:
- source_name (string, required): Name of the source containing the filter
- filter_name (string, required): Name of the filter to toggle
- filter_enabled (bool, optional): Set to true/false to enable/disable; omit to toggle

**Output**:
- source_name: Source name
- filter_name: Filter name
- filter_enabled: New enabled state
- message: Success confirmation

**Example Input** (toggle):
{
  "source_name": "Webcam",
  "filter_name": "Color Correction"
}

**Example Input** (explicit):
{
  "source_name": "Webcam",
  "filter_name": "Color Correction",
  "filter_enabled": false
}

**Use Case**: Quickly enable/disable effects without removing the filter.`,

	"set_source_filter_settings": `# set_source_filter_settings

**Category**: Filters

**Description**: Modify the configuration settings of a filter.

**Input**:
- source_name (string, required): Name of the source containing the filter
- filter_name (string, required): Name of the filter to update
- filter_settings (object, required): Settings to apply to the filter
- overlay (bool, optional): If true, merge with existing settings; if false, replace entirely (default: true)

**Output**:
- source_name: Source name
- filter_name: Filter name
- overlay: Whether overlay mode was used
- message: Success confirmation

**Example Input**:
{
  "source_name": "Webcam",
  "filter_name": "Color Correction",
  "filter_settings": {"brightness": 0.2, "saturation": 0.1},
  "overlay": true
}

**Note**: Use overlay=true to update specific settings while keeping others unchanged.`,

	"list_filter_kinds": `# list_filter_kinds

**Category**: Filters

**Description**: List all available filter types in OBS.

**Input**: None

**Output**:
- filter_kinds: Array of available filter type IDs
- count: Total number of types

**Example Output**:
{
  "filter_kinds": [
    "color_filter_v2",
    "sharpness_filter_v2",
    "noise_suppress_filter_v2",
    "compressor_filter",
    "limiter_filter",
    "gain_filter",
    "chroma_key_filter_v2"
  ],
  "count": 15
}

**Use Case**: Discover available filter types before creating filters with create_source_filter.`,

	// Transitions (FB-24)
	"list_transitions": `# list_transitions

**Category**: Transitions

**Description**: List all available scene transitions and identify the current one.

**Input**: None

**Output**:
- transitions: Array of transition objects (name, kind, fixed, configurable)
- current_transition: Name of currently active transition
- count: Total number of transitions

**Example Output**:
{
  "transitions": [
    {"name": "Cut", "kind": "cut_transition", "fixed": true, "configurable": false},
    {"name": "Fade", "kind": "fade_transition", "fixed": false, "configurable": true},
    {"name": "Swipe", "kind": "swipe_transition", "fixed": false, "configurable": true}
  ],
  "current_transition": "Fade",
  "count": 3
}`,

	"get_current_transition": `# get_current_transition

**Category**: Transitions

**Description**: Get details about the current scene transition including duration and settings.

**Input**: None

**Output**:
- name: Transition name
- kind: Transition type
- duration_ms: Transition duration in milliseconds
- configurable: Whether transition has configurable settings
- settings: Current transition settings (if configurable)

**Example Output**:
{
  "name": "Fade",
  "kind": "fade_transition",
  "duration_ms": 300,
  "configurable": true,
  "settings": {}
}`,

	"set_current_transition": `# set_current_transition

**Category**: Transitions

**Description**: Change the active scene transition (e.g., Cut, Fade, Swipe).

**Input**:
- transition_name (string, required): Name of the transition to set as current

**Output**:
- transition_name: Transition name
- message: Success confirmation

**Example Input**:
{
  "transition_name": "Fade"
}

**Use Case**: Change how scenes transition during scene switches.`,

	"set_transition_duration": `# set_transition_duration

**Category**: Transitions

**Description**: Set the duration of the current scene transition in milliseconds.

**Input**:
- transition_duration (int, required): Duration in milliseconds

**Output**:
- duration_ms: New duration value
- message: Success confirmation

**Example Input**:
{
  "transition_duration": 500
}

**Note**: Typical durations range from 100ms (quick) to 1000ms (slow). Cut transition ignores duration.`,

	"trigger_transition": `# trigger_transition

**Category**: Transitions

**Description**: Trigger the current transition in studio mode (swaps preview and program scenes).

**Input**: None

**Output**:
- message: Success confirmation

**Requirements**: OBS must be in Studio Mode for this to work.

**Use Case**: Manually trigger scene changes in studio mode workflow.

**Error Handling**: Returns error if studio mode is not enabled.`,
}

// GetToolHelpContent returns the help text for a specific tool, or empty if not found.
func GetToolHelpContent(toolName string) (string, bool) {
	help, exists := toolHelpContent[toolName]
	return help, exists
}
