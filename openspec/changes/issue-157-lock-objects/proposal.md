## Why

All objects are always editable, with no way to protect specific objects from accidental modification. Users need the ability to lock individual objects so that CLI and TUI editing operations refuse to modify them.

## What Changes

- Add `locked` as a stored system property (boolean, default false/absent) in object YAML frontmatter
- Register `locked` in the system property registry as an immutable system property
- Add `tmd object lock <id>` and `tmd object unlock <id>` CLI commands
- Guard write operations in `ObjectService`: `SaveObject()` checks `locked` field and returns an error if true
- TUI shows lock icon (🔒) next to locked objects and prevents entering edit mode with a toast message
- Bulk operations encountering a locked object fail immediately with an error identifying the locked object

## Capabilities

### New Capabilities
- `object-lock`: Lock/unlock individual objects to prevent accidental editing, including CLI commands, core service guards, and TUI visual indicators

### Modified Capabilities
- `system-properties`: Adding `locked` as a new immutable system property to the registry
- `inline-property-editing`: Locked objects must refuse inline property editing in the TUI

## Impact

- **core/system_property.go** — register `locked` in SystemProperty registry
- **core/object_service.go** — guard `SaveObject()` and related write operations against locked objects
- **core/object.go** — add `IsLocked()` helper method on Object
- **cmd/** — new `tmd object lock` / `tmd object unlock` commands
- **tui/** — lock icon display, edit prevention, toast notification for locked objects
- **core/projector.go** — index `locked` field for query filtering
