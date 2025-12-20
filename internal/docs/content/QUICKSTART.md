# Quick Start Guide: agentic-obs

**Get AI control of OBS Studio in 10 minutes**

This guide will walk you through installing and configuring agentic-obs, a Model Context Protocol (MCP) server that gives Claude Desktop (and other AI assistants) programmatic control over OBS Studio.

---

## Prerequisites Check

Before you begin, ensure you have the following installed:

### Required Software

1. **Go 1.25.5 or later**
   - Check your version: `go version`
   - Download from: https://go.dev/dl/
   - Required for installing the agentic-obs server

2. **OBS Studio 28 or later**
   - Check: Open OBS → Help → About
   - Download from: https://obsproject.com/download
   - Must have built-in WebSocket server (v28+)

3. **Git**
   - Check your version: `git --version`
   - Download from: https://git-scm.com/downloads
   - Required for `go install` to fetch the package

4. **Claude Desktop**
   - Download from: https://claude.ai/download
   - The AI assistant that will control OBS

### System Requirements

- **Operating System**: Windows, macOS, or Linux
- **Network**: Local WebSocket connection (localhost:4455)
- **Disk Space**: ~50 MB for the binary and SQLite database

---

## Installation

### Step 1: Install agentic-obs

Open your terminal and run:

```bash
go install github.com/ironystock/agentic-obs@latest
```

This will:
- Download the latest version from GitHub
- Compile the binary for your system
- Install it to your `$GOPATH/bin` directory (typically `~/go/bin` or `%USERPROFILE%\go\bin`)

**Verify installation:**

```bash
# Check that the binary is in your PATH
agentic-obs --help
```

If you get "command not found", add Go's bin directory to your PATH:

- **macOS/Linux**: Add to `~/.bashrc` or `~/.zshrc`:
  ```bash
  export PATH=$PATH:$(go env GOPATH)/bin
  ```

- **Windows**: Add `%USERPROFILE%\go\bin` to your system PATH environment variable

---

## OBS WebSocket Setup

agentic-obs communicates with OBS Studio via the built-in WebSocket server. You need to enable and configure it.

### Step 2: Enable OBS WebSocket Server

1. **Open OBS Studio**

2. **Navigate to WebSocket Settings**:
   - Click **Tools** (top menu)
   - Select **WebSocket Server Settings**

3. **Configure the server**:
   - Check **Enable WebSocket server**
   - **Server Port**: `4455` (default, recommended)
   - **Enable Authentication**: Optional but recommended
     - If enabled, set a password (e.g., "my_secret_password")
     - You'll need this later during first run
   - Click **OK** to save

4. **Verify it's running**:
   - You should see "WebSocket server started successfully" in OBS output log
   - The status indicator should be green in the WebSocket settings

**Note**: Keep OBS Studio running whenever you want to use agentic-obs. The WebSocket server only runs when OBS is open.

---

## First Run and Auto-Detection

### Step 3: Run agentic-obs for the First Time

When you first run agentic-obs, it will automatically detect and connect to your OBS WebSocket server.

```bash
agentic-obs
```

**What happens**:

1. **Auto-detection**: The server attempts to connect to `localhost:4455`
2. **Authentication prompt** (if you set a password):
   ```
   Failed to connect to OBS WebSocket
   Please enter OBS WebSocket password:
   ```
   - Enter the password you configured in OBS
   - Press Enter

3. **Success**:
   ```
   Connected to OBS WebSocket at localhost:4455
   MCP server listening on stdio
   ```

4. **Configuration saved**: Your connection settings are stored in SQLite (`agentic-obs.db`) for future runs

**Troubleshooting first run**:

- **"Connection refused"**: Make sure OBS Studio is running and WebSocket server is enabled
- **"Authentication failed"**: Double-check your password in OBS → Tools → WebSocket Server Settings
- **"Port already in use"**: Change the port in OBS settings and provide it when prompted

---

## Claude Desktop Configuration

### Step 4: Add agentic-obs to Claude Desktop

To give Claude access to OBS controls, you need to register the MCP server in Claude Desktop's configuration file.

#### Find your config file:

- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
- **Linux**: `~/.config/Claude/claude_desktop_config.json`

#### Edit the configuration:

Open `claude_desktop_config.json` in a text editor and add the following:

```json
{
  "mcpServers": {
    "agentic-obs": {
      "command": "agentic-obs",
      "args": []
    }
  }
}
```

**If you already have other MCP servers**, add agentic-obs to the existing `mcpServers` object:

```json
{
  "mcpServers": {
    "filesystem": {
      "command": "mcp-server-filesystem",
      "args": ["/Users/yourname/Documents"]
    },
    "agentic-obs": {
      "command": "agentic-obs",
      "args": []
    }
  }
}
```

**Important**: If `agentic-obs` is not in your PATH, use the full path to the binary:

- **macOS/Linux**: `"/Users/yourname/go/bin/agentic-obs"`
- **Windows**: `"C:\\Users\\YourName\\go\\bin\\agentic-obs.exe"`

#### Restart Claude Desktop

Close and reopen Claude Desktop completely for the changes to take effect.

#### Verify the connection

In Claude Desktop, you should see the MCP server indicator (typically in the status bar or settings). Claude will now have access to OBS controls!

---

## First Commands to Try

Now that everything is set up, here are 5 example prompts to test your installation:

