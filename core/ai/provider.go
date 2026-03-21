// Package ai provides AI-powered assistance for typemd vaults.
package ai

import (
	"context"
	"encoding/json"
)

// Provider defines the AI completion contract.
type Provider interface {
	Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)
}

// CompletionRequest holds parameters for an AI completion call.
type CompletionRequest struct {
	SystemPrompt string          `json:"system_prompt"`
	UserPrompt   string          `json:"user_prompt"`
	JSONSchema   json.RawMessage `json:"json_schema,omitempty"`
	Model        string          `json:"model,omitempty"`
}

// CompletionResponse holds the result of an AI completion call.
type CompletionResponse struct {
	Content    string          `json:"content"`
	JSONResult json.RawMessage `json:"json_result,omitempty"`
}
