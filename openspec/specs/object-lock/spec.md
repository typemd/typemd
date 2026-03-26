# object-lock Specification

## Purpose
Lock/unlock individual objects to prevent accidental editing, including CLI commands, core service guards, and TUI visual indicators.
## Requirements
### Requirement: Locked objects cannot be modified through standard operations
The system SHALL prevent modifications to objects that have `locked: true` in their frontmatter. When a write operation (SaveObject, SetProperty, SetPropertyMultiple) is attempted on a locked object, the system SHALL return an `ErrObjectLocked` error without modifying the object.

#### Scenario: SaveObject on locked object returns error
- **WHEN** SaveObject is called on an object with `locked: true`
- **THEN** the system SHALL return an `ErrObjectLocked` error
- **AND** the object file SHALL not be modified

#### Scenario: SetProperty on locked object returns error
- **WHEN** SetProperty is called on an object with `locked: true`
- **THEN** the system SHALL return an `ErrObjectLocked` error
- **AND** the object file SHALL not be modified

#### Scenario: SetPropertyMultiple on locked object returns error
- **WHEN** SetPropertyMultiple is called on an object with `locked: true`
- **THEN** the system SHALL return an `ErrObjectLocked` error
- **AND** the object file SHALL not be modified

#### Scenario: Unlocked object can be modified normally
- **WHEN** SaveObject is called on an object without `locked` or with `locked: false`
- **THEN** the save SHALL proceed normally

### Requirement: Lock and unlock CLI commands
The system SHALL provide `tmd object lock <id>` and `tmd object unlock <id>` commands to toggle the lock state of an object.

#### Scenario: Lock an unlocked object
- **WHEN** `tmd object lock book/clean-code-01abc` is executed
- **AND** the object is not locked
- **THEN** the object's frontmatter SHALL contain `locked: true`
- **AND** the command SHALL print a confirmation message

#### Scenario: Lock an already locked object
- **WHEN** `tmd object lock book/clean-code-01abc` is executed
- **AND** the object already has `locked: true`
- **THEN** the command SHALL print a message indicating the object is already locked
- **AND** no modification SHALL occur

#### Scenario: Unlock a locked object
- **WHEN** `tmd object unlock book/clean-code-01abc` is executed
- **AND** the object has `locked: true`
- **THEN** the `locked` property SHALL be removed from the object's frontmatter
- **AND** the command SHALL print a confirmation message

#### Scenario: Unlock a non-locked object
- **WHEN** `tmd object unlock book/clean-code-01abc` is executed
- **AND** the object does not have `locked: true`
- **THEN** the command SHALL print a message indicating the object is not locked
- **AND** no modification SHALL occur

### Requirement: Object entity provides lock state accessor
The Object struct SHALL have an `IsLocked()` method that returns true if the object has `locked: true` in its properties.

#### Scenario: IsLocked returns true for locked object
- **WHEN** an object has `locked: true` in its properties
- **THEN** `IsLocked()` SHALL return true

#### Scenario: IsLocked returns false for unlocked object
- **WHEN** an object does not have a `locked` property
- **THEN** `IsLocked()` SHALL return false

#### Scenario: IsLocked returns false when locked is false
- **WHEN** an object has `locked: false` in its properties
- **THEN** `IsLocked()` SHALL return false

### Requirement: Locked property stored in frontmatter
The `locked` property SHALL be stored as a boolean in the object's YAML frontmatter. When `locked` is false or absent, it SHALL be omitted from the frontmatter output (omitempty behavior).

#### Scenario: Locked true appears in frontmatter
- **WHEN** an object has `locked: true`
- **THEN** the frontmatter output SHALL include `locked: true`

#### Scenario: Locked false is omitted from frontmatter
- **WHEN** an object has `locked: false`
- **THEN** the frontmatter output SHALL NOT include a `locked` key

#### Scenario: Absent locked is omitted from frontmatter
- **WHEN** an object does not have a `locked` property
- **THEN** the frontmatter output SHALL NOT include a `locked` key

### Requirement: TUI shows lock indicator for locked objects
The TUI SHALL display a 🔒 badge right-aligned in the title panel when a locked object is selected.

#### Scenario: Lock badge in title panel
- **WHEN** the title panel displays an object with `locked: true`
- **THEN** a 🔒 badge SHALL appear right-aligned in the title panel

#### Scenario: No lock badge for unlocked objects
- **WHEN** the title panel displays an object without `locked: true`
- **THEN** no lock badge SHALL be displayed

### Requirement: TUI toggle lock via keybinding
The TUI SHALL allow toggling the lock state of the selected object by pressing `L` (uppercase).

#### Scenario: Toggle lock on unlocked object
- **WHEN** the user presses `L` on an unlocked object in object detail view
- **THEN** the object SHALL become locked
- **AND** a toast notification SHALL display "Locked"

#### Scenario: Toggle lock on locked object
- **WHEN** the user presses `L` on a locked object in object detail view
- **THEN** the object SHALL become unlocked
- **AND** a toast notification SHALL display "Unlocked"

### Requirement: TUI prevents editing locked objects
The TUI SHALL prevent entering property edit mode for locked objects. When a user attempts to edit a locked object, a toast notification SHALL inform them that the object is locked.

#### Scenario: Tab to properties panel on locked object
- **WHEN** the user presses Tab to focus the properties panel on a locked object
- **THEN** the properties panel SHALL display properties in read-only mode
- **AND** no cursor indicator SHALL appear

#### Scenario: Attempt to edit locked object property
- **WHEN** the user attempts to enter edit mode on a locked object
- **THEN** a toast notification SHALL display "Object is locked. Unlock to edit."
- **AND** the property editor SHALL NOT activate

### Requirement: Projector sync bypasses object lock
The Projector sync operation SHALL bypass the lock check when writing back resolved names, expanded wiki-links, or other system-level sync operations. Lock only guards user-initiated edits through ObjectService.

#### Scenario: Sync writes back to locked object
- **WHEN** vault sync resolves a shorthand wiki-link in a locked object
- **THEN** the resolved wiki-link SHALL be written back to the object file
- **AND** no lock error SHALL occur

### Requirement: SetLocked dedicated method
ObjectService SHALL provide a `SetLocked(id string, locked bool)` method that bypasses the standard lock guard to toggle the lock state.

#### Scenario: SetLocked locks an object
- **WHEN** `SetLocked("book/clean-code-01abc", true)` is called
- **THEN** the object's `locked` property SHALL be set to true
- **AND** the object SHALL be saved

#### Scenario: SetLocked unlocks an object
- **WHEN** `SetLocked("book/clean-code-01abc", false)` is called
- **THEN** the object's `locked` property SHALL be removed
- **AND** the object SHALL be saved

### Requirement: LinkObjects and UnlinkObjects respect lock
The system SHALL prevent linking or unlinking relations on locked objects.

#### Scenario: LinkObjects on locked source object
- **WHEN** LinkObjects is called with a locked source object
- **THEN** the system SHALL return an `ErrObjectLocked` error

#### Scenario: UnlinkObjects on locked source object
- **WHEN** UnlinkObjects is called with a locked source object
- **THEN** the system SHALL return an `ErrObjectLocked` error
