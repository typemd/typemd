package tui

import (
	"strings"
	"testing"

	"github.com/typemd/typemd/core"
)

func TestRenderList_GroupHeaderWithEmoji(t *testing.T) {
	groups := []typeGroup{
		{
			Name:     "book",
			Plural:   "books",
			Emoji:    "📚",
			Objects:  []*core.Object{{ID: "book/test", Type: "book", Filename: "test"}},
			Expanded: true,
		},
	}

	result := renderList(groups, 0, 0, true, 40, 10)

	if !strings.Contains(result, "📚 books (1)") {
		t.Errorf("expected group header to contain emoji, got: %s", result)
	}
}

func TestRenderList_GroupHeaderWithoutEmoji(t *testing.T) {
	groups := []typeGroup{
		{
			Name:     "note",
			Plural:   "note",
			Objects:  []*core.Object{{ID: "note/test", Type: "note", Filename: "test"}},
			Expanded: true,
		},
	}

	result := renderList(groups, 0, 0, true, 40, 10)

	if !strings.Contains(result, "▼ note (1)") {
		t.Errorf("expected group header without emoji, got: %s", result)
	}
	// Ensure no double space between arrow and name
	if strings.Contains(result, "▼  note") {
		t.Errorf("unexpected extra space in header without emoji, got: %s", result)
	}
}
