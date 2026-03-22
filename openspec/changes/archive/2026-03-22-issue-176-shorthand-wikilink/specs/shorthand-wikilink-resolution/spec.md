## ADDED Requirements

### Requirement: Type-qualified shorthand links are resolved during sync

The system SHALL resolve wiki-links in `[[type/name]]` format (without ULID suffix) by searching the name index within the specified type. If exactly one object matches, the link SHALL be resolved to the full object ID. If multiple objects match, the link SHALL be treated as ambiguous and left unresolved.

#### Scenario: Type-qualified link resolves to unique match

- **WHEN** an object body contains `[[book/clean-code]]` and exactly one book object has the slug `clean-code`
- **THEN** the link resolves to the full ID `book/clean-code-<ulid>` and `to_id` is set accordingly

#### Scenario: Type-qualified link with no match

- **WHEN** an object body contains `[[book/nonexistent]]` and no book object has the slug `nonexistent`
- **THEN** the link is stored with empty `to_id` (broken link)

#### Scenario: Type-qualified link with ambiguous match

- **WHEN** an object body contains `[[book/golang]]` and two book objects have slugs matching `golang`
- **THEN** the link is stored with empty `to_id` and reported as ambiguous in sync result

### Requirement: Same-type shorthand links are resolved during sync

The system SHALL resolve wiki-links in `[[name]]` format (no type prefix, no ULID suffix) by searching the name index within the source object's own type. Resolution rules are identical to type-qualified resolution.

#### Scenario: Same-type link resolves to unique match

- **WHEN** a `book` object body contains `[[clean-code]]` and exactly one book object has the slug `clean-code`
- **THEN** the link resolves to the full ID `book/clean-code-<ulid>`

#### Scenario: Same-type link with no match in own type

- **WHEN** a `book` object body contains `[[some-person]]` and no book object has the slug `some-person`
- **THEN** the link is stored with empty `to_id` (broken link, no cross-type fallback)

#### Scenario: Same-type link with ambiguous match

- **WHEN** a `note` object body contains `[[meeting]]` and two note objects have slugs matching `meeting`
- **THEN** the link is stored with empty `to_id` and reported as ambiguous

### Requirement: Full ID links continue to work unchanged

The system SHALL continue to resolve wiki-links in `[[type/name-ulid]]` format using exact ID matching against known objects. This existing behavior SHALL NOT change.

#### Scenario: Full ID link resolves exactly

- **WHEN** an object body contains `[[book/clean-code-01abc]]` and that exact object ID exists
- **THEN** the link resolves with `to_id` set to `book/clean-code-01abc`

### Requirement: Resolved shorthand links are written back to source files

When a shorthand wiki-link is successfully resolved during sync, the system SHALL replace the shorthand target with the full object ID in the source file's body content. Display text SHALL be preserved during write-back.

#### Scenario: Type-qualified shorthand is expanded in source file

- **WHEN** an object body contains `[[book/clean-code]]` and it resolves to `book/clean-code-01abc`
- **THEN** the source file body is updated to `[[book/clean-code-01abc]]`

#### Scenario: Same-type shorthand is expanded in source file

- **WHEN** a `book` object body contains `[[clean-code]]` and it resolves to `book/clean-code-01abc`
- **THEN** the source file body is updated to `[[book/clean-code-01abc]]`

#### Scenario: Display text is preserved during write-back

- **WHEN** an object body contains `[[clean-code|Clean Code by Uncle Bob]]` and it resolves to `book/clean-code-01abc`
- **THEN** the source file body is updated to `[[book/clean-code-01abc|Clean Code by Uncle Bob]]`

#### Scenario: Unresolved shorthand is not modified

- **WHEN** an object body contains `[[nonexistent]]` and no match is found
- **THEN** the source file body retains `[[nonexistent]]` unchanged

### Requirement: Shorthand resolution results are reported in sync result

The system SHALL report the count of expanded wiki-links and list unresolved wiki-links (with reason and matches) in the sync result, enabling toast notifications and detailed logging.

#### Scenario: Expanded links are counted

- **WHEN** sync resolves 3 shorthand wiki-links across the vault
- **THEN** the sync result includes `WikiLinksExpanded: 3`

#### Scenario: Ambiguous links are reported with matches

- **WHEN** sync encounters `[[golang]]` matching `book/golang-01abc` and `book/golang-guide-01def`
- **THEN** the sync result includes an unresolved wiki-link entry with reason `ambiguous` and both matches listed

### Requirement: Name resolution uses slug and name property

Shorthand targets SHALL be matched against both the object slug (filename without ULID) and the slugified name property value, consistent with how relation name resolution works.

#### Scenario: Resolution matches by slug

- **WHEN** a wiki-link target is `clean-code` and a book object has filename `clean-code-01abc.md`
- **THEN** the link resolves to `book/clean-code-01abc`

#### Scenario: Resolution matches by name property

- **WHEN** a wiki-link target is `clean-code` and a book object has name property "Clean Code" (slugified to `clean-code`) but a different filename slug
- **THEN** the link resolves to that book object's full ID
