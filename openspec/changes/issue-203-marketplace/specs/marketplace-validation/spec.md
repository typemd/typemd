## ADDED Requirements

### Requirement: Automated PR validation

A GitHub Actions workflow SHALL validate marketplace plugin structure on PRs that modify `marketplace/`.

#### Scenario: PR with valid plugin passes CI
- **WHEN** a PR adds a correctly structured plugin under `marketplace/plugins/`
- **THEN** CI validation SHALL pass

#### Scenario: PR with invalid plugin fails CI
- **WHEN** a PR adds a plugin missing required files or with invalid JSON/frontmatter
- **THEN** CI validation SHALL fail with a descriptive error
