## ADDED Requirements

### Requirement: Headings are rendered

The TUI body panel SHALL render markdown headings (`#` through `######`) by hiding the `#` markers and displaying the heading text with a distinct foreground color and bold weight. Lines that do not start with `#` followed by a space SHALL NOT be treated as headings.

#### Scenario: H1 heading is rendered

- **WHEN** the body contains a line `# Introduction`
- **THEN** `Introduction` is displayed in the heading color with bold (the `#` marker is hidden)

#### Scenario: H4 heading is rendered

- **WHEN** the body contains a line `#### Details`
- **THEN** `Details` is displayed in the heading color with bold (the `####` markers are hidden)

#### Scenario: Hash in non-heading context is not styled

- **WHEN** the body contains a line `Use #hashtag for tagging`
- **THEN** the line is rendered without heading styling

### Requirement: Bold text is rendered

The TUI body panel SHALL render bold markdown syntax (`**text**`) by hiding the `**` markers and displaying the content with bold weight.

#### Scenario: Bold word in a line

- **WHEN** the body contains `This is **important** text`
- **THEN** `important` is displayed with bold styling (the `**` markers are hidden)

#### Scenario: Multiple bold spans in one line

- **WHEN** the body contains `**first** and **second**`
- **THEN** both `first` and `second` are displayed with bold styling

### Requirement: Italic text is rendered

The TUI body panel SHALL render italic markdown syntax (`*text*` and `_text_`) by hiding the markers and displaying the content with italic styling. The italic pattern SHALL NOT match bold markers (`**`). Intra-word underscores (e.g. `snake_case`) SHALL NOT be treated as italic.

#### Scenario: Italic with asterisks

- **WHEN** the body contains `This is *emphasized* text`
- **THEN** `emphasized` is displayed with italic styling (the `*` markers are hidden)

#### Scenario: Italic with underscores

- **WHEN** the body contains `This is _emphasized_ text`
- **THEN** `emphasized` is displayed with italic styling (the `_` markers are hidden)

#### Scenario: Bold markers are not treated as italic

- **WHEN** the body contains `This is **bold** text`
- **THEN** `bold` is displayed with bold styling, not italic

#### Scenario: Snake_case identifiers are not italicized

- **WHEN** the body contains `use my_function_name here`
- **THEN** the line is rendered without italic styling

### Requirement: Inline code is rendered

The TUI body panel SHALL render inline code (`` `code` ``) by hiding the backtick markers and displaying the content with a distinct foreground color.

#### Scenario: Inline code in text

- **WHEN** the body contains ``Use `fmt.Println` for output``
- **THEN** `fmt.Println` is displayed in the inline code color (backticks are hidden)

#### Scenario: Multiple inline code spans

- **WHEN** the body contains ``Both `foo` and `bar` are valid``
- **THEN** both `foo` and `bar` are displayed in the inline code color

### Requirement: Fenced code blocks are rendered

The TUI body panel SHALL render fenced code blocks (lines between ` ``` ` delimiters) by hiding the fence lines and styling the content lines with a distinct foreground color. Inline markdown syntax (bold, italic, etc.) SHALL NOT be applied inside code blocks.

#### Scenario: Code block content is styled

- **WHEN** the body contains a fenced code block with content `x := 1`
- **THEN** the content line is rendered in the code block color and the fence lines are hidden

#### Scenario: Markdown syntax inside code block is not processed

- **WHEN** the body contains a fenced code block with content `**not bold**`
- **THEN** `**not bold**` is rendered in the code block color without bold styling

### Requirement: Links are rendered

The TUI body panel SHALL render markdown links (`[text](url)`) by hiding the URL and bracket syntax, displaying only the link text with a distinct foreground color.

#### Scenario: Link is rendered

- **WHEN** the body contains `See [documentation](https://example.com) for details`
- **THEN** `documentation` is displayed in the link color (the URL and brackets are hidden)

### Requirement: Blockquotes are rendered

The TUI body panel SHALL render blockquote lines (starting with `>`) by hiding the `>` marker, adding a `│` prefix, and applying a distinct foreground color.

#### Scenario: Blockquote line is rendered

- **WHEN** the body contains `> This is a quote`
- **THEN** `│ This is a quote` is displayed in the blockquote color

### Requirement: Horizontal rules are rendered

The TUI body panel SHALL render horizontal rules (`---`, `***`, or `___` on a line by themselves, with at least 3 characters) by replacing them with a styled `────────────────────` line.

#### Scenario: Dash horizontal rule is rendered

- **WHEN** the body contains a line `---`
- **THEN** a styled horizontal line (`────────────────────`) is displayed

#### Scenario: Asterisk horizontal rule is rendered

- **WHEN** the body contains a line `***`
- **THEN** a styled horizontal line is displayed

### Requirement: Theme colors are configurable

All markdown element colors SHALL be configurable via `.typemd/config.yaml` under the `tui.theme` section. When a color is not configured, the default SHALL be used. Existing theme fields (`focus_border`, `wiki_link`) SHALL continue to work unchanged.

#### Scenario: Custom heading color

- **WHEN** `.typemd/config.yaml` contains `tui.theme.heading: "196"`
- **THEN** headings are rendered using ANSI color 196

#### Scenario: Missing config uses defaults

- **WHEN** `.typemd/config.yaml` does not contain a `tui.theme.heading` entry
- **THEN** headings are rendered using the default heading color

#### Scenario: Existing wiki_link config still works

- **WHEN** `.typemd/config.yaml` contains `tui.theme.wiki_link: "42"`
- **THEN** wiki-links are rendered using ANSI color 42

### Requirement: Wiki-link styling is preserved

Markdown rendering SHALL NOT interfere with existing wiki-link styling. Wiki-links SHALL continue to be rendered with their configured color.

#### Scenario: Wiki-link in markdown body

- **WHEN** the body contains `See [[book/example-01abc]] for details`
- **THEN** the wiki-link is rendered with wiki-link styling, not markdown link styling

### Requirement: Markdown styling coexists with pinned properties

Markdown rendering SHALL only apply to the body section of the body panel. Pinned properties displayed above the separator SHALL NOT be affected by markdown rendering.

#### Scenario: Pinned property with markdown-like content

- **WHEN** a pinned property value contains `**bold**`
- **THEN** the property is rendered without markdown styling

### Requirement: Edit mode shows raw markdown

When the user enters edit mode (`e` key), the body panel SHALL display the raw markdown source with all syntax markers visible. Markdown rendering only applies in view mode.

#### Scenario: Edit mode shows markers

- **WHEN** the user presses `e` to enter edit mode on an object with `## Heading`
- **THEN** the textarea displays `## Heading` with the `##` markers visible
