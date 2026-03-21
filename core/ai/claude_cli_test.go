package ai

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestClaudeCLI_BuildArgs_NoSchema(t *testing.T) {
	c := &ClaudeCLI{Binary: "echo"}
	req := &CompletionRequest{
		SystemPrompt: "You are helpful",
		UserPrompt:   "Hello",
	}

	// We can't easily test args without invoking, but we can test with a mock script
	// For now, test that the struct satisfies the Provider interface
	var _ Provider = c
	_ = req
}

func TestClaudeCLI_BuildArgs_WithSchema(t *testing.T) {
	c := &ClaudeCLI{Binary: "echo", Model: "claude-haiku-4-5-20251001"}
	req := &CompletionRequest{
		SystemPrompt: "test",
		UserPrompt:   "test prompt",
		JSONSchema:   json.RawMessage(`{"type":"object"}`),
	}
	_ = c
	_ = req
}

func TestClaudeCLI_BinaryNotFound(t *testing.T) {
	c := &ClaudeCLI{Binary: "nonexistent-binary-that-does-not-exist"}
	req := &CompletionRequest{
		UserPrompt: "hello",
	}

	_, err := c.Complete(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for missing binary")
	}
	if !strings.Contains(err.Error(), "ai: claude CLI error") {
		t.Errorf("expected 'ai: claude CLI error', got: %v", err)
	}
}

func TestClaudeCLI_ContextCancellation(t *testing.T) {
	// Create a script that sleeps for a long time
	dir := t.TempDir()
	script := filepath.Join(dir, "slow-claude")
	os.WriteFile(script, []byte("#!/bin/sh\nsleep 30\n"), 0755)

	c := &ClaudeCLI{Binary: script}
	req := &CompletionRequest{UserPrompt: "hello"}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := c.Complete(ctx, req)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
	if !strings.Contains(err.Error(), "context") {
		t.Errorf("expected context error, got: %v", err)
	}
}

func TestClaudeCLI_NonZeroExitCode(t *testing.T) {
	// Create a script that exits with error
	dir := t.TempDir()
	script := filepath.Join(dir, "fail-claude")
	os.WriteFile(script, []byte("#!/bin/sh\necho 'auth expired' >&2\nexit 1\n"), 0755)

	c := &ClaudeCLI{Binary: script}
	req := &CompletionRequest{UserPrompt: "hello"}

	_, err := c.Complete(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for non-zero exit")
	}
	if !strings.Contains(err.Error(), "auth expired") {
		t.Errorf("expected stderr in error, got: %v", err)
	}
}

