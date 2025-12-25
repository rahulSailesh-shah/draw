# How VAD (Voice Activity Detection) Integration Works

## Overview

The speech service now uses **Silero VAD** to automatically detect speech and silence in real-time audio streams. This eliminates the need for manual mute/unmute events to trigger transcription - transcriptions happen automatically when VAD detects the end of a speech segment (after silence).

## Architecture Flow

```
┌─────────────────┐
│  LiveKit Audio  │
│   (16kHz PCM)   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Voice Handler   │  Sends audio chunks continuously
│  (Go Client)    │  via SendAudioChunk()
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  gRPC Stream    │  Bidirectional streaming
│  (Go → Python)  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Session Manager │  Per-session isolation
│   (Python)      │
└────────┬────────┘
         │
         ├─────────────────┐
         │                 │
         ▼                 ▼
┌──────────────┐   ┌─────────────────┐
│  Silero VAD  │   │  Audio Buffer    │
│  (Real-time) │   │  (Speech chunks) │
└──────┬───────┘   └─────────────────┘
       │
       │ Detects: Speech Start/Stop
       │
       ▼
┌─────────────────┐
│ Background Loop │  Checks silence duration
│  (Thread)       │  every 100ms
└────────┬────────┘
         │
         │ When silence >= threshold:
         │
         ▼
┌─────────────────┐
│ faster-whisper  │  Transcribe buffered audio
│   (STT Model)   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Callback      │  Send transcription
│   (Queue)       │  via gRPC stream
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Go Client      │  Receives transcription
│  (Callback)     │  automatically
└─────────────────┘
```

## How VAD Works (Step by Step)

### 1. **Audio Chunk Processing** (`feed_audio()`)

When an audio chunk arrives from the Go client:

```python
# In session_manager.py:feed_audio()
audio_array = np.frombuffer(audio_chunk, dtype=np.int16)
audio_tensor = torch.from_numpy(audio_array).unsqueeze(0)

# Run VAD on the chunk
speech_prob = self._vad_model(audio_tensor, sample_rate).item()
is_speech = speech_prob >= self.stt_config.silero_sensitivity
```

**VAD Output:**

- Returns a probability (0.0 to 1.0) that speech is present
- Compares against `silero_sensitivity` threshold (default: 0.5)
- `is_speech = True` → Speech detected
- `is_speech = False` → Silence detected

### 2. **Speech State Management**

The session tracks three states:

```python
_is_speaking: bool              # Currently in a speech segment?
_speech_start_time: float       # When did speech start?
_silence_start_time: float      # When did silence start (after speech)?
```

**State Transitions:**

```
[Silence] → [Speech Detected] → Start buffering audio
                                    Set _is_speaking = True
                                    Set _speech_start_time = now()

[Speech] → [Silence Detected] → Set _silence_start_time = now()
                                 (Keep buffering for now)

[Silence] → [Silence Duration >= threshold] → Transcribe!
```

### 3. **Background Processing Loop** (`_process_audio_loop()`)

A background thread runs continuously (checks every 100ms):

```python
while not stopped:
    if _is_speaking and _silence_start_time is not None:
        silence_duration = time.time() - _silence_start_time

        if silence_duration >= post_speech_silence_duration:
            # Check minimum speech duration
            if speech_duration >= min_speech_duration:
                # Transcribe the buffered audio
                _transcribe_and_send(audio_buffer)

            # Reset for next speech segment
            _is_speaking = False
            _speech_buffer = BytesIO()
```

**Key Parameters:**

- `post_speech_silence_duration`: How long to wait after speech ends (default: 0.4s)
- `min_speech_duration`: Minimum speech length to transcribe (default: 0.3s)
- `silero_sensitivity`: VAD threshold for speech detection (default: 0.5)

### 4. **Transcription Trigger**

Transcription happens automatically when:

1. ✅ Speech was detected (`_is_speaking = True`)
2. ✅ Silence detected after speech (`_silence_start_time` is set)
3. ✅ Silence duration >= `post_speech_silence_duration` (default: 0.4s)
4. ✅ Speech duration >= `min_speech_duration` (default: 0.3s)

**No manual trigger needed!** VAD handles everything automatically.

## Key Differences from Old System

### Old System (Manual Finalization)

```
User speaks → Buffer audio → User mutes → Manual finalize() → Transcribe
```

- Required mute/unmute events
- Transcription only on explicit finalization
- No automatic speech detection

### New System (VAD-Based)

```
User speaks → VAD detects → Buffer audio → Silence detected → Auto transcribe
```

- Works continuously (no mute/unmute needed)
- Automatic speech detection
- Transcriptions arrive automatically via callback

## Mute/Unmute Behavior

With VAD integration, mute/unmute events are **optional**:

### OnUnmute()

- Creates a new transcription session
- Starts streaming audio to VAD
- VAD begins detecting speech automatically

### OnMute()

- Closes the transcription session cleanly
- Waits for any pending transcriptions
- Note: Transcriptions may have already arrived via VAD callback

**Important:** Even if user never mutes/unmutes, VAD will still detect speech segments and send transcriptions automatically!

## Per-Session Isolation

Each session maintains its own:

- ✅ Silero VAD model instance
- ✅ faster-whisper model instance
- ✅ Audio buffer (`_speech_buffer`)
- ✅ Background processing thread
- ✅ State tracking (`_is_speaking`, timestamps)

This ensures complete isolation between different users/sessions.

## Configuration

VAD behavior can be tuned via environment variables:

```bash
STT_SILERO_SENSITIVITY=0.5        # VAD threshold (0.0-1.0)
                                  # Lower = more sensitive (detects quieter speech)
                                  # Higher = less sensitive (fewer false positives)

STT_SILENCE_DURATION=0.4          # Seconds of silence before transcribing
                                  # Shorter = faster transcription (may cut off speech)
                                  # Longer = waits more (may include background noise)

STT_MIN_SPEECH_DURATION=0.3       # Minimum speech length to transcribe (seconds)
                                  # Filters out very short utterances/noises
```

## Example Timeline

```
Time    Audio          VAD State        Action
─────────────────────────────────────────────────────
0.0s    [silence]      Silence          Waiting...
0.5s    "Hello"        Speech detected   Start buffering
1.0s    "world"        Speech            Continue buffering
1.5s    [silence]      Silence start    Mark silence_start_time
1.9s    [silence]      Silence          Still waiting...
2.0s    [silence]      Silence >= 0.4s  ✅ TRANSCRIBE!
                          (threshold)
2.1s    "How are"      Speech detected   Start new buffer
2.5s    "you?"         Speech            Continue buffering
3.0s    [silence]      Silence start    Mark silence_start_time
3.4s    [silence]      Silence >= 0.4s  ✅ TRANSCRIBE!
```

## Error Handling

- **VAD initialization fails**: Falls back to buffering all audio (old behavior)
- **VAD processing error**: Falls back to buffering that chunk
- **Transcription error**: Logs error, continues processing
- **Stream disconnection**: Cleanup thread and resources

## Thread Safety

All operations are thread-safe:

- `feed_audio()` uses `_lock` when updating state
- Background loop uses `_lock` when checking/updating state
- Callback execution is safe (called from background thread)

## Summary

**VAD does NOT manually detect silence** - it uses a neural network (Silero VAD) that:

1. Analyzes each audio chunk in real-time
2. Outputs a speech probability score
3. Compares against threshold to determine speech/silence
4. A background loop monitors silence duration
5. Automatically triggers transcription when silence threshold is met

This is **fully automatic** - no manual intervention needed!
