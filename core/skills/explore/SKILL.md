---
name: explore
description: Explore existing markdown files and suggest typemd type schemas and properties. Use when the user wants to analyze their files before importing, plan their vault structure, or get type/property recommendations. Read-only — never creates, modifies, or deletes any files.
---

# Explore

Analyze existing markdown files and suggest type schemas and properties for a typemd vault. This is a **read-only exploration** — never create, modify, or delete any files.

## Prerequisites

If the `vault-guide` plugin is installed, load it first — it provides the complete typemd reference (type schema format, property types, and vault structure) that informs your recommendations.

## Constraints

**This skill is strictly read-only.** You MUST NOT:

- Create any files or directories
- Modify any existing files
- Run any commands that write to disk (no `tmd` commands, no `mkdir`, no file writes)
- Offer to "go ahead and create" anything

Your output is **analysis and recommendations only**. If the user wants to act on your suggestions, they should use the `importer` skill or do it manually.

## Exploration Process

### 1. Identify the Target

Ask the user which files or directories to analyze. Common scenarios:

- A specific directory of markdown files (e.g., `~/notes/`, `./docs/`)
- A set of files the user points to
- The current directory

### 2. Survey the Files

Read the target files and gather:

- **File count and directory structure** — how many files, how they're organized into folders
- **Frontmatter patterns** — do files already have YAML frontmatter? What keys appear?
- **Content patterns** — recurring structures (headings, lists, metadata blocks, link styles)
- **Naming conventions** — how files are named (dates, slugs, categories)
- **Size distribution** — short notes vs. long-form content

### 3. Analyze the Existing Vault (if present)

If a `.typemd/` directory exists, read the current type schemas to understand:

- What types already exist
- What properties are defined
- What shared properties are available in `.typemd/properties.yaml`
- How the existing types relate to the files being explored

### 4. Classify Content into Type Groups

Group the files by content type. For each proposed group:

- **Suggested type name** — a concise, singular noun (e.g., `recipe`, `meeting-note`, `person`)
- **File count** — how many files fall into this group
- **Example files** — 2–3 representative filenames
- **Rationale** — why these files belong together

Use these classification signals:

- Folder structure (e.g., `recipes/` → `recipe` type)
- Frontmatter keys (e.g., `author`, `isbn` → `book` type)
- Content structure (e.g., "Ingredients" + "Steps" headings → `recipe` type)
- Naming patterns (e.g., `2024-01-15-standup.md` → `meeting-note` type)

### 5. Suggest Type Schemas

For each proposed type, recommend a schema with:

- **emoji** — a fitting emoji for the type
- **plural** — the plural form
- **unique** — whether names should be unique (yes for entities like books/people, no for instances like notes/ideas)
- **properties** — extracted from content patterns, using appropriate typemd property types:
  - `string` — free text
  - `number` — numeric values
  - `date` — date only (YYYY-MM-DD)
  - `datetime` — date + time (ISO 8601)
  - `url` — URLs
  - `checkbox` — boolean
  - `select` — single choice from options (include suggested `options` list)
  - `multi_select` — multiple choices (include suggested `options` list)
  - `relation` — link to another type (specify `target`)

Present schemas in YAML format matching typemd's schema format:

```yaml
emoji: 📖
plural: books
unique: true
properties:
  - name: author
    type: string
  - name: genre
    type: select
    options:
      - fiction
      - non-fiction
      - technical
  - name: rating
    type: number
```

### 6. Suggest Relations

Identify potential connections between proposed types:

- Explicit references (e.g., a note mentioning a person's name)
- Structural patterns (e.g., meeting notes always listing attendees)
- Folder co-location (e.g., project files referencing shared resources)

Present relations as property additions to the relevant type schemas.

### 7. Suggest Shared Properties

If multiple types would share the same property definitions, suggest extracting them into `.typemd/properties.yaml`:

```yaml
source:
  type: url
  description: Original source URL

priority:
  type: select
  options:
    - low
    - medium
    - high
```

### 8. Present the Summary

Conclude with a structured overview:

1. **File survey** — total files, directory structure, notable patterns
2. **Proposed types** — table of type name, emoji, file count, key properties
3. **Relations** — how types connect to each other
4. **Shared properties** — properties to define once and reuse
5. **Recommendations** — any additional observations (e.g., "consider splitting the `notes/` folder into `meeting-note` and `idea` types", "some files have no clear type — `page` would work as a catch-all")

## Output Format

Use clear markdown formatting with:

- Tables for type overviews
- YAML code blocks for schema suggestions
- Bullet lists for observations
- File paths for examples

Keep the tone practical and actionable — the user should be able to take your suggestions and directly create type schemas or use the `importer` skill to start importing.
