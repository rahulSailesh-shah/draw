package llm

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewOllamaLLMClient(t *testing.T) {
	tests := []struct {
		name      string
		ollamaHost string
		model     string
		wantError bool
	}{
		{
			name:      "valid initialization",
			ollamaHost: "http://localhost:11434",
			model:     "llama3.2:3b",
			wantError: false,
		},
		{
			name:      "empty host",
			ollamaHost: "",
			model:     "llama3.2:3b",
			wantError: false, // ClientFromEnvironment handles empty host
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewOllamaLLMClient(tt.ollamaHost, tt.model)
			
			if tt.wantError {
				if err == nil {
					t.Errorf("NewOllamaLLMClient() expected error, got nil")
				}
				return
			}

			if err != nil {
				// If Ollama is not running, skip the test
				if strings.Contains(err.Error(), "connection") || 
				   strings.Contains(err.Error(), "refused") ||
				   strings.Contains(err.Error(), "timeout") {
					t.Skipf("Ollama not available: %v", err)
				}
				t.Fatalf("NewOllamaLLMClient() error = %v", err)
			}

			if client == nil {
				t.Fatal("NewOllamaLLMClient() returned nil client")
			}

			if client.model != tt.model {
				t.Errorf("NewOllamaLLMClient() model = %v, want %v", client.model, tt.model)
			}

			if client.requestChan == nil {
				t.Error("NewOllamaLLMClient() requestChan is nil")
			}

			// Cleanup
			client.Close()
		})
	}
}

func TestOllamaLLMClient_GenerateResponse(t *testing.T) {
	// Skip if Ollama is not available
	client, err := NewOllamaLLMClient("http://localhost:11434", "llama3.2:3b")
	if err != nil {
		if strings.Contains(err.Error(), "connection") || 
		   strings.Contains(err.Error(), "refused") ||
		   strings.Contains(err.Error(), "timeout") {
			t.Skipf("Ollama not available: %v", err)
		}
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	tests := []struct {
		name        string
		prompt      string
		ctx         context.Context
		wantError   bool
		errorSubstr string
	}{
		{
			name:      "valid prompt",
			prompt:    "Hello, how are you?",
			ctx:       context.Background(),
			wantError: false,
		},
		{
			name:        "empty prompt",
			prompt:      "",
			ctx:         context.Background(),
			wantError:   true,
			errorSubstr: "empty text",
		},
		{
			name:        "whitespace only prompt",
			prompt:      "   \n\t  ",
			ctx:         context.Background(),
			wantError:   true,
			errorSubstr: "empty text",
		},
		{
			name:      "context cancellation",
			prompt:    "This is a test",
			ctx:       func() context.Context { ctx, cancel := context.WithCancel(context.Background()); cancel(); return ctx }(),
			wantError: true,
		},
		{
			name:      "context timeout",
			prompt:    "This is a test",
			ctx:       func() context.Context { ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond); defer cancel(); return ctx }(),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := client.GenerateResponse(tt.ctx, tt.prompt)

			if tt.wantError {
				if err == nil {
					t.Errorf("GenerateResponse() expected error, got nil")
					return
				}
				if tt.errorSubstr != "" && !strings.Contains(err.Error(), tt.errorSubstr) {
					t.Errorf("GenerateResponse() error = %v, want error containing %q", err, tt.errorSubstr)
				}
				return
			}

			if err != nil {
				// If it's a connection/timeout error, skip
				if strings.Contains(err.Error(), "connection") || 
				   strings.Contains(err.Error(), "refused") ||
				   strings.Contains(err.Error(), "timeout") ||
				   strings.Contains(err.Error(), "context") {
					t.Skipf("Ollama not available or context cancelled: %v", err)
				}
				t.Fatalf("GenerateResponse() error = %v", err)
			}

			if response == nil {
				t.Fatal("GenerateResponse() returned nil response")
			}

			if response.Response == "" {
				t.Error("GenerateResponse() returned empty response")
			}

			if response.Timestamp.IsZero() {
				t.Error("GenerateResponse() returned zero timestamp")
			}

			// Verify timestamp is recent (within last minute)
			if time.Since(response.Timestamp) > time.Minute {
				t.Errorf("GenerateResponse() timestamp = %v, expected recent timestamp", response.Timestamp)
			}
		})
	}
}

func TestOllamaLLMClient_GenerateResponse_Concurrent(t *testing.T) {
	client, err := NewOllamaLLMClient("http://localhost:11434", "llama3.2:3b")
	if err != nil {
		if strings.Contains(err.Error(), "connection") || 
		   strings.Contains(err.Error(), "refused") ||
		   strings.Contains(err.Error(), "timeout") {
			t.Skipf("Ollama not available: %v", err)
		}
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Test concurrent requests
	const numRequests = 5
	results := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func(id int) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			_, err := client.GenerateResponse(ctx, "Say hello")
			results <- err
		}(i)
	}

	// Collect results
	var errors []error
	for i := 0; i < numRequests; i++ {
		err := <-results
		if err != nil {
			// Skip if it's a connection error
			if strings.Contains(err.Error(), "connection") || 
			   strings.Contains(err.Error(), "refused") ||
			   strings.Contains(err.Error(), "timeout") {
				t.Skipf("Ollama not available: %v", err)
			}
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		t.Errorf("GenerateResponse() concurrent requests failed: %v", errors)
	}
}

