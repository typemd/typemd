---
name: release-note
description: Use when preparing a release, writing release notes, creating a GitHub draft release, or updating CHANGELOG.md. Triggers on "write release notes", "prepare release", "create release", "update changelog".
---

# Release Note

Write release notes and create a GitHub draft release for a given version. Also update `CHANGELOG.md`.

## Process

1. **Gather** — Collect all changes for the version
2. **Write** — Draft release notes and CHANGELOG entry
3. **Create** — Create GitHub draft release
4. **Done** — Present the draft URL to the user

## 1. Gather

Identify the version scope:

```bash
# Get milestone info
gh api repos/typemd/typemd/milestones --jq '.[] | select(.title=="v<VERSION>") | {title, closed_issues, open_issues}'

# List closed issues in milestone
gh issue list --milestone "v<VERSION>" --state closed --json number,title,labels --limit 50

# Get commit history since last tag (or all if first release)
git log --oneline <last-tag>..HEAD
```

If the milestone has open issues, warn the user before proceeding.

## 2. Write

### Format

Release notes and CHANGELOG use the same structure. The release note adds a highlights section and installation instructions on top.

```markdown
<1-2 paragraphs highlighting what's notable in this release — written as prose, not a list>

### Installation

```bash
go install github.com/typemd/typemd/cmd/tmd@v<VERSION>
```

Pre-built binaries for macOS, Linux, and Windows are available below.

### Added

- Objects & Types — define typed schemas in YAML, create objects with `tmd object create`
- Relations — bidirectional links via `tmd relation link` / `tmd relation unlink`
- ...

### Changed

- ...

### Fixed

- ...

### Documentation

- Docs: https://docs.typemd.io
- Blog: https://blog.typemd.io

### Thanks

<invite community participation, link to CONTRIBUTING.md>

**Full Changelog:** <milestone URL or compare URL>
```

Uses [Keep a Changelog](https://keepachangelog.com/) conventions:
- `Added` for new features
- `Changed` for changes in existing functionality
- `Fixed` for bug fixes
- `Removed` for removed features
- `Deprecated` for soon-to-be removed features

### Writing guidelines

- **Highlights first** — Open with 1-2 paragraphs of prose describing what's notable in this release. What should users be excited about? What problem does this solve?
- **Group by theme, not by package** — Users care about capabilities (Objects, Relations, TUI), not internal package structure (core, cmd).
- **Bullet format** — `Feature Name — one-sentence description (#issue)`
- **Issue numbers required** — Every CHANGELOG entry must reference the corresponding issue number. This is why we open issues — to trace changes back to decisions.
- **Pre-release** — For `v0.x` releases, mark as pre-release.

### CHANGELOG

Update both `CHANGELOG.md` (English) and `CHANGELOG.zh-TW.md` (Chinese) at the project root. Format:

```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

## [v0.1.0] - YYYY-MM-DD

### Added

- Objects & Types — define typed schemas in YAML, create objects with `tmd object create` (#18)
- Relations — bidirectional links via `tmd relation link` / `tmd relation unlink` (#XX)
- Wiki-links & Backlinks — `[[target]]` syntax with automatic backlink tracking (#10)
- ...

### Changed

- ...

### Fixed

- ...

[v0.1.0]: https://github.com/typemd/typemd/releases/tag/v0.1.0
```

CHANGELOG uses [Keep a Changelog](https://keepachangelog.com/) conventions:
- `Added` for new features
- `Changed` for changes in existing functionality
- `Fixed` for bug fixes
- `Removed` for removed features
- `Deprecated` for soon-to-be removed features

## 3. Create

Create a GitHub draft release:

```bash
gh release create v<VERSION> --draft --prerelease --title "v<VERSION>" --notes "<notes>"
```

- Always create as **draft** — the user decides when to publish
- Mark as **prerelease** for `v0.x` versions

If a draft already exists, update it:

```bash
gh api repos/typemd/typemd/releases/<ID> -X PATCH -f body="<notes>"
```

## 4. Done

Present:
- The draft release URL
- Summary of what's included
- Remind the user to review before publishing
