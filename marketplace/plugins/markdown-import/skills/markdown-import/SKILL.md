---
name: markdown-import
description: Convert existing markdown files into typemd objects. Use when the user wants to import, migrate, or convert markdown files into their typemd vault.
---

# Markdown Import

Help the user convert existing markdown files into typemd objects.

## Prerequisites

If the `vault-guide` plugin is installed, load it first — it provides the complete typemd reference (CLI commands, type schema format, property types, views, and vault structure) that informs every decision in this workflow.

## Before You Start

1. **Read the vault's type schemas** from `.typemd/types/*/schema.yaml` (or legacy `.typemd/types/*.yaml`) to understand available types and their properties
2. **Check `.typemd/properties.yaml`** for shared property definitions (referenced via `use:` in type schemas)
3. **List existing objects** with `tmd object list` to understand what's already in the vault (for relation discovery)
4. **Analyze existing objects to learn vault conventions** (see below)

If no type schemas exist, inform the user and suggest creating basic types first, or proceed with system properties only.

### Learn from Existing Objects

Before classifying new content, study how the vault owner has been organizing their knowledge:

1. **Sample 3–5 objects per type** — read their frontmatter to understand which properties are actually used, what values are common, and how detailed the descriptions tend to be
2. **Identify classification patterns**:
   - Which types have `unique: true`? These represent distinct entities (books, people) vs. instances (notes, ideas)
   - Which properties use `relation`? These reveal how the vault connects concepts
   - Which properties use `select`/`multi_select`? These reveal the vault's taxonomy and categories
   - Which types have `tags` populated? Study the tag vocabulary to reuse existing tags
3. **Note the naming style** — are names formal ("Clean Code: A Handbook of Agile Software Craftsmanship") or casual ("clean code notes")? Follow the established convention
4. **Check for templates** in `templates/` — they reveal the vault owner's expected structure for each type

Use these observations to make classification decisions consistent with the vault's existing style, not just based on content analysis alone.

## Import Process

For each markdown file the user wants to import:

### 1. Analyze and Classify the Content

Read the file and determine:
- What **type** best fits this content
- What **properties** can be extracted from the content
- What **relations** might exist to other objects in the vault

**Type classification priority** (try in this order):

1. **Match existing type schemas** — if the content's structure (properties, vocabulary) aligns with a defined type, use it. Compare against the patterns you learned from existing objects
2. **Match by property overlap** — count how many extractable properties match each type's schema. The type with the highest match wins
3. **Use built-in types as fallback** — `page` (📄) for general content, `tag` (🏷️) only for taxonomy/label concepts
4. **Suggest a new type** — if the content clearly represents a recurring entity class not yet in the vault (e.g., multiple recipe files but no `recipe` type), suggest creating a new type schema before importing

If the content doesn't clearly map to an existing type, ask the user which type to use or suggest creating a new type schema.

### 2. Generate the typemd Object

Use `tmd object create <type> <name>` to create the object file. This generates the proper ULID filename at `objects/<type>/<slug>-<ulid>.md` (ULID: 26-character Crockford Base32, lowercase). If the type has object templates, use `-t <template>` to select one.

Then edit the created file to populate:

**Frontmatter** (YAML between `---` delimiters):
- System properties first, in this order:
  - `name`: display name derived from the content (e.g., book title, person name)
  - `description`: optional one-line summary extracted from the content
  - `created_at`: ISO 8601 timestamp (use the original file's date if available, otherwise current time)
  - `updated_at`: current ISO 8601 timestamp
  - `tags`: array of tag object IDs (e.g., `tag/programming-01kk...`)
- Type-specific properties after system properties, matching the type schema's property definitions and order. Available property types: `string`, `number`, `date`, `datetime`, `url`, `checkbox`, `select`, `multi_select`, `relation`

**Body**: the markdown content, cleaned up and adapted. Preserve the substance of the original content.

### 3. Discover Relations

Look for mentions of existing objects in the vault:
- Names that match existing objects → suggest converting to wiki-links: `[[type/slug-ulid]]`
- Content that implies a relation defined in the type schema (e.g., a book mentioning an author) → populate the relation property in frontmatter

### 4. Confirm Before Writing

Always show the user the generated object before writing it. Present:
- The target path (`objects/<type>/<slug>-<ulid>.md`)
- The full frontmatter
- A preview of the body (first ~10 lines if long)
- Any suggested wiki-links or relations

Wait for the user to confirm or request changes before creating the file.

## Batch Import

When the user has multiple files to import:

1. List all files and suggest types for each
2. Let the user confirm or adjust the type assignments
3. Import files one at a time, maintaining consistency
4. After all imports, summarize what was created and suggest any cross-object relations that could be added

## Important Rules

- **Never modify or delete the original markdown file** — only create new typemd object files
- **Never guess ULIDs** — always use `tmd object create <type> <name>` to generate the object with a proper ULID filename
- **Always respect the type schema** — only include properties defined in the schema
- **System properties are managed by typemd** — set them during import but know that typemd will manage `updated_at` going forward
