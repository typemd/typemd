---
name: skill-retro
description: |
  Retrospective analysis and correction of skills based on execution experience. Analyzes the current conversation to identify what went wrong during skill execution and applies targeted fixes to the SKILL.md and metadata.
disable-model-invocation: true
user-invocable: true
---

# Skill Retro

Analyze the current session's skill execution history, identify corrections that occurred, and apply those learnings back into the skill's SKILL.md and metadata.

## Why This Skill Exists

Skills are written with assumptions about how execution will flow. In practice, three kinds of corrections happen:

1. **User interruption** — the user stops the flow and says "do it this way instead." This means the skill's instructions led to the wrong action or an unnecessary confirmation step.
2. **Self-correction** — the agent hits an error, recovers, and finds the right approach. This means the skill could have guided the agent to the right approach from the start.
3. **Unnecessary pause** — the agent stops to ask or confirm when the user's next instruction shows it could have proceeded directly. This means the skill is overly cautious or has an unnecessary checkpoint.

Each of these is a signal that the skill's instructions can be improved. The goal is to feed these signals back into the SKILL.md so future executions avoid the same problems.

## Step 1: Identify Target Skill

Scan the conversation history for evidence of skill execution. Look for:

- `<command-name>/skill-name</command-name>` tags — these mark skill invocations
- `Skill` tool calls — these show which skills were loaded
- Patterns of tool calls that match a known skill's workflow (e.g., `gh issue view` sequences suggest `resolve-issue`)

If multiple skills were executed, rank them by how many corrections occurred during their execution. Present the findings to the user via AskUserQuestion:

- List the detected skill(s) with a one-line summary of what issues were observed
- Let the user confirm which skill to retro, or specify a different one

If the user passed an argument (e.g., `/skill-retro resolve-issue`), use that directly but still summarize the observed issues for confirmation.

## Step 2: Catalog Corrections

Walk through the conversation chronologically and extract every correction event. For each one, record:

| Field | Description |
|-------|-------------|
| **Type** | `user-interrupt`, `self-correct`, or `unnecessary-pause` |
| **Location** | Which step/phase of the skill was executing |
| **What happened** | The action the agent took or was about to take |
| **What should have happened** | The correct action (from user feedback or the successful recovery) |
| **Root cause in skill** | Which part of the SKILL.md led to the wrong behavior |

### How to Identify Each Type

**User interruption** — look for:
- Tool calls that return "The user doesn't want to proceed with this tool use"
- User messages that redirect the flow (e.g., "不要這樣做", "改成...", "直接...")
- User messages that provide the correct approach after an agent action

**Self-correction** — look for:
- Tool calls that fail, followed by a different approach that succeeds
- Agent messages like "let me try a different approach"
- Sequences where the agent backtracks

**Unnecessary pause** — look for:
- AskUserQuestion calls where the user's response indicates the agent could have just proceeded
- Confirmation requests where the answer was obvious from context
- The agent stopping execution and the user saying "繼續" or "你直接做就好"

### Bash Command Analysis

In addition to correction events, scan the conversation for Bash commands executed during the skill's workflow. Look for opportunities to consolidate:

- **Multi-step sequences** — multiple Bash calls that always run together (e.g., a GraphQL query followed by `--jq` filtering, then another query using the first result). These can become a single script.
- **Complex inline commands** — long one-liners with pipes, jq filters, or string manipulation that are hard to read and error-prone inline.
- **Repeated patterns** — the same command structure used multiple times with different parameters.

For each candidate, note:
- Which Bash commands would be consolidated
- What the script would do (inputs → outputs)
- Where in the skill it would be called

These will be addressed alongside correction fixes in Step 3.

### Present Catalog

Present the catalog to the user as a numbered list. For each item, show:
- The correction type (emoji: 🛑 user-interrupt, 🔄 self-correct, ⏸️ unnecessary-pause, 📦 script-candidate)
- A short description of what happened
- The proposed root cause in the skill (or consolidation rationale for script candidates)

Ask the user to confirm the catalog is accurate before proceeding.

## Step 3: Design Fixes

For each cataloged correction, propose a specific change to the SKILL.md:

### Fix Patterns

**For user interruptions:**
- If the skill instructed the wrong action → rewrite that instruction
- If the skill used an incorrect command/API → fix the command with a comment explaining why
- If the skill had the user choose when a default was obvious → remove the choice or set a sensible default with fallback

**For self-corrections:**
- If the agent found a better approach → document that approach as the primary path
- If the agent hit a known API limitation → add a note/comment in the skill (e.g., "Note: GitHub GraphQL filterBy does NOT support X")
- If the error was due to missing context → add the necessary context or a preliminary check

**For unnecessary pauses:**
- If the skill explicitly asks for confirmation that isn't needed → remove the AskUserQuestion step or make it conditional
- If the skill splits into too many micro-steps → consolidate steps that can flow naturally
- If the agent was overly cautious about a safe operation → add guidance that this step can proceed without confirmation

### Script Extraction

For each 📦 script-candidate identified in Step 2:

1. Write the script to `${CLAUDE_SKILL_DIR}/scripts/` (e.g., `scripts/find-release-issues.sh`). The script should:
   - Accept parameters via arguments or stdin
   - Output structured data (JSON preferred) to stdout
   - Include a brief usage comment at the top
   - Be executable (`chmod +x`)

2. Update the SKILL.md to call the script instead of inline Bash:
   ```bash
   # Before (inline in SKILL.md):
   gh api graphql -f query='...' --jq '...'

   # After (script reference in SKILL.md):
   ${CLAUDE_SKILL_DIR}/scripts/find-release-issues.sh
   ```

3. Update `allowed-tools` in the frontmatter if new Bash patterns are needed.

### Structural Improvements

Beyond individual fixes, look for patterns:

- **Multiple corrections in the same phase** → the phase's instructions may need restructuring
- **Corrections that suggest a missing step** → add the step
- **Corrections caused by stale information** → update the outdated content
- **Skill instructions that contradict CLAUDE.md or project conventions** → align with project standards

Present each proposed fix to the user with:
1. The correction it addresses (reference the catalog number)
2. The specific SKILL.md section to change
3. A before/after preview of the change

If any fix suggests splitting or creating a new skill, discuss with the user first before proceeding.

## Step 4: Apply Fixes

After the user approves the fixes (or a subset):

1. Read the target SKILL.md using the Read tool
2. Apply each approved fix using the Edit tool
3. If metadata changes are needed (description, allowed-tools), update those too
4. Show a summary of all changes made

Do NOT create a new file or rewrite the entire SKILL.md — use targeted edits to preserve the existing structure and only change what needs to change.

## Step 5: Verify

After applying fixes, do a quick sanity check:

- Re-read the modified SKILL.md
- Verify the changes are coherent with the rest of the skill's instructions
- Check that no instructions contradict each other
- Confirm the YAML frontmatter is still valid

Present a final summary to the user:
- Number of corrections addressed
- Sections modified
- Any remaining issues that couldn't be fixed through skill changes alone (e.g., tool limitations, external API behavior)
