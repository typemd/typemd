package tui

import (
	"regexp"
	"strings"

	"charm.land/lipgloss/v2"
)

// Markdown regex patterns.
var (
	reHeading    = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)
	reBold       = regexp.MustCompile(`\*\*(.+?)\*\*`)
	reItalicStar = regexp.MustCompile(`(?:^|[^*])\*([^*]+?)\*(?:[^*]|$)`)
	reItalicUS   = regexp.MustCompile(`(?:^|[^_])_([^_]+?)_(?:[^_]|$)`)
	reInlineCode = regexp.MustCompile("`([^`]+)`")
	reLink       = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	reBlockquote = regexp.MustCompile(`^>\s?(.*)$`)
	reHRule      = regexp.MustCompile(`^(?:-{3,}|\*{3,}|_{3,})\s*$`)
	reFence      = regexp.MustCompile("^```")
)

// renderMarkdown applies syntax highlighting to markdown body text.
// It processes line-by-line with fenced code block state tracking.
func renderMarkdown(body string) string {
	lines := strings.Split(body, "\n")
	inCodeBlock := false
	result := make([]string, 0, len(lines))

	for _, line := range lines {
		// Toggle code block state on fence lines.
		if reFence.MatchString(line) {
			inCodeBlock = !inCodeBlock
			result = append(result, mdCodeBlockStyle.Render(line))
			continue
		}

		// Inside code block: style entire line, no inline processing.
		if inCodeBlock {
			result = append(result, mdCodeBlockStyle.Render(line))
			continue
		}

		result = append(result, styleMarkdownLine(line))
	}

	return strings.Join(result, "\n")
}

// styleMarkdownLine applies markdown styling to a single line outside of code blocks.
func styleMarkdownLine(line string) string {
	// Horizontal rule (full-line match, must check before heading).
	if reHRule.MatchString(line) {
		return mdHRuleStyle.Render(line)
	}

	// Heading (full-line match).
	if reHeading.MatchString(line) {
		return mdHeadingStyle.Render(line)
	}

	// Blockquote (full-line match).
	if reBlockquote.MatchString(line) {
		return mdBlockquoteStyle.Render(line)
	}

	// Inline styles applied in order: inline code first (to protect code content),
	// then bold, italic, and links.
	line = reInlineCode.ReplaceAllStringFunc(line, func(m string) string {
		return mdInlineCodeStyle.Render(m)
	})

	line = reBold.ReplaceAllStringFunc(line, func(m string) string {
		return mdBoldStyle.Render(m)
	})

	line = applyDelimitedStyle(line, reItalicUS, "_", mdItalicStyle)
	line = applyDelimitedStyle(line, reItalicStar, "*", mdItalicStyle)

	line = reLink.ReplaceAllStringFunc(line, func(m string) string {
		return mdLinkStyle.Render(m)
	})

	return line
}

// applyDelimitedStyle applies styling to text between delimiter characters.
// The regex may capture surrounding context characters to avoid false matches,
// so we extract the actual delimited span using the delimiter character.
func applyDelimitedStyle(line string, re *regexp.Regexp, delim string, style lipgloss.Style) string {
	return re.ReplaceAllStringFunc(line, func(m string) string {
		idx := strings.Index(m, delim)
		lastIdx := strings.LastIndex(m, delim)
		if idx == -1 || lastIdx == idx {
			return m
		}
		prefix := m[:idx]
		styled := m[idx : lastIdx+len(delim)]
		suffix := m[lastIdx+len(delim):]
		return prefix + style.Render(styled) + suffix
	})
}
