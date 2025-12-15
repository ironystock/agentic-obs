# Agentic Screenshot Sources

## Why Screenshots Matter for AI

Traditional OBS integrations allow AI to *control* your stream, but they can't *see* it. The screenshot feature bridges this gap, enabling AI assistants to visually observe your OBS output - turning a blind controller into a seeing collaborator.

### The Vision Problem

Without visual context, AI assistants face significant limitations:

- **Scene Verification**: When you say "switch to Gaming scene," the AI can confirm the command succeeded, but it can't verify the scene actually looks correct
- **Source Issues**: Overlays might be mispositioned, webcams frozen, or game capture showing a black screen - the AI has no way to know
- **Quality Monitoring**: Encoding artifacts, incorrect resolutions, or color issues are invisible to a text-based assistant
- **Layout Feedback**: When arranging sources, the AI can't provide design feedback without seeing the result

### The Screenshot Solution

Periodic screenshot capture enables AI to:

1. **Visually verify** that scenes look as expected after changes
2. **Detect problems** like black screens, frozen sources, or missing overlays
3. **Provide feedback** on layouts, composition, and visual quality
4. **Monitor streams** for issues that only manifest visually
5. **Assist with design** by seeing and suggesting improvements

---

## How It Works

### Architecture Overview

```
AI Assistant                agentic-obs                    OBS Studio
     |                           |                              |
     |  create_screenshot_source |                              |
     |-------------------------->|                              |
     |                           |   capture every 5s           |
     |                           |----------------------------->|
     |                           |   <-- base64 PNG             |
     |                           |                              |
     |                           |   Store in SQLite            |
     |                           |   Serve via HTTP             |
     |                           |                              |
     |  GET /screenshot/mystream |                              |
     |<--------------------------|                              |
     |  <-- PNG image            |                              |
```

### Components

1. **Screenshot Sources**: Named configurations that define what to capture and how often
2. **Background Workers**: Goroutines that capture screenshots at the configured cadence
3. **SQLite Storage**: Stores screenshot images with automatic cleanup (keeps latest 10 per source)
4. **HTTP Server**: Serves screenshots at `http://localhost:8765/screenshot/{name}`

---

## Quick Start

### 1. Create a Screenshot Source

Tell the AI:
> "Create a screenshot source called 'stream-monitor' that captures my current scene every 5 seconds"

Or more specifically:
> "Set up a screenshot capture of my Gaming scene at 1080p, saving a PNG every 10 seconds, call it 'game-view'"

### 2. Access Screenshots

The AI can access screenshots via:
- **HTTP URL**: `http://localhost:8765/screenshot/stream-monitor`
- **Direct Tool Call**: Use `list_screenshot_sources` to get URLs

### 3. Use for Verification

After making changes:
> "Switch to my BRB scene and show me a screenshot to verify it looks right"

---

## Use Cases

### Stream Monitoring

**Scenario**: You're streaming but can't see your own output easily.

**Solution**:
> "Create a screenshot source called 'stream-check' that captures every 30 seconds. Let me know if anything looks wrong."

The AI can periodically check:
- Is the webcam visible?
- Are overlays in the right position?
- Is the game capture working?

**Example Prompts**:
- "Check my stream output - does everything look normal?"
- "Take a screenshot and tell me if my webcam is showing"
- "Verify my BRB scene has the timer overlay visible"

---

### Scene Verification

**Scenario**: You've made changes to a scene and want confirmation.

**Solution**:
> "I just rearranged my Gaming scene. Take a screenshot and tell me how it looks"

**Example Workflow**:
1. "Move my webcam source to the bottom right corner"
2. "Take a screenshot of the Gaming scene"
3. "Does my webcam look properly positioned?"

**Example Prompts**:
- "Show me what my current scene looks like"
- "Take a screenshot after switching to Intro scene"
- "Verify my layout changes look correct"

---

### Problem Detection

**Scenario**: Something's wrong but you're not sure what.

**Solution**:
> "My stream looks off according to chat. Take a screenshot and help me figure out what's wrong"

**Common Issues AI Can Detect**:
- Black or frozen game capture
- Missing or mispositioned overlays
- Webcam not visible or in wrong location
- Incorrect aspect ratio or resolution
- Text overlays not rendering

