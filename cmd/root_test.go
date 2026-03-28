package cmd

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/typemd/typemd/core"
)

func TestDebugFlag_EnablesLogging(t *testing.T) {
	buf := &bytes.Buffer{}
	core.InitLogging(slog.LevelDebug, buf)
	slog.Debug("test-debug-message")

	if !strings.Contains(buf.String(), "test-debug-message") {
		t.Errorf("expected debug message in output, got %q", buf.String())
	}
}

func TestDebugFlag_WarnLevelSilent(t *testing.T) {
	buf := &bytes.Buffer{}
	core.InitLogging(slog.LevelWarn, buf)
	slog.Debug("should-not-appear")
	slog.Info("should-not-appear-either")

	if buf.Len() != 0 {
		t.Errorf("expected no output at warn level, got %q", buf.String())
	}
}

func TestDebugFlag_Registered(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("debug")
	if flag == nil {
		t.Fatal("expected --debug persistent flag to be registered")
	}
	if flag.DefValue != "false" {
		t.Errorf("expected default value 'false', got %q", flag.DefValue)
	}
}
