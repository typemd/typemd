---
name: resolve-issue
description: |
  This skill should be used when the user asks to "resolve an issue", "work on issue #N", "fix #N", "implement #N", "close #N", "tackle #N", "pick up #N", "start working on #N", "what should I work on next", or references a specific GitHub issue number they want to work on. Can also accept a version number (e.g., "0.5.0") to select from that release's sub-issues, or auto-select the best issue when no argument is specified.
---

# Resolve Issue

Orchestrate the full lifecycle of resolving a GitHub issue — from reading the issue to opening a PR — with user confirmation at each phase.

All progress is tracked via **OpenSpec changes**, enabling resume from any interruption point.

## Resume Detection

Before starting, check if an OpenSpec change already exists for this issue:

```bash
openspec list --json
```

Look for a change name matching the issue (e.g., `issue-<N>-<slug>`). If found, check its status:

```bash
openspec status --change "<name>" --json
```

Based on artifact completion:

- **No proposal yet** → resume from Phase 0 (Explore). The explore phase is interactive, so the user can skip ahead to Phase 1 if they've already explored.
- **proposal/design/specs done, tasks done** → resume from Phase 2 (Implement), check task progress in `tasks.md`
- **All tasks complete** → resume from Phase 3 (Verify and Ship)

If prior progress is detected, present a summary and ask the user via AskUserQuestion:

- **"Continue from where we left off"** — resume from the detected point
- **"Start over"** — delete the existing change and begin from Preflight

If no matching change exists, start from Preflight.

## Preflight

Preflight covers all lightweight preparation steps before the main phases begin.

### Argument Parsing

The argument can be one of three forms:

| Input | Example | Interpretation |
|-------|---------|---------------|
| Issue number | `42`, `#42` | Resolve that specific issue → skip to **Check Issue State** |
| Version number | `0.5.0`, `v0.5.0` | Expand that version's sub-issues, then select one |
| Empty | *(none)* | List all open Release issues → choose a version → select a sub-issue |

**How to detect:**
- Matches `#?\d+` (with optional `#` prefix) → issue number
- Matches `v?\d+\.\d+\.\d+` → version number
- Otherwise → empty / auto-select

### Issue Selection (when no issue number is specified)

If the argument is a **version number** or **empty**, use the two-step flow below.

**Step 1: Choose a Release**

If no version is specified, list all open Release issues first:

```bash
./scripts/find-release-issues
```

Returns all open Release-type issues (verified via GraphQL `issueType`):

```json
{
  "releases": [
    { "number": 218, "title": "v0.5.0 — Filter, Sort & View" },
    { "number": 219, "title": "v0.6.0 — Navigate & Discover" }
  ]
}
```

Present the releases to the user via AskUserQuestion and let them pick one. If a version was already provided as an argument, skip this step.

**Step 2: Expand sub-issues for the chosen release**

```bash
./scripts/find-release-issues <version>
```

Returns the release's open sub-issues with body (truncated to 300 chars):

```json
{
  "release": { "number": 218, "title": "v0.5.0 — ..." },
  "issues": [
    { "number": 74, "title": "...", "labels": ["core"], "body": "..." }
  ]
}
```

**Step 3: Rank issues by priority**

Evaluate each issue using these criteria (highest priority first):

1. **Blocker** — blocks other issues (look for "blocks #N" or "blocked by #N" references in issue bodies, or issues labeled `blocker` / `priority:critical`)
2. **High value** — labeled `priority:high`, or is a bug affecting core functionality
3. **Low effort, high impact** — small scope issues that unblock progress (labeled `good first issue`, `quick win`, or estimated as small)
4. **Dependencies resolved** — issues whose blockers are already closed

**All issue types are valid candidates**, including `discussion` issues. If an issue belongs to a Release, it needs to be resolved in that timeframe regardless of its label.

**Step 4: Present top 3 candidates**

Ask the user via AskUserQuestion with the top 3 recommended issues:

- **"#N: \<title\>"** — for each candidate, include a one-line reason why it's recommended (e.g., "Blocks 3 other issues", "Critical bug", "Quick win for release X")

The user selects one, then proceed to **Check Issue State** with that issue number.

### Standalone Issue Lookup

When a specific issue number is given (not part of a release), use `get-issue-details` to fetch its details:

```bash
./scripts/get-issue-details <number> [number ...]
```

Returns a JSON array with each issue's number, title, labels, and body (truncated to 300 chars).

### Check Issue State

Verify the issue is actionable:

```bash
gh issue view <number> --json state,closedByPullRequestsReferences
```

- If the issue is **closed**, inform the user and stop.
- If there is already an **open PR** linked to this issue, inform the user and ask whether to continue or stop.

### Understand the Issue

Read the issue and confirm understanding with the user.

```bash
gh issue view <number> --json title,body,labels,assignees
```

Present a summary:

- **Title**
- **Type** (Bug / Feature / Task / Epic)
- **Labels**
- **Release** (which Release issue it belongs to, if any)
- **Key requirements** extracted from the body

Ask the user via AskUserQuestion:

- **"Looks correct"** — proceed
- **"I have additional context"** — let the user add info before proceeding

### Issue Type Routing

After confirming understanding, check the labels already retrieved above:

- If the issue has a **`discussion` label** → invoke the `resolve-discussion` skill and stop.
- Otherwise → continue to **Workspace Setup** below.

## Workspace Setup

Before entering the phases, set up an isolated working environment so that all artifacts (explore notes, design docs, code) live on a feature branch.

If a branch matching `fix/issue-<N>-*` or `feat/issue-<N>-*` already exists, inform the user and ask whether to reuse it or create a new one.

Ask the user how to set up the working environment via AskUserQuestion:

