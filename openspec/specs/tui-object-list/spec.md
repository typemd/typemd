# tui-object-list Specification

## Purpose
TBD - created by archiving change issue-163-tui-type-emoji-headers. Update Purpose after archive.
## Requirements
### Requirement: Group header displays type emoji

The TUI object list panel SHALL display the type's emoji prefix in group headers when the type schema defines an emoji field.

#### Scenario: Type with emoji defined
- **WHEN** a type schema has an emoji field (e.g., book with 📚)
- **THEN** the group header displays as `▼ 📚 book (N)` where N is the object count

#### Scenario: Type without emoji defined
- **WHEN** a type schema does not have an emoji field
- **THEN** the group header displays as `▼ book (N)` with no extra spacing or placeholder

