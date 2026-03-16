package core

import (
	"fmt"
	"os"
	"path/filepath"

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