- **"Worktree (isolated)"** — invoke `superpowers:using-git-worktrees` skill
- **"Branch in current repo"** — create a branch directly

Branch naming convention:

- Bug → `fix/issue-<N>-<slug>`
- Feature / Task / Epic → `feat/issue-<N>-<slug>`

Where `<slug>` is a short kebab-case summary derived from the issue title (max 5 words).

```bash
git checkout -b <branch-name>
```

## Phases

### Phase 0: Explore

Use the `openspec-explore` skill to interactively explore the problem space with the user before committing to a design.

The goal of this phase is to:

- Clarify ambiguous requirements or edge cases in the issue
- Investigate the relevant codebase areas (existing code, data model, dependencies)
- Discuss trade-offs and possible approaches
- Surface hidden complexity or constraints early

The explore session should be grounded in the issue context gathered during Preflight. Pass the issue summary, key requirements, and any additional user context into the explore session.

This phase is interactive — continue the exploration dialogue until the problem is well-understood.

**Gate check before proceeding:** Since Phase 1 (Design) flows directly into Phase 2 (Implement) without a separate review pause, the explore phase must reach sufficient depth. Before leaving this phase, confirm with the user via AskUserQuestion:

- **Scope** — what's in and what's out
- **Approach** — the high-level technical direction (which packages, what patterns)
- **Edge cases** — any tricky scenarios surfaced during exploration

The user may end the exploration explicitly (e.g., "looks good", "let's proceed") or the explore skill may naturally conclude. Either way, ensure the three points above have been addressed.

Once exploration is complete, proceed to Phase 1.

### Phase 1: Design

Use the `openspec-propose` skill to create an OpenSpec change for this issue.

**Change naming convention:** `issue-<N>-<slug>` where `<slug>` is a short kebab-case summary derived from the issue title (max 5 words). Example: `issue-10-wiki-links-backlinks`.

The propose skill will create:
- `proposal.md` — what and why (derived from the issue description)
- `design.md` — how (architecture decisions, approach)
- `specs/<capability>/spec.md` — behavioral requirements with scenarios
- `tasks.md` — implementation steps

**Task ordering must follow test-first (BDD → TDD):**

For each feature group in `tasks.md`, tasks MUST be ordered test-first:

1. **BDD scenario first** — write `.feature` file with Gherkin scenarios (for `core/` and `tui/` changes)
2. **Step definitions** — implement BDD step definitions (initially failing)
3. **Implementation** — write production code to make BDD scenarios pass
4. **Unit tests** — add unit tests for edge cases, exact values, error conditions

Example of correct ordering:
```
## 1. Core: GetName

- [ ] 1.1 Write BDD scenarios for GetName (name present, missing, empty)
- [ ] 1.2 Implement step definitions for GetName scenarios
- [ ] 1.3 Add GetName() method to Object (make scenarios pass)
- [ ] 1.4 Add unit tests for GetName edge cases (whitespace, special chars)
```

Example of **incorrect** ordering (implementation before tests):
```
## 1. Core: GetName

- [ ] 1.1 Add GetName() method to Object    ← WRONG: implementation first
- [ ] 1.2 Write BDD scenarios for GetName   ← tests after implementation
```

For `cmd/` changes, BDD tests cover CLI-layer behavior (argument parsing, output format, error messages) in `cmd/features/`. For `mcp/`, use unit tests. See CLAUDE.md "Testing" section for full guidance.

Once all artifacts are generated, proceed directly to Phase 2.

### Phase 2: Implement

Use the `openspec-apply-change` skill to execute the tasks from the OpenSpec change. The apply skill reads `tasks.md` and implements each task in order.

Choose the appropriate implementation approach:

- **BDD** — the default for `core/`, `cmd/`, and `tui/` changes. BDD scenarios define behaviors and shared vocabulary (what a feature does), not implementation details. Write Gherkin `.feature` files first (in `<package>/features/`), then implement step definitions and production code. Use unit tests for precise validation (edge cases, output formats, exact values). For `mcp/`, use unit tests until BDD scope is decided.
- **Subagent-driven** (`superpowers:subagent-driven-development`) — when the plan has 3+ sequential steps that each produce testable output
- **Parallel agents** (`superpowers:dispatching-parallel-agents`) — when the plan has 2+ independent tasks with no shared state (e.g., separate packages, separate files)

If unsure, default to BDD with sequential implementation.

At key decision points, check with the user before proceeding.

### Post-Implementation Review

After all tasks are complete, run these two review steps before shipping:

1. **Simplify** — invoke the `simplify` skill to review changed code for reuse, quality, and efficiency, then fix any issues found.

2. **Update Documentation** — invoke the `update-doc` skill to compare implementation against all documentation and fix discrepancies.

### Phase 3: Verify and Ship

Execute the following steps in order:

1. **Verify** — invoke `superpowers:verification-before-completion` to confirm all tests pass, no regressions, and implementation matches the plan.

2. **Commit and Push** — invoke `git:commit-push` skill.

3. **Archive** — use the `openspec-archive-change` skill to archive the completed change. This syncs any delta specs to the main `openspec/specs/` directory and moves the change to `openspec/changes/archive/`.

4. **Open PR** — create a pull request using the project's PR template at `.github/pull_request_template.md` as the body structure:

```bash
gh pr create --title "<title>" --body "$(cat <<'EOF'
## Summary

- <bullet point 1>
- <bullet point 2>

## Issue

Closes #<N>

## Test Plan

- [ ] `go test ./...` — all pass
- [ ] `go build ./...` — clean build
- [ ] Manual: <specific manual steps>
EOF
)"
```

### Done

Present the PR URL to the user. The issue will be automatically closed when the PR is merged (via `Closes #N`).
