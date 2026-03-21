package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const configFileName = "config.yaml"

// VaultConfig holds vault-level configuration loaded from .typemd/config.yaml.
type VaultConfig struct {
	CLI CLIConfig `yaml:"cli"`
	TUI TUIConfig `yaml:"tui,omitempty"`
	AI  AIConfig  `yaml:"ai,omitempty"`
}

// AIConfig holds AI-specific configuration.
type AIConfig struct {
	Enabled bool          `yaml:"enabled,omitempty"`
	Model   string        `yaml:"model,omitempty"`
	Prompts PromptsConfig `yaml:"prompts,omitempty"`
	Explore ExploreConfig `yaml:"explore,omitempty"`
}

// PromptsConfig holds customizable system prompts for AI operations.
type PromptsConfig struct {
	Describe string `yaml:"describe,omitempty"`
	Tag      string `yaml:"tag,omitempty"`
	Explore  string `yaml:"explore,omitempty"`
}

// ExploreConfig holds schema explore parameters.
type ExploreConfig struct {
	SampleCount  int `yaml:"sample_count,omitempty"`
	BodyTruncate int `yaml:"body_truncate,omitempty"`
}

// CLIConfig holds CLI-specific configuration.
type CLIConfig struct {
	DefaultType string `yaml:"default_type"`
}

// TUIConfig holds TUI-specific configuration.
type TUIConfig struct {
	DebounceMs      int    `yaml:"debounce_ms,omitempty"`
	StatsTypeLayout string `yaml:"stats_type_layout,omitempty"`
}

// configKeyEntry maps a dot-notation key to getter/setter on VaultConfig.
type configKeyEntry struct {
	Get func(cfg *VaultConfig) string
	Set func(cfg *VaultConfig, value string)
}

// configKeyRegistry maps dot-notation config keys to VaultConfig struct fields.
var configKeyRegistry = map[string]configKeyEntry{
	"cli.default_type": {
		Get: func(cfg *VaultConfig) string { return cfg.CLI.DefaultType },
		Set: func(cfg *VaultConfig, value string) { cfg.CLI.DefaultType = value },
	},
	"tui.debounce_ms": {
		Get: func(cfg *VaultConfig) string {
			if cfg.TUI.DebounceMs == 0 {
				return ""
			}
			return strconv.Itoa(cfg.TUI.DebounceMs)
		},
		Set: func(cfg *VaultConfig, value string) {
			n, _ := strconv.Atoi(value)
			cfg.TUI.DebounceMs = n
		},
	},
	"tui.stats_type_layout": {
		Get: func(cfg *VaultConfig) string { return cfg.TUI.StatsTypeLayout },
		Set: func(cfg *VaultConfig, value string) { cfg.TUI.StatsTypeLayout = value },
	},
	"ai.enabled": {
		Get: func(cfg *VaultConfig) string {
			if cfg.AI.Enabled {
				return "true"
			}
			return ""
		},
		Set: func(cfg *VaultConfig, value string) {
			cfg.AI.Enabled = value == "true"
		},
	},
	"ai.model": {
		Get: func(cfg *VaultConfig) string { return cfg.AI.Model },
		Set: func(cfg *VaultConfig, value string) { cfg.AI.Model = value },
	},
	"ai.prompts.describe": {
		Get: func(cfg *VaultConfig) string { return cfg.AI.Prompts.Describe },
		Set: func(cfg *VaultConfig, value string) { cfg.AI.Prompts.Describe = value },
	},
	"ai.prompts.tag": {
		Get: func(cfg *VaultConfig) string { return cfg.AI.Prompts.Tag },
		Set: func(cfg *VaultConfig, value string) { cfg.AI.Prompts.Tag = value },
	},
	"ai.prompts.explore": {
		Get: func(cfg *VaultConfig) string { return cfg.AI.Prompts.Explore },
		Set: func(cfg *VaultConfig, value string) { cfg.AI.Prompts.Explore = value },
	},
	"ai.explore.sample_count": {
		Get: func(cfg *VaultConfig) string {
			if cfg.AI.Explore.SampleCount == 0 {
				return ""
			}
			return strconv.Itoa(cfg.AI.Explore.SampleCount)
		},
		Set: func(cfg *VaultConfig, value string) {
			n, _ := strconv.Atoi(value)
			cfg.AI.Explore.SampleCount = n
		},
	},
	"ai.explore.body_truncate": {
		Get: func(cfg *VaultConfig) string {
			if cfg.AI.Explore.BodyTruncate == 0 {
				return ""
			}
			return strconv.Itoa(cfg.AI.Explore.BodyTruncate)
		},
		Set: func(cfg *VaultConfig, value string) {
			n, _ := strconv.Atoi(value)
			cfg.AI.Explore.BodyTruncate = n
		},
	},
}

// Config returns the current in-memory VaultConfig.
func (v *Vault) Config() *VaultConfig {
	return v.config
}

// loadVaultConfig reads and parses the vault config file.
// Returns a zero-value VaultConfig if the file does not exist or is empty.
func loadVaultConfig(metaDir string) (*VaultConfig, error) {
	path := filepath.Join(metaDir, configFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &VaultConfig{}, nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg VaultConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &cfg, nil
}

// WriteConfig writes a VaultConfig to the vault's config.yaml file
// and updates the in-memory config.
func (v *Vault) WriteConfig(cfg *VaultConfig) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	if err := os.WriteFile(filepath.Join(v.Dir(), configFileName), data, 0644); err != nil {
		return err
	}
	v.config = cfg
	return nil
}

// GetConfigValue returns the value for a dot-notation config key.
// Returns an error if the key is unknown.
func (v *Vault) GetConfigValue(key string) (string, error) {
	entry, ok := configKeyRegistry[key]
	if !ok {
		return "", fmt.Errorf("unknown config key %q. Known keys: %s", key, strings.Join(ConfigKeys(), ", "))
	}
	cfg := v.config
	if cfg == nil {
		cfg = &VaultConfig{}
	}
	return entry.Get(cfg), nil
}

// SetConfigValue sets a value for a dot-notation config key.
// Returns an error if the key is unknown. Creates config.yaml if it doesn't exist.
func (v *Vault) SetConfigValue(key, value string) error {
	entry, ok := configKeyRegistry[key]
	if !ok {
		return fmt.Errorf("unknown config key %q. Known keys: %s", key, strings.Join(ConfigKeys(), ", "))
	}
	cfg := v.config
	if cfg == nil {
		cfg = &VaultConfig{}
	}
	entry.Set(cfg, value)
	return v.WriteConfig(cfg)
}

// ConfigKeys returns all known config keys sorted alphabetically.
func ConfigKeys() []string {
	keys := make([]string, 0, len(configKeyRegistry))
	for k := range configKeyRegistry {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
