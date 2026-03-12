# markdown-import

Convert existing markdown files into [typemd](https://github.com/typemd/typemd) objects with automatic type detection, frontmatter generation, and relation discovery.

## What It Does

This skill helps you migrate existing markdown notes into your typemd vault. It:

- **Detects the type** — analyzes your markdown content and suggests the best matching type from your vault's schemas
- **Generates frontmatter** — extracts metadata (title, author, dates, tags) and maps them to your type's properties
- **Discovers relations** — finds mentions of existing objects in your vault and suggests wiki-links
- **Supports batch import** — guides you through importing multiple files with consistent type assignment

## Prerequisites

- A typemd vault initialized with `tmd init`
- Type schemas defined in `.typemd/types/*.yaml` (recommended for best results)
- [Claude Code](https://claude.com/code) installed

## Installation

```
/plugin marketplace add typemd/marketplace
/plugin install markdown-import@typemd-marketplace
```

## Usage

### Import a single file

Open Claude Code in your vault directory and ask:

> Import my-notes/book-review.md into the vault

Claude will analyze the file, suggest a type, generate frontmatter, and create the object after your confirmation.

### Import a directory

> Import all markdown files from my-notes/ into the vault

Claude will list the files, suggest types for each, and walk you through importing them one at a time.

### With specific type

> Import meeting-notes.md as a note type object

Claude will use the specified type and map properties accordingly.

## How It Works

1. Reads your vault's type schemas to understand available types and properties
2. Analyzes the markdown file content
3. Suggests the best matching type (or asks you to choose)
4. Extracts metadata and generates YAML frontmatter
5. Identifies potential relations to existing objects
6. Shows you the result and waits for confirmation before writing

The original markdown file is never modified or deleted.
