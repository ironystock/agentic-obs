# HTTP API Reference

The agentic-obs server includes an optional HTTP API for monitoring and configuration. This API is enabled by default and runs on `http://localhost:8765`.

## Configuration

The HTTP server can be configured via the `/api/config` endpoint or through first-run setup:

| Setting | Default | Description |
|---------|---------|-------------|
| `enabled` | `true` | Enable/disable the HTTP server |
| `host` | `localhost` | Bind address (localhost, 127.0.0.1, or 0.0.0.0) |
| `port` | `8765` | Port number (1024-65535) |

## Endpoints

### GET /api/status

Returns server status information.

**Response:**
```json
{
  "status": "ok",
  "server_name": "agentic-obs",
  "version": "0.1.0",
  "uptime": "2h15m30s",
  "http_address": "localhost:8765"
}
```

---

### GET /api/history

Returns action history (MCP tool call log).

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `limit` | int | 50 | Maximum records to return (1-500) |
| `tool` | string | - | Filter by tool name |

**Example:** `GET /api/history?limit=10&tool=set_current_scene`

**Response:**
```json
{
  "count": 2,
  "limit": 10,
  "actions": [
    {
      "id": 42,
      "action": "Set current scene",
      "tool_name": "set_current_scene",
      "input": "{\"scene_name\":\"Gaming\"}",
      "output": "{\"message\":\"Successfully switched to scene: Gaming\"}",
      "success": true,
      "duration_ms": 15,
      "created_at": "2025-01-15T10:30:00Z"
    }
  ]
}
```

---

### GET /api/history/stats

Returns aggregate statistics for action history.

**Response:**
```json
{
  "total_actions": 156,
  "successful_actions": 150,
  "failed_actions": 6,
  "avg_duration_ms": 23.5,
  "top_tools": [
    {"tool_name": "list_scenes", "count": 45},
    {"tool_name": "set_current_scene", "count": 32},
    {"tool_name": "get_obs_status", "count": 28}
  ]
}
```

---

### GET /api/screenshots

Returns list of configured screenshot sources with URLs.

**Response:**
```json
{
  "count": 1,
  "sources": [
    {
      "id": 1,
      "name": "main_view",
      "source_name": "Scene 1",
      "cadence_ms": 5000,
      "image_format": "png",
      "enabled": true,
      "url": "http://localhost:8765/screenshot/main_view",
      "created_at": "2025-01-15T09:00:00Z"
    }
  ]
}
```

---

### GET /api/config

Returns current server configuration (excluding sensitive data).

**Response:**
```json
{
  "obs": {
    "host": "localhost",
    "port": 4455
  },
  "tool_groups": {
    "core": true,
    "visual": true,
    "layout": true,
    "audio": true,
    "sources": true
  },
  "web_server": {
    "enabled": true,
    "host": "localhost",
    "port": 8765
  }
}
```

**Note:** OBS password is intentionally omitted for security.

---

### POST /api/config

Updates server configuration. Changes take effect after server restart.

**Request Body:**
```json
{
  "tool_groups": {
    "core": true,
    "visual": false,
    "layout": true,
    "audio": true,
    "sources": true
  },
  "web_server": {
    "enabled": true,
    "host": "localhost",
    "port": 9000
  }
}
```

**Validation Rules:**

| Field | Constraint | Error |
|-------|------------|-------|
| `web_server.host` | Must be `localhost`, `127.0.0.1`, or `0.0.0.0` | 400 Bad Request |
| `web_server.port` | Must be 1024-65535 | 400 Bad Request |

**Response (Success):**
```json
{
  "status": "ok",
  "message": "Configuration updated. Restart server for changes to take effect."
}
```

**Response (Validation Error):**
```json
{
  "error": "Invalid host: must be localhost, 127.0.0.1, or 0.0.0.0"
}
```

---

### GET /screenshot/{name}

Returns the latest screenshot image for a configured source.

**Path Parameters:**
| Parameter | Description |
|-----------|-------------|
| `name` | Screenshot source name |

**Response:**
- **Content-Type:** `image/png` or `image/jpeg`
- **Body:** Binary image data

**Error Response (404):**
```json
{
  "error": "Screenshot not found"
}
```

---

## Error Responses

All endpoints return JSON error responses for failures:

| Status Code | Description |
|-------------|-------------|
| 400 | Bad Request - Invalid input or validation failure |
| 404 | Not Found - Resource does not exist |
| 405 | Method Not Allowed - Wrong HTTP method |
| 500 | Internal Server Error - Server-side failure |

**Error Format:**
```json
{
  "error": "Description of the error"
}
```

---

## Security Considerations

1. **Local-only binding:** By default, the server binds to `localhost` only
2. **Host validation:** Only `localhost`, `127.0.0.1`, and `0.0.0.0` are allowed
3. **Port restrictions:** Only non-privileged ports (1024-65535) are allowed
4. **Password protection:** OBS WebSocket password is never exposed via API
5. **Request size limits:** POST bodies are limited to 64KB

---

## Dashboard

The HTTP server includes a web dashboard accessible at the root URL:

- **URL:** `http://localhost:8765/`
- **Features:**
  - Server status display
  - Action history viewer
  - Screenshot source list
  - Configuration management

---

## API Summary

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Web dashboard |
| `/api/status` | GET | Server status |
| `/api/history` | GET | Action history |
| `/api/history/stats` | GET | History statistics |
| `/api/screenshots` | GET | Screenshot sources |
| `/api/config` | GET | Get configuration |
| `/api/config` | POST | Update configuration |
| `/screenshot/{name}` | GET | Get screenshot image |

**Total: 8 HTTP API endpoints**