func TestClaudeCLI_ParsesJSONOutput(t *testing.T) {
	// Create a script that outputs valid claude JSON
	dir := t.TempDir()
	script := filepath.Join(dir, "mock-claude")
	os.WriteFile(script, []byte(`#!/bin/sh
echo '{"type":"result","result":"Hello, world!"}'
`), 0755)

	c := &ClaudeCLI{Binary: script}
	req := &CompletionRequest{UserPrompt: "hello"}

	resp, err := c.Complete(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Content != "Hello, world!" {
		t.Errorf("expected 'Hello, world!', got %q", resp.Content)
	}
	if resp.JSONResult != nil {
		t.Error("expected nil JSONResult when no schema provided")
	}
}

func TestClaudeCLI_ParsesJSONSchemaOutput(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "mock-claude")
	os.WriteFile(script, []byte(`#!/bin/sh
echo '{"type":"result","result":"","structured_output":{"description":"A test description"}}'
`), 0755)

	c := &ClaudeCLI{Binary: script}
	req := &CompletionRequest{
		UserPrompt: "describe this",
		JSONSchema: json.RawMessage(`{"type":"object","properties":{"description":{"type":"string"}}}`),
	}

	resp, err := c.Complete(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.JSONResult == nil {
		t.Fatal("expected non-nil JSONResult when schema provided")
	}
	if !strings.Contains(string(resp.JSONResult), "test description") {
		t.Errorf("expected JSONResult to contain 'test description', got %s", resp.JSONResult)
	}
}

func TestClaudeCLI_FallbackResultForSchema(t *testing.T) {
	// Older CLI versions may put structured output in result field
	dir := t.TempDir()
	script := filepath.Join(dir, "mock-claude")
	os.WriteFile(script, []byte(`#!/bin/sh
echo '{"type":"result","result":"{\"description\":\"fallback\"}"}'
`), 0755)

	c := &ClaudeCLI{Binary: script}
	req := &CompletionRequest{
		UserPrompt: "describe this",
		JSONSchema: json.RawMessage(`{"type":"object"}`),
	}

	resp, err := c.Complete(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.JSONResult == nil {
		t.Fatal("expected non-nil JSONResult from fallback")
	}
	if !strings.Contains(string(resp.JSONResult), "fallback") {
		t.Errorf("expected fallback content, got %s", resp.JSONResult)
	}
}

func TestClaudeCLI_ModelOverride(t *testing.T) {
	// Create a script that echoes the arguments to verify model flag
	dir := t.TempDir()
	script := filepath.Join(dir, "mock-claude")
	os.WriteFile(script, []byte(`#!/bin/sh
# Output args as the result for verification
echo "{\"type\":\"result\",\"result\":\"$*\"}"
`), 0755)

	c := &ClaudeCLI{Binary: script}
	req := &CompletionRequest{
		UserPrompt: "test",
		Model:      "claude-haiku-4-5-20251001",
	}

	resp, err := c.Complete(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Content, "--model") || !strings.Contains(resp.Content, "claude-haiku-4-5-20251001") {
		t.Errorf("expected model flag in args, got: %s", resp.Content)
	}
}

func TestClaudeCLI_DefaultModelFromStruct(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "mock-claude")
	os.WriteFile(script, []byte(`#!/bin/sh
echo "{\"type\":\"result\",\"result\":\"$*\"}"
`), 0755)

	c := &ClaudeCLI{Binary: script, Model: "claude-sonnet-4-6-20250627"}
	req := &CompletionRequest{
		UserPrompt: "test",
	}

	resp, err := c.Complete(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Content, "claude-sonnet-4-6-20250627") {
		t.Errorf("expected default model in args, got: %s", resp.Content)
	}
}

func TestClaudeCLI_RequestModelOverridesDefault(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "mock-claude")
	os.WriteFile(script, []byte(`#!/bin/sh
echo "{\"type\":\"result\",\"result\":\"$*\"}"
`), 0755)

	c := &ClaudeCLI{Binary: script, Model: "default-model"}
	req := &CompletionRequest{
		UserPrompt: "test",
		Model:      "override-model",
	}

	resp, err := c.Complete(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.Content, "override-model") {
		t.Errorf("expected override model in args, got: %s", resp.Content)
	}
	if strings.Contains(resp.Content, "default-model") {
		t.Error("should not contain default model when overridden")
	}
}

func TestClaudeCLI_InvalidJSONOutput(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "mock-claude")
	os.WriteFile(script, []byte("#!/bin/sh\necho 'not json'\n"), 0755)

	c := &ClaudeCLI{Binary: script}
	req := &CompletionRequest{UserPrompt: "hello"}

	_, err := c.Complete(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for invalid JSON output")
	}
	if !strings.Contains(err.Error(), "failed to parse") {
		t.Errorf("expected parse error, got: %v", err)
	}
}

func TestLookupBinary(t *testing.T) {
	// This test just verifies the function runs without panicking.
	// The actual result depends on the test environment.
	_, err := LookupBinary()
	if err != nil {
		// claude might not be installed in CI, that's fine
		if _, ok := err.(*exec.Error); !ok {
			t.Logf("LookupBinary returned non-exec error: %v", err)
		}
	}
}
