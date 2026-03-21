## MODIFIED Requirements

### Requirement: Wikilinks and tags rebuilt after incremental sync

- **WHEN** `SyncFiles` completes incremental object sync
- **THEN** the Projector SHALL perform a full wikilink sync, a full tag relation sync, and relation sync for changed objects

#### Scenario: Wikilinks and tags rebuilt after incremental sync

- **WHEN** `SyncFiles` completes incremental object sync
- **THEN** the Projector SHALL perform a full wikilink sync and full tag relation sync

#### Scenario: Relations for changed objects rebuilt after incremental sync

- **WHEN** `SyncFiles` completes incremental object sync for a book whose `author` changed
- **THEN** the Projector SHALL delete existing relations for that book and rebuild them from the updated frontmatter
- **AND** relations for unchanged objects SHALL remain intact
