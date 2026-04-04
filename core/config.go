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
	DateFormat     string    `yaml:"date_format,omitempty"`
	DatetimeFormat string    `yaml:"datetime_format,omitempty"`
	CLI            CLIConfig `yaml:"cli"`
	TUI            TUIConfig `yaml:"tui,omitempty"`
	AI             AIConfig  `yaml:"ai,omitempty"`
	Web            WebConfig `yaml:"web,omitempty"`
}

// WebConfig holds web UI configuration.
type WebConfig struct {
	Theme string `yaml:"theme,omitempty"` // "warm" (default), "dark", or "light"
}

// AIConfig holds AI-specific configuration.
type AIConfig struct {
	Enabled   bool                      `yaml:"enabled,omitempty"`
	Default   string                    `yaml:"default,omitempty"`
	Providers map[string]ProviderConfig `yaml:"providers,omitempty"`
	Language  string                    `yaml:"language,omitempty"`
	Prompts   PromptsConfig             `yaml:"prompts,omitempty"`
	Explore   ExploreConfig             `yaml:"explore,omitempty"`
}

// ProviderConfig holds per-provider settings.
type ProviderConfig struct {
	Type   string `yaml:"type"`
	BaseURL string `yaml:"base_url,omitempty"`
	Model  string `yaml:"model,omitempty"`
	APIKey string `yaml:"api_key,omitempty"`
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
	DebounceMs      int         `yaml:"debounce_ms,omitempty"`
	StatsTypeLayout string      `yaml:"stats_type_layout,omitempty"`
	Toast           ToastConfig `yaml:"toast,omitempty"`
	Theme           ThemeConfig `yaml:"theme,omitempty"`
}

// ThemeConfig holds TUI color theme configuration.
type ThemeConfig struct {
	FocusBorder string `yaml:"focus_border,omitempty"`
	WikiLink    string `yaml:"wiki_link,omitempty"`
	Heading     string `yaml:"heading,omitempty"`
	Bold        string `yaml:"bold,omitempty"`
	Italic      string `yaml:"italic,omitempty"`
	InlineCode  string `yaml:"inline_code,omitempty"`
	CodeBlock   string `yaml:"code_block,omitempty"`
	Link        string `yaml:"link,omitempty"`
	Blockquote  string `yaml:"blockquote,omitempty"`
	HRule       string `yaml:"hrule,omitempty"`
}

// ToastConfig holds toast notification configuration.
type ToastConfig struct {
	Position     string `yaml:"position,omitempty"`       // "bottom-right" (default) or "help-bar"
	DurationMs   int    `yaml:"duration_ms,omitempty"`    // default 3000
	DismissKey   string `yaml:"dismiss_key,omitempty"`    // default "esc"
	ShowWarnings *bool  `yaml:"show_warnings,omitempty"`  // default true (pointer for zero-value detection)
	ShowSuccess  *bool  `yaml:"show_success,omitempty"`   // default false
}

// configKeyEntry maps a dot-notation key to getter/setter on VaultConfig.
type configKeyEntry struct {
	Get         func(cfg *VaultConfig) string
	Set         func(cfg *VaultConfig, value string)
	Default     string // human-readable default value for display
	Description string // human-readable description for display
}

