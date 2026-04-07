## ADDED Requirements

### Requirement: Links property returns outgoing wiki-link targets

The system SHALL expose a `links` computed system property that returns the resolved target object IDs of outgoing wiki-links parsed from the object's markdown body. The property SHALL only include links with a non-empty resolved ID (broken links are excluded).

#### Scenario: Object with wiki-links has links property

- **WHEN** an object body contains `[[person/alice-01abc]]` and `[[book/clean-code-01def]]` and both targets exist
- **THEN** the `links` display properties include entries for `person/alice-01abc` and `book/clean-code-01def`

#### Scenario: Object with no wiki-links has empty links

- **WHEN** an object body contains no wiki-links
- **THEN** no `links` display properties are returned

#### Scenario: Broken links are excluded

- **WHEN** an object body contains `[[person/nobody-01zzz]]` and no such object exists (stored with empty `to_id`)
- **THEN** the `links` display properties do not include an entry for that target

#### Scenario: Duplicate wiki-links produce one entry

- **WHEN** an object body contains `[[book/clean-code-01abc]]` twice
- **THEN** only one `links` display property entry appears for that target

### Requirement: Links property is read-only

The `links` property SHALL be a computed system property. Attempting to set it via `SetProperty` SHALL return an error indicating it is read-only.

#### Scenario: Setting links returns error

- **WHEN** a user calls `SetProperty` with key `links`
- **THEN** the system returns an error containing "computed system property"

### Requirement: Links appear in display properties

The system SHALL include `links` entries in the display properties returned by `BuildDisplayProperties`. Each link SHALL be a separate `DisplayProperty` with `IsLink` set to true and `FromID` set to the link's target object ID. Links SHALL appear after reverse relations and before backlinks.

#### Scenario: Links appear between reverse relations and backlinks

- **WHEN** an object has both outgoing links and incoming backlinks
- **THEN** display properties list links before backlinks

#### Scenario: Link display property has correct fields

- **WHEN** an object has an outgoing link to `person/alice-01abc`
- **THEN** the corresponding `DisplayProperty` has Key `links`, IsLink true, and FromID `person/alice-01abc`

### Requirement: Links display format uses arrow prefix

When formatting a `links` display property for display, the system SHALL prefix the target's display ID with `→ ` (right arrow + space).

#### Scenario: Link format shows arrow and display ID

- **WHEN** a `links` display property has FromID `person/alice-01abc`
- **THEN** `FormatValue()` returns `→ person/alice`
