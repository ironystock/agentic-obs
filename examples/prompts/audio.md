# Audio Control Prompts

Natural language prompts for managing OBS audio through AI assistants.

**MCP Prompts:** `audio-check`
**Tools Used:** `get_input_mute`, `toggle_input_mute`, `set_input_volume`, `get_input_volume`

---

## Checking Mute Status

**Tools Used:** `get_input_mute`

### Basic Mute Checks

**Check if muted:**
- "Is my microphone muted?"
- "Check if my mic is muted"
- "Am I muted right now?"
- "Tell me if my microphone is on or off"
- "What's the mute status of my mic?"

**Check specific sources:**
- "Is my Desktop Audio muted?"
- "Check mute status for my Blue Yeti microphone"
- "Tell me if Game Audio is muted"
- "Is my Music source muted or unmuted?"

### Before Going Live

**Pre-stream checks:**
- "Before I go live, is my microphone muted?"
- "Check that my mic is unmuted before we start"
- "Make sure my microphone isn't muted"
- "Verify my mic is on before the stream"

## Muting Audio Sources

**Tools Used:** `toggle_input_mute`

### Basic Muting

**Mute microphone:**
- "Mute my microphone"
- "Turn off my mic"
- "Can you mute my microphone?"
- "Silence my mic"
- "Mute the mic please"

**Mute other sources:**
- "Mute my Desktop Audio"
- "Turn off Game Audio"
- "Mute the Music source"
- "Silence my Browser Audio"

### Contextual Muting

**During breaks:**
- "Going on break - mute my microphone"
- "BRB, mute my mic"
- "Stepping away, turn off my microphone"

**For interruptions:**
- "Someone's at the door - mute my mic quick"
- "Need to cough - mute microphone"
- "Phone ringing - mute my mic"

**During content:**
- "Playing copyrighted music - mute Desktop Audio"
- "Loud game cutscene coming - mute Game Audio"
- "Taking a call - mute everything"

### Emergency Muting

**Quick mute:**
- "Mute everything!"
- "Kill all audio!"
- "Emergency mute on my microphone"
- "Mute my mic now!"

## Unmuting Audio Sources

### Basic Unmuting

**Unmute microphone:**
- "Unmute my microphone"
- "Turn my mic back on"
- "Can you unmute my microphone?"
- "Enable my mic"
- "Turn on my microphone"

**Unmute other sources:**
- "Unmute Desktop Audio"
- "Turn Game Audio back on"
- "Unmute the Music source"
- "Enable Browser Audio"

### Contextual Unmuting

**After breaks:**
- "Back from break - unmute my microphone"
- "I'm back, turn my mic on"
- "Returning, unmute my mic"

**After interruptions:**
- "All clear - unmute my microphone"
- "Done with that call - unmute mic"
- "Back to normal - turn on my microphone"

### Unmute Before Important Moments

**Pre-content unmute:**
- "About to speak - make sure my mic is unmuted"
- "Starting the show - unmute my microphone"
- "Before I begin, unmute my mic please"

## Toggle Muting

### Basic Toggle

**Simple toggle:**
- "Toggle my microphone mute"
- "Flip the mute on my mic"
- "Switch my microphone mute status"
- "Toggle mute for my mic"

**Toggle other sources:**
- "Toggle Desktop Audio mute"
- "Switch Game Audio mute"
- "Flip the Music source mute"

### Smart Toggle

**Conditional toggle:**
- "If my mic is muted, unmute it. If it's unmuted, mute it"
- "Toggle my microphone - I don't know what state it's in"
- "Switch my mic to the opposite of whatever it is now"

## Volume Control

### Setting Volume Levels

**Set specific volume:**
- "Set my microphone volume to 80%"
- "Change my mic volume to 100%"
- "Set Desktop Audio volume to 50%"
- "Put my microphone at 75% volume"
- "Adjust my mic to 90%"

**Set to common levels:**
- "Set my microphone to full volume"
- "Put my mic at half volume"
- "Max out my microphone volume"
- "Set my mic to 0% volume"

### Adjusting Volume

**Increase volume:**
- "Turn up my microphone"
- "Increase my mic volume to 85%"
- "Make my microphone louder - set it to 95%"
- "Raise my mic volume to 100%"

**Decrease volume:**
- "Turn down my microphone"
- "Lower my mic volume to 60%"
- "Make my microphone quieter - set it to 40%"
- "Reduce my mic volume to 50%"

### Contextual Volume Changes

**For different content:**
- "I'm whispering - turn my mic up to 100%"
- "Loud section coming - lower my mic to 70%"
- "Game is quiet - set Game Audio to 80%"
- "Music is too loud - set it to 30%"

## Combined Audio Operations

### Mute and Volume Together

**Mute then adjust:**
- "Mute my microphone and set the volume to 75% for when I unmute"
- "Turn off my mic and adjust volume to 80%"

**Unmute with volume:**
- "Unmute my microphone and make sure it's at 85%"
- "Turn on my mic at 90% volume"

### Multiple Source Control

**Mute multiple sources:**
- "Mute my microphone and Desktop Audio"
- "Turn off both my mic and Game Audio"
- "Mute everything except my microphone"

