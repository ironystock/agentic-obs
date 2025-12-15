# Examples Directory

This directory contains example prompts and workflows for using agentic-obs with AI assistants like Claude.

## What's Here

- **prompts/** - Natural language prompt examples showing how to interact with OBS through AI
  - Scene management
  - Recording control
  - Audio settings
  - Multi-step workflows

## How to Use These Examples

The prompts in this directory are written exactly as you would speak to an AI assistant. Simply copy and adapt them for your needs, or use them as inspiration for your own commands.

### Quick Start

1. Make sure OBS Studio is running with WebSocket enabled (port 4455)
2. Start the agentic-obs MCP server
3. Connect through Claude Desktop or another MCP-compatible client
4. Try any of the example prompts from the `prompts/` directory

### Example Session

```
You: "Can you show me what scenes I have in OBS?"
AI: [Lists your current OBS scenes]

You: "Switch to my Gaming scene"
AI: [Changes to the Gaming scene]

You: "Start recording"
AI: [Begins recording]
```

## Directory Structure

```
examples/
├── README.md              # This file
└── prompts/
    ├── README.md          # Guide to using prompts
    ├── scenes.md          # Scene management examples
    ├── recording.md       # Recording workflow examples
    ├── audio.md           # Audio control examples
    └── workflows.md       # Multi-step workflow examples
```

## Contributing Examples

If you have useful workflows or prompt patterns, consider contributing them! Good examples:

- Use natural, conversational language
- Include context about when to use them
- Show multiple variations of the same request
- Demonstrate real-world scenarios

## Tips for Working with AI + OBS

1. **Be specific**: "Switch to my Gaming scene" is better than "change scene"
2. **Check status first**: Ask "What's my current scene?" before making changes
3. **Combine actions**: "Start recording and switch to my Talk Show scene"
4. **Use natural language**: Talk to the AI like you would a human assistant

## Resources

- [OBS Studio Documentation](https://obsproject.com/wiki/)
- [Model Context Protocol](https://modelcontextprotocol.io)
- [agentic-obs README](../README.md)

---

**Need Help?** Check the individual prompt files for specific examples, or refer to the main project README.
