## MODIFIED Requirements

### Requirement: Name resolution uses slug, name property, and aliases

Shorthand targets SHALL be matched against the object slug (filename without ULID), the slugified name property value, and each slugified alias in the object's `aliases` system property. This is consistent with how relation name resolution works.

#### Scenario: Resolution matches by slug

- **WHEN** a wiki-link target is `clean-code` and a book object has filename `clean-code-01abc.md`
- **THEN** the link resolves to `book/clean-code-01abc`

#### Scenario: Resolution matches by name property

- **WHEN** a wiki-link target is `clean-code` and a book object has name property "Clean Code" (slugified to `clean-code`) but a different filename slug
- **THEN** the link resolves to that book object's full ID

#### Scenario: Resolution matches by alias

- **WHEN** a wiki-link target is `go-語言` (slugified form of "Go 語言") and a book object has alias "Go 語言"
- **THEN** the link resolves to that book object's full ID

#### Scenario: Alias match is lower priority than exact name match

- **WHEN** object A has alias "Clean Code" and object B has name "Clean Code", and a wiki-link target is `clean-code`
- **THEN** the link resolves to object B (exact name match takes priority over alias match)
