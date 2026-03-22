---
name: resolve-issue
description: |
  This skill should be used when the user asks to "resolve an issue", "work on issue #N", "fix #N", "implement #N", "close #N", "tackle #N", "pick up #N", "start working on #N", "what should I work on next", or references a specific GitHub issue number they want to work on. Can also accept a version number (e.g., "0.5.0") to select from that release's sub-issues, or auto-select the best issue when no argument is specified.
---

# Resolve Issue

Orchestrate the full lifecycle of resolving a GitHub issue — from reading the issue to opening a PR — with minimal interruption.

Progress is tracked via **OpenSpec changes** for user-facing spec changes. Pure technical changes (labeled `tech-only`) skip OpenSpec entirely.

## Resume Detection

Before starting, check if an OpenSpec change already exists for this issue:

```bash
openspec list --json
```

Look for a change name matching the issue (e.g., `issue-<N>-<slug>`). If found, **delete it and start fresh** — always begin from Preflight.

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

If the argument is a **version number** or **empty**, use the flow below.

**Step 1: Choose a Release**

If no version is specified, list all open Release issues and **automatically select the one with the smallest version number** (e.g., `v0.5.0` before `v0.6.0`):

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

If a version was already provided as an argument, use that version directly.

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

**Step 4: Auto-select the best candidate**

Rank issues by the criteria above and **automatically select the top candidate**. However, before proceeding, check if a branch matching `fix/issue-<N>-*` or `feat/issue-<N>-*` already exists for that issue:

```bash
git branch --list "fix/issue-<N>-*" "feat/issue-<N>-*"
```

If a matching branch already exists, **skip that candidate** and try the next one. This avoids picking up issues that are already in progress in another session.

Inform the user which issue was selected and why, then proceed to **Check Issue State**.

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
- If there is already an **open PR** linked to this issue, inform the user and **stop**.

### Understand the Issue

Read the issue and assess whether it's ready for implementation.

```bash
gh issue view <number> --json title,body,labels,assignees
```

**Issue Type Routing:**

- If the issue has a **`discussion` label** → invoke the `resolve-discussion` skill and stop.
- If the issue has a **`tech-only` label** → mark as tech-only. This skips the OpenSpec design phase entirely (Phase 1). After Workspace Setup, proceed directly to **Phase 2: Implement** without OpenSpec.

**Readiness check (AI-judged):**

Review the issue body for completeness. The issue is ready if the Scope, Approach, and Edge Cases sections (for Feature/Task) or the Problem and Steps to Reproduce sections (for Bug) are adequately filled in from the `create-issue` brainstorming phase.

- If the issue is **ready** → inform the user which issue you're working on and proceed directly to **Workspace Setup**.
- If the issue has **gaps or ambiguities** that would block implementation → stop and ask the user for clarification via AskUserQuestion before proceeding.

## Workspace Setup

Always use a **git worktree** for isolated development. Invoke the `superpowers:using-git-worktrees` skill.

Branch naming convention:

- Bug → `fix/issue-<N>-<slug>`
- Feature / Task / Epic → `feat/issue-<N>-<slug>`

Where `<slug>` is a short kebab-case summary derived from the issue title (max 5 words).

## Phases

### Tech-only Fast Path

If the issue is labeled `tech-only` (detected in **Understand the Issue**), **skip Phase 1 entirely**. There is no OpenSpec change — no proposal, design, specs, or tasks artifacts are created.

Instead, after Workspace Setup, proceed directly to **Phase 2: Implement**. For tech-only issues:

- Read the issue body for the technical approach and scope
- Use `superpowers:writing-plans` to create a lightweight implementation plan (no OpenSpec)
- Implement directly using the appropriate approach (BDD, subagent-driven, or parallel agents)

Then continue to **Post-Implementation Review** and **Phase 3** as normal (skipping the Archive step since there is no OpenSpec change to archive).

### Phase 1: Design

> **Note:** Phase 0 (Explore) has been removed. All exploration — Scope, Approach, Edge Cases — is now done during `create-issue` brainstorming. The issue body should already contain this context.

> **Note:** This phase is skipped for `tech-only` issues. See **Tech-only Fast Path** above.

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

For **tech-only issues**: implement based on the lightweight plan created in the Tech-only Fast Path. There is no OpenSpec `tasks.md` — use the plan from `superpowers:writing-plans` instead.

For **regular issues**: use the `openspec-apply-change` skill to execute the tasks from the OpenSpec change. The apply skill reads `tasks.md` and implements each task in order.

Choose the appropriate implementation approach:

- **BDD** — the default for `core/`, `cmd/`, and `tui/` changes. BDD scenarios define behaviors and shared vocabulary (what a feature does), not implementation details. Write Gherkin `.feature` files first (in `<package>/features/`), then implement step definitions and production code. Use unit tests for precise validation (edge cases, output formats, exact values). For `mcp/`, use unit tests until BDD scope is decided.
- **Subagent-driven** (`superpowers:subagent-driven-development`) — when the plan has 3+ sequential steps that each produce testable output
- **Parallel agents** (`superpowers:dispatching-parallel-agents`) — when the plan has 2+ independent tasks with no shared state (e.g., separate packages, separate files)

If unsure, default to BDD with sequential implementation.

**Do not stop to ask the user** unless the planned approach is genuinely blocked (e.g., a required API doesn't exist, a dependency conflict makes the design infeasible). Scope being larger than expected is not a reason to stop — complete the full implementation.

### Post-Implementation Review

After all tasks are complete, run these two review steps before shipping:

1. **Simplify** — invoke the `simplify` skill to review changed code for reuse, quality, and efficiency, then fix any issues found.

2. **Update Documentation** — invoke the `update-doc` skill to compare implementation against all documentation and fix discrepancies.

### Phase 3: Verify and Ship

Execute the following steps in order:

1. **Verify** — invoke `superpowers:verification-before-completion` to confirm all tests pass, no regressions, and implementation matches the plan.

2. **Commit and Push** — invoke `git:commit-push` skill.

3. **Archive** (skip for `tech-only` issues) — use the `openspec-archive-change` skill to archive the completed change. This syncs any delta specs to the main `openspec/specs/` directory and moves the change to `openspec/changes/archive/`.

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
