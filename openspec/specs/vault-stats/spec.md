## Requirements

### Requirement: Vault-wide statistics summary
The system SHALL provide a vault-wide statistics summary showing each type's emoji, display name (plural if available), object count, and most recent update time.

#### Scenario: Vault with multiple types
- **WHEN** user runs `tmd stats` on a vault containing books (3), ideas (2), and tags (1)
- **THEN** the output SHALL list each type with emoji, plural name, count, and last updated date, followed by a total count

#### Scenario: Empty vault
- **WHEN** user runs `tmd stats` on a vault with no objects
- **THEN** the output SHALL display a total of 0 with no type rows

#### Scenario: Built-in types without schema files
- **WHEN** user runs `tmd stats` and the vault contains `tag` objects but no `.typemd/types/tag.yaml`
- **THEN** the output SHALL use the built-in type definition (emoji 🏷️, plural "tags")

### Requirement: Single-type property statistics
The system SHALL provide per-property aggregate statistics for a specified type, showing the property name, type, fill rate, and type-appropriate aggregations.

#### Scenario: Number property aggregation
- **WHEN** user runs `tmd stats --type book` and the `book` type has a `rating` number property with values [3, 4, 5]
- **THEN** the output SHALL show sum (12), average (4.0), min (3), max (5), and filled count (3)

#### Scenario: Select property distribution
- **WHEN** user runs `tmd stats --type book` and the `book` type has a `status` select property
- **THEN** the output SHALL show the count for each option value (e.g., reading: 2, finished: 3)

#### Scenario: Multi-select property distribution
- **WHEN** user runs `tmd stats --type book` and the `book` type has a `genres` multi_select property
- **THEN** the output SHALL show the count for each selected option across all objects

#### Scenario: Checkbox property ratio
- **WHEN** user runs `tmd stats --type book` and the `book` type has a `read` checkbox property
- **THEN** the output SHALL show the count of true and false values

#### Scenario: Date property range
- **WHEN** user runs `tmd stats --type book` and the `book` type has a `published` date property with values
- **THEN** the output SHALL show the earliest and latest dates

#### Scenario: Relation property count
- **WHEN** user runs `tmd stats --type book` and the `book` type has an `author` relation property
- **THEN** the output SHALL show the total number of relation links

#### Scenario: Property with no values
- **WHEN** user runs `tmd stats --type book` and no objects have a value for the `rating` property
- **THEN** the output SHALL show filled 0 out of total and skip type-specific stats

#### Scenario: Non-existent type
- **WHEN** user runs `tmd stats --type nonexistent`
- **THEN** the system SHALL return an error indicating the type does not exist

### Requirement: JSON output format
The system SHALL support a `--json` flag that outputs statistics as structured JSON instead of human-readable text.

#### Scenario: Vault-wide JSON output
- **WHEN** user runs `tmd stats --json`
- **THEN** the output SHALL be valid JSON containing an array of type summaries and a total count

#### Scenario: Single-type JSON output
- **WHEN** user runs `tmd stats --type book --json`
- **THEN** the output SHALL be valid JSON containing the type name, object count, and per-property statistics with type-specific aggregation data

### Requirement: QueryService aggregation methods
The system SHALL expose vault statistics through `QueryService` methods, making aggregation logic reusable across CLI, MCP, and Web consumers. Internally, these methods SHALL use `[]FilterRule` (via `TypeFilter()`) instead of legacy filter strings when querying the index.

#### Scenario: VaultStats method
- **WHEN** a consumer calls `QueryService.VaultStats()`
- **THEN** it SHALL return a `VaultStats` struct with per-type summaries and total count

#### Scenario: TypeStats method
- **WHEN** a consumer calls `QueryService.TypeStats("book")`
- **THEN** it SHALL return a `TypeStats` struct with property-level aggregations for all aggregable properties of the `book` type

#### Scenario: Internal query uses structured filter
- **WHEN** `VaultStats()` or `TypeStats()` queries the index for objects of a specific type
- **THEN** it SHALL pass `TypeFilter(typeName)` (i.e. `[]FilterRule{{Property: "type", Operator: "is", Value: typeName}}`) instead of a filter string

### Requirement: Stats mode entry via global shortcut
The TUI SHALL enter stats mode when the user presses `ctrl+s` from any panel (sidebar, body, or properties). Stats mode SHALL replace the three-panel layout with a fullscreen view.

#### Scenario: Enter stats mode from sidebar
- **WHEN** the user presses `ctrl+s` while focus is on the sidebar
- **THEN** the TUI SHALL switch to stats mode displaying the Vault Overview screen

#### Scenario: Enter stats mode from body panel
- **WHEN** the user presses `ctrl+s` while focus is on the body panel
- **THEN** the TUI SHALL switch to stats mode displaying the Vault Overview screen

#### Scenario: Enter stats mode from properties panel
- **WHEN** the user presses `ctrl+s` while focus is on the properties panel
- **THEN** the TUI SHALL switch to stats mode displaying the Vault Overview screen

### Requirement: Stats mode exit
The TUI SHALL exit stats mode and return to the previous panel state when the user presses `Esc` at the Vault Overview level or presses `q`.

#### Scenario: Exit stats mode via Esc from vault overview
- **WHEN** the user presses `Esc` while on the Vault Overview screen (not in type detail)
- **THEN** the TUI SHALL return to the three-panel layout with the previous panel state restored

#### Scenario: Exit stats mode via q
- **WHEN** the user presses `q` while in stats mode
- **THEN** the TUI SHALL quit entirely (consistent with global quit behavior)

### Requirement: Vault Overview screen
The Vault Overview SHALL display a fullscreen summary of the vault including a header with total object count and total type count, followed by a navigable list of types sorted by object count descending.

