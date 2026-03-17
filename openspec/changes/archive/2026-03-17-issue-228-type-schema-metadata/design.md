## Context

TypeSchema currently has `Name`, `Plural`, `Emoji`, `Unique`, `Version`, and `Properties`. Property has `Name`, `Use`, `Type`, `Emoji`, `Pin`, `Options`, `Target`, `Default`, `Multiple`, `Bidirectional`, `Inverse`, and `Template`. Neither struct has a `description` or `color` field.

The TUI type editor uses `metaFieldCount = 4` for cursor navigation across meta fields (Name, Plural, Emoji, Unique). The `marshalSchema` struct mirrors TypeSchema for YAML serialization.

The `validateUseOverrides()` function currently allows only `pin` and `emoji` overrides on `use` entries.

## Goals / Non-Goals

**Goals:**
- Add `color` and `description` fields to TypeSchema
- Add `description` field to Property
- Validate color values (preset names and hex format)
- Allow `description` override on shared property `use` entries
- Update TUI type editor to display and edit new fields
- Update YAML serialization to include new fields

**Non-Goals:**
- Rendering colors in TUI (lipgloss theming) — that's a separate UI concern
- Adding color to Property — only TypeSchema gets color
- Adding color presets as an enum type — they remain simple strings validated against a list

## Decisions

### Decision 1: Color validation approach

**Choice:** Validate against a fixed preset list + hex regex.

Preset names: `red`, `blue`, `green`, `yellow`, `purple`, `orange`, `pink`, `cyan`, `gray`, `brown`.

Hex format: `#RRGGBB` (6-digit) or `#RGB` (3-digit), case-insensitive.

**Rationale:** Preset names give quick, memorable options. Hex gives full customization. Validation at the schema level catches typos early. The preset list is intentionally small — it covers the most common UI color needs without creating decision paralysis.

**Alternative considered:** No validation (treat as opaque string, let UI interpret). Rejected because typos in hex codes would silently fail in the UI.

### Decision 2: Color field position in YAML

**Choice:** `color` appears after `emoji`, before `unique`. `description` appears after `unique`, before `version`.

```yaml
name: presentation
plural: presentations
emoji: "🖥️"
color: green
unique: false
description: "Slide decks and presentation materials"
version: "1.0"
properties:
  - name: speaker
    type: relation
    target: person
    description: "The person who gave this presentation"
```

**Rationale:** Groups visual metadata together (`emoji`, `color`), then identity/documentation (`unique`, `description`), then versioning. This feels natural when reading a YAML file.

### Decision 3: Property description field position in YAML

**Choice:** `description` appears after `emoji` (or after `name` if no emoji), before `pin`.

**Rationale:** Description is documentation metadata — it belongs near the top of a property definition, after identification fields (`name`, `type`, `emoji`) and before behavioral fields (`pin`, `options`, `target`).

### Decision 4: Use entry override expansion

**Choice:** Add `description` to allowed `use` overrides alongside `pin` and `emoji`.

**Rationale:** A shared property like `due_date` might be described as "Project deadline" in one type and "Payment due date" in another. The description override enables context-specific documentation.

### Decision 5: TUI meta field ordering

**Choice:** `metaFieldCount` increases from 4 to 6. New order: Name (0), Plural (1), Emoji (2), Color (3), Unique (4), Description (5).

**Rationale:** Matches the YAML field order for consistency. Color follows emoji (visual fields), description follows unique (documentation fields).

## Risks / Trade-offs

- **[Risk] TUI type editor complexity increases** → Mitigated by following the existing pattern for meta field editing. The cursor logic is already index-based, so adding two more indices is straightforward.
- **[Risk] Color preset list may need future expansion** → Mitigated by making the preset list a package-level variable that's easy to extend. No breaking change to add more presets later.
- **[Trade-off] Hex validation is format-only** → We validate the format (#RGB or #RRGGBB) but not that the color is visually distinct or accessible. This is intentional — color accessibility is a UI-layer concern.
