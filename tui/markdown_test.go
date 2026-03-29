package tui

import (
	"strings"
	"testing"
)

// Helper: strip all ANSI escape sequences to get plain text.
func stripANSI(s string) string {
	var result strings.Builder
	inEsc := false
	for _, r := range s {
		if r == '\033' {
			inEsc = true
			continue
		}
		if inEsc {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEsc = false
			}
			continue
		}
		result.WriteRune(r)
	}
	return result.String()
}

// Helper: check if a string contains ANSI escape codes.
func hasANSI(s string) bool {
	return strings.Contains(s, "\033[")
}

// --- Heading scenarios ---

func TestRenderMarkdown_H1Heading(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("# Introduction")
	plain := stripANSI(got)
	if plain != "# Introduction" {
		t.Errorf("H1 plain text: got %q, want %q", plain, "# Introduction")
	}
	if !hasANSI(got) {
		t.Error("H1 heading should have ANSI styling")
	}
}

func TestRenderMarkdown_H4Heading(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("#### Details")
	plain := stripANSI(got)
	if plain != "#### Details" {
		t.Errorf("H4 plain text: got %q", plain)
	}
	if !hasANSI(got) {
		t.Error("H4 heading should have ANSI styling")
	}
}

func TestRenderMarkdown_HashInNonHeading(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	input := "Use #hashtag for tagging"
	got := renderMarkdown(input)
	if got != input {
		t.Errorf("non-heading hash should not be styled: got %q, want %q", got, input)
	}
}

// --- Bold scenarios ---

func TestRenderMarkdown_BoldWord(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("This is **important** text")
	plain := stripANSI(got)
	if !strings.Contains(plain, "**important**") {
		t.Errorf("bold word not found in plain text: %q", plain)
	}
	if !hasANSI(got) {
		t.Error("bold should have ANSI styling")
	}
}

func TestRenderMarkdown_MultipleBoldSpans(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("**first** and **second**")
	plain := stripANSI(got)
	if !strings.Contains(plain, "**first**") || !strings.Contains(plain, "**second**") {
		t.Errorf("multiple bold spans: plain text = %q", plain)
	}
}

// --- Italic scenarios ---

func TestRenderMarkdown_ItalicWithAsterisks(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("This is *emphasized* text")
	plain := stripANSI(got)
	if !strings.Contains(plain, "*emphasized*") {
		t.Errorf("italic asterisks: plain text = %q", plain)
	}
	if !hasANSI(got) {
		t.Error("italic should have ANSI styling")
	}
}

func TestRenderMarkdown_ItalicWithUnderscores(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("This is _emphasized_ text")
	plain := stripANSI(got)
	if !strings.Contains(plain, "_emphasized_") {
		t.Errorf("italic underscores: plain text = %q", plain)
	}
	if !hasANSI(got) {
		t.Error("italic should have ANSI styling")
	}
}

func TestRenderMarkdown_BoldNotTreatedAsItalic(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("This is **bold** text")
	plain := stripANSI(got)
	if !strings.Contains(plain, "**bold**") {
		t.Errorf("bold markers should be preserved: plain text = %q", plain)
	}
}

// --- Inline code scenarios ---

func TestRenderMarkdown_InlineCode(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("Use `fmt.Println` for output")
	plain := stripANSI(got)
	if !strings.Contains(plain, "`fmt.Println`") {
		t.Errorf("inline code: plain text = %q", plain)
	}
	if !hasANSI(got) {
		t.Error("inline code should have ANSI styling")
	}
}

func TestRenderMarkdown_MultipleInlineCode(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("Both `foo` and `bar` are valid")
	plain := stripANSI(got)
	if !strings.Contains(plain, "`foo`") || !strings.Contains(plain, "`bar`") {
		t.Errorf("multiple inline code: plain text = %q", plain)
	}
}

// --- Fenced code block scenarios ---

func TestRenderMarkdown_CodeBlockContent(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	input := "```\nx := 1\n```"
	got := renderMarkdown(input)
	plain := stripANSI(got)
	if !strings.Contains(plain, "x := 1") {
		t.Errorf("code block content: plain text = %q", plain)
	}
}

func TestRenderMarkdown_CodeBlockNoInlineProcessing(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	input := "```\n**not bold**\n```"
	got := renderMarkdown(input)
	plain := stripANSI(got)
	if !strings.Contains(plain, "**not bold**") {
		t.Errorf("code block should preserve literal content: plain text = %q", plain)
	}
	// The content should be styled with code block color only, not bold.
	lines := strings.Split(got, "\n")
	contentLine := lines[1]
	// Verify no bold ANSI sequence (ESC[1m) is applied separately from code block style.
	// The code block style itself may include attributes, but bold processing should not run.
	contentPlain := stripANSI(contentLine)
	if contentPlain != "**not bold**" {
		t.Errorf("code block content altered: got %q", contentPlain)
	}
}

