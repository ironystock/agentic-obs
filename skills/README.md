# Claude Skills for agentic-obs

This directory contains Claude Skills that teach Claude how to effectively orchestrate the agentic-obs MCP server. These skills provide domain-specific expertise for managing OBS Studio through AI assistance.

## What are Claude Skills?

Claude Skills are structured markdown files that teach Claude how to use MCP tools effectively for specific domains. Each skill contains:

- **Trigger conditions**: When Claude should activate the skill
- **Tool inventory**: Which MCP tools are relevant for the task
- **Workflow guidance**: Step-by-step procedures for common operations
- **Best practices**: Domain-specific recommendations and tips
- **Examples**: Concrete use cases with input/output patterns

Skills enhance Claude's ability to orchestrate complex multi-tool workflows, provide context-aware recommendations, and apply domain expertise to user requests.

## Installing Skills in Claude Desktop

### Method 1: Copy to Claude Desktop Config

1. Locate your Claude Desktop configuration directory:
   - **macOS**: `~/Library/Application Support/Claude/`
   - **Windows**: `%APPDATA%\Claude\`
   - **Linux**: `~/.config/Claude/`

2. Copy the entire `skills` directory:
   ```bash
   # macOS/Linux
   cp -r skills ~/Library/Application\ Support/Claude/

   # Windows
   xcopy /E /I skills %APPDATA%\Claude\skills
   ```

3. Restart Claude Desktop

### Method 2: Symlink (Development)

For active development, symlink the skills directory:

```bash
# macOS/Linux
ln -s /path/to/agentic-obs/skills ~/Library/Application\ Support/Claude/skills

# Windows (Command Prompt as Administrator)
mklink /D "%APPDATA%\Claude\skills" "E:\code\agentic-obs\skills"
```

### Verification

After installation, Claude will automatically load these skills when appropriate. You can verify by asking:

> "What skills do you have available for OBS Studio?"

Claude should recognize and describe the agentic-obs skills.

## Available Skills

### 1. Streaming Assistant (`streaming-assistant`)

**When to use**: Pre-stream setup, live streaming management, source orchestration

**Key capabilities**:
- Pre-stream checklist execution (audio, video, scene verification)
- Real-time source management during streams
- Audio level monitoring and adjustment
- Scene preset application for quick transitions
- Stream health monitoring and diagnostics

**Tools used**: `get_obs_status`, `list_scenes`, `list_sources`, `toggle_source_visibility`, `get_input_mute`, `get_input_volume`, `list_scene_presets`, `apply_scene_preset`, `start_streaming`, `stop_streaming`, `get_streaming_status`

**Best for**: Users who need AI assistance during live streaming sessions, including pre-stream setup, real-time adjustments, and post-stream cleanup.

---

### 2. Scene Designer (`scene-designer`)

**When to use**: Creating visual layouts, designing stream overlays, source positioning

**Key capabilities**:
- Text overlay creation with custom fonts and styling
- Image source placement and transformation
- Color source backgrounds and accents
- Browser source integration (alerts, chat, widgets)
- Precise source positioning, scaling, and rotation
- Multi-source layout composition
- Crop and visibility control

**Tools used**: All 14 Design tools including `create_text_source`, `create_image_source`, `create_color_source`, `create_browser_source`, `set_source_transform`, `get_source_transform`, `set_source_crop`, `remove_source_from_scene`, `duplicate_source`, `set_source_index`, `set_source_blend_mode`, `set_source_locked`, `set_source_visible`, `get_scene_item_id`

**Best for**: Users designing stream layouts, creating overlays, positioning sources, or building complex visual compositions.

---

### 3. Audio Engineer (`audio-engineer`)

**When to use**: Audio configuration, troubleshooting sound issues, mixing levels

**Key capabilities**:
- Audio source identification and enumeration
- Mute state management (check and toggle)
- Volume level monitoring and adjustment
- Audio troubleshooting workflows
- Multi-source audio mixing guidance
- Audio quality verification

**Tools used**: `get_input_mute`, `toggle_input_mute`, `set_input_volume`, `get_input_volume`, `list_sources`, `get_obs_status`

**Best for**: Users configuring audio setups, troubleshooting sound issues, or managing multiple audio sources during production.

---

### 4. Preset Manager (`preset-manager`)

**When to use**: Managing scene presets, complex source visibility workflows

**Key capabilities**:
- Scene preset creation with source visibility states
- Preset organization and naming strategies
- Preset application and switching
- Preset comparison and diff analysis
- Preset cleanup and maintenance
- Multi-scene preset workflows

**Tools used**: `save_scene_preset`, `list_scene_presets`, `get_preset_details`, `apply_scene_preset`, `rename_scene_preset`, `delete_scene_preset`

**Best for**: Users managing multiple stream configurations, switching between source visibility states, or organizing complex scene setups.

---

## Skill Selection Guide

Claude will automatically select the appropriate skill based on your request. However, you can explicitly invoke a skill:

| User Request Pattern | Recommended Skill |
|---------------------|-------------------|
| "Help me set up for streaming" | `streaming-assistant` |
| "Create a text overlay that says..." | `scene-designer` |
| "Why can't I hear my microphone?" | `audio-engineer` |
| "Save this source configuration as a preset" | `preset-manager` |
| "Position this image in the bottom right" | `scene-designer` |
| "Check all my audio levels" | `audio-engineer` |
| "Switch to my gaming preset" | `preset-manager` |

## Using Skills Effectively

### Best Practices

1. **Be specific**: Describe what you want to achieve rather than which tools to use
   - Good: "Create a centered title that says 'LIVE' in red"
   - Better: "Create a centered title overlay that says 'LIVE' in red, positioned at the top of the screen"

2. **Provide context**: Mention the scene or sources you're working with
   - "In my Gaming scene, hide the webcam and show the full-screen game capture"

3. **Use natural language**: Skills enable Claude to understand domain terminology
   - "Mix down the music and boost my voice" (audio-engineer understands this)

4. **Iterate**: Skills support multi-step refinement
   - "Make that text bigger... no, a bit smaller... and move it left"

### Example Workflows

**Pre-stream Setup** (streaming-assistant):
```
User: "Help me get ready to stream"

