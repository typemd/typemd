## 1. Core: System Property Registration & Object Entity

- [x] 1.1 Write BDD scenarios for locked system property (registration, IsLocked, frontmatter ordering)
- [x] 1.2 Implement step definitions for locked system property scenarios
- [x] 1.3 Register `locked` in systemProperties registry and add `IsLocked()` method to Object
- [x] 1.4 Add unit tests for IsLocked edge cases (nil, false, non-boolean values)

## 2. Core: ObjectService Lock Guard

- [x] 2.1 Write BDD scenarios for lock guard (SaveObject, SetProperty, SetPropertyMultiple, LinkObjects, UnlinkObjects)
- [x] 2.2 Implement step definitions for lock guard scenarios
- [x] 2.3 Add ErrObjectLocked error and lock checks in ObjectService.Save, SetProperty, SetPropertyMultiple
- [x] 2.4 Add lock checks in Vault.LinkObjects and UnlinkObjects
- [x] 2.5 Add unit tests for lock guard edge cases

## 3. Core: SetLocked Dedicated Method

- [x] 3.1 Write BDD scenarios for SetLocked (lock, unlock, already locked, already unlocked)
- [x] 3.2 Implement step definitions for SetLocked scenarios
- [x] 3.3 Add SetLocked method to ObjectService and Vault facade
- [x] 3.4 Add unit tests for SetLocked edge cases

## 4. CLI: Lock and Unlock Commands

- [x] 4.1 Write BDD scenarios for `tmd object lock` and `tmd object unlock` commands
- [x] 4.2 Implement step definitions for lock/unlock CLI scenarios
- [x] 4.3 Add `tmd object lock` and `tmd object unlock` Cobra commands
- [x] 4.4 Add unit tests for CLI edge cases (missing args, invalid ID)

## 5. TUI: Lock Indicator and Edit Prevention

- [x] 5.1 Add 🔒 icon display next to locked objects in sidebar
- [x] 5.2 Prevent property editor activation for locked objects with toast notification
- [x] 5.3 Add unit tests for TUI lock behavior
