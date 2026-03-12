## ADDED Requirements

### Requirement: Automatic publish to typemd/marketplace

A GitHub Actions workflow SHALL push `marketplace/` contents to the `typemd/marketplace` repo when changes merge to main.

#### Scenario: Marketplace changes auto-publish
- **WHEN** a commit touching `marketplace/` is merged to main
- **THEN** the `typemd/marketplace` repo SHALL be updated with the latest marketplace contents at its root
