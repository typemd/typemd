## Context

typemd CLI commands accept object IDs in `type/name-ulid` format (e.g. `book/clean-code-01kk39c30x27ck7ahyc7ct4nyn`). Without tab completion, users must type or copy-paste these long IDs. The codebase already has prefix resolution via `QueryService.Resolve()` and `GlobIDs()` which glob the filesystem — these provide the foundation for completion candidates.

Commands needing completion:
- **Object ID**: `show` (arg 0), `link` (args 0, 2), `unlink` (args 0, 2)
- **Type name**: `migrate` (arg 0), `type show` (arg 0), `stats --type`, `format --type`
- **Relation name**: `link` (arg 1), `unlink` (arg 1) — context-dependent on from-object's type

## Goals / Non-Goals

**Goals:**
- Progressive two-stage object ID completion: type prefix → object name within type
- Type name completion for commands accepting type arguments or `--type` flags
- Relation name completion based on source object's schema
- `tmd completion` subcommand for bash, zsh, fish script generation
- Fast completion using filesystem reads (no SQLite dependency)

**Non-Goals:**
- Completion for object property values
- Completion inside the TUI (already has its own navigation)
- Remote/network-based completion
- Fuzzy matching (Cobra completion uses prefix matching)

## Decisions

### 1. Filesystem-based completion (not SQLite)

Completion functions will read the `objects/` directory structure directly via `Vault.ListTypes()` and `ObjectRepository.GlobIDs()`, not query SQLite.

**Rationale**: Completion must work even when the index is stale or missing. Files are the source of truth. The existing `GlobIDs` function already does exactly the right thing — `filepath.Glob(prefix+"*.md")`.

**Alternative considered**: SQLite queries would be faster for large vaults but add a hard dependency on an up-to-date index, which contradicts the "files are source of truth" principle.

### 2. Two-stage progressive completion for object IDs

When the user types `tmd show <TAB>`:
1. If no `/` in the current word → complete type names with `/` suffix (e.g. `book/`)
2. If word contains `/` → split on first `/`, use type as directory, complete object names within that type

**Rationale**: Mirrors the natural `type/name` path structure. Users discover types first, then objects within a type.

### 3. Shared completion helper functions in `cmd/completion.go`

Create reusable helper functions:
- `completeObjectID(vault, args, toComplete)` — two-stage object ID completion
- `completeTypeName(vault, toComplete)` — type name completion
- `completeRelationName(vault, fromID, toComplete)` — relation name completion based on from-object's schema

Each command's `ValidArgsFunction` delegates to these helpers. The `link`/`unlink` commands use positional dispatch (arg 0 → object, arg 1 → relation, arg 2 → object).

### 4. Cobra's built-in completion command

Use `rootCmd.CompletionOptions` and Cobra's auto-generated `completion` subcommand rather than writing a custom one. Cobra generates correct scripts for bash, zsh, fish, and PowerShell.

**Rationale**: Less code to maintain. Cobra's completion scripts handle edge cases (quoting, escaping) that are hard to get right manually.

### 5. Vault opening strategy for completion

Completion functions need a `Vault` instance. For type-name completion, `resolveVault()` (no Open/index) suffices since `ListTypes()` reads the filesystem. For object ID completion, `resolveVault()` + direct `GlobIDs` is sufficient. For relation name completion, we need `LoadType()` to read the schema — also filesystem-only.

No `vault.Open()` (no SQLite) is needed for any completion function.

## Risks / Trade-offs

- **[Performance on large vaults]** → `GlobIDs` does filesystem glob, which is fast for hundreds of objects but could slow down with thousands. Mitigation: Cobra limits completion to `ShellCompDirectiveNoFileComp` so the shell won't fallback to file completion on timeout.
- **[Vault path resolution]** → Completion runs before the command executes, so `--vault` flag may not be parsed yet. Mitigation: Use Cobra's `Flag("vault").Value.String()` to read the flag value if set, otherwise default to `.`.
- **[Stale completions]** → Filesystem reads may show deleted-but-not-yet-synced objects. Acceptable: completion is best-effort; the command will error on invalid IDs.