func TestOllamaLLMClient_Close(t *testing.T) {
	client, err := NewOllamaLLMClient("http://localhost:11434", "llama3.2:3b")
	if err != nil {
		if strings.Contains(err.Error(), "connection") || 
		   strings.Contains(err.Error(), "refused") ||
		   strings.Contains(err.Error(), "timeout") {
			t.Skipf("Ollama not available: %v", err)
		}
		t.Fatalf("Failed to create client: %v", err)
	}

	// Close should not panic
	err = client.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Verify context is cancelled
	select {
	case <-client.ctx.Done():
		// Expected
	default:
		t.Error("Close() did not cancel context")
	}

	// Verify request channel is closed
	select {
	case _, ok := <-client.requestChan:
		if ok {
			t.Error("Close() did not close request channel")
		}
		// Channel is closed, which is expected
	default:
		t.Error("Close() request channel should be closed")
	}
}

func TestOllamaLLMClient_Worker_ErrorHandling(t *testing.T) {
	client, err := NewOllamaLLMClient("http://localhost:11434", "llama3.2:3b")
	if err != nil {
		if strings.Contains(err.Error(), "connection") || 
		   strings.Contains(err.Error(), "refused") ||
		   strings.Contains(err.Error(), "timeout") {
			t.Skipf("Ollama not available: %v", err)
		}
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Test that worker properly handles requests
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := client.GenerateResponse(ctx, "Test prompt")
	if err != nil {
		if strings.Contains(err.Error(), "connection") || 
		   strings.Contains(err.Error(), "refused") ||
		   strings.Contains(err.Error(), "timeout") {
			t.Skipf("Ollama not available: %v", err)
		}
		t.Fatalf("GenerateResponse() error = %v", err)
	}

	if response == nil {
		t.Fatal("GenerateResponse() returned nil response")
	}
}

func TestOllamaLLMClient_GenerateResponse_PromptFormatting(t *testing.T) {
	client, err := NewOllamaLLMClient("http://localhost:11434", "llama3.2:3b")
	if err != nil {
		if strings.Contains(err.Error(), "connection") || 
		   strings.Contains(err.Error(), "refused") ||
		   strings.Contains(err.Error(), "timeout") {
			t.Skipf("Ollama not available: %v", err)
		}
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Test that the prompt is properly formatted with system message
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	response, err := client.GenerateResponse(ctx, "What is 2+2?")
	if err != nil {
		if strings.Contains(err.Error(), "connection") || 
		   strings.Contains(err.Error(), "refused") ||
		   strings.Contains(err.Error(), "timeout") {
			t.Skipf("Ollama not available: %v", err)
		}
		t.Fatalf("GenerateResponse() error = %v", err)
	}

	if response == nil {
		t.Fatal("GenerateResponse() returned nil response")
	}

	// The response should be non-empty
	if len(response.Response) == 0 {
		t.Error("GenerateResponse() returned empty response")
	}
}

func TestOllamaLLMClient_GenerateResponse_Timeout(t *testing.T) {
	client, err := NewOllamaLLMClient("http://localhost:11434", "llama3.2:3b")
	if err != nil {
		if strings.Contains(err.Error(), "connection") || 
		   strings.Contains(err.Error(), "refused") ||
		   strings.Contains(err.Error(), "timeout") {
			t.Skipf("Ollama not available: %v", err)
		}
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Test with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Wait a bit to ensure timeout
	time.Sleep(1 * time.Millisecond)

	_, err = client.GenerateResponse(ctx, "Test")
	if err == nil {
		t.Error("GenerateResponse() expected timeout error, got nil")
	}

	if !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
		// If it's a connection error, skip
		if strings.Contains(err.Error(), "connection") || 
		   strings.Contains(err.Error(), "refused") {
			t.Skipf("Ollama not available: %v", err)
		}
		// Otherwise, check if error message contains timeout/cancel
		if !strings.Contains(err.Error(), "timeout") && 
		   !strings.Contains(err.Error(), "deadline") &&
		   !strings.Contains(err.Error(), "canceled") {
			t.Errorf("GenerateResponse() error = %v, expected context timeout/cancel", err)
		}
	}
}

func TestOllamaLLMClient_MultipleClose(t *testing.T) {
	client, err := NewOllamaLLMClient("http://localhost:11434", "llama3.2:3b")
	if err != nil {
		if strings.Contains(err.Error(), "connection") || 
		   strings.Contains(err.Error(), "refused") ||
		   strings.Contains(err.Error(), "timeout") {
			t.Skipf("Ollama not available: %v", err)
		}
		t.Fatalf("Failed to create client: %v", err)
	}

	// Close multiple times should not panic
	err1 := client.Close()
	err2 := client.Close()
	err3 := client.Close()

	if err1 != nil {
		t.Errorf("First Close() error = %v", err1)
	}
	if err2 != nil {
		t.Errorf("Second Close() error = %v", err2)
	}
	if err3 != nil {
		t.Errorf("Third Close() error = %v", err3)
	}
}

