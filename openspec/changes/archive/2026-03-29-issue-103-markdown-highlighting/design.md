## Context

The TUI body panel (`tui/detail.go:renderBody`) currently renders markdown as plain text. The only styling applied is wiki-link highlighting via `core.RenderWikiLinksStyled()`. The theme system (`tui/theme.go`) supports two configurable colors (`focus_border`, `wiki_link`) loaded from `.typemd/tui.yaml`.

Markdown content is rendered line-by-line in the body panel, with each line prefixed by a single space. After styling, soft wrapping may be applied in `app_render.go`.

## Goals / Non-Goals

**Goals:**

- Apply lipgloss-based syntax highlighting to markdown elements (headings, bold, italic, inline code, code blocks, links, blockquotes, horizontal rules) in the body panel
- Make all markdown element colors configurable via `.typemd/tui.yaml` under `theme`
- Maintain compatibility with existing wiki-link rendering
- Keep the rendering fast (line-by-line processing, no full AST)

**Non-Goals:**

- Full markdown-to-terminal rendering (e.g., `glamour`-style with reflowing, nested list indentation, or image placeholders)
- Syntax highlighting inside code blocks (language-aware highlighting)
- Markdown rendering in other panels (properties, sidebar)
- Table rendering

## Decisions

### 1. Line-by-line regex approach (not AST parser)

Use regex-based pattern matching applied line-by-line, similar to the existing wiki-link approach.

**Why over a full parser like goldmark/glamour:** The body panel already works line-by-line with simple string operations. A full AST parser would add a dependency, increase complexity, and risk altering whitespace/layout in unexpected ways. The regex approach handles the common markdown elements well enough and is trivially testable.

**Trade-off:** Cannot handle deeply nested or multi-paragraph constructs. Acceptable because the goal is visual highlighting, not semantic rendering.

### 2. New `tui/markdown.go` file

Create a dedicated `markdown.go` file for the rendering pipeline rather than expanding `detail.go`.

**Why:** Keeps markdown rendering logic isolated and testable. The `detail.go` file remains focused on panel composition.

### 3. Rendering order: markdown first, then wiki-links

Apply markdown styling before wiki-link styling. Wiki-links use `[[...]]` syntax which won't collide with markdown syntax.

**Why:** Wiki-link rendering replaces `[[target]]` with styled display text. If markdown styling ran after, it could try to style the already-styled wiki-link output. Running markdown first ensures clean input for wiki-link processing.

### 4. Fenced code block state tracking

Track whether we're inside a fenced code block (``` ``` ```) to avoid applying inline styles within code blocks.

**Why:** Code blocks should render their content literally. Applying heading or bold styles inside a code block would be incorrect.

### 5. Extend themeConfig with markdown colors

Add fields to `themeConfig` struct for each markdown element, with sensible defaults:

```go
type themeConfig struct {
    FocusBorder string `yaml:"focus_border"`
    WikiLink    string `yaml:"wiki_link"`
    Heading     string `yaml:"heading"`
    Bold        string `yaml:"bold"`
    Italic      string `yaml:"italic"`
    InlineCode  string `yaml:"inline_code"`
    CodeBlock   string `yaml:"code_block"`
    Link        string `yaml:"link"`
    Blockquote  string `yaml:"blockquote"`
    HRule       string `yaml:"hrule"`
}
```

Each has a default ANSI color and a corresponding lipgloss style variable in `theme.go`.

### 6. Default color palette

| Element | Default | Rationale |
|---------|---------|-----------|
| Heading | `"3"` (yellow) | Stand out as structural markers |
| Bold | `""` (terminal default, bold attribute) | Bold is about weight, not color |
| Italic | `""` (terminal default, italic attribute) | Italic is about style, not color |
| Inline code | `"245"` (gray) | Subtle distinction from body text |
| Code block | `"245"` (gray) | Same as inline code for consistency |
| Link | `"33"` (cyan) | Match wiki-link color by default |
| Blockquote | `"8"` (dim) | Visually recessed |
| Horizontal rule | `"8"` (dim) | Structural separator, not prominent |

## Risks / Trade-offs

- **[Regex limitations]** → Complex nested markdown (e.g., bold inside heading) may not render perfectly. Mitigation: Handle the most common cases; accept imperfect rendering for edge cases.
- **[ANSI escape interaction with soft wrap]** → Lipgloss-styled text contains ANSI escape codes that affect string width calculation. Mitigation: `softWrapLines()` already handles styled text (it works with wiki-links today); verify it continues to work.
- **[Performance with large files]** → Line-by-line regex on large markdown files could be slow. Mitigation: The viewport only renders visible lines; markdown styling is applied once when content changes, not on every render frame.
