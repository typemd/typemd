## Context

Currently `defaultTypes` contains only `tag`. Users who skip starter types during `tmd init` have no type for general content. The `page` type fills this gap as a minimal, always-available content container.

The existing infrastructure already supports multiple built-in types — `defaultTypes` is a map, `DeleteType()` checks it for protection, `GetSchema()` falls back to it, and `ListSchemas()` merges it with custom types. Adding `page` requires no architectural changes.

## Goals / Non-Goals

**Goals:**
- Add `page` as a second built-in type with minimal properties (system properties only)
- Protect `page` from deletion using existing `defaultTypes` mechanism
- Make `page` available as `cli.default_type` fallback during `tmd init`

**Non-Goals:**
- Merge/override mechanism for built-in type properties (separate issue)
- Adding custom properties to the built-in `page` definition
- Changing `tag` behavior or the `defaultTypes` architecture

## Decisions

### Decision 1: Add `PageTypeName` constant in `system_property.go`

Place it alongside `TagTypeName` since both are built-in type name constants.

**Rationale**: Follows the existing pattern. `TagTypeName` is already in `system_property.go` because it's referenced by the `tags` system property.

### Decision 2: Page definition in `defaultTypes` — no custom properties

```go
PageTypeName: {
    Name:   PageTypeName,
    Plural: "pages",
    Emoji:  "📄",
}
```

No `Unique`, no `Properties`. This keeps page as a pure content container. `Unique` defaults to `false`.

**Rationale**: The issue explicitly states "minimal properties (no opinionated schema beyond system properties)". Users can override via custom `.typemd/types/page.yaml`.

### Decision 3: Update `resolveDefaultType()` with page fallback

Current logic: idea → note → (empty). New logic: idea → note → page.

This ensures vaults always have a usable default type, even when no starters are selected.

**Rationale**: The issue states "Users who skip starter types have no type available for general content." This is the minimal fix.

## Risks / Trade-offs

- **[Risk] Existing tests assume only `tag` in defaultTypes** → Update tests that count or enumerate built-in types. The existing `TestVault_ListTypes_CustomOverridesDefault` may need adjustment.
- **[Risk] Spec says "Only tag is a built-in type"** → Update the spec requirement to include both `tag` and `page`.
- **[Trade-off] `page` has no properties** → This is intentional for minimalism. Users wanting structured pages can override with custom YAML.