**Example Prompts**:
- "Viewers say my game capture is black - can you check?"
- "Take a screenshot and tell me if my alert box is visible"
- "Something looks wrong with my layout - what is it?"

---

### Multi-Scene Comparison

**Scenario**: You want to compare layouts across scenes.

**Solution**:
> "Take screenshots of all my scenes so I can review them"

**Example Workflow**:
1. Create screenshot sources for each scene
2. "List all my screenshot sources with their URLs"
3. "Compare my Gaming and Chatting scenes - are they consistent?"

---

### Design Feedback

**Scenario**: You're designing a new scene and want AI input.

**Solution**:
> "I'm working on a new overlay design. Take a screenshot and give me feedback on the composition"

**Example Prompts**:
- "How does my new scene layout look? Any suggestions?"
- "Take a screenshot - is there too much empty space?"
- "Does my color scheme work well together?"
- "Is my webcam size appropriate for the layout?"

---

### Automated Monitoring

**Scenario**: Long streams where you want periodic checks.

**Solution**:
> "Set up continuous monitoring of my stream at 30-second intervals. Alert me if you notice any issues."

The AI can watch for:
- Sudden scene changes
- Sources appearing/disappearing unexpectedly
- Quality degradation
- Layout shifts

---

### Recording Quality Checks

**Scenario**: Before starting an important recording.

**Solution**:
> "I'm about to record a tutorial. Take a screenshot and verify everything looks professional"

**Pre-Recording Checklist**:
1. All sources visible and positioned correctly
2. No test elements or debug overlays showing
3. Clean, professional appearance
4. Correct resolution and aspect ratio

---

### Browser Source Integration

**Advanced Use Case**: Create a browser source in OBS that displays the screenshot URL, enabling:
- Picture-in-picture monitoring
- Debug views during production
- External dashboard integration

```
Screenshot URL: http://localhost:8765/screenshot/stream-monitor
Add as Browser Source in OBS for live preview monitoring
```

---

## Tool Reference

### create_screenshot_source

Create a new periodic screenshot capture source.

**Parameters**:
| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| name | string | Yes | - | Unique identifier for this source |
| source_name | string | Yes | - | OBS source or scene name to capture |
| cadence_ms | integer | No | 5000 | Capture interval in milliseconds |
| image_format | string | No | "png" | Image format: "png" or "jpg" |
| image_width | integer | No | 0 | Resize width (0 = original) |
| image_height | integer | No | 0 | Resize height (0 = original) |
| quality | integer | No | 80 | Compression quality (1-100) |

**Returns**:
```json
{
  "id": 1,
  "name": "stream-monitor",
  "url": "http://localhost:8765/screenshot/stream-monitor",
  "message": "Screenshot source created successfully"
}
```

**Example Prompts**:
- "Create a screenshot source called 'main-view' for my Gaming scene"
- "Set up 10-second screenshot capture of my current output as 'stream-check'"
- "Create a high-quality PNG screenshot source at 1920x1080 called 'archive'"
- "Make a fast JPG capture at 720p every 2 seconds for monitoring, call it 'quick-check'"

---

### remove_screenshot_source

Stop and delete a screenshot source.

**Parameters**:
| Name | Type | Required | Description |
|------|------|----------|-------------|
| name | string | Yes | Name of the source to remove |

**Returns**:
```json
{
  "message": "Screenshot source 'stream-monitor' removed successfully"
}
```

**Example Prompts**:
- "Remove the screenshot source called 'test-capture'"
- "Stop capturing screenshots for 'old-monitor'"
- "Delete the 'temp-check' screenshot source"

---

### list_screenshot_sources

List all configured screenshot sources with their status and URLs.

**Parameters**: None

**Returns**:
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

**Example Prompts**:
- "What screenshot sources do I have set up?"
- "Show me all my screenshot capture configurations"
- "List screenshot URLs for all my sources"

---

### configure_screenshot_cadence

Update the capture interval for an existing source.

**Parameters**:
| Name | Type | Required | Description |
|------|------|----------|-------------|
| name | string | Yes | Name of the source to update |
| cadence_ms | integer | Yes | New capture interval in milliseconds |

**Returns**:
```json
{
  "message": "Screenshot cadence updated to 10000ms for source 'stream-monitor'"
}
```

