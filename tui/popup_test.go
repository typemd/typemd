package tui

import (
	"strings"
	"testing"
)

func TestRenderPopup_ContainsContent(t *testing.T) {
	resetThemeDefaults()

	out := renderPopup("hello world", 80, 24, 0)
	if !strings.Contains(out, "hello world") {
		t.Error("popup output should contain the content")
	}
}

func TestRenderPopup_HasRoundedBorder(t *testing.T) {
	resetThemeDefaults()

	out := renderPopup("test", 80, 24, 0)
	// Rounded border uses ╭ as top-left corner
	if !strings.Contains(out, "╭") {
		t.Error("popup should use rounded border (╭)")
	}
}

func TestRenderPopup_IsCentered(t *testing.T) {
	resetThemeDefaults()

	out := renderPopup("X", 80, 24, 0)
	lines := strings.Split(out, "\n")

	// With 24 rows of output and a small popup, the first line should be
	// whitespace (centering padding).
	if len(lines) < 24 {
		t.Fatalf("expected at least 24 lines for centering, got %d", len(lines))
	}
	if strings.TrimRight(lines[0], " ") != "" {
		t.Error("first line should be blank (vertical centering)")
	}
}

func TestRenderPopup_RespectsWidth(t *testing.T) {
	resetThemeDefaults()

	narrow := renderPopup("content", 80, 24, 20)
	wide := renderPopup("content", 80, 24, 60)

	// The wide popup should produce a wider rendered box.
	narrowMaxW := maxLineWidth(narrow)
	wideMaxW := maxLineWidth(wide)
	if wideMaxW <= narrowMaxW {
		t.Errorf("wide popup (%d) should be wider than narrow popup (%d)", wideMaxW, narrowMaxW)
	}
}

func maxLineWidth(s string) int {
	max := 0
	for _, line := range strings.Split(s, "\n") {
		w := len(strings.TrimRight(line, " "))
		if w > max {
			max = w
		}
	}
	return max
}
