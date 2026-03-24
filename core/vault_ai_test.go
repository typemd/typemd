package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitAI_MissingBaseURL(t *testing.T) {
	dir := t.TempDir()
	metaDir := filepath.Join(dir, ".typemd")
	os.MkdirAll(metaDir, 0755)
	os.WriteFile(filepath.Join(metaDir, "config.yaml"), []byte(`
ai:
  enabled: true
  default: broken
  providers:
    broken:
      type: openai-compatible
      model: llama3.2
`), 0644)

	v := NewVault(dir)
	if err := v.Open(); err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer v.Close()

	if v.AIService() != nil {
		t.Error("expected nil AI service for openai-compatible without base_url")
	}
}

func TestInitAI_UnknownType(t *testing.T) {
	dir := t.TempDir()
	metaDir := filepath.Join(dir, ".typemd")
	os.MkdirAll(metaDir, 0755)
	os.WriteFile(filepath.Join(metaDir, "config.yaml"), []byte(`
ai:
  enabled: true
  default: custom
  providers:
    custom:
      type: bedrock
      model: some-model
`), 0644)

	v := NewVault(dir)
	if err := v.Open(); err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer v.Close()

	if v.AIService() != nil {
		t.Error("expected nil AI service for unknown provider type")
	}
}

func TestInitAI_DefaultPointsToUndefined(t *testing.T) {
	dir := t.TempDir()
	metaDir := filepath.Join(dir, ".typemd")
	os.MkdirAll(metaDir, 0755)
	os.WriteFile(filepath.Join(metaDir, "config.yaml"), []byte(`
ai:
  enabled: true
  default: nonexistent
  providers:
    claude:
      type: cli
`), 0644)

	v := NewVault(dir)
	if err := v.Open(); err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer v.Close()

	if v.AIService() != nil {
		t.Error("expected nil AI service when default points to undefined provider")
	}
}

func TestInitAI_OpenAICompatibleProvider(t *testing.T) {
	dir := t.TempDir()
	metaDir := filepath.Join(dir, ".typemd")
	os.MkdirAll(metaDir, 0755)
	os.WriteFile(filepath.Join(metaDir, "config.yaml"), []byte(`
ai:
  enabled: true
  default: ollama
  providers:
    ollama:
      type: openai-compatible
      base_url: http://localhost:11434
      model: qwen3-coder:30b
`), 0644)

	v := NewVault(dir)
	if err := v.Open(); err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer v.Close()

	if v.AIService() == nil {
		t.Error("expected non-nil AI service for openai-compatible provider")
	}
}

func TestInitAI_Disabled(t *testing.T) {
	dir := t.TempDir()
	metaDir := filepath.Join(dir, ".typemd")
	os.MkdirAll(metaDir, 0755)
	os.WriteFile(filepath.Join(metaDir, "config.yaml"), []byte(`
ai:
  enabled: false
  default: ollama
  providers:
    ollama:
      type: openai-compatible
      base_url: http://localhost:11434
      model: qwen3-coder:30b
`), 0644)

	v := NewVault(dir)
	if err := v.Open(); err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer v.Close()

	if v.AIService() != nil {
		t.Error("expected nil AI service when disabled")
	}
}
