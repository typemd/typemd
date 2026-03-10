## Why

In the TUI detail view, all properties are displayed equally in the Properties panel. Users have no way to emphasize important metadata — they must scan the full property list to find key information like status or rating. A "pinned properties" feature would let type schema authors mark specific properties for prominent display at the top of the body panel, making important metadata immediately visible.

## What Changes

- Add an optional `pin` field (positive integer) to the `Property` struct in type schemas
- Validate pin uniqueness within a type's properties scope (duplicate pin numbers rejected)
- Render pinned properties as key-value lines at the top of the TUI body panel, separated from body content by a horizontal rule
- Remove pinned properties from the Properties panel to avoid duplication
- Display property emoji (if defined) alongside pinned property values

## Capabilities

### New Capabilities
- `pinned-properties`: Optional integer `pin` field on property definitions with per-type uniqueness validation and prominent TUI body panel display

### Modified Capabilities
- `type-schema`: Property struct gains a `pin` field; `ValidateSchema` adds duplicate pin checking
- `tui-object-list`: Detail view body panel renders pinned properties at top; Properties panel excludes pinned properties

## Impact

- `core/type_schema.go` — `Property` struct, `ValidateSchema()` function
- `core/display.go` — `DisplayProperty` struct, `BuildDisplayProperties()` function
- `tui/detail.go` — `renderBody()`, `renderProperties()` functions
- `.typemd/types/*.yaml` — type schema files gain optional `pin` field
- `examples/` — example vault schemas updated with sample pinned properties
- No breaking changes — `pin` is optional with zero-value default
