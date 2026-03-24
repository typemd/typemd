package ai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOpenAICompatible_StructuredOutput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/v1/chat/completions") {
			t.Errorf("expected /v1/chat/completions, got %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var req chatRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("unmarshal request: %v", err)
		}

		if req.Model != "qwen3-coder:30b" {
			t.Errorf("expected model 'qwen3-coder:30b', got %q", req.Model)
		}
		if len(req.Messages) != 2 {
			t.Errorf("expected 2 messages, got %d", len(req.Messages))
		}
		if req.Messages[0].Role != "system" {
			t.Errorf("expected system message first, got %q", req.Messages[0].Role)
		}
		if req.ResponseFormat == nil {
			t.Fatal("expected response_format to be set")
		}
		if req.ResponseFormat.Type != "json_schema" {
			t.Errorf("expected type 'json_schema', got %q", req.ResponseFormat.Type)
		}
		if req.ResponseFormat.JSONSchema.Name != "response" {
			t.Errorf("expected schema name 'response', got %q", req.ResponseFormat.JSONSchema.Name)
		}
		if !req.ResponseFormat.JSONSchema.Strict {
			t.Error("expected strict: true")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chatResponse{
			Choices: []chatChoice{
				{Message: chatMessage{Role: "assistant", Content: `{"description":"A great book"}`}},
			},
		})
	}))
	defer server.Close()

	provider := &OpenAICompatible{
		BaseURL:    server.URL,
		Model:      "qwen3-coder:30b",
		HTTPClient: server.Client(),
	}

	schema := json.RawMessage(`{"type":"object","properties":{"description":{"type":"string"}},"required":["description"]}`)
	resp, err := provider.Complete(context.Background(), &CompletionRequest{
		SystemPrompt: "You are helpful",
		UserPrompt:   "Describe this book",
		JSONSchema:   schema,
	})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}

	if resp.Content != `{"description":"A great book"}` {
		t.Errorf("unexpected content: %q", resp.Content)
	}
	if resp.JSONResult == nil {
		t.Fatal("expected JSONResult to be set")
	}

	var result struct {
		Description string `json:"description"`
	}
	if err := json.Unmarshal(resp.JSONResult, &result); err != nil {
		t.Fatalf("unmarshal JSONResult: %v", err)
	}
	if result.Description != "A great book" {
		t.Errorf("expected 'A great book', got %q", result.Description)
	}
}

func TestOpenAICompatible_NoSchema(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req chatRequest
		json.Unmarshal(body, &req)

		if req.ResponseFormat != nil {
			t.Error("expected no response_format when no schema")
		}
		if len(req.Messages) != 1 {
			t.Errorf("expected 1 message (no system prompt), got %d", len(req.Messages))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chatResponse{
			Choices: []chatChoice{
				{Message: chatMessage{Role: "assistant", Content: "Hello, world!"}},
			},
		})
	}))
	defer server.Close()

	provider := &OpenAICompatible{
		BaseURL:    server.URL,
		Model:      "llama3.2",
		HTTPClient: server.Client(),
	}

	resp, err := provider.Complete(context.Background(), &CompletionRequest{
		UserPrompt: "Say hello",
	})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}

	if resp.Content != "Hello, world!" {
		t.Errorf("unexpected content: %q", resp.Content)
	}
	if resp.JSONResult != nil {
		t.Error("expected nil JSONResult when no schema")
	}
}

func TestOpenAICompatible_ModelOverride(t *testing.T) {
	var gotModel string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req chatRequest
		json.Unmarshal(body, &req)
		gotModel = req.Model

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chatResponse{
			Choices: []chatChoice{
				{Message: chatMessage{Role: "assistant", Content: "ok"}},
			},
		})
	}))
	defer server.Close()

	provider := &OpenAICompatible{
		BaseURL:    server.URL,
		Model:      "default-model",
		HTTPClient: server.Client(),
	}

	// Override model in request
	_, err := provider.Complete(context.Background(), &CompletionRequest{
		UserPrompt: "test",
		Model:      "override-model",
	})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if gotModel != "override-model" {
		t.Errorf("expected model 'override-model', got %q", gotModel)
	}

	// Use default model
	_, err = provider.Complete(context.Background(), &CompletionRequest{
		UserPrompt: "test",
	})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if gotModel != "default-model" {
		t.Errorf("expected model 'default-model', got %q", gotModel)
	}
}

