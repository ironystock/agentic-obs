# JSON-RPC API Examples

Technical examples of MCP tool calls using the JSON-RPC 2.0 protocol.

## What Are These Examples?

These files contain request/response pairs showing how to interact with the agentic-obs MCP server programmatically. Each example includes:

- **Request**: The JSON-RPC 2.0 formatted tool call
- **Response**: Expected successful response
- **Notes**: Prerequisites, error scenarios, and usage tips

## Who Should Use These?

- Developers building MCP clients
- Testing and debugging tool implementations
- Understanding the technical protocol
- Building automation scripts

## Example Files

- **recording.json** - Recording control with pause/resume
- **scenes.json** - Scene management operations
- **streaming.json** - Streaming control
- **sources.json** - Source management
- **audio.json** - Audio control

## JSON-RPC 2.0 Format

All MCP tool calls follow the JSON-RPC 2.0 specification:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "tool_name",
    "arguments": {
      "param1": "value1"
    }
  }
}
```

## MCP Tool Call Structure

### Request Format
```json
{
  "jsonrpc": "2.0",
  "id": <number>,
  "method": "tools/call",
  "params": {
    "name": "<tool_name>",
    "arguments": {
      // Tool-specific parameters
    }
  }
}
```

### Response Format
```json
{
  "jsonrpc": "2.0",
  "id": <number>,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "<JSON or text response>"
      }
    ]
  }
}
```

### Error Response
```json
{
  "jsonrpc": "2.0",
  "id": <number>,
  "error": {
    "code": -32000,
    "message": "<error description>"
  }
}
```

## Using These Examples

### Testing with curl
```bash
echo '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"list_scenes","arguments":{}}}' | ./agentic-obs
```

### Testing with Node.js
```javascript
const stdin = process.stdin;
const request = {
  jsonrpc: "2.0",
  id: 1,
  method: "tools/call",
  params: {
    name: "list_scenes",
    arguments: {}
  }
};
console.log(JSON.stringify(request));
```

## Tips

1. **IDs are sequential** - Increment for each request
2. **Empty arguments** - Use `{}` not omit for no-parameter tools
3. **Tool names** - Must match exactly (snake_case)
4. **Responses are wrapped** - Content is in `.result.content[0].text`

## More Information

- See [TOOLS.md](../../docs/TOOLS.md) for detailed tool documentation
- See [prompts/](../prompts/) for natural language examples
- Check [OBS WebSocket Protocol](https://github.com/obsproject/obs-websocket) for backend details

---

For conversational AI examples, see the [prompts directory](../prompts/).
