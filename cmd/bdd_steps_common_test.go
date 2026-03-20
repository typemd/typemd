package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cucumber/godog"
	"github.com/typemd/typemd/core"
)

// cmdContext holds shared state across BDD scenarios.
type cmdContext struct {
	vaultDir string
	vault    *core.Vault
	stdout   string
	lastErr  error

	// Track created objects for reference in later steps.
	createdObjectIDs []string
}

func newCmdContext() *cmdContext {
	return &cmdContext{}
}

// setupVault creates a temp dir, initializes a vault, writes a book type
// schema, opens the vault, and syncs the index.
func (cc *cmdContext) setupVault() error {
	dir, err := os.MkdirTemp("", "typemd-cmd-bdd-*")
	if err != nil {
		return err
	}
	cc.vaultDir = dir

	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		return err
	}

	// Write config with default_type: page (like tmd init does).
	cfg := &core.VaultConfig{
		CLI: core.CLIConfig{DefaultType: core.PageTypeName},
	}
	if err := v.WriteConfig(cfg); err != nil {
		return err
	}

	// Write a book type schema for testing.
	if err := os.WriteFile(v.TypesDir()+"/book.yaml", []byte(`name: book
emoji: "\U0001F4DA"
plural: books
properties:
  - name: author
    type: string
  - name: rating
    type: number
`), 0644); err != nil {
		return err
	}

	if err := v.Open(); err != nil {
		return err
	}
	if _, err := v.SyncIndex(); err != nil {
		v.Close()
		return err
	}

	cc.vault = v
	return nil
}

// runCmd executes a CLI command by setting rootCmd args and capturing stdout.
// Note: os.Stdout is redirected because commands use fmt.Print* (not
// cmd.OutOrStdout()), so Cobra's SetOut() cannot capture the output.
func (cc *cmdContext) runCmd(args ...string) {
	resetAllFlags()
	vaultPath = cc.vaultDir

	rootCmd.SetArgs(args)

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		cc.lastErr = fmt.Errorf("os.Pipe: %w", err)
		return
	}
	os.Stdout = w

	cc.lastErr = rootCmd.Execute()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	cc.stdout = buf.String()
}

// runCmdAndTrack executes a CLI command and tracks the created object ID
// from the "Created <id>" output line.
func (cc *cmdContext) runCmdAndTrack(args ...string) {
	cc.runCmd(args...)
	if cc.lastErr == nil {
		line := strings.TrimSpace(cc.stdout)
		if id, ok := strings.CutPrefix(line, "Created "); ok {
			cc.createdObjectIDs = append(cc.createdObjectIDs, id)
		}
	}
}

// resetAllFlags resets all package-level flag variables and Cobra local flags
// to prevent state leakage between scenarios.
func resetAllFlags() {
	vaultPath = ""
	readOnly = false
	reindex = false
	templateFlag = ""
	typeFlag = ""
	noStarters = false
	searchJSON = false

	// Reset Cobra local flags that use cmd.Flags().GetBool() instead of
	// package-level vars. Without this, flags set in one scenario leak
	// into subsequent scenarios.
	objectListCmd.Flags().Set("json", "false")
}

func initializeScenario(sc *godog.ScenarioContext) {
	cc := newCmdContext()

	sc.After(func(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		if cc.vault != nil {
			cc.vault.Close()
		}
		if cc.vaultDir != "" {
			os.RemoveAll(cc.vaultDir)
		}
		return ctx, nil
	})

	initCommonSteps(sc, cc)
	initInitSteps(sc, cc)
	initObjectSteps(sc, cc)
	initSearchSteps(sc, cc)
}

func initCommonSteps(ctx *godog.ScenarioContext, cc *cmdContext) {
	ctx.Step(`^a vault is ready$`, cc.aVaultIsReady)
	ctx.Step(`^the output should contain "([^"]*)"$`, cc.theOutputShouldContain)
	ctx.Step(`^the output should not contain "([^"]*)"$`, cc.theOutputShouldNotContain)
	ctx.Step(`^the output should be empty$`, cc.theOutputShouldBeEmpty)
	ctx.Step(`^the command should succeed$`, cc.theCommandShouldSucceed)
	ctx.Step(`^the command should fail$`, cc.theCommandShouldFail)
	ctx.Step(`^the command should fail with "([^"]*)"$`, cc.theCommandShouldFailWith)
}

func (cc *cmdContext) aVaultIsReady() error {
	return cc.setupVault()
}

func (cc *cmdContext) theOutputShouldContain(expected string) error {
	if !strings.Contains(cc.stdout, expected) {
		return fmt.Errorf("expected output to contain %q, got:\n%s", expected, cc.stdout)
	}
	return nil
}

func (cc *cmdContext) theOutputShouldNotContain(unexpected string) error {
	if strings.Contains(cc.stdout, unexpected) {
		return fmt.Errorf("expected output NOT to contain %q, got:\n%s", unexpected, cc.stdout)
	}
	return nil
}

func (cc *cmdContext) theOutputShouldBeEmpty() error {
	trimmed := strings.TrimSpace(cc.stdout)
	if trimmed != "" {
		return fmt.Errorf("expected empty output, got:\n%s", cc.stdout)
	}
	return nil
}

func (cc *cmdContext) theCommandShouldSucceed() error {
	if cc.lastErr != nil {
		return fmt.Errorf("expected command to succeed, got error: %v", cc.lastErr)
	}
	return nil
}

func (cc *cmdContext) theCommandShouldFail() error {
	if cc.lastErr == nil {
		return fmt.Errorf("expected command to fail, but it succeeded")
	}
	return nil
}

func (cc *cmdContext) theCommandShouldFailWith(expected string) error {
	if cc.lastErr == nil {
		return fmt.Errorf("expected command to fail with %q, but it succeeded", expected)
	}
	if !strings.Contains(cc.lastErr.Error(), expected) {
		return fmt.Errorf("expected error to contain %q, got: %v", expected, cc.lastErr)
	}
	return nil
}