#### Scenario: Vault overview header
- **WHEN** the user enters stats mode on a vault with 43 objects across 4 types
- **THEN** the header SHALL display the total object count (43) and total type count (4)

#### Scenario: Type list content
- **WHEN** the vault contains types book (12 objects, last updated 2026-03-15), idea (8), page (23)
- **THEN** each row SHALL display the type emoji, display name (plural if available), object count, and relative last-updated time

#### Scenario: Type list sorting
- **WHEN** the vault contains types with varying object counts
- **THEN** the type list SHALL be sorted by object count descending (most objects first)

#### Scenario: Empty vault
- **WHEN** the vault contains no objects
- **THEN** the Vault Overview SHALL display a message indicating no objects exist

### Requirement: Vault Overview navigation
The Vault Overview type list SHALL support j/k (or arrow keys) for cursor movement and Enter to drill into type detail.

#### Scenario: Navigate with j/k
- **WHEN** the user presses `j` on the Vault Overview
- **THEN** the cursor SHALL move to the next type in the list

#### Scenario: Navigate with arrow keys
- **WHEN** the user presses the down arrow on the Vault Overview
- **THEN** the cursor SHALL move to the next type in the list

#### Scenario: Enter type detail
- **WHEN** the user presses `Enter` on a type in the list
- **THEN** the TUI SHALL display the Type Detail screen for that type

#### Scenario: Cursor wrapping
- **WHEN** the cursor is on the last type and the user presses `j`
- **THEN** the cursor SHALL remain on the last type (no wrapping)

### Requirement: Type Detail screen displays property statistics
The Type Detail screen SHALL display the type name with emoji and object count as a header, followed by per-property statistics for all properties defined in the type schema.

#### Scenario: Type detail header
- **WHEN** the user enters Type Detail for type "book" with emoji 📚 and 12 objects
- **THEN** the header SHALL display "📚 book (12 objects)"

#### Scenario: Number property stats
- **WHEN** the type has a "pages" number property with values across objects
- **THEN** the detail SHALL display min, max, avg, and fill rate for that property

#### Scenario: Select property distribution
- **WHEN** the type has a "genre" select property with values fiction (8) and non-fiction (4)
- **THEN** the detail SHALL display each option with its count and a proportional bar chart

#### Scenario: Multi-select property distribution
- **WHEN** the type has a "tags" multi_select property
- **THEN** the detail SHALL display each selected option with its frequency count and bar chart

#### Scenario: Checkbox property ratio
- **WHEN** the type has a "read" checkbox property with 8 true and 4 false
- **THEN** the detail SHALL display the true/false counts and fill rate

#### Scenario: Date property range
- **WHEN** the type has a "published" date property
- **THEN** the detail SHALL display the earliest and latest dates and fill rate

#### Scenario: Relation property count
- **WHEN** the type has an "author" relation property
- **THEN** the detail SHALL display the total link count and fill rate

#### Scenario: String/URL property fill rate
- **WHEN** the type has a "url" string property
- **THEN** the detail SHALL display only the fill rate (filled / total)

#### Scenario: Property with no values
- **WHEN** no objects have a value for a given property
- **THEN** the detail SHALL display filled 0/N with no type-specific stats

### Requirement: Type Detail navigation
The Type Detail screen SHALL support Esc to return to the Vault Overview and scrolling for long property lists.

#### Scenario: Return to vault overview
- **WHEN** the user presses `Esc` on the Type Detail screen
- **THEN** the TUI SHALL return to the Vault Overview with the cursor on the same type

#### Scenario: Scroll property list
- **WHEN** the property list is longer than the available screen height
- **THEN** the user SHALL be able to scroll with j/k or arrow keys

### Requirement: Type Detail layout configuration
The Type Detail screen SHALL support two layouts controlled by `tui.stats_type_layout` config: `fullscreen` (replaces vault overview) and `popup` (centered overlay on vault overview). The default SHALL be `fullscreen`.

#### Scenario: Fullscreen layout
- **WHEN** `tui.stats_type_layout` is `fullscreen` or not set
- **AND** the user presses Enter on a type
- **THEN** the Type Detail SHALL render as a fullscreen view replacing the Vault Overview

#### Scenario: Popup layout
- **WHEN** `tui.stats_type_layout` is `popup`
- **AND** the user presses Enter on a type
- **THEN** the Type Detail SHALL render as a centered popup overlay on the Vault Overview

#### Scenario: Toggle layout at runtime
- **WHEN** the user presses `t` while on the Type Detail screen
- **THEN** the layout SHALL toggle between fullscreen and popup immediately

### Requirement: Bar chart rendering
Select and multi-select property distributions SHALL be rendered as horizontal bar charts using Unicode block characters, with bar width proportional to the available terminal width.

#### Scenario: Proportional bar scaling
- **WHEN** a select property has options fiction (8) and non-fiction (4)
- **THEN** the fiction bar SHALL be twice as long as the non-fiction bar
- **AND** the longest bar SHALL scale to fill the available width

#### Scenario: Single option
- **WHEN** a select property has only one option with count 5
- **THEN** that option SHALL display a full-width bar

#### Scenario: Narrow terminal
- **WHEN** the terminal width is less than 60 columns
- **THEN** bar charts SHALL scale down proportionally without truncating labels

### Requirement: Stats refresh
The user SHALL be able to refresh stats data by pressing `r` while in stats mode.

#### Scenario: Refresh from vault overview
- **WHEN** the user presses `r` on the Vault Overview
- **THEN** vault stats SHALL be recomputed from the current index data

#### Scenario: Refresh from type detail
- **WHEN** the user presses `r` on the Type Detail screen
- **THEN** type stats for the current type SHALL be recomputed from the current index data
