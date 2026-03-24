package core

import (
	"testing"
)

func TestMigrateAIConfig_FlatModelToProviders(t *testing.T) {
	ai := &AIConfig{
		Enabled: true,
		Model:   "claude-haiku-4-5-20251001",
	}
	migrateAIConfig(ai)

	if ai.Default != "claude" {
		t.Errorf("expected default 'claude', got %q", ai.Default)
	}
	if len(ai.Providers) != 1 {
		t.Fatalf("expected 1 provider, got %d", len(ai.Providers))
	}
	p := ai.Providers["claude"]
	if p.Type != "cli" {
		t.Errorf("expected type 'cli', got %q", p.Type)
	}
	if p.Model != "claude-haiku-4-5-20251001" {
		t.Errorf("expected model 'claude-haiku-4-5-20251001', got %q", p.Model)
	}
}

func TestMigrateAIConfig_EnabledWithoutModel(t *testing.T) {
	ai := &AIConfig{
		Enabled: true,
	}
	migrateAIConfig(ai)

	if ai.Default != "claude" {
		t.Errorf("expected default 'claude', got %q", ai.Default)
	}
	p := ai.Providers["claude"]
	if p.Type != "cli" {
		t.Errorf("expected type 'cli', got %q", p.Type)
	}
	if p.Model != "" {
		t.Errorf("expected empty model, got %q", p.Model)
	}
}

func TestMigrateAIConfig_ProvidersPresent_NoMigration(t *testing.T) {
	ai := &AIConfig{
		Enabled: true,
		Model:   "old-model",
		Default: "ollama",
		Providers: map[string]ProviderConfig{
			"ollama": {Type: "openai-compatible", BaseURL: "http://localhost:11434", Model: "qwen3"},
		},
	}
	migrateAIConfig(ai)

	// Should NOT modify when providers already exist
	if ai.Default != "ollama" {
		t.Errorf("expected default 'ollama', got %q", ai.Default)
	}
	if len(ai.Providers) != 1 {
		t.Errorf("expected 1 provider, got %d", len(ai.Providers))
	}
	if _, ok := ai.Providers["claude"]; ok {
		t.Error("should not have created claude provider when providers already exist")
	}
}

func TestMigrateAIConfig_Disabled_NoMigration(t *testing.T) {
	ai := &AIConfig{
		Enabled: false,
		Model:   "some-model",
	}
	migrateAIConfig(ai)

	if ai.Default != "" {
		t.Errorf("expected empty default, got %q", ai.Default)
	}
	if ai.Providers != nil {
		t.Errorf("expected nil providers, got %v", ai.Providers)
	}
}

func TestMigrateAIConfig_NoAISection(t *testing.T) {
	ai := &AIConfig{}
	migrateAIConfig(ai)

	if ai.Default != "" {
		t.Errorf("expected empty default, got %q", ai.Default)
	}
	if ai.Providers != nil {
		t.Errorf("expected nil providers, got %v", ai.Providers)
	}
}
