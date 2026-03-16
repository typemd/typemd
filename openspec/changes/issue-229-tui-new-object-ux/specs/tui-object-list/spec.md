## MODIFIED Requirements

### Requirement: Group header displays type emoji

The TUI object list panel SHALL display the type's emoji prefix in group headers when the type schema defines an emoji field.

#### Scenario: Type with emoji defined
- **WHEN** a type schema has an emoji field (e.g., book with 📚)
- **THEN** the group header displays as `▼ 📚 book (N)` where N is the object count

#### Scenario: Type without emoji defined
- **WHEN** a type schema does not have an emoji field
- **THEN** the group header displays as `▼ book (N)` with no extra spacing or placeholder

## ADDED Requirements

### Requirement: Normal mode help bar shows both creation keybindings

When the sidebar is focused in normal mode, the help bar SHALL include both `n` (new) and `N` (quick create) keybinding hints.

#### Scenario: Sidebar focused help bar

- **WHEN** the sidebar is focused in normal mode with a type header or object selected
- **THEN** the help bar SHALL include hints for both `n: new` and `N: quick create`
