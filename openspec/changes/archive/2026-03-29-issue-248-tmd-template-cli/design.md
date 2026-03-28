## Context

Object templates live at `templates/<type>/<name>.md` and are already fully supported by the core layer — `Vault` provides `ListTemplates`, `LoadTemplate`, `SaveTemplate`, and `DeleteTemplate`. The `tmd object create` command uses these to apply templates during object creation. However, there are no CLI commands for managing templates directly; users must navigate the filesystem manually.

The existing CLI follows a consistent pattern: `tmd <noun> <verb>` with Cobra command groups (e.g., `tmd type list`, `tmd object show`). Each command group is a parent `*cobra.Command` in its own file, with subcommands registered via `init()`.

## Goals / Non-Goals

**Goals:**
- Add `tmd template list [type]`, `tmd template show <type/name>`, `tmd template create <type/name>`, and `tmd template delete <type/name>` commands
- Follow existing CLI patterns for output formatting, error handling, and editor integration
- Provide shell completion for type and template names

**Non-Goals:**
- Template validation against type schema (future work)
- Template inheritance or variables
- Modifying existing `tmd object create` template behavior
- TUI or MCP changes

## Decisions

### 1. Argument format: `<type/name>` for show/create/delete

Template identifiers use `type/name` slash format (e.g., `book/review`), consistent with how object IDs are displayed (`type/slug-ulid`). This is parsed by splitting on `/`. The `list` command takes an optional `[type]` positional arg to filter.

**Alternative considered:** Two positional args `<type> <name>` — rejected because the slash format is shorter and matches the existing `type/slug` pattern users already know.

### 2. `create` opens `$EDITOR`

`tmd template create` scaffolds an empty template file (with frontmatter delimiters) and opens it in `$EDITOR` (falling back to `$VISUAL`, then `vi`). This matches the standard Unix CLI convention for content creation.

**Alternative considered:** Interactive prompts for frontmatter properties — rejected as overly complex for a file that's just markdown with optional frontmatter.

### 3. `delete` requires `--force` or interactive confirmation

`tmd template delete` prompts for confirmation in interactive terminals. The `--force` / `-f` flag skips the prompt for scripting use. This is consistent with destructive CLI operations.

### 4. Single file: `cmd/template.go`

All four subcommands fit in a single file (~150 lines), following the `cmd/typecmd.go` + `cmd/type_list.go` pattern but consolidated since the commands are small.

### 5. `list` groups by type, shows template count

`tmd template list` (no args) shows all templates grouped by type. `tmd template list <type>` shows only that type's templates. Output includes the template name, one per line. With `--json`, outputs a structured JSON array.

## Risks / Trade-offs

- **[Risk] `$EDITOR` not set** → Fall back to `$VISUAL`, then `vi`. Print an error message if none are available.
- **[Risk] Template file exists on `create`** → Return a clear error message pointing to `tmd template show` to view existing content.
- **[Risk] Type doesn't exist on `create`** → Allow it (templates can be created before type schemas, since they're just files). Print a warning that the type doesn't exist yet.
