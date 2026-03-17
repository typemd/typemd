## Why

Type schemas and properties currently lack metadata for visual theming and documentation. Types only have `name`, `plural`, `emoji`, and `unique` — there is no way to assign a color for UI theming or describe what a type represents. Properties similarly have no way to document their purpose. This limits the richness of display in TUI and future Web UI, and makes it harder for users to understand the intent behind types and properties in shared vaults.

## What Changes

- **TypeSchema** gains two new optional fields:
  - `color` (string) — accepts 10 preset palette names (`red`, `blue`, `green`, `yellow`, `purple`, `orange`, `pink`, `cyan`, `gray`, `brown`) or custom hex codes (`#RRGGBB` or `#RGB`). Used by TUI/Web UI for visual theming.
  - `description` (string) — free-text description of the type's purpose. No length limit.
- **Property** gains one new optional field:
  - `description` (string) — free-text description of the property's purpose. No length limit.
- **Shared properties `use` entries** allow `description` as an additional override field (alongside existing `pin` and `emoji`).
- **TUI type editor** updated to display and edit `color` and `description` metadata fields.
- **YAML serialization** updated to include new fields in output.
- **Validation** updated to enforce color format (preset name or valid hex).

## Capabilities

### New Capabilities

- `type-color`: Color metadata for type schemas — preset palette names and custom hex validation, YAML serialization, and TUI editing.
- `type-description`: Description metadata for type schemas and properties — free-text field on TypeSchema and Property, `use` entry override support, YAML serialization, and TUI editing.

### Modified Capabilities

- `type-schema`: Use entries expand allowed overrides to include `description` (currently only `pin` and `emoji`).

## Impact

- `core/type_schema.go` — TypeSchema struct, Property struct, marshalSchema, ValidateSchema(), validateUseOverrides(), MarshalTypeSchema()
- `core/type_schema_test.go` — New validation unit tests for color format and description fields
- `core/features/` — New BDD scenarios for color and description behaviors
- `tui/type_editor.go` — metaFieldCount 4→6, new meta field rows for color and description, inline editing support
- No breaking changes — all new fields are optional with zero-value defaults
