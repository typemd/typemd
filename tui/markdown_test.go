package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

// --- Heading scenarios ---

func TestRenderMarkdown_H1Heading(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("# Introduction")
	plain := ansi.Strip(got)
	if plain != "Introduction" {
		t.Errorf("H1 plain text: got %q, want %q", plain, "Introduction")
	}
	if got == plain {
		t.Error("H1 heading should have ANSI styling")
	}
}

func TestRenderMarkdown_H4Heading(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("#### Details")
	plain := ansi.Strip(got)
	if plain != "Details" {
		t.Errorf("H4 plain text: got %q, want %q", plain, "Details")
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
	plain := ansi.Strip(got)
	if !strings.Contains(plain, "important") {
		t.Errorf("bold word not found in plain text: %q", plain)
	}
	if strings.Contains(plain, "**") {
		t.Error("** markers should be stripped")
	}
}

func TestRenderMarkdown_MultipleBoldSpans(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("**first** and **second**")
	plain := ansi.Strip(got)
	if !strings.Contains(plain, "first") || !strings.Contains(plain, "second") {
		t.Errorf("multiple bold spans: plain text = %q", plain)
	}
	if strings.Contains(plain, "**") {
		t.Error("** markers should be stripped")
	}
}

// --- Italic scenarios ---

func TestRenderMarkdown_ItalicWithAsterisks(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("This is *emphasized* text")
	plain := ansi.Strip(got)
	if !strings.Contains(plain, "emphasized") {
		t.Errorf("italic asterisks: plain text = %q", plain)
	}
	if strings.Contains(plain, "*") {
		t.Error("* markers should be stripped")
	}
}

func TestRenderMarkdown_ItalicWithUnderscores(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("This is _emphasized_ text")
	plain := ansi.Strip(got)
	if !strings.Contains(plain, "emphasized") {
		t.Errorf("italic underscores: plain text = %q", plain)
	}
	if strings.Contains(plain, "_") {
		t.Error("_ markers should be stripped")
	}
}

func TestRenderMarkdown_BoldNotTreatedAsItalic(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("This is **bold** text")
	plain := ansi.Strip(got)
	if !strings.Contains(plain, "bold") {
		t.Errorf("bold text should be present: plain text = %q", plain)
	}
	if strings.Contains(plain, "**") {
		t.Error("** markers should be stripped")
	}
}

// --- Inline code scenarios ---

func TestRenderMarkdown_InlineCode(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("Use `fmt.Println` for output")
	plain := ansi.Strip(got)
	if !strings.Contains(plain, "fmt.Println") {
		t.Errorf("inline code: plain text = %q", plain)
	}
	if strings.Contains(plain, "`") {
		t.Error("backtick markers should be stripped")
	}
}

func TestRenderMarkdown_MultipleInlineCode(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("Both `foo` and `bar` are valid")
	plain := ansi.Strip(got)
	if !strings.Contains(plain, "foo") || !strings.Contains(plain, "bar") {
		t.Errorf("multiple inline code: plain text = %q", plain)
	}
	if strings.Contains(plain, "`") {
		t.Error("backtick markers should be stripped")
	}
}

// --- Fenced code block scenarios ---

func TestRenderMarkdown_CodeBlockContent(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	input := "```\nx := 1\n```"
	got := renderMarkdown(input)
	plain := ansi.Strip(got)
	if !strings.Contains(plain, "x := 1") {
		t.Errorf("code block content: plain text = %q", plain)
	}
	if strings.Contains(plain, "```") {
		t.Error("fence markers should be hidden")
	}
}

func TestRenderMarkdown_CodeBlockNoInlineProcessing(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	input := "```\n**not bold**\n```"
	got := renderMarkdown(input)
	plain := ansi.Strip(got)
	if !strings.Contains(plain, "**not bold**") {
		t.Errorf("code block should preserve literal content: plain text = %q", plain)
	}
}

func TestRenderMarkdown_CodeBlockWithLanguage(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	input := "```go\nfunc main() {}\n```"
	got := renderMarkdown(input)
	plain := ansi.Strip(got)
	if !strings.Contains(plain, "func main() {}") {
		t.Errorf("code block with language: plain text = %q", plain)
	}
	if strings.Contains(plain, "```") {
		t.Error("fence markers should be hidden")
	}
}

// --- Link scenarios ---

func TestRenderMarkdown_Link(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("See [documentation](https://example.com) for details")
	plain := ansi.Strip(got)
	if !strings.Contains(plain, "documentation") {
		t.Errorf("link text should be present: %q", plain)
	}
	if strings.Contains(plain, "https://example.com") {
		t.Error("URL should be hidden")
	}
	if strings.Contains(plain, "[") || strings.Contains(plain, "]") {
		t.Error("brackets should be stripped")
	}
}

// --- Blockquote scenarios ---

func TestRenderMarkdown_Blockquote(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("> This is a quote")
	plain := ansi.Strip(got)
	if !strings.Contains(plain, "│ This is a quote") {
		t.Errorf("blockquote should show │ prefix: plain text = %q", plain)
	}
}

// --- Horizontal rule scenarios ---

func TestRenderMarkdown_DashHRule(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("---")
	plain := ansi.Strip(got)
	if !strings.Contains(plain, "────") {
		t.Errorf("dash hrule should render as line: plain text = %q", plain)
	}
}

func TestRenderMarkdown_AsteriskHRule(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("***")
	plain := ansi.Strip(got)
	if !strings.Contains(plain, "────") {
		t.Errorf("asterisk hrule should render as line: plain text = %q", plain)
	}
}

func TestRenderMarkdown_UnderscoreHRule(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("___")
	plain := ansi.Strip(got)
	if !strings.Contains(plain, "────") {
		t.Errorf("underscore hrule should render as line: plain text = %q", plain)
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
	plain := ansi.Strip(got)
	if strings.Contains(plain, "# ") {
		t.Error("heading marker should be stripped")
	}
	if !strings.Contains(plain, "Title") {
		t.Error("heading text should be present")
	}
	if strings.Contains(plain, "**") {
		t.Error("bold markers should be stripped")
	}
	if !strings.Contains(plain, "│ A quote") {
		t.Error("blockquote should have │ prefix")
	}
	if !strings.Contains(plain, "────") {
		t.Error("horizontal rule should be rendered as line")
	}
}

func TestRenderMarkdown_CodeBlockBoundary(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	input := "**bold before**\n```\ncode\n```\n**bold after**"
	got := renderMarkdown(input)
	plain := ansi.Strip(got)
	if !strings.Contains(plain, "bold before") {
		t.Errorf("bold before code block should be present: got %q", plain)
	}
	if !strings.Contains(plain, "bold after") {
		t.Errorf("bold after code block should be present: got %q", plain)
	}
	if strings.Contains(plain, "**") {
		t.Error("bold markers should be stripped")
	}
	if strings.Contains(plain, "```") {
		t.Error("fence markers should be hidden")
	}
}

func TestRenderMarkdown_SoftWrapCompatibility(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	input := "This is **bold** and `code` in a somewhat longer line of text"
	styled := renderMarkdown(input)
	// Soft wrap at 30 chars — should not panic or produce corrupted output.
	wrapped := softWrapLines(styled, 30)
	plain := ansi.Strip(wrapped)
	if !strings.Contains(plain, "bold") {
		t.Errorf("soft wrap lost bold text: plain = %q", plain)
	}
	if !strings.Contains(plain, "code") {
		t.Errorf("soft wrap lost code text: plain = %q", plain)
	}
	if strings.Contains(plain, "**") || strings.Contains(plain, "`") {
		t.Error("markers should be stripped after soft wrap")
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

func TestRenderMarkdown_WikiLinkPreserved(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	input := "see [[book/test]] for details"
	got := renderMarkdown(input)
	if !strings.Contains(got, "[[book/test]]") {
		t.Error("wiki-link syntax should pass through unchanged")
	}
}

func TestRenderMarkdown_HeadingWithInline(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("## **Bold** heading")
	plain := ansi.Strip(got)
	if strings.Contains(plain, "##") {
		t.Error("heading marker should be stripped")
	}
	if strings.Contains(plain, "**") {
		t.Error("bold markers should be stripped")
	}
	if !strings.Contains(plain, "Bold") || !strings.Contains(plain, "heading") {
		t.Errorf("heading content should be present: %q", plain)
	}
}

func TestRenderMarkdown_SnakeCaseNotItalicized(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("use my_function_name here")
	plain := ansi.Strip(got)
	if plain != "use my_function_name here" {
		t.Errorf("snake_case should not be italicized: got %q", plain)
	}
}

func TestRenderMarkdown_ConsecutiveItalicSpans(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)
	got := renderMarkdown("*a* and *b* and *c*")
	plain := ansi.Strip(got)
	if !strings.Contains(plain, "a") || !strings.Contains(plain, "b") || !strings.Contains(plain, "c") {
		t.Errorf("all italic spans should be present: %q", plain)
	}
	if strings.Contains(plain, "*") {
		t.Errorf("no stray * should remain: %q", plain)
	}
}
