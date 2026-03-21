## Context

CLI commands (`show`, `link`, `unlink`) accept object ID prefixes. When a prefix matches multiple objects, `QueryService.Resolve()` returns an `AmbiguousMatchError` containing the prefix and all matching full IDs. Currently the CLI prints this error and exits — users must copy a full ID and re-run the command.

The codebase already uses Bubble Tea for interactive CLI selection in `cmd/starter_picker.go` (used by `tmd init`). All necessary dependencies (`bubbletea/v2`, `lipgloss/v2`) are already in `go.mod`.

## Goals / Non-Goals

**Goals:**

- When a prefix matches multiple objects in an interactive terminal, display a Bubble Tea picker so users can select the intended object without re-running the command
- Show human-readable object names (from `name` property) alongside full IDs for easier identification
- Maintain existing error behavior for non-interactive contexts (piped input, CI)
- Provide a shared helper so all ID-resolving commands get this behavior consistently

**Non-Goals:**

- Changing core resolution logic — `AmbiguousMatchError` already provides everything needed
- Adding disambiguation to TUI — it has its own object selection flow
- Fuzzy matching or search refinement within the picker
- Displaying object properties beyond name in the picker

## Decisions

### 1. CLI-layer interception, core unchanged

**Decision:** Handle disambiguation entirely in `cmd/` by catching `AmbiguousMatchError` and launching the picker. Core logic stays untouched.

**Rationale:** `AmbiguousMatchError` already carries `Matches []string` with full IDs. No need to add callbacks or interfaces to core. The TUI has its own navigation and doesn't need this. Keeping the change in `cmd/` minimizes blast radius.

**Alternatives considered:**
- *Core callback*: Add a `disambiguator func([]string) (string, error)` parameter to `QueryService.Resolve()`. Rejected — over-engineering for a CLI-only concern; TUI doesn't use `Resolve()` this way.

### 2. Shared `resolveIDInteractive` helper in `cmd/helper.go`

**Decision:** Add a `resolveIDInteractive(vault, prefix)` function that wraps `vault.ResolveID()`, checks for `AmbiguousMatchError`, detects TTY, and launches the picker if interactive.

**Rationale:** All three commands (`show`, `link`, `unlink`) need the same behavior. A single helper in the existing `helper.go` avoids duplication. Similarly, `resolveObjectInteractive` wraps `vault.ResolveObject()` for `show` which needs the full `*Object`.

### 3. Bubble Tea picker (consistent with starterPicker)

**Decision:** Build `disambiguatePicker` in a new `cmd/disambiguate_picker.go` following the same Bubble Tea model pattern as `starterPicker`.

**Rationale:** Consistency with existing codebase patterns. Users of `tmd init` will recognize the interaction model (↑↓ to move, Enter to select, Esc to cancel).

**Display format:** Each item shows the object's `name` property (queried from the index via `vault.QueryService`) on the primary line, with the full object ID as a secondary line below it. If the name cannot be resolved, fall back to displaying only the ID.

### 4. TTY detection via `os.Stdin` file descriptor

**Decision:** Use `term.IsTerminal(int(os.Stdin.Fd()))` from `golang.org/x/term` (already an indirect dependency via bubbletea) to detect interactive mode.

**Rationale:** Simple, standard Go approach. When stdin is not a terminal (piped), skip the picker and return the original `AmbiguousMatchError`.

## Risks / Trade-offs

- **[Risk] Bubble Tea captures terminal during picker** → Picker is short-lived (select and quit). Same pattern as `starterPicker` which works reliably.
- **[Risk] Object name lookup adds latency** → Names come from the already-opened SQLite index; lookup is sub-millisecond per candidate. Typical ambiguous matches have 2-5 candidates.
- **[Trade-off] Picker adds visual weight to simple CLI** → Acceptable because it only activates on ambiguity (an error path), and Esc provides immediate cancel with the original error as fallback.
