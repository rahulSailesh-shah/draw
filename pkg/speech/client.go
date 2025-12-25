package speech

import (
	"context"
	"fmt"
	"io"
	"sync"

	pb "draw/pkg/speech/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.SpeechServiceClient
}

func NewClient(host string) (*Client, error) {
	conn, err := grpc.NewClient(
		host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to speech service: %w", err)
	}

	return &Client{
		conn:   conn,
		client: pb.NewSpeechServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// TranscriptionCallback is called whenever a transcription is received from the server.
type TranscriptionCallback func(transcription string, err error)

type TranscribeSession struct {
	client              pb.SpeechServiceClient
	sessionID           string
	stream              pb.SpeechService_StreamTranscribeClient
	mu                  sync.Mutex
	started             bool
	closed              bool
	transcriptionCallback TranscriptionCallback
	receiveDone         chan struct{}
	receiveErr          error
}

func (c *Client) NewTranscribeSession(ctx context.Context, sessionID string, callback TranscriptionCallback) (*TranscribeSession, error) {
	stream, err := c.client.StreamTranscribe(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transcribe stream: %w", err)
	}

	session := &TranscribeSession{
		client:              c.client,
		sessionID:           sessionID,
		stream:              stream,
		started:             true,
		transcriptionCallback: callback,
		receiveDone:         make(chan struct{}),
	}

	go session.receiveTranscriptions()

	return session, nil
}

func (s *TranscribeSession) receiveTranscriptions() {
	defer close(s.receiveDone)

	fmt.Println("Receiving transcriptions")

	for {
		resp, err := s.stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			s.mu.Lock()
			s.receiveErr = err
			s.mu.Unlock()
			if s.transcriptionCallback != nil {
				s.transcriptionCallback("", fmt.Errorf("failed to receive transcription: %w", err))
			}
			return
		}

		if resp.Success {
			if s.transcriptionCallback != nil {
				s.transcriptionCallback(resp.Transcription, nil)
			}
		} else {
			err := fmt.Errorf("transcription failed: %s", resp.Error)
			if s.transcriptionCallback != nil {
				s.transcriptionCallback("", err)
			}
		}
	}
}

func (s *TranscribeSession) SendAudio(audioChunk []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return fmt.Errorf("session is closed")
	}

	return s.stream.Send(&pb.TranscribeRequest{
		SessionId:   s.sessionID,
		AudioChunk:  audioChunk,
		EndOfStream: false,
	})
}

func (s *TranscribeSession) Finalize() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return fmt.Errorf("session is closed")
	}

	if err := s.stream.Send(&pb.TranscribeRequest{
		SessionId:   s.sessionID,
		AudioChunk:  nil,
		EndOfStream: true,
	}); err != nil {
		return fmt.Errorf("failed to send end-of-stream: %w", err)
	}

	if err := s.stream.CloseSend(); err != nil {
		return fmt.Errorf("failed to close send: %w", err)
	}

	s.closed = true

	<-s.receiveDone

	if s.receiveErr != nil && s.receiveErr != io.EOF {
		return fmt.Errorf("receive error: %w", s.receiveErr)
	}

	return nil
}

func (s *TranscribeSession) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true

	if err := s.stream.CloseSend(); err != nil {
		return err
	}

	select {
	case <-s.receiveDone:
		return nil
	default:
		return nil
	}
}

func (c *Client) CleanupSession(ctx context.Context, sessionID string) error {
	resp, err := c.client.CleanupSession(ctx, &pb.CleanupRequest{
		SessionId: sessionID,
	})
	if err != nil {
		return fmt.Errorf("failed to cleanup session: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("cleanup returned false")
	}

	return nil
}
