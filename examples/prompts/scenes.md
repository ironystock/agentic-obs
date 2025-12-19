# Scene Management Prompts

Natural language prompts for managing OBS scenes through AI assistants.

**MCP Prompts:** `quick-status`, `scene-organizer`

---

## Listing Scenes

**Tools Used:** `list_scenes`, `get_obs_status`

### Basic Listing

**Show me all scenes:**
- "Can you show me what scenes I have in OBS?"
- "List all my OBS scenes"
- "What scenes are available?"
- "Show me my scene list"
- "What scenes do I have set up?"

**Get current scene:**
- "Which scene am I currently on?"
- "What's my active scene?"
- "Tell me my current scene"
- "Which scene is showing right now?"

### Contextual Listing

**Before switching:**
- "Before I switch, show me all available scenes"
- "Let me see what scenes I can choose from"
- "I want to change scenes - what are my options?"

**For planning:**
- "I'm setting up for a stream - what scenes do I have?"
- "Show me my scenes so I can plan my transitions"

## Switching Scenes

**Tools Used:** `set_current_scene`, `list_scenes`

### Direct Scene Switching

**Switch to a specific scene:**
- "Switch to my Gaming scene"
- "Change to the Chatting scene"
- "Go to my Be Right Back scene"
- "Set the current scene to Starting Soon"
- "I want to use my Ending scene now"

**Switch with context:**
- "I'm starting my stream - switch to the Starting Soon scene"
- "Time for gameplay, switch to Gaming"
- "Going on break, change to Be Right Back"
- "Stream's over, go to my Ending scene"

### Conditional Switching

**With status check:**
- "If I'm on the Starting Soon scene, switch to Gaming"
- "Change to Chatting unless I'm already there"
- "Switch to Gaming and let me know when it's done"

**With confirmation:**
- "Switch to Gaming and confirm which scene is active"
- "Change to my Be Right Back scene and tell me if it worked"

## Creating Scenes

**Tools Used:** `create_scene`

### Basic Scene Creation

**Create a new scene:**
- "Create a new scene called Studio View"
- "I need a new scene named Break Screen"
- "Make a scene called Technical Difficulties"
- "Add a new scene for me called Workout Camera"
- "Can you create a scene named Just Chatting?"

### Creation with Context

**Planning ahead:**
- "I'm setting up for a podcast - create a scene called Podcast Intro"
- "Before my stream, create a new scene named Guest Interview"
- "Make a scene called Music Break for when I play songs"

**With explanation:**
- "Create a scene called Multi-Cam so I can set up multiple camera angles"
- "I need a clean scene with nothing on it - call it Blank Canvas"

### Multiple Scene Creation

**Create several at once:**
- "Create three scenes: Stream Starting, Main Content, and Stream Ending"
- "I need scenes for Opening, Gameplay, and Closing"
- "Make me a set of scenes: Intro, Gaming, Chatting, and Outro"

## Removing Scenes

**Tools Used:** `remove_scene`, `list_scenes`

### Basic Scene Removal

**Remove a specific scene:**
- "Delete my Old Gaming scene"
- "Remove the scene called Test Scene"
- "I don't need the Break Screen scene anymore - delete it"
- "Get rid of my Unused Scene"
- "Can you remove the scene named Backup Camera?"

### Safe Removal

**With confirmation:**
- "Before you delete it, am I currently using the Test Scene?"
- "If I'm not on it, delete my Old Layout scene"
- "Remove the Deprecated scene but make sure I'm on a different one first"

**Checking before removal:**
- "I want to delete my Test Scene - am I currently on it?"
- "Can I safely remove the Old Gaming scene or am I using it?"

### Cleanup Operations

**Multiple removals:**
- "Delete these scenes: Test1, Test2, and Old Layout"
- "Clean up my unused scenes - remove Test Scene and Backup"
- "I want to remove several scenes: Old Intro, Old Outro, and Deprecated Layout"

**Organized cleanup:**
- "Show me all my scenes, then I'll tell you which ones to delete"
- "List my scenes so I can decide what to remove"

## Combined Scene Operations

### List and Switch

**Two-step process:**
- "Show me my scenes, then switch to Gaming"
- "List all scenes and then go to the one called Chatting"
- "What scenes do I have? Then switch to Be Right Back"

### Create and Switch

**Create then activate:**
- "Create a scene called New Layout and switch to it"
- "Make a new scene named Testing and set it as current"
- "Add a scene called Experiment and make it active"

### Check, Create, Switch

**Complete workflow:**
- "Show me my scenes. If I don't have one called Podcast, create it and switch to it"
- "Check if I have a Guest scene - if not, make one and activate it"

### Switch and Confirm

**With verification:**
- "Switch to Gaming and show me all my scenes to confirm"
- "Change to Chatting and tell me which scene is now active"
- "Go to Be Right Back and verify the switch was successful"

## Real-World Scenarios

### Starting a Stream

**Pre-stream setup:**
- "I'm about to go live - switch to my Starting Soon scene"
- "Stream is starting in 5 minutes, get my Starting Soon scene ready"
- "Put me on the pre-stream scene"

### During Stream Transitions

**Content transitions:**
- "We're moving to gameplay now - switch to Gaming scene"
- "Taking questions from chat - change to my Chatting scene"
- "Time for the main event - go to my Main Content scene"

### Taking Breaks

**Break management:**
- "I'm taking a quick break - switch to Be Right Back"
- "Need a bathroom break, put up my BRB scene"
- "Going AFK for a minute, switch to my Break screen"

### Ending Stream

**Stream conclusion:**
- "Stream's wrapping up - switch to my Ending scene"
- "Time for outro, change to my Goodbye scene"
- "We're done for today - go to the end credits scene"

### Technical Issues

**Problem solving:**
- "Something's wrong with my Gaming scene - create a backup one called Gaming2"
- "My main scene is broken - switch to my backup scene"
- "Create an emergency scene called Technical Difficulties"

### Scene Organization

**Setup and cleanup:**
- "Show me all my scenes so I can organize them"
- "I'm reorganizing - create scenes called Morning Stream and Evening Stream"
- "List my scenes and I'll tell you which old ones to delete"

### Testing and Experimentation

**Testing workflows:**
- "Create a test scene so I can try out new layouts"
- "Make a scene called Experiment for testing new sources"
- "I want to test something - create a scene named Test Layout"

## Tips for Scene Prompts

1. **Use exact scene names**: "Switch to Gaming" not "switch to the game one"
2. **Be specific about timing**: "Switch now" vs "Create for later"
3. **Provide context**: Helps the AI understand your workflow
4. **Combine operations**: "List, then switch" is more efficient
5. **Verify results**: Ask for confirmation on critical changes

## Common Variations

All these mean the same thing:
- "Switch to Gaming" = "Change to Gaming" = "Go to Gaming" = "Set scene to Gaming"
- "Show scenes" = "List scenes" = "What scenes exist" = "Display my scenes"
- "Create scene" = "Make a scene" = "Add a scene" = "New scene"
- "Delete scene" = "Remove scene" = "Get rid of scene" = "Eliminate scene"

---

**Next Steps**: Check out [recording.md](recording.md) for recording control prompts, or [workflows.md](workflows.md) for complete multi-step examples.