// configKeyRegistry maps dot-notation config keys to VaultConfig struct fields.
var configKeyRegistry = map[string]configKeyEntry{
	"date_format": {
		Get:         func(cfg *VaultConfig) string { return cfg.DateFormat },
		Set:         func(cfg *VaultConfig, value string) { cfg.DateFormat = value },
		Default:     DefaultDateFormat,
		Description: "Date display format for properties (Go time layout)",
	},
	"datetime_format": {
		Get:         func(cfg *VaultConfig) string { return cfg.DatetimeFormat },
		Set:         func(cfg *VaultConfig, value string) { cfg.DatetimeFormat = value },
		Default:     DefaultDatetimeFormat,
		Description: "Date-time display format for properties (Go time layout)",
	},
	"cli.default_type": {
		Get:         func(cfg *VaultConfig) string { return cfg.CLI.DefaultType },
		Set:         func(cfg *VaultConfig, value string) { cfg.CLI.DefaultType = value },
		Default:     "",
		Description: "Default type for CLI commands when --type is omitted",
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
		Default:     "200",
		Description: "File watcher debounce delay in milliseconds",
	},
	"tui.stats_type_layout": {
		Get:         func(cfg *VaultConfig) string { return cfg.TUI.StatsTypeLayout },
		Set:         func(cfg *VaultConfig, value string) { cfg.TUI.StatsTypeLayout = value },
		Default:     "fullscreen",
		Description: "Stats type detail layout: fullscreen or popup",
	},
	"tui.toast.position": {
		Get:         func(cfg *VaultConfig) string { return cfg.TUI.Toast.Position },
		Set:         func(cfg *VaultConfig, value string) { cfg.TUI.Toast.Position = value },
		Default:     "bottom-right",
		Description: "Toast notification position (e.g., bottom-right)",
	},
	"tui.toast.duration_ms": {
		Get: func(cfg *VaultConfig) string {
			if cfg.TUI.Toast.DurationMs == 0 {
				return ""
			}
			return strconv.Itoa(cfg.TUI.Toast.DurationMs)
		},
		Set: func(cfg *VaultConfig, value string) {
			n, _ := strconv.Atoi(value)
			cfg.TUI.Toast.DurationMs = n
		},
		Default:     "3000",
		Description: "Toast display duration in milliseconds",
	},
	"tui.toast.dismiss_key": {
		Get:         func(cfg *VaultConfig) string { return cfg.TUI.Toast.DismissKey },
		Set:         func(cfg *VaultConfig, value string) { cfg.TUI.Toast.DismissKey = value },
		Default:     "esc",
		Description: "Key to dismiss toast notifications",
	},
	"tui.toast.show_warnings": {
		Get: func(cfg *VaultConfig) string {
			if cfg.TUI.Toast.ShowWarnings == nil {
				return ""
			}
			if *cfg.TUI.Toast.ShowWarnings {
				return "true"
			}
			return "false"
		},
		Set: func(cfg *VaultConfig, value string) {
			b := value == "true"
			cfg.TUI.Toast.ShowWarnings = &b
		},
		Default:     "true",
		Description: "Show warning toasts (true/false/unset)",
	},
	"tui.toast.show_success": {
		Get: func(cfg *VaultConfig) string {
			if cfg.TUI.Toast.ShowSuccess == nil {
				return ""
			}
			if *cfg.TUI.Toast.ShowSuccess {
				return "true"
			}
			return "false"
		},
		Set: func(cfg *VaultConfig, value string) {
			b := value == "true"
			cfg.TUI.Toast.ShowSuccess = &b
		},
		Default:     "false",
		Description: "Show success toasts (true/false/unset)",
	},
	"tui.theme.focus_border": {
		Get:     func(cfg *VaultConfig) string { return cfg.TUI.Theme.FocusBorder },
		Set:     func(cfg *VaultConfig, value string) { cfg.TUI.Theme.FocusBorder = value },
		Default: "63",
	},
	"tui.theme.wiki_link": {
		Get:     func(cfg *VaultConfig) string { return cfg.TUI.Theme.WikiLink },
		Set:     func(cfg *VaultConfig, value string) { cfg.TUI.Theme.WikiLink = value },
		Default: "33",
	},
	"tui.theme.heading": {
		Get:     func(cfg *VaultConfig) string { return cfg.TUI.Theme.Heading },
		Set:     func(cfg *VaultConfig, value string) { cfg.TUI.Theme.Heading = value },
		Default: "3",
	},
	"tui.theme.bold": {
		Get:     func(cfg *VaultConfig) string { return cfg.TUI.Theme.Bold },
		Set:     func(cfg *VaultConfig, value string) { cfg.TUI.Theme.Bold = value },
		Default: "",
	},
	"tui.theme.italic": {
		Get:     func(cfg *VaultConfig) string { return cfg.TUI.Theme.Italic },
		Set:     func(cfg *VaultConfig, value string) { cfg.TUI.Theme.Italic = value },
		Default: "",
	},
	"tui.theme.inline_code": {
		Get:     func(cfg *VaultConfig) string { return cfg.TUI.Theme.InlineCode },
		Set:     func(cfg *VaultConfig, value string) { cfg.TUI.Theme.InlineCode = value },
		Default: "245",
	},
	"tui.theme.code_block": {
		Get:     func(cfg *VaultConfig) string { return cfg.TUI.Theme.CodeBlock },
		Set:     func(cfg *VaultConfig, value string) { cfg.TUI.Theme.CodeBlock = value },
		Default: "245",
	},
	"tui.theme.link": {
		Get:     func(cfg *VaultConfig) string { return cfg.TUI.Theme.Link },
		Set:     func(cfg *VaultConfig, value string) { cfg.TUI.Theme.Link = value },
		Default: "33",
	},
	"tui.theme.blockquote": {
		Get:     func(cfg *VaultConfig) string { return cfg.TUI.Theme.Blockquote },
		Set:     func(cfg *VaultConfig, value string) { cfg.TUI.Theme.Blockquote = value },
		Default: "8",
	},
	"tui.theme.hrule": {
		Get:     func(cfg *VaultConfig) string { return cfg.TUI.Theme.HRule },
		Set:     func(cfg *VaultConfig, value string) { cfg.TUI.Theme.HRule = value },
		Default: "8",
	},
	"ai.default": {
		Get:         func(cfg *VaultConfig) string { return cfg.AI.Default },
		Set:         func(cfg *VaultConfig, value string) { cfg.AI.Default = value },
		Default:     "",
		Description: "Default AI provider name",
	},
	"ai.language": {
		Get:         func(cfg *VaultConfig) string { return cfg.AI.Language },
		Set:         func(cfg *VaultConfig, value string) { cfg.AI.Language = value },
		Default:     "",
		Description: "Language for AI-generated content",
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
		Default:     "false",
		Description: "Enable AI features (true/false)",
	},
	"ai.prompts.describe": {
		Get:         func(cfg *VaultConfig) string { return cfg.AI.Prompts.Describe },
		Set:         func(cfg *VaultConfig, value string) { cfg.AI.Prompts.Describe = value },
		Default:     "(built-in)",
		Description: "Custom prompt for AI describe action",
	},
	"ai.prompts.tag": {
		Get:         func(cfg *VaultConfig) string { return cfg.AI.Prompts.Tag },
		Set:         func(cfg *VaultConfig, value string) { cfg.AI.Prompts.Tag = value },
		Default:     "(built-in)",
		Description: "Custom prompt for AI tag action",
	},
	"ai.prompts.explore": {
		Get:         func(cfg *VaultConfig) string { return cfg.AI.Prompts.Explore },
		Set:         func(cfg *VaultConfig, value string) { cfg.AI.Prompts.Explore = value },
		Default:     "(built-in)",
		Description: "Custom prompt for AI explore action",
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
		Default:     "10",
		Description: "Number of sample objects for AI schema explore",
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
		Default:     "500",
		Description: "Max body characters for AI schema explore samples",
	},
	"web.theme": {
		Get:         func(cfg *VaultConfig) string { return cfg.Web.Theme },
		Set:         func(cfg *VaultConfig, value string) { cfg.Web.Theme = value },
		Default:     "warm",
		Description: "Web UI theme (warm, dark, or light)",
	},
}

// ConfigKeyInfo holds metadata about a config key for display purposes.
type ConfigKeyInfo struct {
	Key         string
	Description string
	Default     string
	Value       string
}

// ConfigKeysInfo returns metadata for all registered config keys,
// including their current values from the vault config.
func (v *Vault) ConfigKeysInfo() []ConfigKeyInfo {
	keys := ConfigKeys()
	infos := make([]ConfigKeyInfo, 0, len(keys))
	for _, key := range keys {
		entry := configKeyRegistry[key]
		infos = append(infos, ConfigKeyInfo{
			Key:         key,
			Description: entry.Description,
			Default:     entry.Default,
			Value:       entry.Get(v.config),
		})
	}
	return infos
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

// ConfigKeyDefault returns the human-readable default value for a config key.
// Returns an empty string if the key is unknown or has no default.
func ConfigKeyDefault(key string) string {
	entry, ok := configKeyRegistry[key]
	if !ok {
		return ""
	}
	return entry.Default
}
