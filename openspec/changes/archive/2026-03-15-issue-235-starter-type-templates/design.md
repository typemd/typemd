## Context

`tmd init` currently creates the vault directory structure (`.typemd/types/`, `objects/`, SQLite index) and exits. New users face a blank slate with no guidance on how to start defining types. Issue #206 previously removed opinionated built-in types (book, person, note) from `defaultTypes`; this change reintroduces them as opt-in starter templates during initialization.

The `tui/` package already uses Bubble Tea, so the dependency exists. The existing interactive prompt pattern in `cmd/create.go` uses basic `bufio.Scanner`, but for this feature we want a richer checkbox-style selector.

## Goals / Non-Goals

**Goals:**

- Provide 3 starter type templates (idea, note, book) as opt-in during `tmd init`
- Use a Bubble Tea interactive checkbox UI with all items selected by default
- Support `--no-starters` flag for non-interactive/CI use
- Store starter templates as embedded YAML files (`core/starters/*.yaml`)

**Non-Goals:**

- Custom/user-contributed starter template packs
- A "gallery" or marketplace of templates
- Modifying `tmd init` for existing (already initialized) vaults
- Adding starter object templates (`templates/` directory) — only type schemas

## Decisions

### 1. Embedded YAML via `//go:embed` in `core/`

Store starter type definitions as `.yaml` files under `core/starters/` and embed them with `//go:embed`. This keeps the definitions in the same format users will see in `.typemd/types/`, making them self-documenting and easy to maintain.

**Alternative considered:** Go structs in code (like old `defaultTypes`). Rejected because YAML files are more readable, directly testable, and match the user-facing format.

### 2. `core/` exposes starter data; `cmd/` owns UI

- `core/starters.go` provides `StarterTypes() []StarterType` returning metadata (name, emoji, description) and raw YAML bytes for each starter.
- `core/vault.go` `Init()` signature remains unchanged. A new method `WriteStarterTypes(names []string) error` handles writing selected starters to `.typemd/types/`.
- `cmd/init.go` handles the Bubble Tea UI and `--no-starters` flag, then calls `WriteStarterTypes`.

This preserves Clean Architecture: core has no UI dependency, cmd orchestrates interaction.

### 3. Bubble Tea checkbox model in `cmd/`

A small self-contained Bubble Tea model (`cmd/starter_picker.go`) for the multi-select checkbox. Keeps it in `cmd/` since it's init-specific UI, not part of the main TUI app.

Key behaviors:
- Arrow keys to move cursor
- Space to toggle selection
- `a` to select all, `n` to deselect all
- Enter to confirm
- `q`/Esc to skip (select none)
- All items selected by default

### 4. Minimal starter type definitions

Each starter type has a small, opinionated set of properties. The goal is "useful starting point" not "comprehensive schema":

- **idea** (💡): `status` (select: seed/growing/ready)
- **note** (📝): no custom properties (system properties `name`, `description`, `tags` are sufficient)
- **book** (📚): `status` (select: to-read/reading/done), `rating` (number)

### 5. `--no-starters` flag

A simple boolean flag. When set, starter type selection is skipped entirely. The vault is initialized empty (current behavior). This is the safe default for CI/scripts.

## Risks / Trade-offs

- **[Risk] Starter types may not match user preferences** → Mitigation: Types are regular YAML files, fully editable/deletable after creation. The UI defaults to all-selected but allows deselection.
- **[Risk] Bubble Tea UI in `cmd/` may conflict with piped stdin** → Mitigation: `--no-starters` provides a clean non-interactive path. Bubble Tea handles non-TTY gracefully.
- **[Trade-off] Only 3 starter types** → Keeps the selection simple and fast. More types can be added later without architectural changes.
