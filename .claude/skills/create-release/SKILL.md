---
name: create-release
description: Use when preparing a release, writing release notes, creating a GitHub draft release, or updating CHANGELOG.md. Triggers on "write release notes", "prepare release", "create release", "update changelog".
---

# Release Note

Write release notes, update CHANGELOG, write a blog post, and create a GitHub draft release for a given version.

## Input

The user must provide a version number (e.g. `v0.2.0`). If not provided, ask before proceeding.

## Process

1. **Gather** — Collect closed issues and relevant commits for the version
2. **Write** — Draft release notes and CHANGELOG entry
3. **Blog** — Write a blog post in zh-tw (English will be synced via `sync-blog`)
4. **Finish** — Commit, push, and create GitHub draft release

## Language

- Release notes and CHANGELOG: **English**
- Blog post: **Traditional Chinese (zh-tw)**

---

## 1. Gather

```bash
# Get the previous release tag (for commit range)
PREV_TAG=$(git tag --sort=-version:refname | head -2 | tail -1)

# Closed issues in the Release issue (find by title matching version)
# Note: GitHub GraphQL filterBy does NOT support issueType filtering.
# Use --jq to filter Release issues by title pattern instead.
gh api graphql -f query='query {
  repository(owner:"typemd", name:"typemd") {
    issues(first: 50, states: [OPEN, CLOSED], orderBy: {field: CREATED_AT, direction: DESC}) {
      nodes {
        number title state
        subIssues(first: 30) {
          nodes { number title state labels(first: 5) { nodes { name } } }
        }
      }
    }
  }
}' --jq '.data.repository.issues.nodes | map(select(.title | test("^v[0-9]+\\.[0-9]+\\.[0-9]+ —")))'
# Filter for the Release issue matching v<VERSION> and extract its closed sub-issues

# Commits since last tag — exclude chore, skill, and docs-only commits
git log ${PREV_TAG}..HEAD --oneline --no-merges
```

**Commit filtering rules** — exclude commits whose type prefix is:
- `chore` — maintenance, deps, tooling
- `docs` — documentation-only changes
- `skill` or `chore(skill)` — internal skill updates

Include: `feat`, `fix`, `refactor` (if user-visible), `perf`

If the Release issue has open sub-issues, warn the user before proceeding.

---

## 2. Write

### Release notes format

```markdown
<1-2 paragraphs of prose — what's notable, what problem this solves, what users get>

### Installation

```bash
go install github.com/typemd/typemd/cmd/tmd@v<VERSION>
```

Pre-built binaries for macOS, Linux, and Windows are available below.

### ⚠️ Breaking Changes

- Change Name — what breaks and how to migrate (#issue)

### Added

- Feature Name — one-sentence description (#issue)

### Changed

- ...

### Fixed

- ...

### Documentation

- Docs: https://docs.typemd.io
- Blog: https://blog.typemd.io

### Thanks

<invite community participation, link to CONTRIBUTING.md>

**Full Changelog:** https://github.com/typemd/typemd/compare/v<PREV>...v<VERSION>
```

### Writing guidelines

- **Highlight lead** — Open with prose, not a list. What should users be excited about?
- **Group by theme, not package** — Users care about capabilities (Objects, TUI, CLI), not internal packages (core, cmd)
- **Bullet format** — `Feature Name — one-sentence description (#issue)`
- **Issue numbers required** — Every entry must reference the corresponding issue
- **Pre-release** — For `v0.x`, mark as pre-release
- **Breaking changes** — Always check for breaking changes (schema changes, removed commands, renamed flags, changed behavior). If any exist, add a `### ⚠️ Breaking Changes` section at the top with migration instructions. Also add a `### Breaking Changes` section to the CHANGELOG.

### CHANGELOG

Update both `CHANGELOG.md` (English) and `CHANGELOG.zh-TW.md` (Chinese) at the project root.

```markdown
## [v<VERSION>] - YYYY-MM-DD

### Added

- Feature Name — description (#issue)

### Changed / Fixed / Removed

- ...

[v<VERSION>]: https://github.com/typemd/typemd/releases/tag/v<VERSION>
```

---

## 3. Blog

### Read the source before writing examples

Before writing any code or YAML examples, read the actual implementation to verify syntax:

```bash
# Check example vault for correct schema format
cat examples/book-vault/.typemd/types/*.yaml
cat examples/book-vault/.typemd/properties.yaml

# Check feature files for supported types and behaviors
cat core/features/*.feature
```

Do not guess property types, field names, or YAML structure from issue titles or commit messages — always verify against real code.

### Write the post

Write a blog post at:

```
websites/blog/src/content/posts/zh-tw/release-<VERSION-WITH-DASHES>.md
```

e.g. `release-0-2-0.md` (use dashes instead of dots to avoid URL issues)

**Frontmatter:**

```yaml
---
title: "TypeMD <VERSION> — <one-line theme headline>"
description: "<1-2 sentence teaser>"
date: <today's date>
tags: [發布, 開發日誌]
---
```

**Content guidelines:**

- Write in Traditional Chinese (zh-tw)
- Follow the [Capacities](https://capacities.io/whats-new) writing style: conversational, problem-solution framing, user-benefit focused, short punchy section titles
- Use we, not I
- Explain key features with concrete, **verified** examples (code snippets, commands)
- Keep it engaging, not just a feature list — tell a story around the release theme
- **Focus on user-facing value** — explain what users can do now, not how the internals changed. Technical refactors (CQRS, DDD, architecture changes) belong in the CHANGELOG and GitHub release notes, but NOT in the blog post. The blog should only cover features, behaviors, and workflows that users interact with directly.
- End with a forward-looking sentence about what's next

**After writing zh-tw, present the draft to the user for review before writing the file. Only write the file after the user confirms.**

After the user approves the zh-tw post, use the `sync-blog` skill to create the English version.

---

## 4. Finish

Once release notes, CHANGELOG, and blog post are ready:

### Commit and push

```bash
git add CHANGELOG.md CHANGELOG.zh-TW.md \
  websites/blog/src/content/posts/zh-tw/release-<VERSION-WITH-DASHES>.md \
  websites/blog/src/content/posts/en/release-<VERSION-WITH-DASHES>.md

git commit -m "chore(release): prepare v<VERSION> release notes and blog post"
git push origin main
```

### Create GitHub draft release

```bash
gh release create v<VERSION> \
  --draft \
  --prerelease \
  --title "v<VERSION>" \
  --notes "<release notes content>"
```

- Always create as **draft** — the user decides when to publish
- Mark as **prerelease** for `v0.x` versions
- If a draft already exists, update it:

```bash
gh api repos/typemd/typemd/releases/<ID> -X PATCH -f body="<notes>"
```

### Done

Present:
- The draft release URL
- The blog post file path
- Remind the user to review before publishing
