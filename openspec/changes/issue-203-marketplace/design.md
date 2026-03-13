## Context

typemd's `marketplace/` directory is reserved for Claude Code plugins (per CLAUDE.md). Claude Code has a native plugin marketplace system where a git repo with `.claude-plugin/marketplace.json` serves as a plugin catalog. Users add marketplaces via `/plugin marketplace add owner/repo` and install plugins via `/plugin install name@marketplace`.

The `marketplace/` directory lives in the `typemd/typemd` monorepo, but is published to a separate `typemd/marketplace` repo via git subtree push so users can add it as a standard Claude Code marketplace.

## Goals / Non-Goals

**Goals:**

- Establish the marketplace directory structure following Claude Code's native plugin format
- Provide a working example plugin (`markdown-import`) that demonstrates the full plugin structure
- Validate contributed plugins automatically via CI
- Publish the marketplace directory to `typemd/marketplace` repo automatically on merge

**Non-Goals:**

- No browsing website for the marketplace
- No automated quality testing of skill content
- No plugin versioning automation
- No dependency management between plugins
- No changes to Go code, CLI, TUI, or core library

## Decisions

### 1. Monorepo subtree over separate repo

The marketplace source of truth lives in `typemd/typemd` under `marketplace/`. A GitHub Actions workflow uses `git subtree push` to publish to `typemd/marketplace` on every merge to main that touches `marketplace/`.

**Alternative considered**: Separate `typemd/marketplace` repo. Rejected because it adds coordination overhead between repos and makes it harder to review plugin contributions alongside core changes.

### 2. Relative path sources in marketplace.json

All plugins use relative paths (`./plugins/<name>`) in `marketplace.json` rather than GitHub repo sources. This works because users add the marketplace via git (`/plugin marketplace add typemd/marketplace`), which clones the full repo.

**Alternative considered**: Each plugin as a separate GitHub repo with `github` source type. Rejected as unnecessarily complex for a monorepo-hosted marketplace.

### 3. Claude Code plugin validation in CI

The validation workflow uses `claude plugin validate .` directly, since it provides authoritative validation against Claude Code's native plugin format. This ensures plugins are always valid according to the latest plugin specification.

**Alternative considered**: Shell-based validation with `jq` and file existence checks. Rejected because it duplicates logic that Claude Code already implements and could drift out of sync with the actual plugin format.

### 4. Subtree push triggered by path filter

The publish workflow runs only when files under `marketplace/` change, using GitHub Actions path filters. It pushes the `marketplace/` subtree to `typemd/marketplace` main branch.

**Alternative considered**: Manual publish step. Rejected because automation prevents the published repo from drifting out of sync.

## Risks / Trade-offs

- **Subtree push conflicts** → If someone pushes directly to `typemd/marketplace`, subtree push will fail. Mitigation: document that `typemd/marketplace` is read-only and managed by CI.
- **CI validation coverage** → Shell-based validation cannot catch all issues (e.g., broken skill logic). Mitigation: acceptable for MVP; human review covers content quality.
- **marketplace/ directory bloat** → As plugins grow, the monorepo gets larger. Mitigation: not a concern at MVP scale; revisit if plugin count exceeds ~50.
