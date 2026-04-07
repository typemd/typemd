## Context

typemd has two AI-driven skills for getting content into a vault: `explore` (read-only analysis, suggests schemas) and `importer` (single-file conversion, requires schemas to exist). The gap between them — creating type schemas, ordering imports by dependency, resolving cross-object relations — falls entirely on the user.

The onboarding skill bridges this gap with a four-phase workflow (scan → plan → execute → verify) backed by CLI commands that provide structured data for AI orchestration.

### Current state

- `tmd object create <type> <name>` creates a single object via `Vault.NewObject()` → `ObjectService.Create()`
- `Vault.SaveType()` persists type schemas to `types/<name>/schema.yaml`
- `core.ListSkills()` / `core.GetSkillWithOverride()` provide the embedded skill system
- CLI commands in `cmd/` use Cobra with a `rootCmd` → subcommand pattern
- No batch creation, no import-specific CLI commands exist today

## Goals / Non-Goals

**Goals:**

- Provide CLI commands (`tmd import scan/plan/execute`) that output structured JSON for AI consumption
- Keep the CLI layer thin — scan produces analysis, plan produces a conversion plan, execute runs the plan
- Reuse existing `Vault.NewObject()` and `Vault.SaveType()` for object/type creation
- Support dependency-ordered import (tags and referenced objects first)
- Support incremental re-import by detecting existing objects
- Add embedded + marketplace onboarding skills

**Non-Goals:**

- No automatic ingestion without user review (plan must be approved)
- Not replacing `explore` or `importer` skills (they coexist for simple cases)
- No core architecture changes — this is additive
- No MCP write endpoints (that's #377, separate)
- No source provenance tracking via a `source` type (deferred — can be added later without breaking changes)

## Decisions

### 1. CLI command structure: `tmd import` command group

**Decision:** Add a new `tmd import` command group with `scan`, `plan`, and `execute` subcommands.

**Rationale:** Import is a distinct workflow from `tmd object create`. A command group keeps related functionality together and allows future subcommands (e.g., `tmd import status`). The existing `create` command stays for single-object creation.

**Alternative considered:** Adding flags to `tmd object create` (e.g., `--batch`, `--from-dir`). Rejected because the scan/plan/execute workflow is fundamentally different from single-object creation.

### 2. Scan output as structured JSON

**Decision:** `tmd import scan` outputs a `ScanResult` JSON struct containing file inventory, frontmatter analysis, directory structure, and content classification hints.

```go
type ScanResult struct {
    Sources     []SourceInfo     `json:"sources"`
    FileCount   int              `json:"file_count"`
    Directories []DirInfo        `json:"directories"`
    Patterns    FrontmatterStats `json:"patterns"`
}
```

**Rationale:** The CLI provides raw data; the AI skill interprets it (suggests types, maps properties). This keeps the CLI deterministic and testable while leveraging AI for the judgment-heavy classification step.

### 3. Plan file as intermediate artifact

**Decision:** `tmd import plan` produces a JSON plan file that `tmd import execute` consumes. The plan is the contract between the AI's decisions and the CLI's execution.

```go
type ImportPlan struct {
    Types   []TypePlan   `json:"types"`    // schemas to create/modify
    Objects []ObjectPlan `json:"objects"`   // files to import with type + property mapping
    Order   []string     `json:"order"`     // import order (object indices)
}
```

**Rationale:** Separating plan from execution gives the user a review checkpoint. The AI generates the plan, the user approves, the CLI executes deterministically. This also enables retries (re-execute the same plan) and debugging (inspect the plan file).

**Alternative considered:** Direct execution without a plan file. Rejected because the issue explicitly requires user review before execution.

### 4. AI does classification, CLI does execution

**Decision:** The `scan` command collects raw data (file stats, frontmatter keys, directory structure). The AI skill uses this data plus content reading to classify files into types and map properties. The `plan` command is AI-assisted (the skill generates the plan JSON). The `execute` command is purely mechanical.

**Rationale:** Type classification requires judgment (e.g., "these files look like book reviews") which AI excels at. File I/O, ULID generation, and schema validation are mechanical and belong in the CLI.

### 5. Dependency ordering via topological sort

**Decision:** The plan's `order` field lists objects in dependency order. Tags first, then objects without relations, then objects with relations to already-created objects.

**Rationale:** Wiki-links and relations need target objects to exist. A two-pass approach (create objects → resolve links) is simpler than sorting, but sorting avoids a separate reconciliation step for frontmatter relations. We still do a second pass for body wiki-links since those are handled by the reconciler.

### 6. Incremental import via existing object detection

**Decision:** When scanning, check existing vault objects. If a source file maps to an object that already exists (by name + type match), mark it as `skip` in the plan by default.

**Rationale:** Simpler than provenance tracking (which requires a `source` type and relations). Name + type matching covers the common case. Source provenance is deferred as a non-goal.

### 7. Conflict handling via plan flags

**Decision:** Each object in the plan has a `conflict` field: `skip` (default), `overwrite`, or `merge`. The skill sets these based on user preference (per-file or global).

**Rationale:** Conflict resolution is a user decision, not a technical one. The plan file captures the decision; execute honors it.

### 8. Embedded skill + marketplace skill pattern

**Decision:** Follow the existing two-tier pattern:
- `core/skills/onboarding/SKILL.md` — full instructions for the four-phase workflow
- `marketplace/plugins/typemd/skills/onboarding/SKILL.md` — discovery layer that loads vault-guide + calls `tmd instructions onboarding`

**Rationale:** Consistent with how `explore` and `importer` are structured.

## Risks / Trade-offs

**[Large source collections] →** `tmd import scan` reads frontmatter from every file. For 1000+ files this could be slow. Mitigation: scan only reads frontmatter (not full content), which is fast. AI content analysis (for classification) is batched by the skill, not the CLI.

**[Schema evolution during import] →** If the plan becomes stale (e.g., a file doesn't fit the planned type), `execute` will fail for that file and continue. Mitigation: report failures in the verify phase; user can re-scan and re-plan for failed files.

**[Name collision on incremental import] →** Name + type matching for duplicate detection is imperfect (same content, different name = duplicate not detected). Mitigation: acceptable for v1. Source provenance tracking (#378) will improve this later.

**[Circular relations] →** If A references B and B references A, topological sort fails. Mitigation: break cycles arbitrarily (import A first, B second, then reconciler resolves A's link to B).
