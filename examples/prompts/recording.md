# Recording Control Prompts

Natural language prompts for managing OBS recordings through AI assistants.

**MCP Prompts:** `recording-workflow`
**Tools Used:** `start_recording`, `stop_recording`, `get_recording_status`, `pause_recording`, `resume_recording`

---

## Starting Recordings

**Tools Used:** `start_recording`, `get_recording_status`

### Basic Start Commands

**Simple start:**
- "Start recording"
- "Begin recording"
- "Start a recording now"
- "Can you start recording for me?"
- "Hit record"

### Contextual Starts

**With scene context:**
- "Switch to my Gaming scene and start recording"
- "Start recording on my current scene"
- "I'm on the Podcast scene - start recording now"

**With timing context:**
- "I'm ready to go - start recording"
- "Content is starting, begin recording"
- "We're live, start the recording"
- "Starting in 3, 2, 1 - hit record"

### Pre-Start Checks

**Verify before recording:**
- "Am I already recording? If not, start now"
- "Check if recording is active, then start if it isn't"
- "Make sure I'm not recording, then start a new one"

**Status-aware starts:**
- "What's my recording status? Then start if I'm not recording"
- "If recording is stopped, please start it"

## Stopping Recordings

**Tools Used:** `stop_recording`

### Basic Stop Commands

**Simple stop:**
- "Stop recording"
- "End the recording"
- "Stop recording now"
- "Can you stop the recording?"
- "Finish recording"

### Contextual Stops

**With completion context:**
- "That's a wrap - stop recording"
- "Content is done, end the recording"
- "Stream's over, stop recording"
- "We're finished, stop the recording"

### Safe Stops

**Verify before stopping:**
- "Am I recording? If so, stop it"
- "Check recording status and stop if active"
- "Stop recording if I'm currently recording"

**With confirmation:**
- "Stop recording and confirm it stopped"
- "End the recording and let me know when it's done"

## Pausing Recordings

### Basic Pause Commands

**Simple pause:**
- "Pause the recording"
- "Pause recording"
- "Hold the recording"
- "Can you pause recording for me?"
- "Put recording on hold"

### Contextual Pauses

**During breaks:**
- "Taking a break - pause the recording"
- "Need to pause, someone's at the door"
- "Hold on, pause the recording for a minute"
- "Pause recording while I deal with this"

**For interruptions:**
- "Someone's calling - pause recording quick"
- "Need to answer this - pause the recording"
- "Emergency pause on the recording"

### Pause with Context

**Explaining why:**
- "I need to sneeze - pause recording"
- "Doorbell rang, pause the recording"
- "Taking a water break, pause recording"
- "Technical issue, pause recording while I fix it"

## Resuming Recordings

### Basic Resume Commands

**Simple resume:**
- "Resume recording"
- "Continue recording"
- "Unpause the recording"
- "Can you resume recording?"
- "Keep recording"

### Contextual Resumes

**After breaks:**
- "Break's over - resume recording"
- "Back from break, continue recording"
- "I'm ready again - resume recording"
- "All good now, unpause the recording"

**After interruptions:**
- "Interruption handled - resume recording"
- "Issue fixed, continue recording"
- "Back to it - resume recording"

### Resume with Countdown

**Preparation before resume:**
- "Give me a second, then resume recording"
- "Resume recording in a moment"
- "I'm ready - resume recording now"

## Checking Recording Status

### Basic Status Checks

**Simple status queries:**
- "Am I recording?"
- "What's my recording status?"
- "Is recording active?"
- "Check if I'm recording"
- "Tell me the recording status"

### Detailed Status Checks

**Full status information:**
- "Give me full recording status - am I recording, paused, or stopped?"
- "What's the current state of recording?"
- "Tell me everything about the recording status"
- "Am I recording, and if so, for how long?"

### Contextual Status Checks

**Before actions:**
- "Before I start, am I already recording?"
- "Check recording status before I begin"
- "Make sure I'm not recording before we start"

**During stream:**
- "Quick check - is recording running?"
- "Verify recording is active"
- "Confirm I'm still recording"

## Combined Recording Operations

### Start with Scene Change

**Coordinated start:**
- "Switch to Gaming and start recording"
- "Change to Podcast scene and begin recording"
- "Go to my Main Content scene and hit record"

### Stop with Scene Change

**Coordinated stop:**
- "Stop recording and switch to my Ending scene"
- "End recording and go to Be Right Back"
- "Finish recording and change to my Idle scene"

### Status Check Before Action

**Verify then act:**
- "Check if I'm recording, then stop if I am"
- "Am I recording? If not, please start"
- "Verify recording status, then pause if it's running"

