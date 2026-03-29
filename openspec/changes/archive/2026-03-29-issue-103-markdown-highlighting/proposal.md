## Why

The TUI body panel renders markdown as plain text with no visual distinction for syntax elements. Only wiki-links are currently highlighted. Adding syntax highlighting for headings, bold, italic, inline code, code blocks, and links would significantly improve readability and make the TUI feel like a proper markdown viewer.

## What Changes

- Add a markdown rendering pipeline that applies lipgloss styles to markdown syntax elements in the body panel
- Support headings (`#`–`####`), bold (`**`), italic (`*`/`_`), inline code (`` ` ``), fenced code blocks (`` ``` ``), links (`[text](url)`), horizontal rules (`---`), and blockquotes (`>`)
- Add configurable theme colors for each markdown element in `.typemd/tui.yaml` under `theme`
- Integrate markdown styling with the existing wiki-link rendering pipeline

## Capabilities

### New Capabilities

- `markdown-rendering`: Markdown syntax highlighting in the TUI body panel with configurable theme colors

### Modified Capabilities

_(none — wiki-link rendering behavior is unchanged; markdown rendering is applied alongside it)_

## Impact

- **tui/**: New markdown rendering function in `detail.go` or a dedicated `markdown.go`; extended `themeConfig` struct in `theme.go` with new color fields; `loadTheme()` updated to load new colors
- **No core/ changes**: Rendering is purely a TUI presentation concern
- **No breaking changes**: Default colors are applied when no theme config exists; existing `tui.yaml` files continue to work
