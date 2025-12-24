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

type TranscribeSession struct {
	client    pb.SpeechServiceClient
	sessionID string
	stream    pb.SpeechService_StreamTranscribeClient
	mu        sync.Mutex
	started   bool
	closed    bool
}

func (c *Client) NewTranscribeSession(ctx context.Context, sessionID string) (*TranscribeSession, error) {
	stream, err := c.client.StreamTranscribe(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transcribe stream: %w", err)
	}

	return &TranscribeSession{
		client:    c.client,
		sessionID: sessionID,
		stream:    stream,
		started:   true,
	}, nil
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

func (s *TranscribeSession) Finalize() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return "", fmt.Errorf("session is closed")
	}

	if err := s.stream.Send(&pb.TranscribeRequest{
		SessionId:   s.sessionID,
		AudioChunk:  nil,
		EndOfStream: true,
	}); err != nil {
		return "", fmt.Errorf("failed to send end-of-stream: %w", err)
	}

	resp, err := s.stream.CloseAndRecv()
	if err != nil {
		if err == io.EOF {
			return "", fmt.Errorf("stream closed without response")
		}
		return "", fmt.Errorf("failed to receive transcription: %w", err)
	}

	s.closed = true

	if !resp.Success {
		return "", fmt.Errorf("transcription failed: %s", resp.Error)
	}

	return resp.Transcription, nil
}

func (s *TranscribeSession) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true
	return s.stream.CloseSend()
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

