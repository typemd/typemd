package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

// ClaudeCLI implements Provider by invoking the claude CLI binary.
type ClaudeCLI struct {
	Binary string // path to claude binary (default: "claude")
	Model  string // optional default model override
}

// claudeJSONResponse is the JSON output format from `claude -p --output-format json`.
type claudeJSONResponse struct {
	Type             string          `json:"type"`
	Result           string          `json:"result"`
	StructuredOutput json.RawMessage `json:"structured_output,omitempty"`
}

// Complete invokes the claude CLI as a subprocess and returns the result.
func (c *ClaudeCLI) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	binary := c.Binary
	if binary == "" {
		binary = "claude"
	}

	args := []string{
		"-p",
		"--output-format", "json",
		"--tools", "",
	}

	if req.SystemPrompt != "" {
		args = append(args, "--system-prompt", req.SystemPrompt)
	}

	model := req.Model
	if model == "" {
		model = c.Model
	}
	if model != "" {
		args = append(args, "--model", model)
	}

	if len(req.JSONSchema) > 0 {
		args = append(args, "--json-schema", string(req.JSONSchema))
	}

	args = append(args, req.UserPrompt)

	cmd := exec.CommandContext(ctx, binary, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("ai: context cancelled: %w", ctx.Err())
		}
		return nil, fmt.Errorf("ai: claude CLI error: %s: %w", stderr.String(), err)
	}

	var resp claudeJSONResponse
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		return nil, fmt.Errorf("ai: failed to parse claude output: %w", err)
	}

	result := &CompletionResponse{
		Content: resp.Result,
	}

	if len(req.JSONSchema) > 0 && len(resp.StructuredOutput) > 0 {
		result.JSONResult = resp.StructuredOutput
	} else if len(req.JSONSchema) > 0 && resp.Result != "" {
		// Fallback: some versions put structured output in result
		result.JSONResult = json.RawMessage(resp.Result)
	}

	return result, nil
}

// LookupBinary checks if the claude binary is available in PATH.
func LookupBinary() (string, error) {
	return exec.LookPath("claude")
}
