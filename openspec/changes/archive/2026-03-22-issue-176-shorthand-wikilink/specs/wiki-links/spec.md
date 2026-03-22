## MODIFIED Requirements

### Requirement: Broken links are detected by validation

The `tmd type validate` command SHALL report links whose targets do not resolve to existing objects. For shorthand targets that match multiple objects (ambiguous), validation SHALL report the ambiguous target along with the list of matching full IDs as suggestions.

#### Scenario: Broken link is reported

- **WHEN** an object contains a link `[[person/nobody-01jjjjjjjjjjjjjjjjjjjjjjjj]]` that does not resolve
- **THEN** validation reports `<object-id>: broken wiki-link [[person/nobody-01jjjjjjjjjjjjjjjjjjjjjjjj]]`

#### Scenario: Valid links pass validation

- **WHEN** all links in the vault resolve to existing objects
- **THEN** link validation reports no errors

#### Scenario: Ambiguous shorthand link is reported with suggestions

- **WHEN** an object contains `[[golang]]` matching `book/golang-01abc` and `book/golang-guide-01def`
- **THEN** validation reports the ambiguous link and lists both matching full IDs as suggestions

#### Scenario: Resolvable shorthand link passes validation

- **WHEN** an object contains `[[clean-code]]` and exactly one object matches within the source type
- **THEN** link validation reports no error for that link
