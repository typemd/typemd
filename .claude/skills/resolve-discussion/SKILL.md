---
name: resolve-discussion
description: Use when resolving a GitHub issue labeled `discussion` — facilitates decision-making on open questions, documents conclusions as an issue comment, creates follow-up issues if needed, and removes the discussion label. Triggered by resolve-issue when it detects a discussion-labeled issue, or directly when user asks to "resolve discussion #N", "close discussion", "wrap up discussion". Can also accept a version number (e.g., "0.5.0") to select from that release's discussion sub-issues, or auto-select the best discussion issue when no argument is specified.
allowed-tools:
  - Bash(gh issue view:*)
  - Bash(gh issue comment:*)
  - Bash(gh issue edit:*)
  - Bash(gh issue list:*)
  - Bash(./scripts/find-release-issues:*)
  - Bash(./scripts/get-issue-details:*)
---

# Resolve Discussion

Facilitate decision-making for discussion issues, document conclusions, and remove the discussion label.

Discussion issues don't produce code — they produce **decisions**. The goal is to reach conclusions on open questions, document them, and create actionable follow-up issues if needed.

**Important:** This skill does NOT invoke `superpowers:brainstorming`. Brainstorming produces specs and implementation plans — discussions produce decisions. Keep the process lightweight: research the codebase, ask clarifying questions, propose approaches, capture decisions, and document them in the issue.

## Preflight

### Argument Parsing

The argument can be one of three forms:

| Input | Example | Interpretation |
|-------|---------|---------------|
| Issue number | `42`, `#42` | Resolve that specific discussion issue → skip to **Check Issue State** |
| Version number | `0.5.0`, `v0.5.0` | Expand that version's sub-issues, filter to discussions, then select one |
| Empty | *(none)* | List all open `💬 discussion` issues → select one |

**How to detect:**
- Matches `#?\d+` (with optional `#` prefix) → issue number
- Matches `v?\d+\.\d+\.\d+` → version number
- Otherwise → empty / auto-select

### Issue Selection (when no issue number is specified)

Discussion issues require user interaction, so **always let the user choose** — never auto-select.

**Step 1: Find discussion issues**

If a **version number** is provided:

```bash
./scripts/find-release-issues <version>
```

From the returned sub-issues, filter to only those with the `💬 discussion` label. If none found, inform the user and stop.

If **no argument** is provided, first check release-scoped discussions, then fall back to unscoped:

1. List all open Release issues and auto-select the smallest version:

```bash
./scripts/find-release-issues
```

2. Expand the selected release's sub-issues and filter to `💬 discussion` label.

3. If no release-scoped discussions found, fall back to all open discussion issues:

```bash
gh issue list --label "💬 discussion" --state open --json number,title,labels --limit 20
```

If no open discussion issues exist at all, inform the user and stop.

**Step 2: Present candidates to the user**

Always present the list via AskUserQuestion and let the user choose which discussion to resolve. If only one exists, still confirm with the user before proceeding (discussions need user buy-in).

### Standalone Issue Lookup

When a specific issue number is given, fetch its details:

```bash
./scripts/get-issue-details <number>
```

### Check Issue State

Verify the issue is actionable:

```bash
gh issue view <number> --json state,labels
```

- If the issue is **closed**, inform the user and stop.
- If the issue does **not** have the `💬 discussion` label, inform the user this is not a discussion issue and suggest using `resolve-issue` instead. Stop.

### Understand the Issue

Read the full issue context:

```bash
gh issue view <number> --json title,body,labels,assignees
```

## Phase 1: Research and Explore

Explore the discussion topic directly (do NOT invoke `superpowers:brainstorming`):

- Research relevant codebase areas, prior art, and constraints
- Check status of any sub-issues referenced in the issue body
- Summarize the issue's open questions and current state to the user
- Ask clarifying questions one at a time to understand the user's direction
- Propose 2-3 approaches with trade-offs when applicable

Keep it lightweight — the goal is to reach decisions, not to produce a spec or design document.

## Phase 2: Facilitate Decisions

Before formalizing decisions, **critically assess whether the listed open questions are complete**. The questions recorded in the issue body were identified at creation time — new questions may have emerged since then, or the original analysis may have missed something. Specifically:

- Review codebase findings from Phase 1 — do they reveal edge cases or constraints not covered by the listed questions?
- Consider interactions with other features or recent changes that weren't anticipated when the issue was created
- Think about failure modes, migration paths, backwards compatibility, or performance implications that may not have been called out

If you identify additional open questions, add them to the list. All questions — whether from the issue body or newly discovered — are simply open questions that need decisions.

Then formalize decisions on all open questions:

1. **List all open questions**
2. **For each question**, present:
   - Summary of what was explored during Phase 1
   - Options with trade-offs (use AskUserQuestion)
   - Recommendation if evidence supports one
3. **Capture decisions** — record the user's choice for each question

## Phase 3: Document and Resolve

Once all questions are resolved:

### 1. Post summary comment

```bash
gh issue comment <number> --body "$(cat <<'EOF'
## Discussion Summary

### Decisions

- <decision 1>
- <decision 2>

### Follow-up

- <follow-up issue description, if any>

---
Resolved via discussion.
EOF
)"
```

### 2. Create follow-up issues

If the discussion concluded with a large feature or multiple work streams, use the `break-down-epic` skill to decompose it into actionable sub-issues.

For simpler follow-ups (individual tasks or bugs), use the `create-issue` skill for each.

### 3. Remove the discussion label

Do NOT close the issue — the issue itself may still need implementation. Only remove the `💬 discussion` label to indicate the discussion is resolved:

```bash
gh issue edit <number> --remove-label "💬 discussion"
```

Present the summary to the user. Done.
