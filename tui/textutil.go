package tui

import (
	"strings"

	"github.com/mattn/go-runewidth"

	"charm.land/lipgloss/v2"
)

// softWrapLines wraps each line individually, preserving leading indentation on continuation lines.
func softWrapLines(content string, width int) string {
	lines := strings.Split(content, "\n")
	var result []string
	for _, line := range lines {
		if lipgloss.Width(line) <= width {
			result = append(result, line)
			continue
		}
		// Detect leading whitespace
		trimmed := strings.TrimLeft(line, " ")
		indent := line[:len(line)-len(trimmed)]
		wrapped := lipgloss.NewStyle().Width(width - lipgloss.Width(indent)).Render(trimmed)
		for _, wl := range strings.Split(wrapped, "\n") {
			result = append(result, indent+wl)
		}
	}
	return strings.Join(result, "\n")
}

// truncate shortens a string to fit within maxLen cells, adding ellipsis if needed.
func truncate(s string, maxLen int) string {
	return runewidth.Truncate(s, maxLen, "…")
}

// padRight pads a string with spaces to fill exactly width display cells.
// Delegates to runewidth.FillRight for CJK/emoji-aware padding.
func padRight(s string, width int) string {
	return runewidth.FillRight(s, width)
}

// padEmoji strips the variation selector (U+FE0F) and pads the emoji
// to a consistent 2-cell display width for terminal alignment.
func padEmoji(emoji string) string {
	display := strings.ReplaceAll(emoji, "\uFE0F", "")
	ew := runewidth.StringWidth(display)
	if ew < 2 {
		return display + strings.Repeat(" ", 2-ew)
	}
	return display
}
