package livekit

import (
	"github.com/livekit/media-sdk"
)

// LivekitHandler is the interface for handling audio from LiveKit.
// Implementations stream audio to the speech service and handle transcription.
type LivekitHandler interface {
	// SendAudioChunk sends a PCM16 audio sample to the speech service.
	SendAudioChunk(sample media.PCM16Sample) error

	// OnMute is called when the user mutes their microphone.
	// This signals end of speech and triggers transcription finalization.
	// Returns the transcription result.
	OnMute() (string, error)

	// OnUnmute is called when the user unmutes their microphone.
	// This starts a new transcription session.
	OnUnmute() error

	// Close cleans up resources.
	Close() error
}

// TranscriptionCallback is called when transcription is complete.
type TranscriptionCallback func(sessionID string, transcription string, err error)
