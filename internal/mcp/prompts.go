package mcp

import (
	"context"
	"fmt"
	"log"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerPrompts registers all MCP prompt handlers with the server
func (s *Server) registerPrompts() {
	// Prompt 1: Stream Launch Checklist
	s.mcpServer.AddPrompt(
		&mcpsdk.Prompt{
			Name:        "stream-launch",
			Description: "Pre-stream checklist to verify OBS is ready for streaming",
		},
		s.handleStreamLaunch,
	)

	// Prompt 2: Stream Teardown
	s.mcpServer.AddPrompt(
		&mcpsdk.Prompt{
			Name:        "stream-teardown",
			Description: "Post-stream cleanup to stop streaming/recording and switch to offline scene",
		},
		s.handleStreamTeardown,
	)

	// Prompt 3: Audio Check
	s.mcpServer.AddPrompt(
		&mcpsdk.Prompt{
			Name:        "audio-check",
			Description: "Verify all audio inputs are configured correctly with proper mute states and volumes",
		},
		s.handleAudioCheck,
	)

	// Prompt 4: Visual Check
	s.mcpServer.AddPrompt(
		&mcpsdk.Prompt{
			Name:        "visual-check",
			Description: "Analyze a screenshot to verify the stream layout and identify any visual issues",
			Arguments: []*mcpsdk.PromptArgument{{
				Name:        "screenshot_source",
				Description: "Name of the screenshot source to analyze",
				Required:    true,
			}},
		},
		s.handleVisualCheck,
	)

	// Prompt 5: Health Check
	s.mcpServer.AddPrompt(
		&mcpsdk.Prompt{
			Name:        "health-check",
			Description: "Run a comprehensive diagnostic of OBS connection, scenes, sources, and streaming state",
		},
		s.handleHealthCheck,
	)

	// Prompt 6: Problem Detection
	s.mcpServer.AddPrompt(
		&mcpsdk.Prompt{
			Name:        "problem-detection",
			Description: "Detect potential stream issues like black screens, frozen frames, or incorrect scenes",
			Arguments: []*mcpsdk.PromptArgument{{
				Name:        "screenshot_source",
				Description: "Name of the screenshot source to analyze for problems",
				Required:    true,
			}},
		},
		s.handleProblemDetection,
	)

	// Prompt 7: Preset Switcher
	s.mcpServer.AddPrompt(
		&mcpsdk.Prompt{
			Name:        "preset-switcher",
			Description: "Manage scene presets: list available presets and optionally apply one",
			Arguments: []*mcpsdk.PromptArgument{{
				Name:        "preset_name",
				Description: "Optional name of the preset to apply",
				Required:    false,
			}},
		},
		s.handlePresetSwitcher,
	)

	// Prompt 8: Recording Workflow
	s.mcpServer.AddPrompt(
		&mcpsdk.Prompt{
			Name:        "recording-workflow",
			Description: "Guide through recording operations: check status, start/stop, verify scene and audio",
		},
		s.handleRecordingWorkflow,
	)

	// Prompt 9: Scene Organizer
	s.mcpServer.AddPrompt(
		&mcpsdk.Prompt{
			Name:        "scene-organizer",
			Description: "Analyze scene structure and organization to suggest improvements",
		},
		s.handleSceneOrganizer,
	)

	// Prompt 10: Quick Status
	s.mcpServer.AddPrompt(
		&mcpsdk.Prompt{
			Name:        "quick-status",
			Description: "Get a brief summary of current OBS state (scene, recording, streaming)",
		},
		s.handleQuickStatus,
	)

	// Prompt 11: Scene Designer
	s.mcpServer.AddPrompt(
		&mcpsdk.Prompt{
			Name:        "scene-designer",
			Description: "Guide through creating and manipulating visual layouts with sources, transforms, and positioning",
			Arguments: []*mcpsdk.PromptArgument{
				{
					Name:        "scene_name",
					Description: "Name of the scene to design",
					Required:    true,
				},
				{
					Name:        "action",
					Description: "Optional action: 'add_elements', 'reposition', or 'adjust_layout'",
					Required:    false,
				},
			},
		},
		s.handleSceneDesigner,
	)

	// Prompt 12: Source Management
	s.mcpServer.AddPrompt(
		&mcpsdk.Prompt{
			Name:        "source-management",
			Description: "Manage existing sources in a scene: visibility, properties, duplication, and removal",
			Arguments: []*mcpsdk.PromptArgument{{
				Name:        "scene_name",
				Description: "Name of the scene containing sources to manage",
				Required:    true,
			}},
		},
		s.handleSourceManagement,
	)

	// Prompt 13: Visual Setup
	s.mcpServer.AddPrompt(
		&mcpsdk.Prompt{
			Name:        "visual-setup",
			Description: "Create and configure screenshot sources for AI visual monitoring of stream output",
			Arguments: []*mcpsdk.PromptArgument{{
				Name:        "monitor_scene",
				Description: "Optional scene name to monitor with screenshot source",
				Required:    false,
			}},
		},
		s.handleVisualSetup,
	)

	log.Println("Prompt handlers registered successfully")
}

// Prompt handler implementations

// handleStreamLaunch provides a pre-stream checklist workflow
func (s *Server) handleStreamLaunch(ctx context.Context, req *mcpsdk.GetPromptRequest) (*mcpsdk.GetPromptResult, error) {
	log.Println("Handling stream-launch prompt")

	promptText := `Help me launch my stream with these comprehensive checks:

1. **Verify OBS Connection**
   - Use get_obs_status to confirm OBS is connected and responsive
   - Report OBS version and connection state

2. **List All Available Scenes**
   - Use list_scenes to get all scenes
   - Identify the current active scene
   - Confirm the correct scene is selected for stream start

3. **Audio Source Verification**
   - Use list_sources to find all audio inputs
   - For each audio input, check mute states with get_input_mute
   - Check volume levels with get_input_volume
   - Report any audio sources that are muted or have unusual volume levels

4. **Check Recording and Streaming Status**
   - Use get_streaming_status to verify streaming is not already active
   - Use get_recording_status to check if recording is active
   - Report current state of both

5. **Provide Pre-Stream Recommendations**
   - Identify any issues that should be resolved before going live
   - Suggest adjustments to scene selection, audio levels, or mute states
   - Confirm when everything is ready for streaming

After completing all checks, provide a clear summary of OBS readiness for streaming.`

	return &mcpsdk.GetPromptResult{
		Description: "Pre-stream checklist to verify OBS is ready for streaming",
		Messages: []*mcpsdk.PromptMessage{{
			Role: "user",
			Content: &mcpsdk.TextContent{
				Text: promptText,
			},
		}},
	}, nil
}

// handleStreamTeardown provides post-stream cleanup workflow
func (s *Server) handleStreamTeardown(ctx context.Context, req *mcpsdk.GetPromptRequest) (*mcpsdk.GetPromptResult, error) {
	log.Println("Handling stream-teardown prompt")

	promptText := `Help me properly end my stream with these cleanup steps:

1. **Stop Active Streaming**
   - Use get_streaming_status to check if streaming is active
   - If streaming, use stop_streaming to end the stream
   - Confirm streaming has stopped successfully

2. **Stop Active Recording**
   - Use get_recording_status to check if recording is active
   - If recording, use stop_recording to save the recording
   - Report the output file path where the recording was saved

3. **Switch to Offline Scene**
   - Use list_scenes to identify available scenes
   - Look for scenes named "Offline", "BRB", "End Screen", or similar
   - If an offline scene exists, use set_current_scene to switch to it
   - If no offline scene exists, suggest creating one for future use

4. **Mute All Audio Inputs**
   - Use list_sources to find all audio inputs
   - For each audio input, check mute state with get_input_mute
   - If any inputs are unmuted, use toggle_input_mute to mute them
   - Confirm all audio is muted

5. **Final Status Confirmation**
   - Verify streaming is stopped
   - Verify recording is stopped
   - Confirm scene is set to offline/end screen
   - Confirm all audio is muted

Provide a summary of all teardown actions completed.`

	return &mcpsdk.GetPromptResult{
		Description: "Post-stream cleanup to stop streaming/recording and switch to offline scene",
		Messages: []*mcpsdk.PromptMessage{{
			Role: "user",
			Content: &mcpsdk.TextContent{
				Text: promptText,
			},
		}},
	}, nil
}

// handleAudioCheck provides audio verification workflow
func (s *Server) handleAudioCheck(ctx context.Context, req *mcpsdk.GetPromptRequest) (*mcpsdk.GetPromptResult, error) {
	log.Println("Handling audio-check prompt")

	promptText := `Help me verify all audio inputs are configured correctly:

1. **List All Audio Inputs**
   - Use list_sources to get all available sources
   - Filter to identify audio inputs (microphone, desktop audio, music, etc.)
   - Report the total count of audio inputs found

2. **Check Mute States**
   - For each audio input, use get_input_mute to check if it's muted
   - Report which inputs are muted and which are unmuted
   - Identify any unexpected mute states (e.g., microphone muted when it should be live)

3. **Check Volume Levels**
   - For each audio input, use get_input_volume to get current volume
   - Report volume in both dB and multiplier format
   - Identify any inputs with unusual volumes (too quiet, too loud, or at 0)

4. **Identify Audio Issues**
   - Flag any microphones that are muted
   - Flag any desktop audio sources that are unmuted (if unexpected)
   - Flag any volume levels outside normal ranges (-30dB to 0dB)
   - Flag any inputs set to 0% volume

5. **Provide Audio Recommendations**
   - Suggest which inputs should be unmuted for streaming/recording
   - Recommend volume adjustments if levels are too low or too high
   - Confirm when audio configuration is optimal

Provide a detailed audio configuration report with clear recommendations.`

	return &mcpsdk.GetPromptResult{
		Description: "Verify all audio inputs are configured correctly with proper mute states and volumes",
		Messages: []*mcpsdk.PromptMessage{{
			Role: "user",
			Content: &mcpsdk.TextContent{
				Text: promptText,
			},
		}},
	}, nil
}

// handleVisualCheck provides screenshot-based visual analysis workflow
func (s *Server) handleVisualCheck(ctx context.Context, req *mcpsdk.GetPromptRequest) (*mcpsdk.GetPromptResult, error) {
	log.Println("Handling visual-check prompt")

	screenshotSource := ""
	if req != nil && req.Params.Arguments != nil {
		if val, ok := req.Params.Arguments["screenshot_source"]; ok {
			screenshotSource = val
		}
	}

	if screenshotSource == "" {
		return nil, fmt.Errorf("screenshot_source argument is required")
	}

	promptText := fmt.Sprintf(`Help me verify the visual layout of my stream using screenshot analysis:

1. **Capture Current Screenshot**
   - Use list_screenshot_sources to verify screenshot source '%s' exists
   - If it doesn't exist, guide me to create it with create_screenshot_source
   - Access the screenshot at the HTTP URL provided by the screenshot source

2. **Analyze Visual Layout**
   - Examine the screenshot to identify all visible elements
   - Describe the scene composition (camera position, overlays, text, graphics)
   - Check if all expected visual elements are present and visible

3. **Check for Visual Issues**
   - Look for black screens or blank areas
   - Identify any frozen frames or static images
   - Check for visual glitches, artifacts, or rendering issues
   - Verify text is readable and not cut off
   - Check for proper aspect ratios and scaling

4. **Evaluate Stream Quality**
   - Assess overall visual quality and professionalism
   - Check color balance and brightness
   - Verify overlays and graphics are positioned correctly
   - Identify any visual elements that overlap inappropriately

5. **Provide Visual Recommendations**
   - Suggest improvements to layout and composition
   - Recommend adjustments to source positioning or sizing
   - Flag any critical issues that need immediate attention
   - Confirm when the visual setup is stream-ready

Provide a detailed visual analysis report with actionable recommendations.

Screenshot source to analyze: %s`, screenshotSource, screenshotSource)

	return &mcpsdk.GetPromptResult{
		Description: "Analyze a screenshot to verify the stream layout and identify any visual issues",
		Messages: []*mcpsdk.PromptMessage{{
			Role: "user",
			Content: &mcpsdk.TextContent{
				Text: promptText,
			},
		}},
	}, nil
}

// handleHealthCheck provides comprehensive OBS diagnostic workflow
func (s *Server) handleHealthCheck(ctx context.Context, req *mcpsdk.GetPromptRequest) (*mcpsdk.GetPromptResult, error) {
	log.Println("Handling health-check prompt")

	promptText := `Help me run a comprehensive diagnostic of OBS to ensure everything is working correctly:

1. **OBS Connection Status**
   - Use get_obs_status to verify connection to OBS
   - Report OBS version, connection state, and WebSocket status
   - Flag any connection issues

2. **Scene Inventory**
   - Use list_scenes to get all available scenes
   - Report total scene count
   - Identify the current active scene
   - Check if critical scenes exist (Main, Offline, BRB, etc.)

3. **Source Count and Status**
   - Use list_sources to get all input sources
   - Report total source count by type (audio, video, image, text, etc.)
   - Identify any disabled or problematic sources

4. **Recording State**
   - Use get_recording_status to check recording state
   - Report if recording is active, paused, or stopped
   - If recording is active, report duration and output path

5. **Streaming State**
   - Use get_streaming_status to check streaming state
   - Report if streaming is active or stopped
   - If streaming is active, report duration and connection status

6. **Screenshot Sources**
   - Use list_screenshot_sources to check configured screenshot sources
   - Report which screenshot sources are active
   - Verify screenshot HTTP endpoints are accessible

7. **Overall Health Assessment**
   - Identify any warning signs or issues detected
   - Provide recommendations for optimization
   - Confirm overall OBS health status (Healthy, Warning, Critical)

Provide a comprehensive health report with clear status indicators for each category.`

	return &mcpsdk.GetPromptResult{
		Description: "Run a comprehensive diagnostic of OBS connection, scenes, sources, and streaming state",
		Messages: []*mcpsdk.PromptMessage{{
			Role: "user",
			Content: &mcpsdk.TextContent{
				Text: promptText,
			},
		}},
	}, nil
}

// handleProblemDetection provides screenshot-based issue detection workflow
func (s *Server) handleProblemDetection(ctx context.Context, req *mcpsdk.GetPromptRequest) (*mcpsdk.GetPromptResult, error) {
	log.Println("Handling problem-detection prompt")

	screenshotSource := ""
	if req != nil && req.Params.Arguments != nil {
		if val, ok := req.Params.Arguments["screenshot_source"]; ok {
			screenshotSource = val
		}
	}

	if screenshotSource == "" {
		return nil, fmt.Errorf("screenshot_source argument is required")
	}

	promptText := fmt.Sprintf(`Help me detect potential streaming issues by analyzing the current screenshot:

1. **Capture and Access Screenshot**
   - Use list_screenshot_sources to verify screenshot source '%s' exists
   - Access the latest screenshot from the HTTP URL
   - Confirm the screenshot loaded successfully

2. **Black Screen Detection**
   - Analyze the screenshot for completely black or blank screens
   - Check if the image is mostly black (>90%% black pixels)
   - If detected, report as CRITICAL issue

3. **Frozen Frame Detection**
   - Compare with previous screenshots if available
   - Look for static elements that should be dynamic (timers, animations)
   - Report if the scene appears frozen or static

4. **Wrong Scene Detection**
   - Use list_scenes and identify the current active scene
   - Analyze screenshot content to verify it matches the expected scene
   - Check if an "Offline" or "BRB" scene is active when it shouldn't be
   - Flag if the wrong scene appears to be active

5. **Visual Artifact Detection**
   - Look for visual glitches, corruption, or rendering errors
   - Identify any overlapping elements that obscure important content
   - Check for missing sources that should be visible

6. **Audio-Visual Sync Issues**
   - Check if audio sources are present in the scene
   - Verify expected audio-visual sources are visible (e.g., waveform displays)

7. **Problem Report and Recommendations**
   - List all detected issues with severity (CRITICAL, WARNING, INFO)
   - Provide specific recommendations to fix each issue
   - If no issues found, confirm stream appears healthy

Provide a detailed problem detection report with prioritized issues and solutions.

Screenshot source to analyze: %s`, screenshotSource, screenshotSource)

	return &mcpsdk.GetPromptResult{
		Description: "Detect potential stream issues like black screens, frozen frames, or incorrect scenes",
		Messages: []*mcpsdk.PromptMessage{{
			Role: "user",
			Content: &mcpsdk.TextContent{
				Text: promptText,
			},
		}},
	}, nil
}

// handlePresetSwitcher provides scene preset management workflow
func (s *Server) handlePresetSwitcher(ctx context.Context, req *mcpsdk.GetPromptRequest) (*mcpsdk.GetPromptResult, error) {
	log.Println("Handling preset-switcher prompt")

	presetName := ""
	if req != nil && req.Params.Arguments != nil {
		if val, ok := req.Params.Arguments["preset_name"]; ok {
			presetName = val
		}
	}

	promptText := `Help me manage scene presets for quick configuration switching:

1. **List Available Presets**
   - Use list_scene_presets to get all saved presets
   - Display preset names, associated scene names, and creation dates
   - Report total count of available presets

2. **Show Current Scene State**
   - Use list_scenes to identify the current active scene
   - Use the scene resource (obs://scene/{name}) to get current source visibility states
   - Describe which sources are currently visible/hidden`

	if presetName != "" {
		promptText += fmt.Sprintf(`

3. **Apply Requested Preset**
   - Use get_preset_details to view preset '%s' configuration
   - Use apply_scene_preset to apply preset '%s'
   - Confirm which sources were enabled/disabled
   - Report success or any issues encountered`, presetName, presetName)
	} else {
		promptText += `

3. **Preset Selection Guidance**
   - Ask which preset I want to apply (if any)
   - Explain what each preset does based on its scene and source configuration
   - If I want to apply a preset, use apply_scene_preset with the chosen name`
	}

	promptText += `

4. **Preset Recommendations**
   - Suggest which preset might be appropriate for my current streaming activity
   - Offer to create a new preset to save the current scene configuration
   - Explain the benefits of using presets for quick scene setup

Provide a clear preset management interface with actionable options.`

	return &mcpsdk.GetPromptResult{
		Description: "Manage scene presets: list available presets and optionally apply one",
		Messages: []*mcpsdk.PromptMessage{{
			Role: "user",
			Content: &mcpsdk.TextContent{
				Text: promptText,
			},
		}},
	}, nil
}

// handleRecordingWorkflow provides recording management workflow
func (s *Server) handleRecordingWorkflow(ctx context.Context, req *mcpsdk.GetPromptRequest) (*mcpsdk.GetPromptResult, error) {
	log.Println("Handling recording-workflow prompt")

	promptText := `Help me manage my recording session with proper workflow:

1. **Check Current Recording Status**
   - Use get_recording_status to check if recording is active, paused, or stopped
   - If recording is active, report duration and output path
   - If recording is paused, show pause duration

2. **Verify Scene is Ready**
   - Use list_scenes to identify the current scene
   - Confirm the correct scene is active for recording
   - If wrong scene, suggest switching with set_current_scene

3. **Verify Audio Configuration**
   - Use list_sources to identify audio inputs
   - For each audio input, check mute state with get_input_mute
   - Check audio levels with get_input_volume
   - Flag any audio sources that should be unmuted but aren't

4. **Recording Control Guidance**
   - If not recording: offer to start_recording after confirming setup is ready
   - If recording: offer to pause_recording or stop_recording
   - If paused: offer to resume_recording or stop_recording
   - Explain the implications of each action

5. **Post-Recording Actions**
   - When stopping recording, report the final output file path
   - Confirm recording duration and file size (if available)
   - Suggest next steps (review recording, start new recording, etc.)

6. **Recording Best Practices**
   - Recommend checking disk space before starting long recordings
   - Suggest testing audio levels before recording important content
   - Remind about scene setup and preset usage

Provide a guided recording management experience with clear options.`

	return &mcpsdk.GetPromptResult{
		Description: "Guide through recording operations: check status, start/stop, verify scene and audio",
		Messages: []*mcpsdk.PromptMessage{{
			Role: "user",
			Content: &mcpsdk.TextContent{
				Text: promptText,
			},
		}},
	}, nil
}

// handleSceneOrganizer provides scene analysis and organization workflow
func (s *Server) handleSceneOrganizer(ctx context.Context, req *mcpsdk.GetPromptRequest) (*mcpsdk.GetPromptResult, error) {
	log.Println("Handling scene-organizer prompt")

	promptText := `Help me analyze and organize my OBS scenes for better workflow:

1. **List All Scenes**
   - Use list_scenes to get all available scenes
   - Report total scene count
   - Identify the current active scene

2. **Analyze Each Scene Structure**
   - For each scene, use the scene resource (obs://scene/{name}) to get details
   - Report the number of sources in each scene
   - List source types (image, video, audio, text, browser, etc.)
   - Identify which sources are enabled vs disabled

3. **Identify Scene Naming Patterns**
   - Analyze scene names for organization patterns
   - Check for common scene types (Main, Gaming, Chatting, Offline, BRB, etc.)
   - Flag any scenes with unclear or generic names

4. **Check for Scene Redundancy**
   - Identify scenes with very similar source configurations
   - Suggest consolidating redundant scenes
   - Recommend using scene presets instead of duplicate scenes

5. **Evaluate Scene Completeness**
   - Check if essential scenes exist (Main scene, Offline scene, BRB scene)
   - Identify missing scene types that might be useful
   - Flag empty scenes or scenes with no visible sources

6. **Organization Recommendations**
   - Suggest better naming conventions for clarity
   - Recommend scene ordering or grouping strategies
   - Propose creating scene presets for frequently used configurations
   - Suggest which scenes could be deleted or consolidated

7. **Scene Workflow Optimization**
   - Identify commonly used scene transitions
   - Suggest preset creation for quick scene setup
   - Recommend scene organization best practices

Provide a detailed scene organization report with actionable improvement suggestions.`

	return &mcpsdk.GetPromptResult{
		Description: "Analyze scene structure and organization to suggest improvements",
		Messages: []*mcpsdk.PromptMessage{{
			Role: "user",
			Content: &mcpsdk.TextContent{
				Text: promptText,
			},
		}},
	}, nil
}

// handleQuickStatus provides brief OBS status overview
func (s *Server) handleQuickStatus(ctx context.Context, req *mcpsdk.GetPromptRequest) (*mcpsdk.GetPromptResult, error) {
	log.Println("Handling quick-status prompt")

	promptText := `Give me a brief summary of my current OBS state:

1. **Current Scene**
   - Use list_scenes to identify the active scene
   - Report just the scene name

2. **Recording Status**
   - Use get_recording_status to check recording state
   - Report: "Recording: Yes (duration)" or "Recording: No"
   - If recording is paused, report: "Recording: Paused (duration)"

3. **Streaming Status**
   - Use get_streaming_status to check streaming state
   - Report: "Streaming: Yes (duration)" or "Streaming: No"

4. **Brief Format**
   - Present all information in a concise, easy-to-read format
   - Use simple yes/no answers with minimal details
   - Total response should be 3-4 lines maximum

Example output format:
---
Current Scene: Gaming
Recording: Yes (15:23)
Streaming: Yes (12:45)
---

Keep it short and clear - this is a quick status check, not a detailed report.`

	return &mcpsdk.GetPromptResult{
		Description: "Get a brief summary of current OBS state (scene, recording, streaming)",
		Messages: []*mcpsdk.PromptMessage{{
			Role: "user",
			Content: &mcpsdk.TextContent{
				Text: promptText,
			},
		}},
	}, nil
}

// handleSceneDesigner provides visual layout creation and manipulation workflow
func (s *Server) handleSceneDesigner(ctx context.Context, req *mcpsdk.GetPromptRequest) (*mcpsdk.GetPromptResult, error) {
	log.Println("Handling scene-designer prompt")

	sceneName := ""
	action := ""
	if req != nil && req.Params.Arguments != nil {
		if val, ok := req.Params.Arguments["scene_name"]; ok {
			sceneName = val
		}
		if val, ok := req.Params.Arguments["action"]; ok {
			action = val
		}
	}

	if sceneName == "" {
		return nil, fmt.Errorf("scene_name argument is required")
	}

	promptText := fmt.Sprintf(`Help me design and manipulate the visual layout for scene '%s':

1. **Verify Scene Exists**
   - Use list_scenes to check if scene '%s' exists
   - If the scene doesn't exist, use create_scene to create it
   - Use set_current_scene to make it active for preview
   - Access scene details via resource: obs://scene/%s

2. **List Existing Sources**
   - Use list_sources to see all available input sources
   - Identify sources currently in the scene
   - Check which sources are visible and enabled
   - Report current source count and types`, sceneName, sceneName, sceneName)

	switch action {
	case "add_elements":
		promptText += `

3. **Add New Visual Elements**
   - Use create_text_source to add text overlays (titles, labels, timers)
   - Use create_image_source to add static images (logos, backgrounds)
   - Use create_color_source to add solid color backgrounds or separators
   - Use create_browser_source to add web content (alerts, widgets)
   - Use create_media_source to add video files or animations
   - Use list_input_kinds to discover available source types

4. **Position New Sources**
   - Use set_source_transform to position each source (x, y coordinates)
   - Set source dimensions with width and height
   - Apply scaling with scale_x and scale_y
   - Rotate sources if needed with rotation angle
   - Use get_source_transform to verify current positioning

5. **Configure Source Properties**
   - Use set_source_bounds to control scaling behavior (fit, stretch, crop)
   - Use set_source_crop to crop edges (top, bottom, left, right)
   - Use set_source_order to arrange layer ordering (front to back)
   - Use set_source_locked to prevent accidental movement
   - Use toggle_source_visibility to show/hide sources`

	case "reposition":
		promptText += `

3. **Analyze Current Layout**
   - Use get_source_transform for each source to see current positions
   - Identify sources that overlap or are positioned incorrectly
   - Check for sources outside visible canvas area
   - Verify source dimensions and scaling

4. **Reposition Sources**
   - Use set_source_transform to move sources to new coordinates
   - Adjust source dimensions (width, height) as needed
   - Update scaling factors (scale_x, scale_y) for better fit
   - Change rotation angles if needed
   - Ensure no critical elements are cut off

5. **Verify and Fine-Tune**
   - Use get_source_transform to confirm new positions
   - Check for overlapping elements that shouldn't overlap
   - Verify all text is readable and not obscured
   - Use set_source_order to adjust layering if needed`

	case "adjust_layout":
		promptText += `

3. **Evaluate Current Layout**
   - Use get_source_transform to analyze all source positions and sizes
   - Use set_source_bounds to review scaling behavior
   - Use get_source_crop to check if any sources are cropped
   - Identify layout issues (overlaps, gaps, poor alignment)

4. **Adjust Layout Properties**
   - Use set_source_bounds to change how sources scale (fit, stretch, crop to bounds)
   - Use set_source_crop to trim unwanted edges (crop pixels from top/bottom/left/right)
   - Use set_source_transform to fine-tune positions and dimensions
   - Adjust layer ordering with set_source_order for proper visibility

5. **Lock and Finalize**
   - Use set_source_locked to lock sources in place once positioned
   - Save the layout as a preset with save_scene_preset
   - Document the layout design for future reference`

	default:
		// No specific action - general design workflow
		promptText += `

3. **Design Actions Available**

   **Adding Elements:**
   - create_text_source: Add text overlays, titles, labels, or timers
   - create_image_source: Add logos, backgrounds, or static graphics
   - create_color_source: Add solid color backgrounds or separators
   - create_browser_source: Add web-based content, alerts, or widgets
   - create_media_source: Add video files or animated content
   - list_input_kinds: Discover all available source types for creation

   **Positioning & Transform:**
   - set_source_transform: Position, scale, and rotate sources
   - get_source_transform: View current transform properties
   - set_source_bounds: Control how sources scale (fit, stretch, crop)
   - set_source_crop: Crop edges of sources (top, bottom, left, right)

   **Organization & Management:**
   - set_source_order: Arrange layer ordering (bring to front, send to back)
   - set_source_locked: Lock sources to prevent accidental changes
   - duplicate_source: Copy a source to reuse its configuration
   - remove_source: Delete sources from the scene
   - toggle_source_visibility: Show or hide sources

4. **Design Workflow Recommendations**
   - Start with background elements (colors, images)
   - Add mid-layer content (video sources, browser sources)
   - Add foreground elements (text, overlays)
   - Position and scale each element carefully
   - Use set_source_locked to protect finalized sources
   - Save layout as preset for easy restoration

5. **Best Practices**
   - Use descriptive source names for easy identification
   - Group related sources using consistent naming (e.g., "Title_Main", "Title_Subtitle")
   - Leave margins around text for readability
   - Ensure contrast between text and background
   - Test layout at target resolution before going live`
	}

	promptText += `

6. **Final Verification**
   - Review all sources with list_sources
   - Verify transforms with get_source_transform for each element
   - Check source visibility states with toggle_source_visibility
   - Use visual-check prompt with a screenshot source to see the final result
   - Save the design as a preset with save_scene_preset

Provide step-by-step guidance for designing the scene layout with clear recommendations.

Scene to design: ` + sceneName

	if action != "" {
		promptText += fmt.Sprintf("\nAction focus: %s", action)
	}

	return &mcpsdk.GetPromptResult{
		Description: "Guide through creating and manipulating visual layouts with sources, transforms, and positioning",
		Messages: []*mcpsdk.PromptMessage{{
			Role: "user",
			Content: &mcpsdk.TextContent{
				Text: promptText,
			},
		}},
	}, nil
}

// handleSourceManagement provides source management workflow
func (s *Server) handleSourceManagement(ctx context.Context, req *mcpsdk.GetPromptRequest) (*mcpsdk.GetPromptResult, error) {
	log.Println("Handling source-management prompt")

	sceneName := ""
	if req != nil && req.Params.Arguments != nil {
		if val, ok := req.Params.Arguments["scene_name"]; ok {
			sceneName = val
		}
	}

	if sceneName == "" {
		return nil, fmt.Errorf("scene_name argument is required")
	}

	promptText := fmt.Sprintf(`Help me manage sources in scene '%s':

1. **Inventory Current Sources**
   - Use list_scenes to verify scene '%s' exists
   - Use set_current_scene to make it active (if needed)
   - Use list_sources to get all sources in the scene
   - Access detailed scene info via resource: obs://scene/%s
   - Report total source count and types (video, audio, image, text, browser, etc.)

2. **Analyze Source States**
   - For each source, check visibility status (visible vs hidden)
   - Use toggle_source_visibility to review which sources are shown
   - Identify sources that are enabled but not visible
   - Check for sources that might be hidden accidentally
   - Flag any sources with unusual configurations

3. **Source Visibility Management**
   - Use toggle_source_visibility to show/hide specific sources
   - Identify which sources should be visible for different scenarios
   - Suggest creating scene presets for common visibility configurations
   - Report current visibility state for all sources

4. **Source Properties Review**
   - Use get_source_settings to view detailed properties for each source
   - Check resolution, format, and quality settings
   - Verify file paths for image and media sources
   - Review URLs for browser sources
   - Identify any sources with problematic settings

5. **Source Organization**
   - Use duplicate_source to copy sources that need variants
   - Use remove_source to delete unused or redundant sources
   - Use set_source_locked to protect important sources from changes
   - Suggest renaming sources for better organization (via source settings)

6. **Cleanup Recommendations**
   - Identify sources that are never used or always hidden
   - Suggest removing duplicate sources with identical settings
   - Flag sources that consume resources unnecessarily
   - Recommend consolidating similar sources

7. **Source Management Best Practices**
   - Keep source count reasonable (remove unused sources)
   - Use descriptive names for easy identification
   - Lock critical sources to prevent accidental deletion
   - Create presets to save visibility configurations
   - Document purpose of each source

8. **Action Recommendations**
   - Which sources should be removed to declutter?
   - Which sources should be duplicated for variants?
   - Which sources should be locked for protection?
   - Should any visibility states be saved as presets?

Provide a detailed source management report with actionable recommendations.

Scene to manage: %s`, sceneName, sceneName, sceneName, sceneName)

	return &mcpsdk.GetPromptResult{
		Description: "Manage existing sources in a scene: visibility, properties, duplication, and removal",
		Messages: []*mcpsdk.PromptMessage{{
			Role: "user",
			Content: &mcpsdk.TextContent{
				Text: promptText,
			},
		}},
	}, nil
}

// handleVisualSetup provides screenshot source setup workflow for AI monitoring
func (s *Server) handleVisualSetup(ctx context.Context, req *mcpsdk.GetPromptRequest) (*mcpsdk.GetPromptResult, error) {
	log.Println("Handling visual-setup prompt")

	monitorScene := ""
	if req != nil && req.Params.Arguments != nil {
		if val, ok := req.Params.Arguments["monitor_scene"]; ok {
			monitorScene = val
		}
	}

	promptText := `Help me set up screenshot sources for AI visual monitoring of my stream:

1. **Understand Screenshot Sources**
   - Screenshot sources enable AI to "see" your stream output periodically
   - Screenshots are captured at configurable intervals (cadence)
   - Images are served via HTTP at http://localhost:8765/screenshot/{name}
   - Screenshot resources available at obs://screenshot/{name}
   - Essential for automated visual monitoring and problem detection

2. **List Existing Screenshot Sources**
   - Use list_screenshot_sources to see currently configured sources
   - Review source names, scenes, cadence, and format (PNG/JPEG)
   - Check which screenshot sources are active
   - Identify gaps in visual monitoring coverage`

	if monitorScene != "" {
		promptText += fmt.Sprintf(`

3. **Create Screenshot Source for Scene '%s'**
   - Use create_screenshot_source with scene_name='%s'
   - Choose a descriptive source_name (e.g., '%s-monitor')
   - Set cadence to appropriate interval:
     * 5000ms (5s) for active monitoring during streaming
     * 10000ms (10s) for periodic checks
     * 30000ms (30s) for low-frequency monitoring
   - Select format: 'png' for quality or 'jpeg' for smaller size
   - Verify creation with list_screenshot_sources`, monitorScene, monitorScene, monitorScene)
	} else {
		promptText += `

3. **Create Screenshot Sources**
   - Identify scenes that need visual monitoring
   - Use create_screenshot_source for each scene:
     * source_name: Descriptive name (e.g., "main-stream-monitor")
     * scene_name: Scene to capture
     * cadence: Capture interval in milliseconds (default: 5000ms)
     * format: 'png' or 'jpeg' (default: 'png')
   - Common monitoring scenarios:
     * Main scene: For stream output verification
     * All active scenes: For comprehensive monitoring
     * Critical scenes: For problem detection`
	}

	promptText += `

4. **Configure Screenshot Cadence**
   - Use configure_screenshot_cadence to adjust capture frequency
   - Cadence recommendations:
     * 1000-3000ms: Real-time monitoring (high resource usage)
     * 5000ms: Standard monitoring (balanced)
     * 10000-15000ms: Periodic checks (lower resource usage)
     * 30000ms+: Infrequent monitoring (minimal overhead)
   - Balance monitoring needs with system performance
   - Faster cadence = more responsive but higher CPU/disk usage

5. **Verify Screenshot Access**
   - Screenshots available at: http://localhost:8765/screenshot/{source_name}
   - Access via MCP resource: obs://screenshot/{source_name}
   - Use visual-check prompt with screenshot_source to analyze content
   - Use problem-detection prompt for automated issue detection
   - Test HTTP endpoint accessibility

6. **Screenshot Source Management**
   - Use list_screenshot_sources to verify all sources are running
   - Use configure_screenshot_cadence to adjust intervals as needed
   - Use remove_screenshot_source to delete unused sources
   - Monitor disk usage (old screenshots are auto-cleaned)

7. **AI Monitoring Workflows**
   - **Visual Check:** Use visual-check prompt to analyze layout and quality
   - **Problem Detection:** Use problem-detection prompt to find issues
   - **Health Monitoring:** Periodically verify scenes look correct
   - **Automated Alerts:** Detect black screens, frozen frames, wrong scenes

8. **Best Practices**
   - Create screenshot sources for all scenes used during streaming
   - Use descriptive names: "{scene}-monitor" or "{purpose}-check"
   - Start with 5000ms cadence, adjust based on needs
   - Use PNG for quality, JPEG for lower disk usage
   - Clean up screenshot sources for deleted scenes
   - Combine with visual-check and problem-detection prompts

9. **Example Setup**
   - Main streaming scene: create_screenshot_source(source_name="main-monitor", scene_name="Gaming", cadence=5000)
   - Offline scene: create_screenshot_source(source_name="offline-monitor", scene_name="Offline", cadence=10000)
   - Verify: list_screenshot_sources
   - Test: Access http://localhost:8765/screenshot/main-monitor

Provide guidance on setting up comprehensive visual monitoring for AI-driven stream management.`

	if monitorScene != "" {
		promptText += fmt.Sprintf("\n\nTarget scene for monitoring: %s", monitorScene)
	}

	return &mcpsdk.GetPromptResult{
		Description: "Create and configure screenshot sources for AI visual monitoring of stream output",
		Messages: []*mcpsdk.PromptMessage{{
			Role: "user",
			Content: &mcpsdk.TextContent{
				Text: promptText,
			},
		}},
	}, nil
}
