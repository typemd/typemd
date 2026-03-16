## MODIFIED Requirements

### Requirement: Template selection when multiple templates exist

When a type has two or more object templates, the creation flow SHALL display a template cycling selector in the title panel alongside the name input. The user SHALL cycle through templates with `←`/`→` (or `↑`/`↓`) keys when the template field is focused. The template list SHALL include all available templates plus a "(none)" option.

#### Scenario: Type with multiple templates

- **WHEN** the user initiates object creation for type "book"
- **AND** type "book" has templates "review" and "summary"
- **THEN** the title panel SHALL display both a name input and a template selector (e.g., `📚 book · [name█] 📝 review`)
- **AND** `Tab` SHALL switch focus between name and template fields
- **AND** `←`/`→` on the template field SHALL cycle through "review", "summary", "(none)"

#### Scenario: Template selection cancelled

- **WHEN** the creation form is displayed in the title panel
- **AND** the user presses Escape
- **THEN** the creation flow SHALL be cancelled entirely

#### Scenario: Template none selected

- **WHEN** the user cycles to "(none)" in the template selector
- **THEN** the name input remains active with no template selected

### Requirement: Help bar reflects creation mode

The help bar SHALL display context-appropriate keybinding hints during object creation.

#### Scenario: Help bar during name input in Create & Edit mode

- **WHEN** the name field is focused in Create & Edit mode
- **THEN** the help bar SHALL display hints including "tab: template  enter: create & edit  esc: cancel"

#### Scenario: Help bar during template selector

- **WHEN** the template field is focused
- **THEN** the help bar SHALL display hints including "◀▶: switch  tab: name  enter: create  esc: cancel"

#### Scenario: Help bar during name input in Quick Create mode

- **WHEN** the name field is focused in Quick Create mode
- **THEN** the help bar SHALL display "tab: template  enter: create  esc: done"

### Requirement: Name template auto-skip in Create & Edit mode

In Create & Edit mode (`n`), when the type defines a name template, the name input SHALL be pre-filled with the evaluated template value instead of being skipped entirely. The user MAY edit the pre-filled name or press Enter to accept it.

#### Scenario: Name template pre-fills name input

- **WHEN** the user presses `n` on type "journal" which has name template `{{ date:YYYY-MM-DD }}`
- **THEN** the title panel SHALL display the name input pre-filled with the evaluated name (e.g., "2026-03-16")
- **AND** the user MAY edit the name or press Enter to create with the pre-filled value

#### Scenario: Name template with template selection

- **WHEN** the user presses `n` on type "journal" which has name template `{{ date:YYYY-MM-DD }}`
- **AND** type "journal" has multiple templates
- **THEN** the title panel SHALL display both name input (pre-filled) and template selector

## ADDED Requirements

### Requirement: Creation form renders in title panel

When object creation is active (`n` or `N`), the title panel SHALL transform into an inline creation form showing the type emoji, type name, a name text input, and (when applicable) a template cycling selector. The sidebar SHALL NOT render any creation-specific UI.

#### Scenario: Title panel shows creation form

- **WHEN** the user presses `n` on a type header for type "book" with emoji "📚"
- **THEN** the title panel SHALL display `📚 book · [name input] 📝 <template>`
- **AND** the sidebar SHALL show the normal object list without any creation UI appended

#### Scenario: Title panel without templates

- **WHEN** the user presses `n` on a type with no templates
- **THEN** the title panel SHALL display only the name input (no template selector)

#### Scenario: Title panel with single template

- **WHEN** the user presses `n` on a type with exactly one template "default"
- **THEN** the title panel SHALL display the name input and a static template label `📝 default` (not interactive)

#### Scenario: Title panel visible during creation even without prior selection

- **WHEN** no object was previously selected (title panel was hidden)
- **AND** the user presses `n` on a type header
- **THEN** the title panel SHALL appear with the creation form

### Requirement: Live template preview in body and properties panels

During object creation, switching the template selector SHALL update the body panel and properties panel in real time to show a preview of the selected template's content.

#### Scenario: Template preview shows body content

- **WHEN** the user cycles the template selector to "review"
- **AND** the "review" template has body content `## Review Notes`
- **THEN** the body panel SHALL display `## Review Notes` as a preview

#### Scenario: Template preview shows properties

- **WHEN** the user cycles the template selector to "review"
- **AND** the "review" template has frontmatter `status: draft`
- **THEN** the properties panel SHALL display `status: draft` as a preview

#### Scenario: None template shows empty preview

- **WHEN** the user cycles the template selector to "(none)"
- **THEN** the body panel SHALL display `(empty)` placeholder
- **AND** the properties panel SHALL display schema default values

#### Scenario: Preview updates on each template switch

- **WHEN** the user switches from "review" to "summary"
- **THEN** the body and properties panels SHALL update to show "summary" template content

### Requirement: Tab switches between name and template fields

During creation, `Tab` SHALL cycle focus between the name input field and the template selector field. The focused field SHALL be visually indicated.

#### Scenario: Tab from name to template

- **WHEN** the name field is focused
- **AND** the user presses Tab
- **THEN** the template selector SHALL become focused
- **AND** `←`/`→` SHALL cycle templates

#### Scenario: Tab from template to name

- **WHEN** the template selector is focused
- **AND** the user presses Tab
- **THEN** the name field SHALL become focused
- **AND** text input SHALL be active

#### Scenario: Tab with no templates

- **WHEN** the type has no templates
- **AND** the user presses Tab
- **THEN** nothing SHALL change (no template field to switch to)

### Requirement: Batch mode flash in title panel

In Quick Create mode, the success flash message SHALL appear in the title panel area, replacing or overlaying the name input briefly.

#### Scenario: Flash in title panel

- **WHEN** an object is created in Quick Create mode
- **THEN** the title panel SHALL briefly show the flash message (e.g., `✓ Created: my-book`)
- **AND** the name input SHALL be cleared and refocused after the flash
