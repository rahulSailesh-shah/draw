package livekit

import (
	"github.com/livekit/media-sdk"
)

type LivekitHandler interface {
	SendAudioChunk(sample media.PCM16Sample) error
	Close() error
}