Claude uses streaming-assistant skill to:
1. Check OBS connection status
2. List all scenes and sources
3. Verify audio inputs are unmuted
4. Check audio levels are appropriate
5. Confirm streaming settings
6. Provide pre-stream checklist
```

**Creating a Lower Third** (scene-designer):
```
User: "Create a lower third with my name and title"

Claude uses scene-designer skill to:
1. Create a color source for the background bar
2. Create a text source for the name
3. Create a text source for the title
4. Position and scale all elements
5. Apply appropriate styling
```

**Audio Troubleshooting** (audio-engineer):
```
User: "My microphone sounds quiet"

Claude uses audio-engineer skill to:
1. List all audio sources
2. Check microphone mute status
3. Get current volume level
4. Recommend volume adjustment
5. Verify changes took effect
```

**Preset Switching** (preset-manager):
```
User: "Switch to my 'BRB' preset"

Claude uses preset-manager skill to:
1. List available presets to confirm 'BRB' exists
2. Get preset details to explain what will change
3. Apply the preset
4. Confirm application success
```

## Combining Skills

Claude can use multiple skills in a single workflow:

**Example**: "Set up my Gaming scene with the webcam in the corner, then save it as a preset"
- Uses `scene-designer` to position the webcam
- Uses `preset-manager` to save the configuration

**Example**: "Prepare my stream and check the audio"
- Uses `streaming-assistant` for pre-stream setup
- Uses `audio-engineer` for detailed audio verification

## Troubleshooting

### Skills Not Loading

1. Verify the skills directory is in the correct location
2. Check that all SKILL.md files have valid YAML frontmatter
3. Restart Claude Desktop
4. Check Claude Desktop logs for skill loading errors

### Skills Not Activating

1. Be more explicit about your request (include keywords from skill descriptions)
2. Explicitly mention OBS or streaming in your request
3. Ask Claude: "Which skill would help me with [task]?"

## Contributing

To add or modify skills:

1. Follow the SKILL.md format with YAML frontmatter
2. Keep descriptions concise and trigger-focused
3. Include concrete examples and workflows
4. Test skills with Claude Desktop before submitting
5. Submit pull requests to the agentic-obs repository

## Additional Resources

- **agentic-obs GitHub**: [https://github.com/ironystock/agentic-obs](https://github.com/ironystock/agentic-obs)
- **MCP Documentation**: [https://modelcontextprotocol.io](https://modelcontextprotocol.io)
- **Claude Skills Guide**: Claude Desktop documentation
- **OBS Studio**: [https://obsproject.com](https://obsproject.com)

## Support

For issues with skills or agentic-obs:
- Open an issue on GitHub
- Check the agentic-obs README.md
- Review the CLAUDE.md for technical context

---

**Last Updated**: 2025-12-18
**agentic-obs Version**: 1.0.0
**Skills Version**: 1.0.0
