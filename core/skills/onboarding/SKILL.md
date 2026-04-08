---
name: onboarding
description: Import existing markdown collections into a structured typemd vault. Use when the user wants to onboard, import, or migrate content from external directories into their vault.
---

# Onboarding

Help the user import existing markdown collections into a structured typemd vault. This is a guided four-phase workflow: scan, plan, execute, verify.

## Prerequisites

If the `typemd` plugin is installed, load it first — it provides the complete typemd reference (CLI commands, type schema format, property types, views, and vault structure) that informs every decision in this workflow.

## Before You Start

1. **Read the vault's type schemas** from `types/*/schema.yaml` to understand available types and their properties
2. **Check `properties/*.yaml`** for shared property definitions (referenced via `use:` in type schemas)
3. **List existing objects** with `tmd object list` to understand what's already in the vault (for relation discovery and conflict detection)
4. **Analyze existing objects to learn vault conventions** — sample 3–5 objects per type to understand naming style, property usage, and tag vocabulary

If no type schemas exist, that's fine — the plan phase will suggest new types based on scan results.

## Phase 1: Scan & Analyze

Accept source directories or files as input from the user. Run the scan command to get structured analysis:

```bash
tmd import scan <paths...>
```

Review the scan output, which includes:

- **File count** — total markdown files found
- **Directory structure** — how files are organized into folders
- **Frontmatter patterns** — which YAML keys appear, how frequently, with sample values
- **Existing vault types** — types already defined in the vault, with their properties
- **Files without frontmatter** — count of plain markdown files with no YAML metadata

Present a summary to the user:

1. Total files and directory breakdown
2. Most common frontmatter keys and their frequency
3. Files with vs. without frontmatter
4. Current vault types that might match the scanned content

## Phase 2: Conversion Plan

Analyze the scan results and classify each source file into a target type.

### Classification Strategy

For each source file, determine:

- **Target type** — which vault type best fits this content
- **Name** — the display name for the object
- **Property mapping** — how frontmatter keys map to type properties
- **Dependencies** — cross-file references (e.g., a book referencing an author)

**Type classification priority** (try in this order):

1. **Match existing type schemas** — if the content's frontmatter keys or structure align with a defined type, use it
2. **Match by property overlap** — count how many extractable properties match each type's schema; the type with the highest match wins
3. **Group by folder** — files in a `recipes/` folder likely belong to a `recipe` type
4. **Group by frontmatter pattern** — files sharing the same set of keys likely belong to the same type
5. **Use built-in types as fallback** — `page` (📄) for general content, `tag` (🏷️) only for taxonomy/label concepts
6. **Suggest a new type** — if a clear recurring pattern doesn't match any existing type, propose creating one

### Property Type Inference

When frontmatter values need property types, use these heuristics:

| Value pattern | Inferred type |
|---|---|
| `true` / `false` | `checkbox` |
| ISO 8601 date (`2024-01-15`) | `date` |
| ISO 8601 datetime (`2024-01-15T10:30:00Z`) | `datetime` |
| Integer or float (`42`, `3.14`) | `number` |
| URL (`https://...`, `http://...`) | `url` |
| Array of strings | `multi_select` (if values repeat across files) or `tags` relation |
| Single string from a small set of values | `select` (collect unique values as options) |
| Free text string | `string` |
| Reference to another file or name | `relation` (specify `target` type) |

### Relation Discovery

Look for cross-file references:

- **Frontmatter references** — a property whose value matches another file's name or title (e.g., `author: "John Doe"` where `john-doe.md` exists)
- **Wiki-link patterns** — `[[...]]` syntax in markdown body
- **Shared tags/categories** — frontmatter arrays that could become `tag` objects
- **Folder co-location** — files organized together that reference each other

### Building the Plan

Group files by type and construct the conversion plan:

1. **New types** — types that need to be created with suggested schemas (emoji, plural, properties)
2. **Object mappings** — for each file: source path, target type, name, property values, conflict status
3. **Import order** — tags first, then independent objects, then objects with dependencies

Run `tmd import plan` or construct the `ImportPlan` JSON directly. The plan JSON structure:

```json
{
  "types": [
    {"name": "book", "emoji": "📖", "plural": "books", "properties": [...]}
  ],
  "objects": [
    {
      "source_path": "books/clean-code.md",
      "type_name": "book",
      "name": "Clean Code",
      "properties": {"author": "Robert C. Martin", "year": 2008},
      "conflict": "none",
      "depends_on": []
    }
  ]
}
```

### Conflict Resolution

For each object, set the conflict field:

- `"none"` — no existing object with the same name and type; proceed normally
- `"skip"` — an object with the same name and type already exists; skip this import
- `"overwrite"` — an object exists but the user wants to replace it

**Default to `"skip"` for conflicts.** Only use `"overwrite"` when the user explicitly requests it.

### User Review

Present the complete plan to the user before proceeding:

1. **New types to create** — with full schema details
2. **Objects to import** — grouped by type, showing name and key properties
3. **Conflicts detected** — which files will be skipped and why
4. **Import order** — the sequence in which objects will be created

Wait for the user to confirm, adjust type assignments, or modify property mappings before proceeding to execution.

## Phase 3: Execute

After user approval, execute the import plan:

1. Write the `ImportPlan` JSON to a temporary file
2. Run the execute command:

```bash
tmd import execute <plan-file>
```

3. Monitor the progress output — the command creates types first, then objects in dependency order, and finally runs reconciliation to resolve wiki-links

If the CLI command is not available, fall back to manual execution:

1. Create type schemas with `tmd type create <name>` and populate `types/<name>/schema.yaml`
2. Create objects with `tmd object create <type> <name>` in dependency order
3. Set properties and body content on each created object
4. Run `tmd fix` to reconcile relations and wiki-links

## Phase 4: Verify

Review the import results:

1. **Check the ImportReport** — types created, objects created/skipped/failed
2. **Review failures** — if any objects failed to import, diagnose and suggest fixes
3. **Check unresolved references** — wiki-links or relations that couldn't be resolved
4. **Suggest follow-up actions**:
   - Create objects for unresolved references
   - Add missing relations between imported objects
   - Review and refine type schemas based on actual data
   - Run `tmd fix` again if manual adjustments were made

### Post-Import Checklist

- [ ] All expected objects appear in `tmd object list`
- [ ] Type schemas have correct properties and property types
- [ ] Relations between objects are properly linked
- [ ] Tags are created and assigned
- [ ] No unresolved wiki-links remain

## Incremental Additions

When re-running on an existing vault with previously imported content:

1. **Scan only new paths** — or re-scan and let conflict detection handle duplicates
2. **Conflict detection is automatic** — objects with matching name and type are flagged as `"skip"` by default
3. **New types merge safely** — if a type already exists, its schema is preserved; only new types are created
4. **Re-run reconciliation** — after incremental imports, run `tmd fix` to update all cross-references

## Important Rules

- **Never modify or delete the original source files** — only create new typemd objects
- **Never guess ULIDs** — always use `tmd object create <type> <name>` or `tmd import execute` to generate proper ULID filenames
- **Always respect existing type schemas** — only include properties defined in the schema
- **Present the plan before executing** — never auto-execute without user confirmation
- **Default to skip for conflicts** — don't overwrite existing objects unless explicitly requested
