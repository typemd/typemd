## MODIFIED Requirements

### Requirement: Group header displays type emoji

The TUI object list panel SHALL display the type's emoji prefix in group headers when the type schema defines an emoji field. Object list loading SHALL use `[]FilterRule` parameters when calling `Vault.QueryObjects()`.

#### Scenario: Type with emoji defined
- **WHEN** a type schema has an emoji field (e.g., book with 📚)
- **THEN** the group header displays as `▼ 📚 book (N)` where N is the object count

#### Scenario: Type without emoji defined
- **WHEN** a type schema does not have an emoji field
- **THEN** the group header displays as `▼ book (N)` with no extra spacing or placeholder

#### Scenario: Object list loading uses structured filter
- **WHEN** the TUI loads the object list (via `app.go` or `view_mode.go`)
- **THEN** it SHALL call `Vault.QueryObjects([]FilterRule{...})` instead of passing a filter string
