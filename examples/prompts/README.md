# Prompt Examples Guide

This directory contains natural language prompt examples for controlling OBS Studio through AI assistants using agentic-obs.

## What Are These Examples?

Each file contains conversational prompts you can use to interact with OBS through an AI assistant like Claude. These aren't code snippets or API calls - they're actual phrases you would type or speak to the AI.

## Available Prompt Collections

### [scenes.md](scenes.md)
Scene management prompts including:
- Listing available scenes
- Switching between scenes
- Creating new scenes
- Removing scenes

### [recording.md](recording.md)
Recording control prompts for:
- Starting and stopping recordings
- Pausing and resuming
- Checking recording status
- Managing recording state

### [audio.md](audio.md)
Audio control examples:
- Muting and unmuting sources
- Adjusting volume levels
- Checking audio status

### [screenshots.md](screenshots.md)
Visual monitoring examples:
- Setting up AI visual observation
- Verifying scene changes
- Problem detection and diagnosis
- Layout design feedback

### [workflows.md](workflows.md)
Multi-step workflow examples:
- Stream preparation sequences
- Recording session management
- Scene transitions with audio
- Visual verification workflows
- Complete broadcast workflows

## How to Use These Prompts

### Direct Usage
Simply copy a prompt and use it with your AI assistant:

```
You: "Can you show me all my OBS scenes?"
```

### Adaptation
Modify prompts to match your setup:

**Example prompt:** "Switch to my Gaming scene"
**Your adaptation:** "Switch to my Podcast scene"

### Combination
Combine multiple actions:

```
You: "Switch to my Be Right Back scene and mute my microphone"
```

## Prompt Writing Tips

### 1. Be Conversational
❌ "Execute scene change to Gaming"
✅ "Switch to my Gaming scene"

### 2. Provide Context When Needed
❌ "Start it"
✅ "Start recording"

### 3. Use Natural Variations
All of these work:
- "Show me my scenes"
- "What scenes do I have?"
- "List all my OBS scenes"

### 4. Specify Details
❌ "Change the volume"
✅ "Set my microphone volume to 80%"

### 5. Combine Related Actions
✅ "Start recording and switch to my Gaming scene"
✅ "Mute my microphone and turn off my camera source"

## Understanding AI Responses

The AI will typically:

1. **Confirm the action**: "I'll switch to your Gaming scene"
2. **Execute the command**: [Uses OBS tools in the background]
3. **Report results**: "Successfully switched to Gaming scene"
4. **Handle errors**: "I couldn't find a scene named 'Gamng' - did you mean 'Gaming'?"

## Common Patterns

### Checking Before Acting
```
You: "Am I currently recording?"
AI: "No, you're not recording right now."
You: "Okay, start recording"
```

### Sequential Operations
```
You: "First show me my scenes, then switch to the one called Gaming"
```

### Conditional Actions
```
You: "If I'm not recording, please start recording"
```

### Status Queries
```
You: "What's my current scene and am I recording?"
```

## Troubleshooting Prompts

If something doesn't work, try:

1. **Check connection**: "Can you check if OBS is connected?"
2. **List available options**: "What scenes are available?"
3. **Get status**: "What's the current status of OBS?"
4. **Be more specific**: Instead of "turn that on", use "unmute my microphone"

## Advanced Usage

### Referencing Previous Context
```
You: "Show me my scenes"
AI: [Lists: Gaming, Chatting, Be Right Back]
You: "Switch to the second one"
AI: [Switches to Chatting]
```

### Asking for Verification
```
You: "Switch to my Gaming scene and let me know when it's done"
```

### Requesting Information First
```
You: "Before we start, what's my current setup? Which scene am I on and am I recording?"
```

## Best Practices

1. **Start Simple**: Use basic prompts first, then combine actions
2. **Verify State**: Check current status before making changes
3. **Use Clear Names**: Reference scenes and sources by their exact names
4. **Ask Questions**: If unsure, ask the AI what's possible
5. **Iterate**: Refine your prompts based on responses

## Example Conversation Flow

```
You: "Hey, I'm about to start my stream. Can you help me get set up?"

AI: "Of course! I can help you prepare for streaming. What would you like me to do?"

You: "First, show me what scenes I have available"

AI: [Lists your scenes: Starting Soon, Gaming, Chatting, Be Right Back, Ending]

You: "Great. Switch to the Starting Soon scene and make sure I'm not recording yet"

AI: "I've switched to your Starting Soon scene. You're not currently recording."

You: "Perfect. Can you start streaming for me?"

AI: "I've started streaming. You're now live with the Starting Soon scene."

You: "Thanks! I'll let you know when I'm ready to switch to Gaming"
```

## Need More Help?

- Check the specific prompt files for detailed examples
- Refer to the [main examples README](../README.md)
- See the [agentic-obs documentation](../../README.md)

---

**Remember**: These are examples and starting points. Feel free to use your own words and phrasing - the AI is designed to understand natural language!
