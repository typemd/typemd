## Why

Object templates (`templates/<type>/<name>.md`) are supported during `tmd object create`, but there are no dedicated CLI commands for managing templates themselves. Users must manually navigate to `templates/<type>/` and create/edit/delete `.md` files. Adding `tmd template` commands closes this gap and provides a consistent CLI experience for template lifecycle management.

## What Changes

- Add `tmd template` command group with four subcommands:
  - `tmd template list [type]` — list available templates, optionally filtered by type
  - `tmd template show <type>/<name>` — display a template's frontmatter and body
  - `tmd template create <type>/<name>` — create a new template file and open in `$EDITOR`
  - `tmd template delete <type>/<name>` — delete a template file with confirmation prompt

## Capabilities

### New Capabilities

- `cli-template-management`: CLI commands for listing, viewing, creating, and deleting object templates

### Modified Capabilities

_(none — existing `object-templates` spec covers core Vault methods; this change only adds CLI surface)_

## Impact

- **Code**: New `cmd/template.go` file with Cobra command group, registered on `rootCmd`
- **Dependencies**: Delegates to existing `Vault.ListTemplates()`, `Vault.LoadTemplate()`, `Vault.SaveTemplate()`, `Vault.DeleteTemplate()` — no new core methods needed
- **APIs**: No API changes (web API already has template endpoints)
- **Systems**: CLI only — no TUI, MCP, or web changes