**Unmute multiple sources:**
- "Unmute my microphone and Desktop Audio"
- "Turn on both my mic and Music"

### Scene Change with Audio

**Coordinated audio and scene:**
- "Switch to BRB scene and mute my microphone"
- "Change to Gaming scene and unmute Game Audio"
- "Go to Chatting scene and make sure my mic is unmuted"

## Audio Source Management

### Listing Audio Sources

**See available sources:**
- "What audio sources do I have?"
- "List all my audio inputs"
- "Show me my audio sources"
- "What microphones and audio do I have in OBS?"

### Checking Source Settings

**Get source information:**
- "What are the settings for my microphone?"
- "Show me my Blue Yeti settings"
- "Tell me about my Desktop Audio configuration"
- "What's configured for my Game Audio?"

### Source-Specific Control

**Control by name:**
- "Mute the source named 'Blue Yeti'"
- "Set 'Desktop Audio' volume to 60%"
- "Unmute my 'Music Player' source"
- "Toggle mute on 'Discord Audio'"

## Real-World Scenarios

### Starting a Stream

**Pre-stream audio check:**
- "Before I go live, check that my microphone is unmuted and at 85%"
- "Make sure my mic is on and Desktop Audio is at 40%"
- "Pre-stream check: is my microphone unmuted?"

### During Gameplay

**Game audio management:**
- "Game's getting loud - lower Desktop Audio to 30%"
- "Cutscene is playing - mute Game Audio for a minute"
- "Gameplay starting - unmute Game Audio and set to 50%"

### During Breaks

**Break audio setup:**
- "Going on break - mute my microphone and play music"
- "BRB - mute my mic and set Music to 60%"
- "Taking a break, mute everything except Music"

**Returning from break:**
- "Back from break - unmute my mic and mute the Music"
- "I'm back - mic on, music off"

### Podcast Recording

**Guest management:**
- "Guest is talking - lower my mic to 70% and theirs to 90%"
- "My turn to speak - my mic to 90%, guest to 60%"

**Recording preparation:**
- "Starting podcast - both mics unmuted and at 85%"
- "Check that my microphone and guest microphone are both on"

### Music Stream

**Music playback:**
- "Playing music - mute my mic and set Music to 100%"
- "Song's over - unmute my mic and mute Music"
- "Music break - mic off, music on at full volume"

**Commentary during music:**
- "Talking over music - my mic to 90%, music to 25%"

### Technical Troubleshooting

**Audio testing:**
- "Testing my mic - is it unmuted and at 100%?"
- "Can't hear myself - what's my microphone volume at?"
- "Audio seems off - check all my audio source settings"

**Quick fixes:**
- "My mic is too quiet - set it to 100%"
- "Game is too loud - lower Game Audio to 20%"
- "Everything is too loud - set all volumes to 50%"

### Multi-Source Management

**Balancing audio:**
- "Set my mic to 85%, Game Audio to 40%, and Music to 20%"
- "Balance my audio: microphone at 90%, everything else at 30%"

**Selective muting:**
- "Mute everything except my microphone"
- "Mute all game audio but keep my mic on"
- "Turn off Desktop Audio and Music, keep mic unmuted"

## Advanced Audio Scenarios

### Conditional Audio Control

**Status-based actions:**
- "If my mic is muted, unmute it and set to 85%"
- "Check if Desktop Audio is muted - if so, unmute and set to 50%"

### Audio Profiles

**Quick setups:**
- "Set up gaming audio: mic 85%, game 40%, music off"
- "Chatting setup: mic 90%, desktop 30%, everything else muted"
- "Music stream setup: mic off, music 100%, game off"

### Emergency Audio

**Panic controls:**
- "Mute everything right now!"
- "Kill all audio immediately!"
- "Emergency: turn off all sound sources!"

**Recovery:**
- "Unmute my microphone only"
- "Restore normal audio: mic 85%, desktop 40%"

## Tips for Audio Prompts

1. **Be specific about source names**: Use exact names from OBS
2. **Include volume levels**: "Set to 80%" is clearer than "turn up"
3. **Check before important moments**: Verify mute status before going live
4. **Use percentages**: "85%" is more precise than "loud"
5. **Combine related actions**: "Mute mic and lower game audio" is efficient

## Common Variations

All these mean the same thing:
- "Mute my mic" = "Turn off my microphone" = "Silence my mic"
- "Unmute my mic" = "Turn on my microphone" = "Enable my mic"
- "Set volume to 80%" = "Adjust volume to 80%" = "Change volume to 80%"
- "Toggle mute" = "Switch mute" = "Flip mute status"

## Common Audio Source Names

Typical audio sources you might have:
- Microphone, Mic, Blue Yeti, AT2020, etc.
- Desktop Audio, System Audio
- Game Audio, Game Capture Audio
- Music, Music Player, Spotify
- Browser Audio, Chrome Audio
- Discord, Voice Chat
- Alerts, Notifications

---

**Next Steps**: Check out [workflows.md](workflows.md) for complete multi-step workflows that include audio control, or [scenes.md](scenes.md) for scene management.
