package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"draw/pkg/config"
)

// Client is a client for the LLM service (Ollama).
type Client struct {
	httpClient *http.Client
	config     *config.LLMConfig
}

// NewClient creates a new LLM client for Ollama.
func NewClient(cfg *config.LLMConfig) (*Client, error) {
	if cfg.Host == "" {
		cfg.Host = "http://localhost:11434"
	}
	if cfg.Model == "" {
		cfg.Model = "llama3.2"
	}

	return &Client{
		httpClient: &http.Client{},
		config:     cfg,
	}, nil
}

// Close closes the LLM client.
func (c *Client) Close() error {
	return nil
}

// CanvasRequest represents a request to generate canvas elements.
type CanvasRequest struct {
	UserTranscript string          `json:"user_transcript"`
	BoardState     json.RawMessage `json:"board_state"` // Current canvas elements as JSON
	BoardID        string          `json:"board_id"`
	UserID         string          `json:"user_id"`
}

// CanvasResponse represents the LLM response with canvas updates.
type CanvasResponse struct {
	Elements    json.RawMessage `json:"elements"`    // Canvas elements to add/update
	Explanation string          `json:"explanation"` // Brief explanation for TTS
}

// ollamaRequest represents the Ollama API request.
type ollamaRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
	Format   string          `json:"format,omitempty"`
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ollamaResponse represents the Ollama API response.
type ollamaResponse struct {
	Model     string        `json:"model"`
	Message   ollamaMessage `json:"message"`
	Done      bool          `json:"done"`
	CreatedAt string        `json:"created_at"`
}

// GenerateCanvasUpdate generates canvas updates based on user voice command.
func (c *Client) GenerateCanvasUpdate(ctx context.Context, req CanvasRequest) (*CanvasResponse, error) {
	prompt := buildPrompt(req)

	ollamaReq := ollamaRequest{
		Model: c.config.Model,
		Messages: []ollamaMessage{
			{
				Role:    "system",
				Content: getSystemPrompt(),
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
		Format: "json",
	}

	reqBody, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.config.Host+"/api/chat", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama request failed: %s - %s", resp.Status, string(body))
	}

	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Parse the JSON response from the LLM
	var canvasResp CanvasResponse
	if err := json.Unmarshal([]byte(ollamaResp.Message.Content), &canvasResp); err != nil {
		// If parsing fails, try to extract from the response
		return nil, fmt.Errorf("failed to parse LLM response as JSON: %w\nResponse: %s", err, ollamaResp.Message.Content)
	}

	return &canvasResp, nil
}

// getSystemPrompt returns the system prompt for the LLM.
func getSystemPrompt() string {
	return `You are a UI component generator for a collaborative canvas application.
You MUST respond with valid JSON only, no markdown, no explanation outside JSON.

Your response format MUST be:
{
  "elements": [
    {
      "id": "unique-id",
      "type": "button|input|text|container|form|card|image",
      "props": { ... component specific props ... },
      "position": { "x": number, "y": number },
      "size": { "width": number, "height": number }
    }
  ],
  "explanation": "Brief explanation of what was done"
}

Common element types and their props:
- button: { "label": "text", "variant": "primary|secondary|outline" }
- input: { "placeholder": "text", "type": "text|email|password" }
- text: { "content": "text", "fontSize": number, "fontWeight": "normal|bold" }
- container: { "backgroundColor": "#hex", "borderRadius": number }
- form: { "title": "text" }
- card: { "title": "text", "description": "text" }
- image: { "src": "url", "alt": "text" }`
}

// buildPrompt creates the user prompt for the LLM.
func buildPrompt(req CanvasRequest) string {
	boardStateStr := "[]"
	if len(req.BoardState) > 0 {
		boardStateStr = string(req.BoardState)
	}

	return fmt.Sprintf(`CURRENT CANVAS STATE:
%s

USER REQUEST:
"%s"

Generate the JSON response with canvas elements to add/update based on the user's request.`, boardStateStr, req.UserTranscript)
}

// StreamCanvasUpdate generates canvas updates with streaming response.
func (c *Client) StreamCanvasUpdate(ctx context.Context, req CanvasRequest) (<-chan string, <-chan error) {
	resultChan := make(chan string, 10)
	errChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errChan)

		prompt := buildPrompt(req)

		ollamaReq := ollamaRequest{
			Model: c.config.Model,
			Messages: []ollamaMessage{
				{
					Role:    "system",
					Content: getSystemPrompt(),
				},
				{
					Role:    "user",
					Content: prompt,
				},
			},
			Stream: true,
			Format: "json",
		}

		reqBody, err := json.Marshal(ollamaReq)
		if err != nil {
			errChan <- fmt.Errorf("failed to marshal request: %w", err)
			return
		}

		httpReq, err := http.NewRequestWithContext(ctx, "POST", c.config.Host+"/api/chat", bytes.NewReader(reqBody))
		if err != nil {
			errChan <- fmt.Errorf("failed to create request: %w", err)
			return
		}
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			errChan <- fmt.Errorf("failed to send request: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			errChan <- fmt.Errorf("ollama request failed: %s - %s", resp.Status, string(body))
			return
		}

		decoder := json.NewDecoder(resp.Body)
		for {
			var ollamaResp ollamaResponse
			if err := decoder.Decode(&ollamaResp); err != nil {
				if err == io.EOF {
					return
				}
				errChan <- fmt.Errorf("failed to decode stream: %w", err)
				return
			}

			if ollamaResp.Message.Content != "" {
				resultChan <- ollamaResp.Message.Content
			}

			if ollamaResp.Done {
				return
			}
		}
	}()

	return resultChan, errChan
}
