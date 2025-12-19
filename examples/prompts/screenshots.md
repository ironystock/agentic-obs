# Screenshot Prompts

Natural language examples for AI-powered visual monitoring of your OBS stream.

**MCP Prompts:** `visual-check`, `visual-setup`, `problem-detection`
**Tools Used:** `create_screenshot_source`, `remove_screenshot_source`, `list_screenshot_sources`, `configure_screenshot_cadence`

---

## Why Screenshots?

Screenshots give your AI assistant **eyes** on your stream. Instead of blindly executing commands, the AI can:
- **See** what your stream actually looks like
- **Verify** changes were applied correctly
- **Detect** problems like black screens or missing overlays
- **Provide feedback** on layouts and composition

---

## Setting Up Screenshot Capture

**Tools Used:** `create_screenshot_source`

### Basic Setup

**Creating your first screenshot source:**
- "Set up a screenshot capture of my stream so you can see what I'm broadcasting"
- "Create a screenshot source called 'stream-view' that captures every 5 seconds"
- "I want you to be able to see my OBS output - set that up for me"
- "Enable visual monitoring of my current scene"

**Quick monitoring setup:**
- "Start monitoring my stream visually"
- "I need you to be able to see what's on screen - create a screenshot capture"
- "Set up a way to check what my stream looks like"

### Customized Setup

**With specific settings:**
- "Create a screenshot source called 'hd-monitor' at 1920x1080, capturing as PNG every 10 seconds"
- "Set up fast screenshot capture at 2-second intervals called 'quick-check' using JPG format"
- "Create a low-res monitoring source at 640x360 for quick status checks"
- "I need a high-quality screenshot source for my Gaming scene - call it 'game-capture'"

**For specific scenes:**
- "Create a screenshot source that watches my BRB scene"
- "Set up monitoring specifically for my Intro scene"
- "I want to capture screenshots of my Webcam-Only scene"

---

## Checking Your Stream

### Visual Verification

**After making changes:**
- "I just switched scenes - take a screenshot and tell me if it looks right"
- "Can you see my stream? How does it look?"
- "Check my current output and describe what you see"
- "Take a screenshot and verify everything is in place"

**Confirming specific elements:**
- "Is my webcam visible in the current scene?"
- "Can you see my chat overlay?"
- "Is my game capture showing anything?"
- "Check if my alerts are displaying correctly"
- "Is the donation goal visible?"

**Quality checks:**
- "How does my stream quality look?"
- "Take a screenshot - is everything rendering clearly?"
- "Check if my overlays look sharp or blurry"
- "Does my layout look professional?"

### Problem Detection

**When something seems wrong:**
- "Viewers are saying something looks off - can you check?"
- "My game capture might be broken - take a look"
- "Is my stream showing a black screen?"
- "Something's wrong with my layout - what is it?"

**Specific issue checking:**
- "Is my webcam frozen or is it live?"
- "Check if my browser source is rendering"
- "Is there a black bar on my stream?"
- "Are any of my sources showing error messages?"

---

## Stream Monitoring

### Active Monitoring

**During streams:**
- "Keep an eye on my stream and let me know if you notice any issues"
- "Monitor my output and alert me if anything looks wrong"
- "Check my stream periodically while I play"
- "Watch for any visual problems during this session"

**Periodic checks:**
- "How does my stream look right now?"
- "Quick visual check - anything I should know about?"
- "Take a look at my current output"
- "Status check on my stream appearance"

### Pre-Stream Checks

**Before going live:**
- "I'm about to stream - take a screenshot and make sure everything looks good"
- "Pre-stream check: how does my Starting Soon scene look?"
- "Verify my stream layout before I go live"
- "Screenshot my current setup and tell me if I'm ready to stream"

**Comprehensive pre-flight:**
- "Run a full visual check: are all my overlays visible, is my webcam working, and does the layout look professional?"
- "Before I start streaming, check my Gaming scene, BRB scene, and Outro scene - do they all look correct?"

---

## Scene Verification

### After Scene Changes

**Confirming switches:**
- "I just switched to my Gaming scene - does it look right?"
- "Verify the BRB scene is showing correctly"
- "Take a screenshot after the scene change"
- "Did the transition to Chatting scene work properly?"

**Layout verification:**
- "I rearranged my sources - how does it look now?"
- "Check if my webcam moved to the right position"
- "Is my overlay now in the bottom left corner?"
- "Verify the chat box is properly sized"

### Multi-Scene Reviews

**Comparing scenes:**
- "Take screenshots of all my scenes so I can review them"
- "Show me what each of my scenes looks like"
- "Compare my Gaming scene to my Chatting scene - are they consistent?"
- "Cycle through my scenes and check each one"

