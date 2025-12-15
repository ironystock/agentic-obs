# MCP Resources & Notifications

This document describes how OBS entities are exposed as MCP resources and how notifications work.

## Overview

**MCP Resources** represent queryable, server-owned state that can change over time. Unlike tools (which perform actions), resources provide read-only access to dynamic data with server-initiated notifications when that data changes.

## Scenes as Resources

### Resource Schema

**URI Pattern:** `obs://scene/{scene_name}`

**Resource Type:** `obs-scene`

**MIME Type:** `application/json`

### Resource Operations

#### 1. List Resources (`resources/list`)

Returns all available scenes as resources.

**Response:**

```json
{
  "resources": [
    {
      "uri": "obs://scene/Gaming",
      "name": "Gaming",
      "mimeType": "application/json",
      "description": "Gaming scene with webcam and game capture"
    },
    {
      "uri": "obs://scene/Chat",
      "name": "Chat",
      "mimeType": "application/json",
      "description": "Chat-focused scene"
    }
  ]
}
```

#### 2. Read Resource (`resources/read`)

Returns detailed configuration for a specific scene.

**Request:**

```json
{
  "uri": "obs://scene/Gaming"
}
```

**Response:**

```json
{
  "contents": [
    {
      "uri": "obs://scene/Gaming",
      "mimeType": "application/json",
      "text": "{\"name\":\"Gaming\",\"isActive\":true,\"sources\":[...]}"
    }
  ]
}
```

**Scene JSON Structure:**

```json
{
  "name": "Gaming",
  "isActive": true,
  "index": 0,
  "sources": [
    {
      "name": "Webcam",
      "type": "input",
      "visible": true
    },
    {
      "name": "Game Capture",
      "type": "input",
      "visible": true
    }
  ]
}
```

## Notification Events

### When to Send Notifications

The MCP server monitors OBS events and sends notifications to keep clients synchronized:

| OBS Event | MCP Notification | Trigger |
|-----------|------------------|---------|
| `SceneCreated` | `notifications/resources/list_changed` | New scene added to OBS |
| `SceneRemoved` | `notifications/resources/list_changed` | Scene deleted from OBS |
| `CurrentProgramSceneChanged` | `notifications/resources/updated` | Active scene switched |
| `SceneItemEnableStateChanged` | `notifications/resources/updated` | Source visibility changed in scene |

### Notification Message Format

#### List Changed Notification

Sent when the list of available scenes changes (scene added/removed).

```json
{
  "jsonrpc": "2.0",
  "method": "notifications/resources/list_changed",
  "params": {}
}
```

**Client Response:** Client should call `resources/list` again to get updated list.

#### Resource Updated Notification

Sent when a specific scene's content or state changes.

```json
{
  "jsonrpc": "2.0",
  "method": "notifications/resources/updated",
  "params": {
    "uri": "obs://scene/Gaming"
  }
}
```

**Client Response:** Client should call `resources/read` for the specified URI to get updated content.

## Implementation Architecture

### Event Flow

```
┌─────────────────┐
│   OBS Studio    │
└────────┬────────┘
         │ WebSocket Event (SceneCreated, etc.)
         ↓
┌─────────────────┐
│  OBS Client     │
│  (events.go)    │
└────────┬────────┘
         │ handleSceneCreated()
         │ handleCurrentProgramSceneChanged()
         ↓
┌─────────────────┐
│  MCP Server     │
│  (server.go)    │
└────────┬────────┘
         │ SendResourceListChanged()
         │ SendResourceUpdated(uri)
         ↓
┌─────────────────┐
│  MCP Client     │
│  (Claude/AI)    │
└─────────────────┘
```

### Code Structure

**`internal/obs/events.go`** - OBS Event Handlers

```go
func (c *Client) handleSceneCreated(event *events.SceneCreated) {
    log.Printf("Scene created: %s", event.SceneName)
    c.mcpServer.SendResourceListChanged()
}

func (c *Client) handleSceneRemoved(event *events.SceneRemoved) {
    log.Printf("Scene removed: %s", event.SceneName)
    c.mcpServer.SendResourceListChanged()
}

func (c *Client) handleCurrentProgramSceneChanged(event *events.CurrentProgramSceneChanged) {
    log.Printf("Scene changed to: %s", event.SceneName)
    uri := fmt.Sprintf("obs://scene/%s", event.SceneName)
    c.mcpServer.SendResourceUpdated(uri)
}
```

**`internal/mcp/resources.go`** - Resource Handlers