### 1. Check OBS Status

**Prompt**: "What's the current status of my OBS Studio?"

**Expected**: Claude will report OBS version, current scene, recording/streaming status, and connection health.

---

### 2. List Available Scenes

**Prompt**: "Show me all my OBS scenes"

**Expected**: Claude will list all scenes in your OBS configuration, including the currently active scene.

---

### 3. Switch Scenes

**Prompt**: "Switch to the 'Gaming' scene in OBS"

**Expected**: Your OBS should immediately change to the Gaming scene. Claude will confirm the switch.

---

### 4. Start/Stop Recording

**Prompt**: "Start recording in OBS"

**Expected**: OBS will begin recording. Claude will confirm recording started successfully.

**Follow-up**: "Stop recording" to end the recording.

---

### 5. Toggle Source Visibility

**Prompt**: "Hide my webcam source in the current OBS scene"

**Expected**: The webcam source will be hidden. Claude will confirm the visibility change.

---

## Common First-Time Issues

### Issue 1: "MCP server not found" in Claude Desktop

**Symptoms**: Claude doesn't recognize OBS commands, no MCP indicator visible.

**Solutions**:
- Verify the config file path is correct for your OS
- Check that `agentic-obs` is in your PATH (`which agentic-obs` on macOS/Linux, `where agentic-obs` on Windows)
- Use the full binary path in `claude_desktop_config.json` instead of just "agentic-obs"
- Restart Claude Desktop completely (quit, not just close window)

---

### Issue 2: "Failed to connect to OBS WebSocket"

**Symptoms**: MCP server starts but can't communicate with OBS.

**Solutions**:
- Ensure OBS Studio is running
- Verify WebSocket server is enabled: Tools → WebSocket Server Settings
- Check the port is 4455 (or update your config if different)
- If using a password, make sure it was entered correctly during first run
- Check OBS logs for WebSocket connection attempts: View → Logs → View Current Log

---

### Issue 3: "Authentication failed"

**Symptoms**: Connection attempts but authentication rejected.

**Solutions**:
- Delete `agentic-obs.db` to reset stored credentials
- Run `agentic-obs` manually to re-enter the password
- Verify the password in OBS → Tools → WebSocket Server Settings
- Consider disabling authentication temporarily for testing

---

### Issue 4: "Scene not found" or "Source not found"

**Symptoms**: Commands fail with "not found" errors.

**Solutions**:
- Check exact scene/source names in OBS (case-sensitive!)
- Ask Claude to "list all scenes" or "list all sources" first
- Scene names must match exactly (e.g., "Gaming Setup" vs "Gaming")

---

### Issue 5: Changes work but OBS doesn't update visually

**Symptoms**: Commands succeed but OBS UI doesn't reflect changes.

**Solutions**:
- This is rare but can happen with OBS UI refresh delays
- Click in OBS to trigger a UI refresh
- Check OBS logs to confirm WebSocket commands were received
- Restart OBS if the issue persists

---

## Quick Reference Card

### Essential Commands

| Task | Example Prompt |
|------|----------------|
| Get OBS status | "What's my OBS status?" |
| List scenes | "Show me all OBS scenes" |
| Switch scene | "Switch to [Scene Name]" |
| Start recording | "Start recording in OBS" |
| Stop recording | "Stop recording" |
| Start streaming | "Start streaming on OBS" |
| Stop streaming | "Stop streaming" |
| List sources | "List all sources in the current scene" |
| Hide/show source | "Hide [Source Name]" or "Show [Source Name]" |
| Mute audio | "Mute [Audio Input Name]" |
| Unmute audio | "Unmute [Audio Input Name]" |
| Set volume | "Set [Audio Input] volume to 75%" |

### File Locations

| Item | Location |
|------|----------|
| Binary | `$GOPATH/bin/agentic-obs` (or `%GOPATH%\bin\agentic-obs.exe`) |
| Config DB | `./agentic-obs.db` (created in working directory) |
| Claude Config | `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS)<br>`%APPDATA%\Claude\claude_desktop_config.json` (Windows)<br>`~/.config/Claude/claude_desktop_config.json` (Linux) |

### Important URLs

- **GitHub Repository**: https://github.com/ironystock/agentic-obs
- **OBS WebSocket Protocol**: https://github.com/obsproject/obs-websocket/blob/master/docs/generated/protocol.md
- **MCP Documentation**: https://modelcontextprotocol.io
- **Issue Tracker**: https://github.com/ironystock/agentic-obs/issues

### Default Settings

- **OBS WebSocket Port**: 4455
- **Host**: localhost
- **Transport**: stdio (standard input/output)
- **Database**: SQLite (`agentic-obs.db`)

---

## Next Steps

Congratulations! You now have AI-powered control over OBS Studio. Here's what to explore next:

1. **Experiment with complex workflows**: Try chaining multiple commands ("Start recording, switch to Gaming scene, unmute my mic")

2. **Create scene presets**: Future versions will support saving and recalling scene configurations

3. **Monitor OBS events**: The server automatically detects scene changes and can notify Claude

4. **Read the full documentation**: Check out `README.md` and `CLAUDE.md` for advanced features

5. **Join the community**: Report issues, request features, or contribute on GitHub

---

**Need help?** Open an issue at https://github.com/ironystock/agentic-obs/issues

**Last Updated**: 2025-12-14
**Version**: 1.0.0
