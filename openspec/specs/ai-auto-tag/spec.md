### Requirement: ctrl+g on tags field triggers AI tag suggestion

In TUI object detail view (`panelObject`), when the cursor is on the `tags` property and the user presses `ctrl+g`, the system SHALL invoke `AIService.SuggestTags()` with the current object, its type schema, and the list of all existing tags in the vault.

#### Scenario: AI suggests tags for an object

- **WHEN** user presses `ctrl+g` on the `tags` field
- **THEN** the system SHALL show a loading indicator
- **AND** SHALL call `AIService.SuggestTags()` with the object, schema, and all existing tags
- **AND** SHALL display results as a selectable popup list

#### Scenario: ctrl+g on tags when AI is unavailable

- **WHEN** user presses `ctrl+g` on `tags` and `vault.AIService()` is nil
- **THEN** nothing SHALL happen (silent no-op)

### Requirement: Tag suggestion popup displays existing and new tags

The tag suggestion popup SHALL display each suggested tag with its classification: `existing` (tag already exists in the vault) or `new` (tag would be created). Each tag SHALL have a checkbox for selection. Tags already assigned to the object SHALL be excluded from suggestions.

#### Scenario: Popup shows mixed existing and new tags

- **WHEN** AI suggests tags ["go", "concurrency", "parallel-computing"]
- **AND** "go" and "concurrency" exist in the vault but "parallel-computing" does not
- **THEN** the popup SHALL show "go" and "concurrency" marked as `existing`
- **AND** "parallel-computing" marked as `new`

#### Scenario: Already-assigned tags are excluded

- **WHEN** the object already has tag "go" assigned
- **AND** AI suggests ["go", "concurrency"]
- **THEN** only "concurrency" SHALL appear in the popup

### Requirement: User confirms tag selection from popup

The user SHALL navigate the popup with up/down keys, toggle selection with Space, confirm with Enter, and cancel with Esc.

#### Scenario: User selects and confirms tags

- **WHEN** user selects "concurrency" and "parallel-computing" in the popup
- **AND** presses Enter
- **THEN** "concurrency" SHALL be linked to the object as a tag relation
- **AND** "parallel-computing" SHALL be created as a new tag object and linked
- **AND** the object SHALL be marked as dirty

#### Scenario: User cancels tag selection

- **WHEN** user presses Esc in the tag popup
- **THEN** the popup SHALL close
- **AND** the object SHALL remain unchanged

#### Scenario: AI call fails

- **WHEN** `AIService.SuggestTags()` returns an error
- **THEN** the loading indicator SHALL be replaced with an error message
- **AND** the error message SHALL be dismissable with any key

### Requirement: Tag suggestion prompt includes existing tag list

The prompt sent to the AI for tag suggestion SHALL include: the object's name, properties, body content, the type schema context, and a complete list of existing tag names (with descriptions if available) so the AI can prefer existing tags over inventing new ones.

#### Scenario: Existing tags with descriptions influence suggestions

- **WHEN** the vault has tags [{"name": "go", "description": "Go programming language"}, {"name": "rust", "description": "Rust programming language"}]
- **AND** the object is about Go concurrency patterns
- **THEN** the prompt SHALL include both tag names and descriptions
- **AND** the AI SHALL be more likely to suggest "go" than create a new "golang" tag

### Requirement: TagSuggestion response structure

`AIService.SuggestTags()` SHALL return a `*TagSuggestion` containing a list of `SuggestedTag` entries. Each `SuggestedTag` SHALL have `Name` (string), `IsNew` (bool — true if the tag does not exist in the vault), and `Reason` (string — brief explanation of why this tag is relevant).

#### Scenario: Structured tag response

- **WHEN** `SuggestTags` returns successfully
- **THEN** each entry in the result SHALL have a non-empty `Name` and `Reason`
- **AND** `IsNew` SHALL be `true` only for tags not found in the vault's existing tag list