---

## Design Feedback

### Layout Feedback

**Getting AI input on design:**
- "How does my current layout look? Any suggestions?"
- "Take a screenshot and give me design feedback"
- "Is my scene composition balanced?"
- "Does my stream layout look professional?"

**Specific design questions:**
- "Is my webcam too big or too small?"
- "Is there too much empty space in my scene?"
- "Do my colors work well together?"
- "Is the text on my overlays readable?"

### Improvement Suggestions

**Asking for recommendations:**
- "Look at my Gaming scene - how could I improve it?"
- "What would make my stream look more professional?"
- "Any suggestions for my layout?"
- "How could I make better use of screen space?"

---

## Managing Screenshot Sources

### Listing Sources

**Checking what's configured:**
- "What screenshot sources do I have set up?"
- "Show me all my visual monitoring configurations"
- "List my screenshot capture sources"
- "What streams am I currently monitoring?"

### Modifying Sources

**Adjusting capture frequency:**
- "Make my screenshot capture faster - every 2 seconds"
- "Slow down the monitoring to every 30 seconds"
- "Change stream-view to capture every 10 seconds"
- "Speed up the screenshot rate for quick-check"

### Removing Sources

**Cleaning up:**
- "Remove the test screenshot source"
- "Stop the screenshot monitoring for 'old-monitor'"
- "Delete all my screenshot sources"
- "Turn off visual monitoring"

---

## Advanced Use Cases

### Multi-Source Monitoring

**Watching multiple things:**
- "Set up screenshot sources for my Gaming, BRB, and Chatting scenes"
- "Create monitors for each of my main scenes"
- "I want to track multiple views - set up sources for each"

### Recording Preparation

**Pre-recording checks:**
- "I'm about to record a tutorial - verify my scene looks professional"
- "Screenshot my recording setup and make sure nothing is out of place"
- "Check that my screen recording area looks correct"
- "Verify there are no test elements visible before I record"

### Troubleshooting Sessions

**Debugging visual issues:**
- "Take rapid screenshots to help me debug this flickering issue"
- "Set up 1-second capture so we can analyze what's happening"
- "Monitor closely while I adjust settings"
- "Capture before and after I make this change"

---

## Real-World Scenarios

### Starting a Stream

```
You: "I'm about to go live. Set up visual monitoring and check my Starting Soon scene"
AI: [Creates screenshot source, captures image]
AI: "Your Starting Soon scene looks good - I can see your countdown timer, your logo is centered, and the animated background is running smoothly."

You: "Perfect. Switch to Gaming scene and verify that too"
AI: [Switches scene, captures new screenshot]
AI: "Gaming scene verified. Your webcam is in the bottom right, game capture is active showing your game, and chat overlay is visible on the left side."
```

### Diagnosing a Problem

```
You: "Chat says my webcam is frozen - can you check?"
AI: [Captures screenshot]
AI: "Looking at your stream, your webcam does appear static. The game is moving but your webcam shows the same frame as the last capture 30 seconds ago."

You: "Okay, I'll restart it. Tell me when it's working again"
[You restart webcam source]
AI: [Captures new screenshot after a few seconds]
AI: "Your webcam is now showing live movement. I can see you waving - the feed is working normally again."
```

### Design Session

```
You: "I'm redesigning my Gaming scene. Set up frequent screenshots so you can help"
AI: [Creates 2-second capture source]

You: "I moved my webcam - how does it look?"
AI: "The webcam is now in the upper left. It overlaps slightly with your health bar HUD. You might want to shift it down about 50 pixels or move it to the upper right."

You: "Good point. I moved it. Better?"
AI: "Much better. The webcam is now clear of the game UI, and the position feels more balanced with your chat overlay on the opposite side."
```

---

## Tips for Effective Visual Monitoring

1. **Use appropriate intervals**: 5-10 seconds for general monitoring, 1-2 seconds for active debugging
2. **Create named sources**: Use descriptive names like 'stream-monitor' or 'game-view'
3. **Check after changes**: Always verify visual changes with a screenshot
4. **Describe what you're looking for**: "Check if my webcam is visible" is more useful than "take a screenshot"
5. **Use for pre-flight checks**: Verify everything looks right before going live
6. **Clean up when done**: Remove temporary screenshot sources after debugging sessions

---

**Related Prompts**:
- [Scene Management](scenes.md) - Switching and managing scenes
- [Workflows](workflows.md) - Complete multi-step processes
- [Audio Control](audio.md) - Managing audio sources

**Documentation**: [Screenshot Feature Guide](../../docs/SCREENSHOTS.md)
