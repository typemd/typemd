package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const configFileName = "config.yaml"

// VaultConfig holds vault-level configuration loaded from .typemd/config.yaml.
type VaultConfig struct {
	CLI CLIConfig `yaml:"cli"`
}

// CLIConfig holds CLI-specific configuration.
type CLIConfig struct {
	DefaultType string `yaml:"default_type"`
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
}

// loadVaultConfig reads and parses the vault config file.
// Returns a zero-value VaultConfig if the file does not exist or is empty.
func loadVaultConfig(metaDir string) (*VaultConfig, error) {
	path := filepath.Join(metaDir, configFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
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
