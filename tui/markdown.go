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
	reItalicStar = regexp.MustCompile(`(^|[^*\w])\*([^*]+?)\*([^*\w]|$)`)
	reItalicUS   = regexp.MustCompile(`(^|[^_\w])_([^_]+?)_([^_\w]|$)`)
	reInlineCode = regexp.MustCompile("`([^`]+)`")
	reLink       = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	reBlockquote = regexp.MustCompile(`^>\s?(.*)$`)
	reHRule      = regexp.MustCompile(`^(?:-{3,}|\*{3,}|_{3,})\s*$`)
	reFence      = regexp.MustCompile("^```")
)

// renderMarkdown transforms markdown body text for view mode:
// syntax markers are hidden and lipgloss styles are applied.
func renderMarkdown(body string) string {
	lines := strings.Split(body, "\n")
	inCodeBlock := false
	result := make([]string, 0, len(lines))

	for _, line := range lines {
		// Toggle code block state on fence lines.
		if reFence.MatchString(line) {
			inCodeBlock = !inCodeBlock
			// Hide fence lines in view mode.
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
// Syntax markers are stripped and content is styled.
func styleMarkdownLine(line string) string {
	// Horizontal rule (full-line match, must check before heading).
	if reHRule.MatchString(line) {
		return mdHRuleStyle.Render("────────────────────")
	}

	// Heading: strip # markers, style content.
	if m := reHeading.FindStringSubmatch(line); m != nil {
		content := styleInline(m[2])
		return mdHeadingStyle.Render(content)
	}

	// Blockquote: strip > marker, add │ prefix.
	if m := reBlockquote.FindStringSubmatch(line); m != nil {
		content := styleInline(m[1])
		return mdBlockquoteStyle.Render("│ " + content)
	}

	return styleInline(line)
}

// styleInline applies inline markdown styling, stripping syntax markers.
// Processing order: inline code → bold → italic → links.
func styleInline(line string) string {
	// Inline code: strip backticks, style content.
	line = reInlineCode.ReplaceAllStringFunc(line, func(m string) string {
		sub := reInlineCode.FindStringSubmatch(m)
		return mdInlineCodeStyle.Render(sub[1])
	})

	// Bold: strip ** markers, style content.
	line = reBold.ReplaceAllStringFunc(line, func(m string) string {
		sub := reBold.FindStringSubmatch(m)
		return mdBoldStyle.Render(sub[1])
	})

	// Italic: strip delimiters, preserve boundary characters ($1, $3).
	line = applyItalicStyle(line, reItalicUS, mdItalicStyle)
	line = applyItalicStyle(line, reItalicStar, mdItalicStyle)

	// Links: strip [](url) syntax, show only display text.
	line = reLink.ReplaceAllStringFunc(line, func(m string) string {
		sub := reLink.FindStringSubmatch(m)
		return mdLinkStyle.Render(sub[1])
	})

	return line
}

// applyItalicStyle replaces italic matches using capturing groups.
// The regex captures boundary chars in groups 1 and 3, content in group 2.
// Replacements are applied in reverse order to preserve indices.
func applyItalicStyle(line string, re *regexp.Regexp, style lipgloss.Style) string {
	locs := re.FindAllStringSubmatchIndex(line, -1)
	if len(locs) == 0 {
		return line
	}
	var b strings.Builder
	b.Grow(len(line))
	last := 0
	for _, loc := range locs {
		b.WriteString(line[last:loc[0]])
		prefix := line[loc[2]:loc[3]]   // group 1: boundary before
		content := line[loc[4]:loc[5]]  // group 2: italic content
		suffix := line[loc[6]:loc[7]]   // group 3: boundary after
		b.WriteString(prefix)
		b.WriteString(style.Render(content))
		b.WriteString(suffix)
		last = loc[1]
	}
	b.WriteString(line[last:])
	return b.String()
}
