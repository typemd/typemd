---
name: create-issue
description: This skill should be used when the user asks to "create an issue", "open an issue", "file a bug", "request a feature", "add a task", or discusses a problem or idea that should be tracked as a GitHub issue.
---

# Create GitHub Issue

Create a GitHub issue for the TypeMD project through iterative Q&A.

Do NOT create the issue until the user explicitly confirms. Ask questions one at a time to refine the idea before writing anything.

## Language

All issue content (title, body) MUST be written in **English**.

## Process

### Step 1: Understand the idea

Use the `superpowers:brainstorming` skill to explore the user's idea. The brainstorming session helps shape a vague idea into a concrete, actionable issue by:

- Understanding the motivation and problem space
- Clarifying what "done" looks like
- Surfacing constraints, related issues, and technical considerations
- Narrowing scope if the idea is too broad

**Technical exploration:** During brainstorming, also investigate the codebase to ground the discussion:

- Identify which packages and files are affected (`core/`, `cmd/`, `tui/`, etc.)
- Check existing patterns — how similar features are already implemented
- Surface potential edge cases or conflicts with existing behavior

**IMPORTANT:** This brainstorming phase **replaces** `resolve-issue` Phase 0 (Explore). The issue body must contain enough technical context for `resolve-issue` to skip exploration and proceed directly to design and implementation. Specifically, the following sections must be adequately filled:

- **Scope** — what's in and what's out
- **Approach** — high-level technical direction (which packages, what patterns, key design decisions)
- **Edge Cases** — tricky scenarios, conflicts with existing behavior, or open questions

If any of these cannot be determined during brainstorming (e.g., the user wants to defer a decision), add the `discussion` label so that `resolve-issue` knows to pause and clarify before implementing.

Once brainstorming concludes with a clear direction, proceed to the next step.

### Step 2: Check for duplicates

Before proceeding, search existing issues for potential duplicates:

```bash
gh issue list --state all --json number,title,state,labels --limit 100
```

Compare the new idea against existing issues by title and topic. If a similar issue exists, present it to the user via AskUserQuestion with options:

- **"It's a duplicate"** — stop and link to the existing issue
- **"Related but different"** — continue creating, and reference the related issue in the body
- **"Not related"** — continue creating as normal

Skip this step only if the user has already referenced a specific issue number in their request.

### Step 3: Determine issue type

Each issue gets exactly **one issue type**. Issue types replace the old type labels (`bug`, `enhancement`, `epic`, `chore`).

| Issue Type | ID | When to use |
|---|---|---|
| Task | `IT_kwDOD9OO7M4B4kKu` | CI, refactoring, dependencies, project configuration, general tasks |
| Bug | `IT_kwDOD9OO7M4B4kKv` | Something isn't working |
| Feature | `IT_kwDOD9OO7M4B4kKw` | New feature or improvement with clear scope |
| Epic | `IT_kwDOD9OO7M4B4rT8` | High-level feature plan, will be broken into sub-issues |

Suggest the type based on context. Ask for confirmation if ambiguous.

### Step 4: Determine labels

After setting the issue type, assign labels based on the feature's **impact area** in TypeMD. Multiple labels are allowed.

**Top-level labels** (pick one or more):

| Label | Scope | When to use |
|---|---|---|
| `🧱 core` | Core library (`core/`) | Affects core design but no specific sub-category fits |
| `💻 cli` | CLI commands (`cmd/`) | Affects CLI commands and argument design |
| `🖥️ tui` | Terminal UI (`tui/`) | Affects TUI interface and interaction flows |
| `🤖 mcp` | MCP server (`mcp/`) | Affects AI agent integration |
| `🌐 web` | Web UI (`web/`) | Affects web interface design |

**`core` sub-labels** (use together with `🧱 core`):

| Label | Scope | When to use |
|---|---|---|
| `🏗️ core/schema` | Type system, inheritance, namespace | Affects object type definitions, schema design |
| `🔍 core/query` | Query, filter, sort | Affects query syntax, dynamic views, saved queries |
| `🔗 core/relation` | Relations, backlinks, wiki-links | Affects inter-object relations, reference system |
| `🔄 core/sync` | Sync, collaboration, import/export | Affects multi-device sync, cross-vault, collaboration |
| `⚙️ core/automation` | Computed properties, scheduling, reminders | Affects automated computation, scheduled tasks |

**Labeling principles:**
- If a core sub-category fits → use `🧱 core` + the sub-label (e.g. `🧱 core` + `🔍 core/query`)
- If it affects core but no sub-category fits → use `🧱 core` alone
- If it spans multiple areas → use multiple labels (e.g. `🧱 core` + `🖥️ tui`)

**Optional extra labels**:

| Label | When to use |
|---|---|
| `💥 breaking-change` | Introduces breaking changes to public API or CLI |
| `💬 discussion` | Needs discussion before implementation |
| `📝 documentation` | Docs changes |
| `🔩 tech-only` | Pure technical change — refactoring, dependency updates, CI fixes, internal restructuring — no user-facing spec changes |

**Tech-only detection:** During the brainstorming phase (Step 1), assess whether the issue is a **pure technical change** — one that does not alter user-facing behavior, APIs, or specs. Examples:

- Refactoring internal code structure
- Dependency version bumps
- CI/CD pipeline fixes
- Code style / lint fixes
- Performance optimizations with no behavior change
- Internal naming or file organization changes

