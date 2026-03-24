package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OpenAICompatible implements Provider by calling an OpenAI-compatible HTTP API.
type OpenAICompatible struct {
	BaseURL    string       // required: e.g. "http://localhost:11434"
	Model      string       // required: model identifier
	APIKey     string       // optional: Bearer token for Authorization header
	HTTPClient *http.Client // optional: custom client for testing; defaults to 60s timeout
}

// Provider type constants.
const (
	ProviderTypeCLI              = "cli"
	ProviderTypeOpenAICompatible = "openai-compatible"
)

type chatRequest struct {
	Model          string        `json:"model"`
	Messages       []chatMessage `json:"messages"`
	ResponseFormat *respFormat   `json:"response_format,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type respFormat struct {
	Type       string          `json:"type"`
	JSONSchema *respJSONSchema `json:"json_schema,omitempty"`
}

type respJSONSchema struct {
	Name   string          `json:"name"`
	Strict bool            `json:"strict"`
	Schema json.RawMessage `json:"schema"`
}

type chatResponse struct {
	Choices []chatChoice `json:"choices"`
}

type chatChoice struct {
	Message chatMessage `json:"message"`
}

// defaultHTTPClient is reused across OpenAICompatible instances that don't provide a custom client.
var defaultHTTPClient = &http.Client{Timeout: 60 * time.Second}

// Complete sends a completion request to the OpenAI-compatible API.
func (o *OpenAICompatible) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	client := o.HTTPClient
	if client == nil {
		client = defaultHTTPClient
	}

	model := req.Model
	if model == "" {
		model = o.Model
	}

	messages := make([]chatMessage, 0, 2)
	if req.SystemPrompt != "" {
		messages = append(messages, chatMessage{Role: "system", Content: req.SystemPrompt})
	}
	messages = append(messages, chatMessage{Role: "user", Content: req.UserPrompt})

	chatReq := chatRequest{
		Model:    model,
		Messages: messages,
	}

	if len(req.JSONSchema) > 0 {
		chatReq.ResponseFormat = &respFormat{
			Type: "json_schema",
			JSONSchema: &respJSONSchema{
				Name:   "response",
				Strict: true,
				Schema: req.JSONSchema,
			},
		}
	}

	body, err := json.Marshal(chatReq)
	if err != nil {
		return nil, fmt.Errorf("ai: marshal request: %w", err)
	}

	url := strings.TrimRight(o.BaseURL, "/") + "/v1/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("ai: create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	if o.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+o.APIKey)
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("ai: context cancelled: %w", ctx.Err())
		}
		return nil, fmt.Errorf("ai: request to %s failed: %w", o.BaseURL, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ai: read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		truncated := string(respBody)
		if len(truncated) > 500 {
			truncated = truncated[:500]
		}
		return nil, fmt.Errorf("ai: HTTP %d from %s: %s", resp.StatusCode, o.BaseURL, truncated)
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("ai: parse response JSON: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("ai: no choices in response from %s", o.BaseURL)
	}

	content := chatResp.Choices[0].Message.Content

	result := &CompletionResponse{
		Content: content,
	}

	if len(req.JSONSchema) > 0 {
		result.JSONResult = json.RawMessage(content)
	}

	return result, nil
}