```go
func handleResourcesList(ctx context.Context, obsClient *obs.Client) ([]*mcp.Resource, error) {
    scenes, err := obsClient.GetSceneList()
    if err != nil {
        return nil, fmt.Errorf("failed to get scenes: %w", err)
    }

    resources := make([]*mcp.Resource, len(scenes))
    for i, scene := range scenes {
        resources[i] = &mcp.Resource{
            URI:         fmt.Sprintf("obs://scene/%s", scene.Name),
            Name:        scene.Name,
            MimeType:    "application/json",
            Description: fmt.Sprintf("OBS Scene: %s", scene.Name),
        }
    }
    return resources, nil
}

func handleResourceRead(ctx context.Context, uri string, obsClient *obs.Client) (*mcp.ResourceContent, error) {
    // Extract scene name from URI
    sceneName := strings.TrimPrefix(uri, "obs://scene/")

    // Get scene details from OBS
    scene, err := obsClient.GetSceneByName(sceneName)
    if err != nil {
        return nil, fmt.Errorf("failed to get scene %s: %w", sceneName, err)
    }

    // Serialize to JSON
    jsonData, err := json.Marshal(scene)
    if err != nil {
        return nil, fmt.Errorf("failed to serialize scene: %w", err)
    }

    return &mcp.ResourceContent{
        URI:      uri,
        MimeType: "application/json",
        Text:     string(jsonData),
    }, nil
}
```

## Client Interaction Pattern

### Typical Workflow

1. **Initial Query**: Client calls `resources/list` to get all scenes
2. **Server Monitors**: MCP server subscribes to OBS events
3. **State Change**: User switches scene in OBS
4. **Notification**: Server sends `notifications/resources/updated` with scene URI
5. **Client Refresh**: Client calls `resources/read` for that URI to get updated state

### Example AI Interaction

**User:** "What scenes do I have?"

**AI → Server:** `resources/list`

**Server → AI:** `[{uri: "obs://scene/Gaming", ...}, {uri: "obs://scene/Chat", ...}]`

**AI → User:** "You have two scenes: Gaming and Chat."

*[User switches to Gaming scene in OBS]*

**Server → AI:** `notifications/resources/updated` with `uri: "obs://scene/Gaming"`

**AI → Server:** `resources/read` for `obs://scene/Gaming`

**Server → AI:** `{name: "Gaming", isActive: true, ...}`

**AI → User:** "I noticed you switched to the Gaming scene."

## Future Resource Types

### Sources as Resources

**URI Pattern:** `obs://source/{scene_name}/{source_name}`

**Use Cases:**

- Monitor source visibility changes
- Track source settings modifications
- Detect source additions/removals

**Notifications:**

- `notifications/resources/list_changed` - Source added/removed from scene
- `notifications/resources/updated` - Source settings changed

### Audio Inputs as Resources

**URI Pattern:** `obs://audio-input/{input_name}`

**Use Cases:**

- Monitor mute state changes
- Track volume adjustments
- Detect audio device changes

**Notifications:**

- `notifications/resources/updated` - Mute state or volume changed

### Filters as Resources

**URI Pattern:** `obs://filter/{source_name}/{filter_name}`

**Use Cases:**

- Monitor filter enable/disable state
- Track filter settings changes

**Notifications:**

- `notifications/resources/list_changed` - Filter added/removed
- `notifications/resources/updated` - Filter settings changed

## Best Practices

### 1. Resource Naming

- Use descriptive URI patterns with clear hierarchy
- Use URL-safe characters (encode spaces and special chars)
- Keep URIs stable (don't change schema without versioning)

### 2. Notification Frequency

- Batch rapid changes when possible (debouncing)
- Only send notifications for meaningful state changes
- Consider client subscription preferences (future)

### 3. Error Handling

- Return contextual errors for invalid URIs
- Handle OBS disconnections gracefully
- Log notification send failures

### 4. Performance

- Cache scene list to reduce OBS queries
- Use efficient JSON serialization
- Monitor notification queue depth

## Testing Notifications

### Manual Testing

1. Start the MCP server connected to OBS
2. Connect MCP client (e.g., Claude Desktop)
3. Query `resources/list` to get scenes
4. In OBS, create a new scene
5. Verify client receives `notifications/resources/list_changed`
6. In OBS, switch scenes
7. Verify client receives `notifications/resources/updated`

### Automated Testing

```go
func TestSceneCreatedNotification(t *testing.T) {
    // Mock OBS client
    mockOBS := &MockOBSClient{}

    // Create MCP server
    server := mcp.NewServer(mockOBS)

    // Subscribe to notifications
    notifications := make(chan *Notification, 10)
    server.OnNotification(func(n *Notification) {
        notifications <- n
    })

    // Trigger OBS event
    mockOBS.TriggerEvent(&events.SceneCreated{SceneName: "Test"})

    // Verify notification received
    select {
    case n := <-notifications:
        assert.Equal(t, "notifications/resources/list_changed", n.Method)
    case <-time.After(1 * time.Second):
        t.Fatal("Notification not received")
    }
}
```

---

**Document Version:** 1.0
**Last Updated:** 2025-12-14
**Related:** CLAUDE.MD, PROJECT_PLAN.md
