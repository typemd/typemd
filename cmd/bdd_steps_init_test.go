package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
	"github.com/typemd/typemd/core"
)

func initInitSteps(ctx *godog.ScenarioContext, cc *cmdContext) {
	ctx.Step(`^an empty directory$`, cc.anEmptyDirectory)
	ctx.Step(`^I run init(?:| with no-starters)$`, cc.iRunInitNoStarters)
	ctx.Step(`^the \.typemd directory should exist$`, cc.theTypemdDirectoryShouldExist)
	ctx.Step(`^the config should have default_type "([^"]*)"$`, cc.theConfigShouldHaveDefaultType)
	ctx.Step(`^a vault is already initialized$`, cc.aVaultIsAlreadyInitialized)
}

func (cc *cmdContext) anEmptyDirectory() error {
	dir, err := os.MkdirTemp("", "typemd-cmd-bdd-*")
	if err != nil {
		return err
	}
	cc.vaultDir = dir
	return nil
}

// iRunInitNoStarters always passes --no-starters to skip the interactive
// Bubble Tea picker which cannot run in test environments.
func (cc *cmdContext) iRunInitNoStarters() {
	cc.runCmd("init", "--no-starters")
}

func (cc *cmdContext) theTypemdDirectoryShouldExist() error {
	path := filepath.Join(cc.vaultDir, ".typemd")
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("expected .typemd directory to exist: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf(".typemd exists but is not a directory")
	}
	return nil
}

func (cc *cmdContext) theConfigShouldHaveDefaultType(expected string) error {
	path := filepath.Join(cc.vaultDir, ".typemd", "config.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config.yaml: %w", err)
	}
	needle := "default_type: " + expected
	if !strings.Contains(string(data), needle) {
		return fmt.Errorf("expected config to contain %q, got:\n%s", needle, data)
	}
	return nil
}

func (cc *cmdContext) aVaultIsAlreadyInitialized() error {
	dir, err := os.MkdirTemp("", "typemd-cmd-bdd-*")
	if err != nil {
		return err
	}
	cc.vaultDir = dir

	v := core.NewVault(dir)
	return v.Init()
}
