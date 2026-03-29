## Context

typemd manages objects as Markdown files and relies on git for version control. Existing CLI commands (`show`, `search`, `object list`) resolve object IDs to domain entities via the Vault facade, but none expose git history. Users must manually determine the file path and run `git log` themselves.

The `cmd/` package follows a consistent pattern: each command file registers a Cobra command, uses `openVault()` + `resolveObjectInteractive()` for object resolution, and registers itself in `init()`.

## Goals / Non-Goals

**Goals:**

- Provide `tmd log <object-id>` that displays git commit history for a specific object file
- Support prefix matching and interactive disambiguation (consistent with `show`, `link`, etc.)
- Support `--oneline` flag for compact output
- Detect and report when the vault is not inside a git repository

**Non-Goals:**

- Diff display (showing actual file changes per commit)
- Interactive log browsing (paging, scrolling)
- Non-git version tracking
- Support for `git log` flags beyond `--oneline` (e.g., `--since`, `--author`)
- Log for types or other non-object entities

## Decisions

### 1. Subprocess `git log` vs Go git library

**Decision:** Use `os/exec` to run `git log` as a subprocess.

**Rationale:** typemd already requires git to be installed (it's a git-backed tool). A subprocess call is simpler, has no dependency overhead, and automatically benefits from the user's git configuration (aliases, date formats, pager). A Go git library would add a large dependency for minimal benefit.

### 2. Command placement: `tmd log` (root) vs `tmd object log` (subcommand)

**Decision:** Register as `tmd log` on `rootCmd`, not under `objectCmd`.

**Rationale:** `git log` is a top-level command. `tmd log` follows the same ergonomic pattern. The command is used frequently enough to justify a short path. The issue spec explicitly says `tmd log <type>/<name>`.

### 3. Output formatting

**Decision:** Pass through git's native output formatting. Use `--oneline` for compact mode, default `git log --follow` for standard mode.

**Rationale:** Git's output is well-understood by developers. Reformatting would lose features like color coding and pagination. The `--follow` flag ensures renames are tracked.

### 4. Git repository detection

**Decision:** Run `git rev-parse --git-dir` before executing `git log` to verify the vault is inside a git repository.

**Rationale:** Provides a clear, typemd-specific error message rather than a confusing git error.

## Risks / Trade-offs

- **[Risk] git not installed** → The command returns a clear error ("git is not installed or not in PATH"). This is acceptable since typemd is a git-backed tool.
- **[Risk] Object file not yet committed** → `git log` returns empty output. The command should handle this gracefully with a message like "no commits found for this object".
- **[Trade-off] No custom formatting** → We pass through git's native output rather than parsing and reformatting. This means output consistency depends on the user's git config, but avoids reimplementing git's rich formatting.
