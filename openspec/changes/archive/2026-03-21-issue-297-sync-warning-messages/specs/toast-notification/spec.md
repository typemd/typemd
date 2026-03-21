## MODIFIED Requirements

### Requirement: Sync warnings display as toast

When `Projector.Sync()` produces unresolved references, the TUI SHALL display them as a Warning-level toast with group-based aggregation. The group key SHALL be derived from the `UnresolvedRelation.Reason` field to distinguish between "not found" and "ambiguous" references.

#### Scenario: Sync produces only not-found references

- **WHEN** a sync operation completes with 2 unresolved references with reason "not_found"
- **THEN** a Warning-level toast SHALL be shown
- **AND** the items SHALL use group key "not found"
- **AND** the toast SHALL display "⚠ 2 not found"

#### Scenario: Sync produces only ambiguous references

- **WHEN** a sync operation completes with 3 unresolved references with reason "ambiguous"
- **THEN** a Warning-level toast SHALL be shown
- **AND** the items SHALL use group key "ambiguous"
- **AND** the toast SHALL display "⚠ 3 ambiguous"

#### Scenario: Sync produces mixed not-found and ambiguous references

- **WHEN** a sync operation completes with 2 not-found and 1 ambiguous unresolved references
- **THEN** a Warning-level toast SHALL be shown
- **AND** the toast SHALL display two lines: "⚠ 2 not found" and "⚠ 1 ambiguous"

#### Scenario: Unknown reason falls back to generic group

- **WHEN** a sync operation completes with an unresolved reference whose reason is neither "not_found" nor "ambiguous"
- **THEN** the item SHALL use group key "unresolved"
