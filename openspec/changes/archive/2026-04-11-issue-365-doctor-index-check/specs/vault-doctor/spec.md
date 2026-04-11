## MODIFIED Requirements

### Requirement: Doctor checks index-disk synchronization
The system SHALL check whether the SQLite index is in sync with files on disk, and auto-rebuild if out of sync. The check MUST be reported as a category named "Index" and MUST appear between the "Files" and "Orphans" categories in the doctor report. Auto-fixed rebuilds MUST be recorded as `AutoFixed` on the category (not as issues) so that they do not cause a non-zero exit code.

#### Scenario: Index is in sync
- **WHEN** the SQLite index matches the current state of files on disk
- **THEN** doctor shows the "Index" category as passing with 0 issues and 0 auto-fixed

#### Scenario: Index is out of sync — new file on disk
- **WHEN** an object file exists on disk that is not yet recorded in the index
- **THEN** doctor automatically reconciles and projects the new object into the index
- **AND** reports the fix as auto-fixed under the "Index" category
- **AND** the "Index" category has 0 issues

#### Scenario: Doctor report includes Index category
- **WHEN** doctor runs on any vault
- **THEN** the report contains exactly 8 categories including "Index"