If the issue is tech-only, automatically suggest adding the `tech-only` label. This label tells `resolve-issue` to skip the OpenSpec design phase and go straight to implementation.

Suggest labels based on context. Ask for confirmation if ambiguous.

### Step 5: Assign to Release

Fetch open Release issues:

```bash
# Note: GitHub GraphQL filterBy does NOT support issueType filtering.
# Fetch issues and filter by title pattern or issueType field in results.
gh issue list --state open --json number,title --limit 50 --jq '.[] | select(.title | test("^v[0-9]+\\.[0-9]+\\.[0-9]+ —"))'
```

Present the Release issues as options using AskUserQuestion. Always include a "None" option for issues that don't belong to any release.

**Note:** The issue will be linked as a sub-issue of the selected Release issue in Step 8. GitHub sub-issues only allow **one parent per issue**:

- If the issue has **no parent** → add as sub-issue of the Release directly.
- If the issue has a parent that is **CLOSED** (e.g. a completed epic or block issue) → suggest removing the stale parent first, then add as sub-issue of the Release.
- If the issue has a parent that is **OPEN** (e.g. an active Epic) → cannot add as sub-issue of the Release. Reference it in the Release issue body instead.

### Step 6: Relationships (optional)

Proactively analyze existing issues to suggest relationships. Do NOT simply ask the user — do the research yourself and present findings.

Fetch open issues with `gh issue list --state open --json number,title,labels,issueType --limit 100`, then compare the new issue against them. Look for:

- **Potential parent (epic)**: Is there an open Epic that this issue logically belongs under? Match by topic, component, or feature area. Note: a parent epic is an **organizational tracker only** — it does not block its children. Children of the same epic may span multiple releases.
- **Potential blockers**: Are there open issues that must be resolved before this one can start? For issues under the same epic, look for **sibling dependencies** (e.g. "implement X" blocks "build Y on top of X"). Blocker relationships are between siblings, not between parent and child.
- **Related issues**: Issues in the same area that aren't parent/blocker but worth cross-referencing.

Present your findings to the user via AskUserQuestion. Format:

> **建議的關聯：**
>
> - **Parent**: #42 "Web UI storage interface" (Epic) — 這個 issue 屬於 Web UI 的範疇
> - **Blocked by**: #38 "Add VaultStorage abstraction" — 需要先完成 storage 介面
> - **Related**: #45 "React component library" — 同為 Web UI 元件
>
> 或者沒有找到明顯關聯。

Options:
- **"Accept all"** — apply all suggested relationships
- **"Let me pick"** — user selects which to keep
- **"No relationships"** — skip all

If the user wants to pick, present each suggestion individually for confirmation.

For confirmed relationships, look up the issue node ID with `gh issue view <number> --json id --jq '.id'`. Multiple relationships can be set. After the user confirms, proceed to the next step.

### Step 7: Draft and confirm

Present the full issue draft to the user, then use AskUserQuestion to confirm:

- **Title** — concise, plain language, no prefix
- **Type** — issue type name
- **Labels** — component + optional extra labels
- **Release** — selected Release issue or none
- **Relationships** — parent issue or blocking issues, if any
- **Body** — using the body template matching the issue type (see below)

Use AskUserQuestion with options like "Create" and "Needs changes" to get user confirmation. Only proceed after the user confirms.

### Step 8: Create issue

Use the `scripts/create-issue` helper script. It handles repo/label/issue ID resolution, GraphQL mutations, sub-issue linking, and blocker relationships in a single call.

```bash
bash scripts/create-issue \
  --title "Issue title" \
  --body "Issue body..." \
  --type task \
  --label "🧱 core" \
  --label "🔩 tech-only" \
  --release 220 \
  --blocked-by 310
```

**Options:**

| Flag | Description |
|---|---|
| `--title TEXT` | Issue title (required) |
| `--body TEXT` | Issue body (required) |
| `--type TYPE` | `task`, `bug`, `feature`, or `epic` (required) |
| `--label NAME` | Label name, repeatable |
| `--release NUMBER` | Release issue number — adds as sub-issue |
| `--parent NUMBER` | Parent issue number (mutually exclusive with `--release` for sub-issue) |
| `--blocked-by NUMBER` | Blocking issue number, repeatable |

**Output:** JSON `{ "number": 325, "url": "https://..." }`

**Notes:**
- If `--release` is set and no `--parent`, the issue is added as a sub-issue of the release
- If `--parent` is set, it's used as the parent and `--release` is skipped for sub-issue linking (reference the release in the body instead)
- The `--body` value can contain newlines — the script handles escaping via temp files and `jq`

**Removing a blocking relationship** (manual, not covered by the script):

```bash
ISSUE_ID=$(gh issue view <number> --json id --jq '.id')
BLOCKER_ID=$(gh issue view <blocking_number> --json id --jq '.id')
gh api graphql -f query="mutation { removeBlockedBy(input: { issueId: \"$ISSUE_ID\", blockingIssueId: \"$BLOCKER_ID\" }) { issue { number } } }"
```

### Body Templates

Read the matching template from `.github/ISSUE_TEMPLATE/<type>.yml` to determine the body structure. Each template defines sections with labels, descriptions, and required/optional validations — use those as-is to compose the issue body.

Omit optional sections if the user didn't provide relevant content during brainstorming.

### Step 9: Confirm

Return the issue URL to the user.
