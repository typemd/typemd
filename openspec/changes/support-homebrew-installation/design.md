## Context

typemd currently distributes binaries via GitHub Releases, powered by GoReleaser and a release workflow (`.github/workflows/release.yml`). Users install via `go install` or manual binary download. Adding Homebrew tap support automates formula generation on each release with minimal ongoing maintenance.

The `typemd/homebrew-tap` repository will follow Homebrew's naming convention: a repo named `homebrew-tap` under the `typemd` org enables `brew tap typemd/tap`.

## Goals / Non-Goals

**Goals:**

- Fully automated formula publishing on each `v*` tag release
- Zero manual formula maintenance — GoReleaser handles everything
- Formula name `typemd-cli`, binary name `tmd`
- macOS support (Intel + Apple Silicon)

**Non-Goals:**

- Homebrew Cask for desktop app (future, when `app/` is ready)
- Submission to `homebrew-core`
- Linux Homebrew (Linuxbrew) — not explicitly targeted but GoReleaser produces Linux binaries so it may work

## Decisions

### 1. GoReleaser `brews` section over manual formula

**Decision**: Use GoReleaser's built-in `brews` configuration.

**Rationale**: GoReleaser already runs on every release. Its `brews` feature auto-generates the Ruby formula with correct SHA256 checksums and version, then pushes it to the tap repo. Manual formula maintenance would require a separate workflow and is error-prone.

**Alternative considered**: GitHub Actions workflow that generates formula from release assets — rejected as redundant since GoReleaser does this natively.

### 2. Fine-grained PAT for cross-repo access

**Decision**: Use a fine-grained Personal Access Token scoped to `typemd/homebrew-tap` with `contents: write` permission, stored as `HOMEBREW_TAP_TOKEN` repository secret.

**Rationale**: The default `GITHUB_TOKEN` can only write to the repository where the workflow runs. A fine-grained PAT limits blast radius to just the tap repo, following least-privilege principle.

**Alternative considered**: Classic PAT with `repo` scope — rejected as overly broad. GitHub App token — viable but higher setup complexity for minimal benefit.

### 3. Formula naming strategy

**Decision**: Formula name `typemd-cli`, reserving `typemd` for future Cask.

**Rationale**: Separates the CLI tool from the future desktop app namespace. Users run `brew install typemd/tap/typemd-cli` for the CLI. When the desktop app ships, `brew install --cask typemd/tap/typemd` won't conflict.

## Risks / Trade-offs

- **[PAT expiration]** → Fine-grained PATs have a max expiration. Set a reminder to rotate. Consider switching to GitHub App token if this becomes burdensome.
- **[tap repo is public]** → Homebrew taps must be public repos. No sensitive data in the formula, so this is acceptable.
- **[GoReleaser version changes]** → `brews` config syntax may change across major versions. Pin GoReleaser v2 in the workflow (already done).
