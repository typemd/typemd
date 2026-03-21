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

This technical context feeds directly into the issue body (see Scope and Approach sections in body templates) and reduces rework during `resolve-issue` Phase 0 (Explore).

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

After setting the issue type, assign **one or more component labels**. Optionally add extra labels if applicable.

**Component label** (pick one or more):

| Label | Scope |
|---|---|
| `core` | Core library — objects, types, relations, index (`core/`) |
| `cli` | CLI commands (`cmd/`) |
| `tui` | Terminal UI (`tui/`) |
| `mcp` | MCP server (`mcp/`) |
| `web` | Web UI (`web/`) |
| `app` | Desktop app via Wails (`app/`) |

**Optional extra labels**:

| Label | When to use |
|---|---|
| `breaking-change` | Introduces breaking changes to public API or CLI |
| `discussion` | Needs discussion before implementation |
| `documentation` | Docs changes |

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

Use GraphQL to create the issue with the issue type set.

**IMPORTANT — use JSON file + `--input`**: The `gh api graphql -f` flag cannot pass array variables (like `labelIds`) correctly — it treats the JSON array as a single string, causing `NOT_FOUND` errors. Shell escaping of `!` in GraphQL types (e.g. `ID!`) also causes issues. Always write a JSON file with heredoc and pass it via `--input`.

```bash
# Get repo ID
REPO_ID=$(gh api repos/typemd/typemd --jq '.node_id')

# Get label IDs (one per label)
LABEL_ID_1=$(gh api repos/typemd/typemd/labels/<name1> --jq '.node_id')
LABEL_ID_2=$(gh api repos/typemd/typemd/labels/<name2> --jq '.node_id')

# Write GraphQL request as JSON file using heredoc (avoids shell escaping issues with `!`)
cat > /tmp/create_issue.json << 'EOF'
{
  "query": "mutation($repoId: ID!, $title: String!, $body: String!, $typeId: ID!, $labelIds: [ID!], $parentId: ID) { createIssue(input: { repositoryId: $repoId, title: $title, body: $body, issueTypeId: $typeId, labelIds: $labelIds, parentIssueId: $parentId }) { issue { number url } } }",
  "variables": {
    "repoId": "<REPO_ID>",
    "title": "<title>",
    "body": "<body with \\n for newlines>",
    "typeId": "<issue_type_id>",
    "labelIds": ["<LABEL_ID_1>", "<LABEL_ID_2>"],
    "parentId": null
  }
}
EOF

gh api graphql --input /tmp/create_issue.json
```

Omit or set to `null` the `labelIds` or `parentId` fields if not applicable. Clean up the temp file after use with `rm -f /tmp/create_issue.json`.

**After creation**, if a Release issue was selected and the new issue has no other parent, add it as a sub-issue:

```bash
RELEASE_ID=$(gh issue view <release_number> --json id --jq '.id')
ISSUE_ID=$(gh issue view <new_number> --json id --jq '.id')
gh api graphql -f query="mutation { addSubIssue(input: { issueId: \"$RELEASE_ID\", subIssueId: \"$ISSUE_ID\" }) { subIssue { number } } }"
```

If the issue already has a parent (e.g. an Epic), skip the sub-issue link and instead edit the Release issue body to reference it.

**After creation**, if there are blocking relationships, add them (one per blocking issue):

```bash
# Get the newly created issue's node ID
ISSUE_ID=$(gh issue view <number> --json id --jq '.id')
BLOCKER_ID=$(gh issue view <blocking_number> --json id --jq '.id')

# addBlockedBy: issueId = the blocked issue, blockingIssueId = the blocker
gh api graphql -f query="mutation { addBlockedBy(input: { issueId: \"$ISSUE_ID\", blockingIssueId: \"$BLOCKER_ID\" }) { issue { number } } }"
```

Repeat for each blocking issue. To remove a blocking relationship:

```bash
gh api graphql -f query="mutation { removeBlockedBy(input: { issueId: \"$ISSUE_ID\", blockingIssueId: \"$BLOCKER_ID\" }) { issue { number } } }"
```

### Body Templates

Read the matching template from `.github/ISSUE_TEMPLATE/<type>.yml` to determine the body structure. Each template defines sections with labels, descriptions, and required/optional validations — use those as-is to compose the issue body.

Omit optional sections if the user didn't provide relevant content during brainstorming.

### Step 9: Confirm

Return the issue URL to the user.
