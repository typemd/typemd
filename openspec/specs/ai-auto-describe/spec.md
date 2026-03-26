### Requirement: g on description field triggers AI description generation

In TUI object detail view (`panelObject`), when the cursor is on the `description` property and the user presses `g`, the system SHALL invoke `AIService.Describe()` with the current object and its type schema. The AI feature SHALL only be available when `vault.AIService()` is non-nil.

#### Scenario: AI generates description for object with body content

- **WHEN** user presses `g` on the `description` field of an object that has body content
- **THEN** the system SHALL show a loading indicator
- **AND** SHALL call `AIService.Describe()` with the object and type schema
- **AND** SHALL display the result as an inline preview

#### Scenario: g on description when AI is unavailable

- **WHEN** user presses `g` on `description` and `vault.AIService()` is nil
- **THEN** nothing SHALL happen (silent no-op)

#### Scenario: g on a non-supported field

- **WHEN** user presses `g` on a field that is not `description` or `tags`
- **THEN** nothing SHALL happen (silent no-op)

### Requirement: Inline preview shows AI-generated description before applying

The AI-generated description SHALL be displayed as an inline preview (ghost text style) in the description field area. The user SHALL accept with `Tab` or reject with `Esc`.

#### Scenario: User accepts AI description

- **WHEN** an AI-generated description is shown as inline preview
- **AND** user presses `Tab`
- **THEN** the description SHALL be written to the object's `description` property
- **AND** the object SHALL be marked as dirty (unsaved changes)
- **AND** the inline preview SHALL be dismissed

#### Scenario: User rejects AI description

- **WHEN** an AI-generated description is shown as inline preview
- **AND** user presses `Esc`
- **THEN** the preview SHALL be dismissed
- **AND** the object SHALL remain unchanged

#### Scenario: AI call fails

- **WHEN** `AIService.Describe()` returns an error
- **THEN** the loading indicator SHALL be replaced with an error message
- **AND** the error message SHALL be dismissable with any key

### Requirement: Describe prompt includes semantic context

The prompt sent to the AI for description generation SHALL include: the object's name, all property key-value pairs, the object's body content, and the type schema's property descriptions as semantic context.

#### Scenario: Object with properties and body

- **WHEN** generating a description for an object with name "Go Concurrency", properties {author: "Rob Pike", year: 2012}, body "A book about goroutines...", and type schema with property descriptions
- **THEN** the prompt SHALL include all of: name, properties, body, and type/property descriptions
