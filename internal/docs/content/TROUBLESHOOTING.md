# Troubleshooting

Common issues and solutions for agentic-obs.

## Connection Issues

### "Failed to connect to OBS"

**Symptoms:**
- Server fails to start
- "connection refused" errors
- Timeout on startup

**Solutions:**

1. **Verify OBS is running**
   - OBS Studio must be open before starting agentic-obs

2. **Enable WebSocket server in OBS**
   - Go to: Tools → WebSocket Server Settings
   - Check "Enable WebSocket server"
   - Note the port (default: 4455)

3. **Check firewall settings**
   - Ensure port 4455 (or your configured port) is not blocked
   - Windows: Check Windows Defender Firewall
   - Linux: Check iptables/ufw

4. **Verify password (if enabled)**
   - OBS WebSocket can require authentication
   - Password is stored in SQLite after first successful connection
   - Delete `~/.agentic-obs/db.sqlite` to reset and re-prompt

### Connection Drops / Auto-Reconnect

**Symptoms:**
- "OBS disconnected" in logs
- Tools fail with connection errors
- Reconnection messages

**Solutions:**

1. **OBS may have restarted** - agentic-obs auto-reconnects
2. **Network instability** - Check if running on different machines
3. **OBS crashed** - Check OBS logs for crash reports

---

## Web UI Issues (localhost:8765)

### Scenes Not Loading / Empty Scene Grid

**Cause:** OBS not connected or connection lost

**Check:**
- Visit `http://localhost:8765/api/status`
- Verify `"connected": true` in response

**Fix:**
- Ensure OBS is running with WebSocket server enabled
- Check OBS password matches config (stored in SQLite database)

### Thumbnails Showing Placeholders (SVG with scene name)

**Cause:** Screenshot capture failing for scene

**Check:**
- Look for `[Thumbnail] Cache miss` or errors in server logs

**Fix:**
- Verify OBS is connected and scenes exist
- Scene may have no visible sources - add a source to the scene
- Note: Placeholders are expected when OBS first connects (cache warming)

### Scene Click Does Nothing / Flickers But Doesn't Switch

**Cause:** Action executor not configured or JavaScript error

**Check:**
- Browser DevTools Console (F12) for JavaScript errors
- Server logs for `[UI Action]` messages

**Fix:**
- Ensure `SetStatusProvider()` was called before `Start()` in server code
- Clear browser cache and reload

### Audio Sliders at Wrong Position

**Cause:** Volume scale mismatch (linear vs logarithmic)

**Note:** Audio uses logarithmic curve matching human perception:
- 0dB = 100%
- -9dB ≈ 50%
- -18dB ≈ 25%
- -30dB ≈ 10%

Sliders snap to 0.5dB increments on release.

### "Template error" or Blank Page

**Cause:** Template parsing or embedding failed

**Check:**
- Server logs for `Template load error`

**Fix:**
- Ensure `templates/*.html` files exist and are valid HTML
- Rebuild with `go build` to re-embed templates

### UI Actions Return "action executor not configured"

**Cause:** HTTP server started before MCP server set as provider

**Fix:**
- Ensure initialization order: `httpServer.SetStatusProvider(mcpServer)` before `httpServer.Start()`

### Port 8765 Already in Use

**Cause:** Previous instance still running or another service using port

**Fix (Windows):**
```cmd
netstat -ano | findstr :8765
taskkill /PID <pid> /F
```

**Fix (Linux/macOS):**
```bash
lsof -i :8765
kill -9 <pid>
```

**Alternative:** Configure different port in settings

---

## Database Issues

### Database Write Errors

**Symptoms:**
- "database is locked" errors
- Failed to save configuration

**Solutions:**

1. **Check file permissions**
   - Ensure directory is writable
   - On Linux: `chmod 755 ~/.agentic-obs/`

2. **Check disk space**
   - SQLite needs space for journal files

3. **Close other instances**
   - Only one agentic-obs instance should access the database

### Reset Configuration

To start fresh, delete the database:

```bash
# Linux/macOS
rm ~/.agentic-obs/db.sqlite

# Windows
del %USERPROFILE%\.agentic-obs\db.sqlite
```

Next startup will re-prompt for OBS connection details.

---

## Screenshot Issues

### Screenshots Not Capturing

**Symptoms:**
- `list_screenshot_sources` shows sources but no images
- HTTP endpoint returns 404

**Check:**
- Server logs for capture errors
- OBS source actually exists and is visible

**Fix:**
- Ensure source name matches exactly (case-sensitive)
- Source must be in current scene or have "Show in all scenes" enabled

### Screenshots Outdated / Stale

**Cause:** Capture cadence too slow or background worker stopped

**Fix:**
- Use `configure_screenshot_cadence` to adjust interval
- Default is 5000ms (5 seconds)
- Minimum is 1000ms (1 second)

### High Memory Usage from Screenshots

**Cause:** Many screenshot sources with high cadence

**Mitigation:**
- Increase cadence (less frequent captures)
- Remove unused screenshot sources
- Screenshots auto-cleanup keeps only 10 per source

---

## MCP / AI Integration Issues

### Tools Not Appearing in Claude

**Cause:** Tool groups may be disabled

**Check:**
- First-run setup enables/disables tool groups
- Database stores tool group preferences

**Fix:**
- Delete database to re-run setup
- Or manually edit config in database

### "Tool not found" Errors

**Cause:** Tool may be in disabled group

**Tool Groups:**
| Group | Tools |
|-------|-------|
| Core | 13 (scenes, recording, streaming) |
| Sources | 3 (visibility, settings) |
| Audio | 4 (mute, volume) |
| Layout | 6 (presets) |
| Visual | 4 (screenshots) |
| Design | 14 (source creation/transforms) |

### Resource Notifications Not Received

**Cause:** Client may not support notifications

**Check:**
- MCP client capabilities
- Server logs for notification dispatch

**Note:** Resource notifications require persistent OBS connection

---

## TUI Dashboard Issues

### TUI Not Displaying Correctly

**Cause:** Terminal encoding or size issues

**Fix:**
- Ensure terminal supports UTF-8
- Terminal must be at least 80x24 characters
- Try different terminal emulator

### TUI Colors Wrong

**Cause:** Terminal color scheme conflicts

**Fix:**
- TUI adapts to terminal theme
- Try setting `TERM=xterm-256color`

---

## Getting More Help

1. **Check server logs** - Most issues are logged with context
2. **Use `help` tool** - `help` tool has topic-based guidance
3. **GitHub Issues** - Report bugs at repository
4. **API Status** - Visit `/api/status` for connection diagnostics
