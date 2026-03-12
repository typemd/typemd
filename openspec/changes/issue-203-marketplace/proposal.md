## Why

There is no centralized place for the community to discover, share, or install typemd-related Claude Code skills. Users must manually copy skill files between projects. A plugin marketplace using Claude Code's native format enables community contribution while leveraging the existing plugin infrastructure.

## What Changes

- Add `marketplace/` directory with Claude Code native plugin marketplace structure (`.claude-plugin/marketplace.json`)
- Include a `markdown-import` example plugin that converts existing markdown files into typemd objects
- Add `CONTRIBUTING.md` with naming rules, quality requirements, and submission process
- Add GitHub Actions workflow to validate marketplace PRs (JSON syntax, plugin structure, SKILL.md frontmatter, README presence)
- Add GitHub Actions workflow to subtree-push `marketplace/` to the `typemd/marketplace` repo on merge to main

## Capabilities

### New Capabilities

- `marketplace-structure`: Claude Code plugin marketplace directory layout, marketplace.json catalog, and plugin organization conventions
- `marketplace-validation`: GitHub Actions CI pipeline for validating plugin structure and metadata on PRs
- `marketplace-publish`: GitHub Actions workflow for subtree-pushing marketplace/ to typemd/marketplace repo
- `markdown-import-plugin`: Example plugin with a Claude skill that converts existing markdown files into typemd objects

### Modified Capabilities

(none)

## Impact

- New `marketplace/` directory at repo root
- New GitHub Actions workflows in `.github/workflows/`
- Requires creating the `typemd/marketplace` repo on GitHub (empty, as a subtree target)
- No changes to existing Go code, CLI, TUI, or core library
