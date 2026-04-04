package cmd

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/cucumber/godog"
)

func initLogSteps(ctx *godog.ScenarioContext, cc *cmdContext) {
	ctx.Step(`^I run log "([^"]*)"$`, cc.iRunLog)
	ctx.Step(`^I run log "([^"]*)" with oneline$`, cc.iRunLogOneline)
	ctx.Step(`^I run log with the first created object$`, cc.iRunLogWithFirstObject)
	ctx.Step(`^the vault is a git repository with committed objects$`, cc.theVaultIsGitRepoWithCommits)
	ctx.Step(`^the vault is a git repository$`, cc.theVaultIsGitRepo)
}

func (cc *cmdContext) iRunLog(objectID string) {
	cc.runCmd("log", objectID)
}

func (cc *cmdContext) iRunLogOneline(objectID string) {
	cc.runCmd("log", "--oneline", objectID)
}

func (cc *cmdContext) iRunLogWithFirstObject() {
	objs, err := cc.vault.QueryObjects(nil)
	if err != nil || len(objs) == 0 {
		cc.lastErr = fmt.Errorf("no objects in vault")
		return
	}
	cc.runCmd("log", objs[0].ID)
}

func (cc *cmdContext) theVaultIsGitRepo() error {
	cmd := exec.Command("git", "init")
	cmd.Dir = cc.vaultDir
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	return cmd.Run()
}

func (cc *cmdContext) theVaultIsGitRepoWithCommits() error {
	if err := cc.theVaultIsGitRepo(); err != nil {
		return err
	}

	// Configure git user for commits.
	for _, args := range [][]string{
		{"config", "user.email", "test@test.com"},
		{"config", "user.name", "Test"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Dir = cc.vaultDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git config: %w", err)
		}
	}

	// Add and commit all objects.
	cmd := exec.Command("git", "add", "objects/")
	cmd.Dir = cc.vaultDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git add: %w", err)
	}

	cmd = exec.Command("git", "commit", "-m", "add objects")
	cmd.Dir = cc.vaultDir
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	return cmd.Run()
}
