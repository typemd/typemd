---
name: break-down-epic
description: Use when the user asks to "break down an epic", "split an epic into sub-issues", "create sub-issues for an epic", or wants to decompose a large epic issue into smaller, actionable sub-issues on GitHub.
---

# Break Down Epic into Sub-Issues

Decompose a GitHub epic issue into sub-issues through analysis and user confirmation.

Do NOT create any sub-issues until the user explicitly confirms the full breakdown plan.

## Language

All issue content (title, body) MUST be written in **English**.

## Process

### Step 1: Read the epic

Fetch the epic issue content:

```bash
gh issue view <number> --json number,title,body,labels
```

Verify the issue is an Epic type. If not, warn the user and ask if they want to proceed anyway.

### Step 2: Extract planned features

Parse the epic body's **"Planned Features"** section. Each bullet point typically represents one sub-issue.

If the epic doesn't have a clear "Planned Features" section, analyze the Summary and body to propose a breakdown. Present the proposed breakdown to the user for feedback before proceeding.

### Step 3: Draft sub-issues

For each planned feature, draft a sub-issue with:

- **Title** — concise, descriptive, no prefix
- **Type** — typically Feature, but could be Task depending on scope
- **Labels** — inherit from parent epic, adjust per sub-issue if needed
- **Parent** — the epic issue number
- **Body** — use the Feature template (see create-issue skill for templates)

**Issue type IDs:**

| Issue Type | ID |
|---|---|
| Task | `IT_kwDOD9OO7M4B4kKu` |
| Feature | `IT_kwDOD9OO7M4B4kKw` |

**Ordering:** Order sub-issues by implementation dependency — foundational work first, dependent features later.

### Step 4: Present and confirm

Present ALL drafted sub-issues in a numbered list showing title, type, and labels for each.

Use AskUserQuestion to ask the user to review. Options:

- **"Create all"** — proceed to create all sub-issues
- **"Needs changes"** — user wants to modify the breakdown

If the user wants changes, iterate on the draft before proceeding.

### Step 5: Create sub-issues

Fetch IDs needed for creation:

```bash
REPO_ID=$(gh api repos/typemd/typemd --jq '.node_id')
PARENT_ID=$(gh issue view <epic_number> --json id --jq '.id')

# For each label
LABEL_ID=$(gh api repos/typemd/typemd/labels/<name> --jq '.node_id')
```

Create each sub-issue using GraphQL with `--input` (same pattern as create-issue):

```bash
cat > /tmp/create_issue.json << 'EOF'
{
  "query": "mutation($repoId: ID!, $title: String!, $body: String!, $typeId: ID!, $labelIds: [ID!], $parentId: ID) { createIssue(input: { repositoryId: $repoId, title: $title, body: $body, issueTypeId: $typeId, labelIds: $labelIds, parentIssueId: $parentId }) { issue { number url } } }",
  "variables": {
    "repoId": "<REPO_ID>",
    "title": "<title>",
    "body": "<body with \\n for newlines>",
    "typeId": "<issue_type_id>",
    "labelIds": ["<LABEL_ID_1>"],
    "parentId": "<PARENT_ID>"
  }
}
EOF

gh api graphql --input /tmp/create_issue.json
rm -f /tmp/create_issue.json
```

Create sub-issues **sequentially** (not in parallel) to ensure consistent ordering.

### Step 6: Set blocking relationships (optional)

If sub-issues have dependencies on each other, ask the user whether to set blocking relationships between them.

```bash
# addBlockedBy: issueId = the blocked issue, blockingIssueId = the blocker
BLOCKED_ID=$(gh issue view <number> --json id --jq '.id')
BLOCKER_ID=$(gh issue view <blocking_number> --json id --jq '.id')
gh api graphql -f query="mutation { addBlockedBy(input: { issueId: \"$BLOCKED_ID\", blockingIssueId: \"$BLOCKER_ID\" }) { issue { number } } }"
```

### Step 7: Report

Return a summary table of all created sub-issues with issue numbers and URLs.