func TestOpenAICompatible_APIKeyHeader(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chatResponse{
			Choices: []chatChoice{
				{Message: chatMessage{Role: "assistant", Content: "ok"}},
			},
		})
	}))
	defer server.Close()

	// With API key
	provider := &OpenAICompatible{
		BaseURL:    server.URL,
		Model:      "model",
		APIKey:     "sk-test-key",
		HTTPClient: server.Client(),
	}
	_, err := provider.Complete(context.Background(), &CompletionRequest{UserPrompt: "test"})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if gotAuth != "Bearer sk-test-key" {
		t.Errorf("expected 'Bearer sk-test-key', got %q", gotAuth)
	}

	// Without API key
	provider.APIKey = ""
	gotAuth = ""
	_, err = provider.Complete(context.Background(), &CompletionRequest{UserPrompt: "test"})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if gotAuth != "" {
		t.Errorf("expected no auth header, got %q", gotAuth)
	}
}

func TestOpenAICompatible_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "model not found"}`))
	}))
	defer server.Close()

	provider := &OpenAICompatible{
		BaseURL:    server.URL,
		Model:      "nonexistent",
		HTTPClient: server.Client(),
	}

	_, err := provider.Complete(context.Background(), &CompletionRequest{UserPrompt: "test"})
	if err == nil {
		t.Fatal("expected error for HTTP 400")
	}
	if !strings.Contains(err.Error(), "HTTP 400") {
		t.Errorf("expected error to contain 'HTTP 400', got: %v", err)
	}
	if !strings.Contains(err.Error(), "model not found") {
		t.Errorf("expected error to contain response body, got: %v", err)
	}
}

func TestOpenAICompatible_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`not json`))
	}))
	defer server.Close()

	provider := &OpenAICompatible{
		BaseURL:    server.URL,
		Model:      "model",
		HTTPClient: server.Client(),
	}

	_, err := provider.Complete(context.Background(), &CompletionRequest{UserPrompt: "test"})
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "parse response JSON") {
		t.Errorf("expected JSON parse error, got: %v", err)
	}
}

func TestOpenAICompatible_ConnectionRefused(t *testing.T) {
	provider := &OpenAICompatible{
		BaseURL: "http://localhost:1", // port 1 should be refused
		Model:   "model",
		HTTPClient: &http.Client{
			Timeout: 1, // 1ns timeout to fail fast
		},
	}

	_, err := provider.Complete(context.Background(), &CompletionRequest{UserPrompt: "test"})
	if err == nil {
		t.Fatal("expected error for connection refused")
	}
	if !strings.Contains(err.Error(), "localhost:1") {
		t.Errorf("expected error to contain base URL, got: %v", err)
	}
}

func TestOpenAICompatible_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Block until context is done
		<-r.Context().Done()
	}))
	defer server.Close()

	provider := &OpenAICompatible{
		BaseURL:    server.URL,
		Model:      "model",
		HTTPClient: server.Client(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := provider.Complete(ctx, &CompletionRequest{UserPrompt: "test"})
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
	if !strings.Contains(err.Error(), "context cancelled") {
		t.Errorf("expected context cancelled error, got: %v", err)
	}
}

func TestOpenAICompatible_EmptyChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chatResponse{Choices: []chatChoice{}})
	}))
	defer server.Close()

	provider := &OpenAICompatible{
		BaseURL:    server.URL,
		Model:      "model",
		HTTPClient: server.Client(),
	}

	_, err := provider.Complete(context.Background(), &CompletionRequest{UserPrompt: "test"})
	if err == nil {
		t.Fatal("expected error for empty choices")
	}
	if !strings.Contains(err.Error(), "no choices") {
		t.Errorf("expected 'no choices' error, got: %v", err)
	}
}

func TestOpenAICompatible_LongResponseBodyTruncation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(strings.Repeat("x", 1000)))
	}))
	defer server.Close()

	provider := &OpenAICompatible{
		BaseURL:    server.URL,
		Model:      "model",
		HTTPClient: server.Client(),
	}

	_, err := provider.Complete(context.Background(), &CompletionRequest{UserPrompt: "test"})
	if err == nil {
		t.Fatal("expected error")
	}
	// Error should be truncated to 500 chars
	if len(err.Error()) > 600 { // some overhead for prefix
		t.Errorf("expected truncated error, got length %d", len(err.Error()))
	}
}
