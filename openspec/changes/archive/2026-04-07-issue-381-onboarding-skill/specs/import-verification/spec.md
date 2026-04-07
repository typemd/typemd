## ADDED Requirements

### Requirement: Report import results
The system SHALL produce a summary report after plan execution with counts of created, failed, and skipped items.

#### Scenario: Successful import report
- **WHEN** a plan with 10 objects executes with 8 created, 1 skipped, 1 failed
- **THEN** the report shows `created: 8, skipped: 1, failed: 1` with details for each

#### Scenario: All objects created successfully
- **WHEN** a plan executes with no failures or skips
- **THEN** the report shows all items as created with no error details

### Requirement: List unresolved references
The report SHALL include a list of wiki-links or relations that could not be resolved after the reconciliation pass.

#### Scenario: Unresolved wiki-links
- **WHEN** an imported object body contains `[[unknown-entity]]` and no matching object exists
- **THEN** the report lists `unknown-entity` as an unresolved reference with the source object ID

#### Scenario: All references resolved
- **WHEN** all wiki-links and relations in imported objects match existing objects
- **THEN** the unresolved references list is empty

### Requirement: Suggest follow-up actions
The report SHALL suggest actionable next steps based on the import results.

#### Scenario: Suggest creating missing relation targets
- **WHEN** the report contains unresolved references
- **THEN** the suggestions include creating objects for the unresolved references

#### Scenario: Suggest reviewing failed imports
- **WHEN** some objects failed to import
- **THEN** the suggestions include reviewing and re-importing the failed files
