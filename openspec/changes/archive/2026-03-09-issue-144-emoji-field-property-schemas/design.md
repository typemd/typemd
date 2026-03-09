## Context

The `Property` struct in `core/type_schema.go` currently has no `emoji` field. The `TypeSchema` struct already has an `emoji` field (type-level emoji). This change adds a similar concept at the property level.

`ValidateSchema()` currently checks for duplicate property names but has no emoji-related validation.

Global/predefined properties are not yet implemented (spec exists at `openspec/specs/predefined-property/`), so global-scope emoji uniqueness validation will be deferred until that feature lands.

## Goals / Non-Goals

**Goals:**
- Add optional `emoji` field to `Property` struct
- Validate emoji uniqueness within a type's properties
- Keep the change minimal and non-breaking

**Non-Goals:**
- UI rendering of property emojis (TUI/web changes deferred to consumers)
- Global property emoji uniqueness (predefined properties not yet implemented)
- Emoji format validation (any non-empty string accepted; no Unicode emoji detection)
- Fallback behavior when emoji is not set (no auto-generation from property name)

## Decisions

### 1. Field type: `string` with `omitempty`

Add `Emoji string \`yaml:"emoji,omitempty"\`` to the `Property` struct, matching the pattern already used by `TypeSchema.Emoji`.

**Alternative**: A dedicated emoji type with Unicode validation. Rejected — over-engineering for an optional display hint. Users may want to use text characters or custom symbols.

### 2. Uniqueness scope: per-type only

Validate that no two properties within the same `TypeSchema.Properties` slice share the same non-empty emoji. This catches typos and copy-paste errors.

**Alternative**: Cross-type global uniqueness. Rejected — different types naturally reuse the same emoji (e.g., both `book` and `person` might use 📝 for a `notes` property).

### 3. Validation in `ValidateSchema()`

Add emoji duplicate checking to the existing `ValidateSchema()` function, alongside the existing duplicate property name check. Empty emojis are skipped (not all properties need one).

### 4. Defer global property scope

The issue mentions "global properties form their own separate uniqueness scope." Since predefined/global properties are not yet implemented, this validation will be added when that feature lands.

## Risks / Trade-offs

- [Low] Users might set non-emoji strings (e.g., "abc") as the emoji value → Accepted. No format enforcement keeps the field flexible.
- [Low] No cross-type duplicate detection → Acceptable. Types are independent schemas.
