## Context

Currently, `tmd object create <type> <name>` requires a name argument that becomes the object's slug and `name` property. For types with predictable naming patterns (journals, meeting notes), users must manually type a name every time. Issue #187 established `GetName()` as the centralized name access method, deliberately designed to accommodate future template logic.

The `name` system property is reserved — `ValidateSchema()` rejects any property named `name` in the `properties` array. This must be relaxed to allow `name` entries with only the `template` field.

## Goals / Non-Goals

**Goals:**

- Allow type schemas to define a name template via `- name: name` with `template` field in the `properties` array
- Evaluate templates at create time, writing the result as a static `name` property value
- Support `{{ date:FORMAT }}` placeholder with user-friendly format syntax (YYYY-MM-DD)
- Make the CLI `name` argument optional when a template is available

**Non-Goals:**

- Dynamic/display-time template evaluation
- Property reference placeholders (e.g., `{{ prop:author }}`) — future work
- Sequence number placeholders — future work
- Body templates (#173 — separate issue)

## Decisions

### D1: Template lives in the `properties` array as a `name` entry

**Decision:** Allow `- name: name` in the `properties` array with only the `template` field permitted. No `type`, `options`, or other fields allowed.

```yaml
properties:
  - name: name
    template: "日記 {{ date:YYYY-MM-DD }}"
  - name: content
    type: string
```

**Rationale:** The `properties` array is the natural place for property configuration. Future system property customization (ordering, display hints) can follow the same pattern. This avoids adding ad-hoc top-level fields to the schema.

**Alternative considered:** Top-level `name_template` field — simpler but establishes a precedent of scattering property config across the schema. Rejected for consistency.

### D2: Special validation path for `name` in properties

**Decision:** When `ValidateSchema()` encounters `name: name`, instead of rejecting it as a reserved system property, it enters a special validation path that only allows the `template` field. All other fields (`type`, `options`, `pin`, `emoji`, etc.) are rejected.

**Rationale:** This is the minimal change to the validation logic. The `name` property remains a system property — the schema can only configure it, not redefine it.

### D3: Template evaluation at create time only

**Decision:** `NewObject()` evaluates the template when the `name` argument is empty, writing the result as a static string to the `name` property. The template string is not stored in the object.

**Rationale:** Create-time evaluation is simple, predictable, and allows users to edit the name afterward. The `{{ date }}` placeholder semantically refers to the creation date, which is fixed.

### D4: User-friendly date format syntax

**Decision:** Templates use `{{ date:YYYY-MM-DD }}` syntax. The template engine converts common tokens to Go reference time format internally:

| Token | Go equivalent | Example |
|-------|--------------|---------|
| YYYY  | 2006         | 2026    |
| MM    | 01           | 03      |
| DD    | 02           | 14      |
| HH    | 15           | 09      |
| mm    | 04           | 30      |
| ss    | 05           | 00      |

**Rationale:** Go's reference time format (`2006-01-02`) is unintuitive for non-Go developers. `YYYY-MM-DD` is universally understood.

### D5: Template stored on TypeSchema, not Property

**Decision:** Add a `NameTemplate` field to the `TypeSchema` struct. During `LoadType()`, when a `name` property entry with `template` is found in the YAML, extract the template value into `TypeSchema.NameTemplate`. The `name` entry is not added to the resolved `Properties` slice.

**Rationale:** The template is a type-level configuration, not an object-level property. Storing it on `TypeSchema` keeps the `Properties` slice clean for schema-defined properties only. `NewObject()` can access it directly via `schema.NameTemplate`.

### D6: CLI argument handling

**Decision:** Change `cobra.ExactArgs(2)` to `cobra.RangeArgs(1, 2)`. When only 1 arg is provided:
1. Load the type schema
2. If `NameTemplate` is set, evaluate it and use as the name
3. If `NameTemplate` is not set, return an error asking for the name argument

**Rationale:** Minimal CLI change. The type arg is always required; only the name becomes conditional.

## Risks / Trade-offs

- **[Risk] Template syntax is limited to date only** → Acceptable for v1. The `{{ placeholder:arg }}` syntax is extensible for future placeholder types (property refs, sequences).
- **[Risk] Relaxing system property validation could be misused** → Mitigated by strict validation: only `template` field is allowed on `name` entries. Other system properties still fully rejected.
- **[Risk] Slug generation from template output** → Template output becomes the slug (used in filename). Must ensure the output produces valid slugs. `Slugify()` already handles this.
