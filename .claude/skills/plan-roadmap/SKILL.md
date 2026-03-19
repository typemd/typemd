---
name: plan-roadmap
description: Use when the user asks to "plan a roadmap", "review release scope", "rebalance releases", "check if release is too large", or wants to adjust issue distribution across versions to fit a one-week release cadence.
---

# Plan Roadmap

Review all Release issues' scope, analyze issue relationships, and rebalance to fit a one-week release cadence.

Version planning is managed through **Release issues** (GitHub issues with the `Release` type). Each Release issue serves as a version container, with child issues linked via GitHub's sub-issue relationship. Milestones are no longer used for new versions.

## Language

All issue comments and Release issue content MUST be written in **English**.

## Principles

- One minor version = one week of work
- Every release must deliver visible user value — avoid pure tech-debt or pure testing releases
- Epic issues are organizational trackers — place them in the same release as their first child
- Children of an epic may span multiple releases; dependency order is between siblings, not parent → child
- Blocked issues cannot ship before their blockers (check `blockedBy` between sibling issues)
- Discussion-tagged issues need design first — don't schedule for immediate delivery
- Themes drive issue selection — not just effort balancing
- When a release is under budget, prefer non-business tasks (testing, infra, docs) over adding more features

## Process

### Step 1: Gather

Fetch all Release issues and their sub-issues, plus unassigned issues:

```bash
# All Release issues (open and closed)
# Note: GitHub GraphQL filterBy does NOT support issueType filtering.
# Fetch all issues with issueType field, then filter by issueType.name == "Release" in the results.
gh api graphql -f query='query {
  repository(owner:"typemd", name:"typemd") {
    issues(first: 100, states: [OPEN, CLOSED], orderBy: {field: CREATED_AT, direction: ASC}) {
      nodes {
        number title state
        issueType { name }
        subIssues(first: 30) {
          nodes { number title state labels(first: 5) { nodes { name } }
            parent { number title state }
          }
        }
      }
    }
  }
}' --jq '.data.repository.issues.nodes | map(select(.issueType.name == "Release"))'

# All open issues (to find unassigned ones)
gh issue list --state open --json number,title,labels --limit 100
```

Filter out issues that are already sub-issues of a Release to identify unassigned issues.

### Step 2: Analyze relationships

For all open issues (not just release sub-issues), query parent, sub-issues, and blocking relationships. Batch into groups of ~20 per GraphQL query:

```bash
gh api graphql -f query='query {
  repository(owner:"typemd", name:"typemd") {
    i<N>: issue(number: <N>) {
      number title state
      parent { number title state }
      subIssues(first: 20) { nodes { number title state } }
      blockedBy(first: 10) { nodes { number title state } }
      blocking(first: 10) { nodes { number title state } }
    }
  }
}'
```

Build a dependency map and flag issues that need cleanup:

- **Epic groups** — parent + children should move together
- **Block chains** — if A blocks B, A must ship first (same or earlier release)
- **Cross-release blockers** — flag as problems
- **Orphaned children** — issues whose parent is CLOSED (completed epic, finished block issue, etc.). These are leftover sub-issues that were never unlinked. **Flag these for cleanup**: suggest removing the stale parent so the issue can be freely added as a sub-issue of a Release.
- **Stale blockedBy** — issues blocked by a CLOSED issue. These don't affect functionality but add noise. Report them for awareness.
- **Missing blocking relationships** — review issue descriptions for implicit dependencies that aren't set as formal `blockedBy` relationships (e.g. "infrastructure" issues that must precede "feature" issues in the same domain).

Present all findings to the user before proceeding. Execute cleanup (unparent orphaned children, add/remove blocking relationships) with user confirmation.

#### Blocking relationship mutations

```bash
# Add: issueId = the blocked issue, blockingIssueId = the blocker
BLOCKED_ID=$(gh issue view <blocked_number> --json id --jq '.id')
BLOCKER_ID=$(gh issue view <blocker_number> --json id --jq '.id')
gh api graphql -f query="mutation { addBlockedBy(input: { issueId: \"$BLOCKED_ID\", blockingIssueId: \"$BLOCKER_ID\" }) { issue { number } } }"

# Remove a blocking relationship
gh api graphql -f query="mutation { removeBlockedBy(input: { issueId: \"$BLOCKED_ID\", blockingIssueId: \"$BLOCKER_ID\" }) { issue { number } } }"
```

