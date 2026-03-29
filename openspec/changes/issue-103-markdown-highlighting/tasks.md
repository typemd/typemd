## 1. Theme Configuration

- [x] 1.1 Extend `themeConfig` struct in `tui/theme.go` with markdown color fields (heading, bold, italic, inline_code, code_block, link, blockquote, hrule)
- [x] 1.2 Add default color constants and lipgloss style variables for each markdown element
- [x] 1.3 Update `loadTheme()` to read and apply markdown color overrides from `tui.yaml`
- [x] 1.4 Update `resetThemeDefaults()` to reset markdown styles
- [x] 1.5 Add unit tests for theme loading with markdown color overrides

## 2. Markdown Rendering

- [x] 2.1 Write BDD scenarios for markdown rendering in `tui/features/` (headings, bold, italic, inline code, code blocks, links, blockquotes, horizontal rules)
- [x] 2.2 Implement BDD step definitions for markdown rendering scenarios
- [x] 2.3 Create `tui/markdown.go` with `renderMarkdown(body string) string` function: line-by-line regex styling with code block state tracking
- [x] 2.4 Add unit tests for markdown rendering edge cases (nested syntax, empty lines, code block boundaries, bold vs italic disambiguation)

## 3. Integration

- [x] 3.1 Integrate `renderMarkdown()` into `renderBody()` in `tui/detail.go` — apply markdown styling before wiki-link styling
- [x] 3.2 Add unit tests verifying markdown styling does not apply to pinned properties
- [x] 3.3 Verify soft wrapping compatibility with styled markdown output
