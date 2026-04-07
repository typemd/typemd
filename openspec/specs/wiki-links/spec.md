## Requirements

### Requirement: Link syntax is parsed from markdown body

The system SHALL parse `[[target]]` and `[[target|display text]]` syntax from object markdown body content. Target SHALL be a full object ID in `type/name-ulid` format. Duplicate targets within the same body SHALL be deduplicated, keeping the first occurrence.

#### Scenario: Simple link is parsed

- **WHEN** an object body contains `[[person/bob-01kk3gqm8zrrbjjwkx90f727y6]]`
- **THEN** the parser extracts one link with target `person/bob-01kk3gqm8zrrbjjwkx90f727y6` and empty display text

#### Scenario: Link with display text is parsed

- **WHEN** an object body contains `[[person/bob-01kk3gqm8zrrbjjwkx90f727y6|Uncle Bob]]`
- **THEN** the parser extracts one link with target `person/bob-01kk3gqm8zrrbjjwkx90f727y6` and display text `Uncle Bob`

#### Scenario: Duplicate targets are deduplicated

- **WHEN** an object body contains the same target `[[book/clean-code-01abc]]` twice
- **THEN** the parser returns only one link for that target

### Requirement: Links are stored in the database on sync

The system SHALL extract links from each object body during index sync and store them in the `wikilinks` table. Each sync SHALL replace all existing links for that object (delete + insert).

#### Scenario: Links are created on first sync

- **WHEN** an object body contains a link to an existing object and the index is synced
- **THEN** a link record is stored with the source object as `from_id` and the resolved target as `to_id`

#### Scenario: Links are updated on re-sync

- **WHEN** an object's link target changes and the index is re-synced
- **THEN** old link records for that object are removed and new ones are inserted

#### Scenario: Links to deleted objects are cleaned up

- **WHEN** an object that is the source of links is deleted and the index is synced
- **THEN** all link records with that object as `from_id` are removed

### Requirement: Broken links have empty resolved ID

When a wiki-link target does not match any existing object ID, the system SHALL store the link with an empty `to_id` field, preserving the original target text.

#### Scenario: Link to non-existent object

- **WHEN** an object body contains `[[person/nobody-01jjjjjjjjjjjjjjjjjjjjjjjj]]` and no such object exists
- **THEN** the link record has an empty `to_id` and target `person/nobody-01jjjjjjjjjjjjjjjjjjjjjjjj`

### Requirement: Backlinks are queryable

The system SHALL provide a way to query all objects that link to a given object (backlinks) via the `wikilinks` table.

#### Scenario: Single backlink

- **WHEN** object A contains a link to object B
- **THEN** querying backlinks for object B returns object A

#### Scenario: Multiple backlinks

- **WHEN** objects A and C both contain links to object B
- **THEN** querying backlinks for object B returns both A and C

### Requirement: Outgoing links are displayed as a built-in property

The system SHALL display outgoing wiki-links as a system-level `links` property in object detail views. Only resolved links (non-empty `to_id`) SHALL be included. This property SHALL appear after reverse relations and before backlinks.

#### Scenario: Object with outgoing links shows them in display properties

- **WHEN** object A has resolved outgoing links to objects B and C
- **THEN** object A's display properties include `links` entries for B and C

#### Scenario: Object without outgoing links omits the property

- **WHEN** object A has no wiki-links in its body
- **THEN** object A's display properties do not include `links` entries

#### Scenario: Broken links are excluded from display

- **WHEN** object A has an outgoing link with empty resolved ID
- **THEN** object A's display properties do not include a `links` entry for that target

### Requirement: Backlinks are displayed as a built-in property

The system SHALL display backlinks as a system-level `backlinks` property in object detail views. This property SHALL appear after schema-defined properties and reverse relations.

#### Scenario: Object with backlinks shows them in display properties

- **WHEN** object B has backlinks from objects A and C
- **THEN** object B's display properties include a `backlinks` entry listing A and C

#### Scenario: Object without backlinks omits the property

- **WHEN** object B has no backlinks
- **THEN** object B's display properties do not include a `backlinks` entry

### Requirement: Links are rendered with display text

When rendering markdown body for display, the system SHALL replace wiki-link syntax with human-readable text. `[[target|text]]` SHALL render as the display text. `[[target]]` SHALL render as the DisplayID (target with ULID suffix stripped).

#### Scenario: Render link with display text

- **WHEN** body contains `[[person/bob-01kk3gqm8zrrbjjwkx90f727y6|Uncle Bob]]`
- **THEN** it renders as `Uncle Bob`

#### Scenario: Render link without display text

- **WHEN** body contains `[[person/bob-01kk3gqm8zrrbjjwkx90f727y6]]`
- **THEN** it renders as `person/bob` (ULID stripped)

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

### Requirement: tmd fix wikilinks expands shorthand links

The `tmd fix wikilinks` command SHALL walk all objects in the vault, resolve shorthand wiki-links to full IDs, and write back the expanded targets to source files.

#### Scenario: Shorthand links are expanded

- **WHEN** the vault contains objects with shorthand wiki-links and `tmd fix wikilinks` is run
- **THEN** all resolvable shorthand links are replaced with full IDs in the source files

#### Scenario: Summary is displayed after fix

- **WHEN** `tmd fix wikilinks` completes
- **THEN** the output shows the count of expanded links and the count of unresolved links

#### Scenario: Ambiguous links are reported

- **WHEN** a shorthand link matches multiple objects
- **THEN** the command reports the ambiguous link with the list of matching full IDs

#### Scenario: No changes needed

- **WHEN** all wiki-links in the vault are already full IDs
- **THEN** the command reports that no changes were needed
