package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cucumber/godog"
)

// watchResult holds output from a brief watch run.
type watchResult struct {
	stdout string
	err    error
}

func initWatchValidateSteps(ctx *godog.ScenarioContext, cc *cmdContext) {
	var wr watchResult

	ctx.Step(`^I run watch validate briefly$`, func() {
		wr = runWatchBriefly(cc)
	})
	ctx.Step(`^the watch output should contain "([^"]*)"$`, func(expected string) error {
		if !strings.Contains(wr.stdout, expected) {
			return fmt.Errorf("expected watch output to contain %q, got:\n%s", expected, wr.stdout)
		}
		return nil
	})
	ctx.Step(`^a broken schema exists$`, func() error {
		dir := cc.vaultDir + "/types/broken"
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		return os.WriteFile(dir+"/schema.yaml", []byte(`name: broken
properties:
  - name: foo
    type: invalid_type
`), 0644)
	})
}

// runWatchBriefly runs the validate --watch command with a short context timeout
// to capture initial output without blocking indefinitely.
func runWatchBriefly(cc *cmdContext) watchResult {
	resetAllFlags()
	vaultPath = cc.vaultDir
	watchFlag = true

	rootCmd.SetArgs([]string{"type", "validate", "--watch"})

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return watchResult{err: err}
	}
	os.Stdout = w

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	rootCmd.SetContext(ctx)
	cmdErr := rootCmd.Execute()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Context timeout is expected — treat as success
	if cmdErr != nil && ctx.Err() == nil {
		return watchResult{stdout: output, err: cmdErr}
	}

	return watchResult{stdout: output}
}
