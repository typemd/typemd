## Context

The `Property` struct in `core/type_schema.go` currently has fields: `Name`, `Type`, `Emoji`, `Options`, `Target`, `Default`, `Multiple`, `Bidirectional`, `Inverse`. The `ValidateSchema()` function checks property names, types, options, and emoji uniqueness.

The TUI detail view has three panels: title (top), body (center), and properties (right). `renderBody()` renders markdown content; `renderProperties()` renders formatted `DisplayProperty` items. `DisplayProperty` carries `Key`, `Value`, `Type`, `IsRelation`, `IsReverse`, `IsBacklink` fields.

Property emoji was added in #144 but is not yet rendered in the TUI.

## Goals / Non-Goals

**Goals:**
- Add optional `pin` integer field to `Property` struct
- Validate pin uniqueness within a type's properties (duplicate numbers rejected)
- Render pinned properties at top of body panel with key-value format
- Exclude pinned properties from Properties panel
- Use property emoji in pinned display when available

**Non-Goals:**
- Pin support for predefined/global properties (not yet implemented)
- Interactive pin/unpin from TUI (schema-only configuration)
- Pin support in web UI (future, separate change)
- Negative or zero pin values (positive integers only)

## Decisions

### 1. Field type: `int` with `omitempty`

Add `Pin int \`yaml:"pin,omitempty"\`` to the `Property` struct. Zero value means "not pinned". Positive integers define display order (lower = higher priority).

**Alternative**: `bool` with schema-order sorting. Rejected — user chose explicit ordering for flexibility.

### 2. Uniqueness scope: per-type only

Validate that no two properties within the same `TypeSchema.Properties` slice share the same non-zero pin value. Mirrors the emoji uniqueness pattern.

### 3. Validation in `ValidateSchema()`

Add pin duplicate checking alongside existing duplicate name and emoji checks. Pin value must be positive (> 0) when set.

### 4. Pin ordering in display

Pinned properties sorted by pin value ascending (pin: 1 first, pin: 2 second, etc.).

### 5. Visual format: key-value lines with separator

Pinned properties rendered as `emoji key: value` lines at the top of the body panel, followed by a horizontal separator (`────`) before the markdown body content. Uses the existing `DisplayProperty.Format()` pattern for value formatting.

### 6. Pinned properties excluded from Properties panel

`renderProperties()` filters out any `DisplayProperty` whose corresponding schema property has a non-zero pin value. This avoids information duplication.

### 7. DisplayProperty gains Pin field

Add `Pin int` to `DisplayProperty` struct. `BuildDisplayProperties()` populates it from the schema. This lets both `renderBody()` and `renderProperties()` make decisions based on pin status without re-loading the schema.

## Risks / Trade-offs

- [Low] Pin values are schema-level only — users cannot pin/unpin per-object. Acceptable for v1; interactive pinning is a separate feature.
- [Low] Property emoji not yet rendered in TUI — this change will be the first to display property emoji, specifically for pinned properties. Non-pinned property emoji display remains deferred.