### Step 3: Define themes for the next 5 releases

Look at the next 5 open Release issues. For each one, propose a **release theme** — a short, user-facing headline that:

- Describes the user-visible value in a concrete, exciting way (e.g. "Type System is Here", "Built-in Types", "TUI Editing")
- Follows naturally from the previous release
- Respects dependency order

**Headline style:** Think product update cards like Capacities — feature-forward, verb-driven, not abstract. "Navigation & Discovery" beats "Usability Improvements". Present all 5 themes together and discuss with the user before proceeding. Adjust based on feedback.

**If a Release issue has 0 open sub-issues**, flag it as ready to release — don't force issues into it.

### Step 4: Populate each themed release

Work through each release **one at a time** with the user:

1. **Start from existing sub-issues** already linked to that Release issue
2. **Scan all unassigned issues** for relevance to the theme — propose good candidates
3. **Move out** issues that don't fit the theme (to a later release or unlink)
4. **Check effort budget** — target 70–80% of weekly capacity (3.5–4 days)

**When under budget:** Look for non-business tasks first — testing, infra, documentation — before adding more features. These are good fillers that don't dilute the release's theme.

**Effort sizing:**

| Size | Criteria | Estimate |
|------|----------|----------|
| **Small** | Single file change, clear scope | < 1 day |
| **Medium** | Multiple files, some design, tests needed | 1–2 days |
| **Large** | Cross-package, significant design, new patterns | 3–5 days |

Heuristics: `discussion` label → add 1 day; new command → medium; epic tracker → sum children.

**Issue fitness criteria** (for adding from backlog):
- Complements or extends the release's theme
- Small or medium effort
- Not blocked by unresolved discussions
- Dependencies resolved in same or earlier release
- Epic tracker (parent): place in same release as its first child
- Sibling dependencies: if child A blocks child B, A must be in same or earlier release than B

### Step 5: Confirm and execute

Confirm each release's issue list with the user before executing moves. Work one release at a time — confirm theme and issues together, then execute immediately before moving to the next.

```bash
# Add a sub-issue to a Release issue
PARENT_ID=$(gh issue view <release_number> --json id --jq '.id')
CHILD_ID=$(gh issue view <number> --json id --jq '.id')
gh api graphql -f query="mutation { addSubIssue(input: { issueId: \"$PARENT_ID\", subIssueId: \"$CHILD_ID\" }) { subIssue { number } } }"

# Remove a sub-issue from a Release issue
gh api graphql -f query="mutation { removeSubIssue(input: { issueId: \"$PARENT_ID\", subIssueId: \"$CHILD_ID\" }) { subIssue { number } } }"

# Close obsolete issues
gh issue close <number> --comment "<reason>"

# Create new Release issue if needed (use create-milestone skill)
```

**GitHub sub-issues only allow one parent per issue.** When an issue already has a parent, handle it as follows (in order of preference):

1. **Parent is CLOSED** → remove the stale parent first (`removeSubIssue`), then add as sub-issue of the Release. This is the preferred approach — completed epics/block issues should not hold orphaned children.
2. **Parent is an OPEN Epic** → cannot add as sub-issue of the Release. Reference it in the Release issue body instead as a last resort.

Never leave issues orphaned under closed parents when they need to be tracked in a Release.

### Step 6: Report

Present final state of all affected Release issues:

```bash
gh api graphql -f query='query {
  repository(owner:"typemd", name:"typemd") {
    issues(first: 100, states: OPEN, orderBy: {field: CREATED_AT, direction: ASC}) {
      nodes {
        number title
        issueType { name }
        subIssues(first: 30) {
          totalCount
          nodes { state }
        }
      }
    }
  }
}' --jq '.data.repository.issues.nodes | map(select(.issueType.name == "Release"))'
```
