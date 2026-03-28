package core

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"strings"
	"testing"
)

func TestInitLogging_NilWriter(t *testing.T) {
	// io.Discard should work as output
	InitLogging(slog.LevelWarn, io.Discard)
	// Should not panic
	slog.Info("test")
}

func TestInitLogging_MultipleCalls(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	InitLogging(slog.LevelDebug, buf1)
	slog.Debug("first")

	InitLogging(slog.LevelDebug, buf2)
	slog.Debug("second")

	if !strings.Contains(buf1.String(), "first") {
		t.Error("expected buf1 to contain 'first'")
	}
	if strings.Contains(buf1.String(), "second") {
		t.Error("expected buf1 NOT to contain 'second'")
	}
	if !strings.Contains(buf2.String(), "second") {
		t.Error("expected buf2 to contain 'second'")
	}
}

func TestInitLogging_JSONFormat(t *testing.T) {
	buf := &bytes.Buffer{}
	InitLogging(slog.LevelDebug, buf)
	slog.Debug("test-msg", "key1", "val1")

	var m map[string]any
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("log output is not valid JSON: %v", err)
	}

	for _, field := range []string{"time", "level", "msg", "key1"} {
		if _, ok := m[field]; !ok {
			t.Errorf("expected field %q in JSON output", field)
		}
	}

	if m["msg"] != "test-msg" {
		t.Errorf("expected msg 'test-msg', got %v", m["msg"])
	}
	if m["key1"] != "val1" {
		t.Errorf("expected key1 'val1', got %v", m["key1"])
	}
}

func TestInitLogging_LevelFiltering(t *testing.T) {
	buf := &bytes.Buffer{}
	InitLogging(slog.LevelWarn, buf)

	slog.Debug("should-not-appear")
	slog.Info("should-not-appear")
	if buf.Len() != 0 {
		t.Errorf("expected no output at warn level for debug/info, got %q", buf.String())
	}

	slog.Warn("should-appear")
	if !strings.Contains(buf.String(), "should-appear") {
		t.Error("expected warn message to appear")
	}
}
