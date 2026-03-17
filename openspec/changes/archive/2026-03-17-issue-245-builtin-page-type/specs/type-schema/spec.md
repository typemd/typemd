## MODIFIED Requirements

### Requirement: Custom type emoji overrides built-in default

When a custom type schema defines its own emoji, it SHALL override the built-in default emoji for that type. This applies to any built-in type (`tag`, `page`) when a custom schema is defined.

#### Scenario: Custom tag type with different emoji

- **WHEN** a custom `tag.yaml` defines `emoji: 🔖`
- **THEN** the loaded tag type SHALL have emoji "🔖" instead of the built-in "🏷️"

## REMOVED Requirements

### Requirement: Only tag is a built-in type

**Reason**: Replaced by the `builtin-type` spec which covers all built-in types (tag and page).
**Migration**: See `builtin-type/spec.md` for the authoritative built-in type definitions.
