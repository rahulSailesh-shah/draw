package livekit

import (
	"github.com/livekit/media-sdk"
)

// LivekitHandler is the interface for handling audio from LiveKit.
// Implementations stream audio to the speech service and handle transcription.
// With VAD integration, transcriptions arrive automatically via callback when
// silence is detected - no need to wait for mute/unmute events.
type LivekitHandler interface {
	// SendAudioChunk sends a PCM16 audio sample to the speech service.
	// Audio chunks are continuously fed to VAD for speech detection.
	SendAudioChunk(sample media.PCM16Sample) error

	// OnMute is called when the user mutes their microphone.
	// This closes the transcription session and waits for any pending transcriptions.
	// Transcriptions are received automatically via callback when VAD detects silence.
	OnMute() error

	// OnUnmute is called when the user unmutes their microphone.
	// This starts a new transcription session with VAD enabled.
	// VAD will automatically detect speech and send transcriptions via callback.
	OnUnmute() error

	// Close cleans up resources.
	Close() error
}

// TranscriptionCallback is called when transcription is complete.
type TranscriptionCallback func(sessionID string, transcription string, err error)