**Example Prompts**:
- "Change my stream-monitor to capture every 10 seconds"
- "Speed up the screenshot capture to every 2 seconds"
- "Slow down 'quick-check' to 30-second intervals"

---

## Best Practices

### Choosing Capture Intervals

| Use Case | Recommended Interval | Rationale |
|----------|---------------------|-----------|
| Real-time monitoring | 1-2 seconds | Quick feedback, higher resource use |
| General verification | 5 seconds | Good balance of responsiveness and efficiency |
| Long-term monitoring | 30-60 seconds | Minimal overhead, catches major issues |
| Archival/logging | 5-10 minutes | Documentation purposes |

### Image Format Selection

| Format | Best For | Trade-offs |
|--------|----------|------------|
| PNG | Quality-critical captures, text overlays | Larger files (~1-2MB at 1080p) |
| JPG | Frequent captures, monitoring | Smaller files, some quality loss |

### Resolution Considerations

- **Full resolution (0x0)**: Best quality, larger storage
- **1920x1080**: Standard HD, good balance
- **1280x720**: Faster captures, less storage
- **640x360**: Minimal footprint for quick checks

### Storage Management

The system automatically keeps only the **10 most recent screenshots per source**. With typical 1080p PNG screenshots:
- ~1-2MB per image
- ~10-20MB storage per source
- Cleanup runs every minute

For long-running sessions with many sources, monitor database size with:
```
Database location: ./agentic-obs.db
```

---

## Technical Details

### HTTP Endpoint

**URL**: `http://localhost:8765/screenshot/{source_name}`

**Method**: GET

**Response**: Binary image with appropriate Content-Type header

**Headers**:
- `Content-Type`: `image/png` or `image/jpeg`
- `X-Screenshot-Source`: Source name
- `X-Screenshot-Captured`: ISO 8601 timestamp

### Error Responses

| Status | Meaning |
|--------|---------|
| 400 | Invalid source name |
| 404 | Source not found or no screenshots available |
| 500 | Internal error (check server logs) |

### Security Notes

- HTTP server binds to `localhost` by default (not exposed externally)
- Source names are validated to prevent path traversal
- Consider rate limiting if exposing externally via reverse proxy

---

## Troubleshooting

### No Screenshots Available

**Symptoms**: 404 error when accessing screenshot URL

**Causes**:
- Source just created, first capture pending
- OBS not connected
- Invalid source name in OBS

**Solutions**:
1. Wait a few seconds for first capture
2. Verify OBS is running and connected
3. Check source name matches exactly (case-sensitive)

### Screenshots Not Updating

**Symptoms**: Same screenshot despite time passing

**Causes**:
- Screenshot source disabled
- OBS source is frozen/static
- Capture cadence is very long

**Solutions**:
1. Check source is enabled with `list_screenshot_sources`
2. Verify OBS source is active
3. Reduce cadence for more frequent captures

### High Resource Usage

**Symptoms**: Lag or high CPU during captures

**Causes**:
- Very short capture intervals (< 1 second)
- High resolution captures
- Many simultaneous sources

**Solutions**:
1. Increase cadence (5+ seconds for monitoring)
2. Reduce resolution if full quality not needed
3. Use JPG format for smaller files
4. Remove unused screenshot sources

---

## FAQ

**Q: Does this affect stream performance?**
A: Minimal impact. Screenshots are captured using OBS's native screenshot API, which is optimized for efficiency. With reasonable intervals (5+ seconds), impact is negligible.

**Q: Can I capture specific sources or only full scenes?**
A: You can capture any named source or scene in OBS. Specify the exact source name when creating the screenshot source.

**Q: How long are screenshots kept?**
A: The 10 most recent screenshots per source are kept. Older ones are automatically deleted during cleanup (runs every minute).

**Q: Can I access screenshots externally?**
A: By default, the HTTP server only accepts connections from localhost. For external access, use a reverse proxy with appropriate security measures.

**Q: What happens if OBS disconnects?**
A: Screenshot capture pauses automatically when OBS is disconnected and resumes when reconnected. No errors are generated for missed captures.

---

**Document Version**: 1.0
**Last Updated**: 2025-12-15
**Related**: [TOOLS.md](TOOLS.md) | [Workflows](../examples/prompts/workflows.md)
