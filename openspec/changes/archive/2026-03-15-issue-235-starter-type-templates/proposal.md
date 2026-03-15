## Why

`tmd init` creates an empty vault with no type schemas. New users must manually create all types from scratch, which is a high-friction onboarding experience. Offering opt-in starter type templates during initialization lets users start with useful defaults (idea, note, book) while keeping full ownership of their schemas.

## What Changes

- Add embedded starter type YAML files (`core/starters/*.yaml`) for `idea`, `note`, and `book`
- Add a Bubble Tea interactive checkbox selector to `tmd init` for choosing starter types (all selected by default)
- Add `--no-starters` flag to skip starter type selection in non-interactive/CI scenarios
- Selected starter types are copied to `.typemd/types/*.yaml` during vault initialization

## Capabilities

### New Capabilities

- `starter-type-templates`: Embedded YAML type definitions offered during `tmd init`, with interactive Bubble Tea multi-select and `--no-starters` flag

### Modified Capabilities

_(none — vault initialization is not yet covered by an existing spec)_

## Impact

- **core/**: New `starters/` directory with embedded YAML files; `Vault.Init()` gains an option to write starter types
- **cmd/**: `init` command gains `--no-starters` flag and Bubble Tea selector integration
- **Dependencies**: Bubble Tea is already a dependency (used by `tui/`)
- **Existing behavior**: Without `--no-starters`, interactive mode changes from immediate completion to showing a type selector. Non-interactive pipelines should use `--no-starters` to preserve current behavior.
