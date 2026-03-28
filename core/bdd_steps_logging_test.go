package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/cucumber/godog"
)

type loggingContext struct {
	buf *bytes.Buffer
}

func initLoggingSteps(ctx *godog.ScenarioContext, lc *loggingContext) {
	ctx.Step(`^I initialize logging at debug level$`, lc.initDebug)
	ctx.Step(`^I initialize logging at warn level$`, lc.initWarn)
	ctx.Step(`^I log a debug message "([^"]*)"$`, lc.logDebug)
	ctx.Step(`^I log a warn message "([^"]*)"$`, lc.logWarn)
	ctx.Step(`^I log a debug message "([^"]*)" with attribute "([^"]*)" value "([^"]*)"$`, lc.logDebugWithAttr)
	ctx.Step(`^the log output should contain "([^"]*)"$`, lc.outputContains)
	ctx.Step(`^the log output should be empty$`, lc.outputEmpty)
	ctx.Step(`^the log output should be valid JSON$`, lc.outputValidJSON)
	ctx.Step(`^the log JSON should have field "([^"]*)" with value "([^"]*)"$`, lc.jsonFieldValue)
	ctx.Step(`^the log JSON should have field "([^"]*)"$`, lc.jsonFieldExists)
}

func (lc *loggingContext) initDebug() {
	lc.buf = &bytes.Buffer{}
	InitLogging(slog.LevelDebug, lc.buf)
}

func (lc *loggingContext) initWarn() {
	lc.buf = &bytes.Buffer{}
	InitLogging(slog.LevelWarn, lc.buf)
}

func (lc *loggingContext) logDebug(msg string) {
	slog.Debug(msg)
}

func (lc *loggingContext) logWarn(msg string) {
	slog.Warn(msg)
}

func (lc *loggingContext) logDebugWithAttr(msg, key, val string) {
	slog.Debug(msg, key, val)
}

func (lc *loggingContext) outputContains(expected string) error {
	if !strings.Contains(lc.buf.String(), expected) {
		return fmt.Errorf("expected log output to contain %q, got %q", expected, lc.buf.String())
	}
	return nil
}

func (lc *loggingContext) outputEmpty() error {
	if lc.buf.Len() != 0 {
		return fmt.Errorf("expected empty log output, got %q", lc.buf.String())
	}
	return nil
}

func (lc *loggingContext) outputValidJSON() error {
	lines := strings.TrimSpace(lc.buf.String())
	for _, line := range strings.Split(lines, "\n") {
		if line == "" {
			continue
		}
		if !json.Valid([]byte(line)) {
			return fmt.Errorf("log output is not valid JSON: %q", line)
		}
	}
	return nil
}

func (lc *loggingContext) jsonFieldValue(field, expected string) error {
	lines := strings.TrimSpace(lc.buf.String())
	lastLine := lines
	if idx := strings.LastIndex(lines, "\n"); idx >= 0 {
		lastLine = lines[idx+1:]
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(lastLine), &m); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	val, ok := m[field]
	if !ok {
		return fmt.Errorf("field %q not found in log JSON", field)
	}
	got := fmt.Sprintf("%v", val)
	if got != expected {
		return fmt.Errorf("expected field %q to be %q, got %q", field, expected, got)
	}
	return nil
}

func (lc *loggingContext) jsonFieldExists(field string) error {
	lines := strings.TrimSpace(lc.buf.String())
	lastLine := lines
	if idx := strings.LastIndex(lines, "\n"); idx >= 0 {
		lastLine = lines[idx+1:]
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(lastLine), &m); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	if _, ok := m[field]; !ok {
		return fmt.Errorf("field %q not found in log JSON", field)
	}
	return nil
}
