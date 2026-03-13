## MODIFIED Requirements

### Requirement: Group header displays type emoji

The TUI object list panel SHALL display the type's emoji prefix in group headers when the type schema defines an emoji field. The group header SHALL use the type's plural display name (from `PluralName()`) instead of the singular type name.

#### Scenario: Type with emoji and plural defined

- **WHEN** a type schema has emoji "📚" and plural "books"
- **THEN** the group header displays as `▼ 📚 books (N)` where N is the object count

#### Scenario: Type with emoji but no plural defined

- **WHEN** a type schema has emoji "📚" but no plural field
- **THEN** the group header displays as `▼ 📚 book (N)` using the singular name as fallback

#### Scenario: Type without emoji defined

- **WHEN** a type schema does not have an emoji field but has plural "notes"
- **THEN** the group header displays as `▼ notes (N)` with no extra spacing or placeholder