### Pause and Resume Workflow

**Complete pause cycle:**
- "Pause recording while I fix my camera, then I'll tell you to resume"
- "Pause now, I'll let you know when to resume"
- "Hold recording, I need a minute"

## Recording Workflows

### Session Start

**Beginning a recording session:**
- "I'm starting my recording session - switch to Intro and start recording"
- "New video starting - get recording going"
- "Let's begin: start recording on my current scene"

**Pre-session check:**
- "Before I start, check that I'm not already recording"
- "Status check: am I recording or are we good to start?"
- "Make sure recording is off before we begin fresh"

### Session End

**Ending a recording session:**
- "That's everything - stop recording"
- "Video is complete, end the recording"
- "We got it all, stop recording now"

**Post-session confirmation:**
- "Stop recording and confirm it's saved"
- "End recording and let me know the status"

### Multi-Segment Recording

**Pausing between segments:**
- "End of segment one - pause recording"
- "Taking a break between takes - pause recording"
- "Pause while I set up the next segment"

**Resuming next segment:**
- "Ready for segment two - resume recording"
- "Next take is ready, continue recording"
- "All set for the next part, unpause recording"

### Emergency Procedures

**Emergency stop:**
- "Stop everything - end recording now!"
- "Cancel this take, stop recording"
- "No good - stop recording immediately"

**Recovery:**
- "That was a false alarm - start recording again"
- "Okay, we're good - resume recording"

## Real-World Scenarios

### Podcast Recording

**Starting podcast:**
- "We're about to start the podcast - switch to Podcast scene and start recording"
- "Intro is done, hit record for the main episode"

**During podcast:**
- "Guest is having tech issues - pause recording"
- "Back with the guest - resume recording"

**Ending podcast:**
- "That's a wrap on the episode - stop recording"

### Tutorial/Educational Content

**Starting tutorial:**
- "Beginning the tutorial - start recording on my Tutorial scene"
- "Ready to teach - start recording"

**Between sections:**
- "Pause recording while I set up the next demo"
- "Next section is ready - resume recording"

**Completing tutorial:**
- "Tutorial complete - stop recording and switch to Outro scene"

### Gaming Content

**Starting gameplay:**
- "About to start the game - switch to Gaming and start recording"
- "Game loaded, begin recording"

**Mid-game pause:**
- "Loading screen - pause recording for a sec"
- "Back in game - resume recording"

**Ending session:**
- "Game session over - stop recording"

### Live Stream Recording

**Stream start:**
- "Starting stream - begin recording"
- "We're live - make sure recording is running"

**Stream break:**
- "On break screen - pause recording or keep going?"
- "BRB for a minute - should I pause recording?"

**Stream end:**
- "Stream ending - stop recording"
- "That's all folks - end the recording"

### Multiple Takes

**First take:**
- "Take one - start recording"

**Between takes:**
- "Cut! Stop recording for a moment"
- "That wasn't good - stop and we'll start fresh"

**New take:**
- "Take two - start recording again"

**Final take:**
- "Perfect take - stop recording, we got it"

## Recording + Streaming Combined

### Dual Operation

**Start both:**
- "Start recording and start streaming"
- "Go live and make sure recording is on"
- "Begin both streaming and recording"

**Stop both:**
- "Stop streaming and stop recording"
- "End everything - stream and recording"
- "We're done - stop both stream and recording"

**Status check both:**
- "Am I recording and streaming?"
- "What's the status of both recording and streaming?"
- "Check if I'm live and recording"

### Conditional Operations

**Recording without streaming:**
- "I'm recording offline - check that streaming is off"
- "Start recording but don't stream"

**Streaming without recording:**
- "Just streaming today - make sure recording is off"
- "Start stream but no recording"

## Tips for Recording Prompts

1. **Be explicit**: "Start recording" is clearer than "start it"
2. **Check status**: Always verify before important recordings
3. **Confirm stops**: Make sure recording stopped before walking away
4. **Use context**: Explain why you're pausing (helps with workflow)
5. **Combine actions**: Save time with "switch and record" commands

## Common Variations

All these mean the same thing:
- "Start recording" = "Begin recording" = "Hit record" = "Start a recording"
- "Stop recording" = "End recording" = "Finish recording" = "Stop the recording"
- "Pause recording" = "Hold recording" = "Pause the recording"
- "Resume recording" = "Continue recording" = "Unpause recording"
- "Am I recording?" = "Is recording active?" = "Check recording status"

---

**Next Steps**: Check out [audio.md](audio.md) for audio control prompts, or [workflows.md](workflows.md) for complete multi-step recording workflows.
