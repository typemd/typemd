package core

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestProviderConfig_YAMLRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  ProviderConfig
	}{
		{
			name:  "cli provider",
			input: "type: cli\nmodel: claude-sonnet-4-20250514\n",
			want:  ProviderConfig{Type: "cli", Model: "claude-sonnet-4-20250514"},
		},
		{
			name:  "openai-compatible provider",
			input: "type: openai-compatible\nbase_url: http://localhost:11434\nmodel: qwen3-coder:30b\n",
			want:  ProviderConfig{Type: "openai-compatible", BaseURL: "http://localhost:11434", Model: "qwen3-coder:30b"},
		},
		{
			name:  "provider with api_key",
			input: "type: openai-compatible\nbase_url: http://example.com\nmodel: llama3.2\napi_key: sk-test\n",
			want:  ProviderConfig{Type: "openai-compatible", BaseURL: "http://example.com", Model: "llama3.2", APIKey: "sk-test"},
		},
		{
			name:  "minimal cli provider",
			input: "type: cli\n",
			want:  ProviderConfig{Type: "cli"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got ProviderConfig
			if err := yaml.Unmarshal([]byte(tt.input), &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got != tt.want {
				t.Errorf("unmarshal:\ngot  %+v\nwant %+v", got, tt.want)
			}

			// Round-trip: marshal and unmarshal again
			data, err := yaml.Marshal(&got)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			var roundTrip ProviderConfig
			if err := yaml.Unmarshal(data, &roundTrip); err != nil {
				t.Fatalf("unmarshal round-trip: %v", err)
			}
			if roundTrip != got {
				t.Errorf("round-trip:\ngot  %+v\nwant %+v", roundTrip, got)
			}
		})
	}
}

func TestAIConfig_ProvidersMap(t *testing.T) {
	input := `
enabled: true
default: ollama
providers:
  claude:
    type: cli
    model: claude-sonnet-4-20250514
  ollama:
    type: openai-compatible
    base_url: http://localhost:11434
    model: qwen3-coder:30b
prompts:
  describe: "custom"
`
	var cfg AIConfig
	if err := yaml.Unmarshal([]byte(input), &cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !cfg.Enabled {
		t.Error("expected enabled true")
	}
	if cfg.Default != "ollama" {
		t.Errorf("expected default 'ollama', got %q", cfg.Default)
	}
	if len(cfg.Providers) != 2 {
		t.Fatalf("expected 2 providers, got %d", len(cfg.Providers))
	}
	claude := cfg.Providers["claude"]
	if claude.Type != "cli" || claude.Model != "claude-sonnet-4-20250514" {
		t.Errorf("claude provider: %+v", claude)
	}
	ollama := cfg.Providers["ollama"]
	if ollama.Type != "openai-compatible" || ollama.BaseURL != "http://localhost:11434" || ollama.Model != "qwen3-coder:30b" {
		t.Errorf("ollama provider: %+v", ollama)
	}
	if cfg.Prompts.Describe != "custom" {
		t.Errorf("expected prompts.describe 'custom', got %q", cfg.Prompts.Describe)
	}
}

func TestAIConfig_EmptyProviders(t *testing.T) {
	input := `
enabled: true
model: claude-haiku-4-5-20251001
`
	var cfg AIConfig
	if err := yaml.Unmarshal([]byte(input), &cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if cfg.Providers != nil {
		t.Errorf("expected nil providers, got %v", cfg.Providers)
	}
	if cfg.Default != "" {
		t.Errorf("expected empty default, got %q", cfg.Default)
	}
	if cfg.Model != "claude-haiku-4-5-20251001" {
		t.Errorf("expected model 'claude-haiku-4-5-20251001', got %q", cfg.Model)
	}
}