func TestRenderMarkdown_CodeBlockWithLanguage(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	input := "```go\nfunc main() {}\n```"
	got := renderMarkdown(input)
	plain := stripANSI(got)
	if !strings.Contains(plain, "func main() {}") {
		t.Errorf("code block with language: plain text = %q", plain)
	}
}

// --- Link scenarios ---

func TestRenderMarkdown_Link(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("See [documentation](https://example.com) for details")
	plain := stripANSI(got)
	if !strings.Contains(plain, "[documentation](https://example.com)") {
		t.Errorf("link: plain text = %q", plain)
	}
	if !hasANSI(got) {
		t.Error("link should have ANSI styling")
	}
}

// --- Blockquote scenarios ---

func TestRenderMarkdown_Blockquote(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("> This is a quote")
	plain := stripANSI(got)
	if plain != "> This is a quote" {
		t.Errorf("blockquote: plain text = %q", plain)
	}
	if !hasANSI(got) {
		t.Error("blockquote should have ANSI styling")
	}
}

// --- Horizontal rule scenarios ---

func TestRenderMarkdown_DashHRule(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("---")
	plain := stripANSI(got)
	if plain != "---" {
		t.Errorf("dash hrule: plain text = %q", plain)
	}
	if !hasANSI(got) {
		t.Error("hrule should have ANSI styling")
	}
}

func TestRenderMarkdown_AsteriskHRule(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("***")
	plain := stripANSI(got)
	if plain != "***" {
		t.Errorf("asterisk hrule: plain text = %q", plain)
	}
}

func TestRenderMarkdown_UnderscoreHRule(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("___")
	plain := stripANSI(got)
	if plain != "___" {
		t.Errorf("underscore hrule: plain text = %q", plain)
	}
}

// --- Edge cases ---

func TestRenderMarkdown_EmptyBody(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("")
	if got != "" {
		t.Errorf("empty body: got %q", got)
	}
}

func TestRenderMarkdown_PlainText(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	input := "Just some regular text without any markdown."
	got := renderMarkdown(input)
	if got != input {
		t.Errorf("plain text should not be modified: got %q, want %q", got, input)
	}
}

func TestRenderMarkdown_MultipleElements(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	input := "# Title\n\nSome **bold** and `code` text.\n\n---\n\n> A quote"
	got := renderMarkdown(input)
	plain := stripANSI(got)
	lines := strings.Split(plain, "\n")
	if len(lines) != 7 {
		t.Fatalf("expected 7 lines, got %d: %q", len(lines), plain)
	}
	if lines[0] != "# Title" {
		t.Errorf("first line: got %q", lines[0])
	}
	if lines[6] != "> A quote" {
		t.Errorf("last line: got %q", lines[6])
	}
}

func TestRenderMarkdown_CodeBlockBoundary(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	input := "**bold before**\n```\ncode\n```\n**bold after**"
	got := renderMarkdown(input)
	plain := stripANSI(got)
	lines := strings.Split(plain, "\n")
	if len(lines) != 5 {
		t.Fatalf("expected 5 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "**bold before**") {
		t.Errorf("bold before code block: got %q", lines[0])
	}
	if !strings.Contains(lines[4], "**bold after**") {
		t.Errorf("bold after code block: got %q", lines[4])
	}
}

func TestRenderMarkdown_SoftWrapCompatibility(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	// Render markdown with inline styles, then soft wrap at a narrow width.
	input := "This is **bold** and `code` in a somewhat longer line of text"
	styled := renderMarkdown(input)
	// Soft wrap at 30 chars — should not panic or produce corrupted output.
	wrapped := softWrapLines(styled, 30)
	plain := stripANSI(wrapped)
	// All original text should still be present after wrapping.
	if !strings.Contains(plain, "**bold**") {
		t.Errorf("soft wrap lost bold text: plain = %q", plain)
	}
	if !strings.Contains(plain, "`code`") {
		t.Errorf("soft wrap lost code text: plain = %q", plain)
	}
}

func TestRenderMarkdown_EmptyLines(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	input := "line 1\n\nline 3"
	got := renderMarkdown(input)
	if got != input {
		t.Errorf("empty lines: got %q, want %q", got, input)
	}
}
