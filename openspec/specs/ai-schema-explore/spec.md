### Requirement: ctrl+e from sidebar enters schema explore mode

When the user presses `ctrl+e` from the sidebar and `vault.AIService()` is non-nil, the TUI SHALL prompt the user to select a type to analyze, then transition to a full-width `panelSchemaExplore` panel mode.

#### Scenario: Enter schema explore with available types

- **WHEN** user presses `ctrl+e` from the sidebar
- **AND** AI is available
- **THEN** the system SHALL display a type selection prompt
- **AND** after selection, SHALL transition to schema explore panel

#### Scenario: ctrl+e when AI is unavailable

- **WHEN** user presses `ctrl+e` and `vault.AIService()` is nil
- **THEN** nothing SHALL happen (silent no-op)

#### Scenario: Exit schema explore

- **WHEN** user presses `Esc` in the schema explore panel
- **THEN** the TUI SHALL return to the previous panel mode (sidebar focus)

### Requirement: Schema explore samples objects for analysis

The schema explore feature SHALL sample objects of the selected type for AI analysis. The number of objects sampled SHALL be configurable via `ai.explore.sample_count` (default 10). Object bodies SHALL be truncated to `ai.explore.body_truncate` characters (default 500). If the type has fewer objects than `sample_count`, all objects SHALL be used.

#### Scenario: Type with fewer objects than sample count

- **WHEN** type "book" has 5 objects and `sample_count` is 10
- **THEN** all 5 objects SHALL be included in the analysis

#### Scenario: Type with more objects than sample count

- **WHEN** type "note" has 50 objects and `sample_count` is 10
- **THEN** exactly 10 objects SHALL be sampled

#### Scenario: Body truncation

- **WHEN** an object's body is 1200 characters and `body_truncate` is 500
- **THEN** only the first 500 characters of the body SHALL be included in the prompt

### Requirement: Schema explore displays structured suggestions

The AI analysis SHALL return structured suggestions, each with a `Type` (one of `add`, `modify`, `remove`), `PropertyName` (string), `PropertyType` (string, for `add`/`modify`), `Reason` (string explanation), and `Description` (string, suggested property description). The explore panel SHALL display each suggestion with its type, property details, and reason.

#### Scenario: AI suggests adding a new property

- **WHEN** AI analyzes book objects and finds many mention a "publisher" in their body
- **THEN** the suggestion SHALL have `Type: "add"`, `PropertyName: "publisher"`, `PropertyType: "string"`, and a `Reason` explaining why

#### Scenario: AI suggests removing an unused property

- **WHEN** AI finds that no sampled objects use the "isbn" property
- **THEN** the suggestion SHALL have `Type: "remove"`, `PropertyName: "isbn"`, and a `Reason` explaining it appears unused

#### Scenario: AI suggests modifying a property type

- **WHEN** AI finds that the "rating" property contains numeric values but is typed as "string"
- **THEN** the suggestion SHALL have `Type: "modify"`, `PropertyName: "rating"`, `PropertyType: "number"`, and a `Reason`

### Requirement: User accepts or skips individual suggestions

Each suggestion in the explore panel SHALL have Accept and Skip actions. The user SHALL navigate between suggestions with up/down keys and press Enter to accept or `s` to skip the current suggestion.

#### Scenario: User accepts an add-property suggestion

- **WHEN** user presses Enter on an "add publisher (string)" suggestion
- **THEN** the type schema SHALL be modified to include the new property
- **AND** the schema file SHALL be saved
- **AND** the suggestion SHALL be marked as accepted in the UI

#### Scenario: User skips a suggestion

- **WHEN** user presses `s` on a suggestion
- **THEN** the suggestion SHALL be marked as skipped in the UI
- **AND** the type schema SHALL NOT be modified

#### Scenario: All suggestions processed

- **WHEN** all suggestions have been accepted or skipped
- **THEN** the explore panel SHALL show a summary of applied changes
- **AND** the user can press Esc to return to the sidebar

### Requirement: SchemaExploration response structure

`AIService.ExploreSchema()` SHALL return a `*SchemaExploration` containing a list of `SchemaSuggestion` entries. Each `SchemaSuggestion` SHALL have `Type` (string: "add", "modify", "remove"), `PropertyName` (string), `PropertyType` (string, empty for "remove"), `Reason` (string), and `Description` (string, suggested property description text).

#### Scenario: Structured exploration response

- **WHEN** `ExploreSchema` returns successfully
- **THEN** each suggestion SHALL have a non-empty `Type`, `PropertyName`, and `Reason`
- **AND** `Type` SHALL be one of "add", "modify", or "remove"
