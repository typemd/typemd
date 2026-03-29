## ADDED Requirements

### Requirement: Headings are styled

The TUI body panel SHALL render markdown headings (`#` through `####`) with a distinct foreground color. The heading marker (`#` characters) and heading text SHALL both be styled. Lines that do not start with `#` followed by a space SHALL NOT be treated as headings.

#### Scenario: H1 heading is styled

- **WHEN** the body contains a line `# Introduction`
- **THEN** the entire line is rendered in the heading color

#### Scenario: H4 heading is styled

- **WHEN** the body contains a line `#### Details`
- **THEN** the entire line is rendered in the heading color

#### Scenario: Hash in non-heading context is not styled

- **WHEN** the body contains a line `Use #hashtag for tagging`
- **THEN** the line is rendered without heading styling

### Requirement: Bold text is styled

The TUI body panel SHALL render bold markdown syntax (`**text**`) with bold weight. The surrounding `**` markers SHALL be included in the styled output.

#### Scenario: Bold word in a line

- **WHEN** the body contains `This is **important** text`
- **THEN** `**important**` is rendered with bold styling

#### Scenario: Multiple bold spans in one line

- **WHEN** the body contains `**first** and **second**`
- **THEN** both `**first**` and `**second**` are rendered with bold styling

### Requirement: Italic text is styled

The TUI body panel SHALL render italic markdown syntax (`*text*` and `_text_`) with italic styling. The surrounding markers SHALL be included in the styled output. The italic pattern SHALL NOT match bold markers (`**`).

#### Scenario: Italic with asterisks

- **WHEN** the body contains `This is *emphasized* text`
- **THEN** `*emphasized*` is rendered with italic styling

#### Scenario: Italic with underscores

- **WHEN** the body contains `This is _emphasized_ text`
- **THEN** `_emphasized_` is rendered with italic styling

#### Scenario: Bold markers are not treated as italic

- **WHEN** the body contains `This is **bold** text`
- **THEN** `**bold**` is rendered with bold styling, not italic

### Requirement: Inline code is styled

The TUI body panel SHALL render inline code (`` `code` ``) with a distinct foreground color. The surrounding backtick markers SHALL be included in the styled output.

#### Scenario: Inline code in text

- **WHEN** the body contains ``Use `fmt.Println` for output``
- **THEN** `` `fmt.Println` `` is rendered in the inline code color

#### Scenario: Multiple inline code spans

- **WHEN** the body contains ``Both `foo` and `bar` are valid``
- **THEN** both `` `foo` `` and `` `bar` `` are rendered in the inline code color

### Requirement: Fenced code blocks are styled

The TUI body panel SHALL render fenced code blocks (lines between ` ``` ` delimiters) with a distinct foreground color. Both the delimiter lines and the content lines SHALL be styled. Inline markdown syntax (bold, italic, etc.) SHALL NOT be applied inside code blocks.

#### Scenario: Code block content is styled

- **WHEN** the body contains a fenced code block with content `x := 1`
- **THEN** the content line is rendered in the code block color

#### Scenario: Markdown syntax inside code block is not processed

- **WHEN** the body contains a fenced code block with content `**not bold**`
- **THEN** `**not bold**` is rendered in the code block color without bold styling

### Requirement: Links are styled

The TUI body panel SHALL render markdown links (`[text](url)`) with a distinct foreground color applied to the entire link syntax.

#### Scenario: Link is styled

- **WHEN** the body contains `See [documentation](https://example.com) for details`
- **THEN** `[documentation](https://example.com)` is rendered in the link color

### Requirement: Blockquotes are styled

The TUI body panel SHALL render blockquote lines (starting with `>`) with a distinct foreground color applied to the entire line.

#### Scenario: Blockquote line is styled

- **WHEN** the body contains `> This is a quote`
- **THEN** the entire line is rendered in the blockquote color

### Requirement: Horizontal rules are styled

The TUI body panel SHALL render horizontal rules (`---`, `***`, or `___` on a line by themselves, with at least 3 characters) with a distinct foreground color.

#### Scenario: Dash horizontal rule is styled

- **WHEN** the body contains a line `---`
- **THEN** the line is rendered in the horizontal rule color

#### Scenario: Asterisk horizontal rule is styled

- **WHEN** the body contains a line `***`
- **THEN** the line is rendered in the horizontal rule color

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

Markdown styling SHALL only apply to the body section of the body panel. Pinned properties displayed above the separator SHALL NOT be affected by markdown rendering.

#### Scenario: Pinned property with markdown-like content

- **WHEN** a pinned property value contains `**bold**`
- **THEN** the property is rendered without markdown styling
