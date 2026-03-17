## ADDED Requirements

### Requirement: Type editor displays template section

The type editor SHALL display a "Templates" section after the properties section, listing all templates available for the current type. Each template SHALL be displayed with a 📝 prefix and its name. If no templates exist, the section SHALL show "(none)". The section SHALL include a "+ Add Template" action row.

#### Scenario: Type with multiple templates

- **WHEN** the type editor is open for type "book"
- **AND** "book" has templates "review" and "summary"
- **THEN** the Templates section SHALL display "📝 review" and "📝 summary"

#### Scenario: Type with no templates

- **WHEN** the type editor is open for type "idea"
- **AND** "idea" has no templates
- **THEN** the Templates section SHALL display "(none)"

#### Scenario: Templates section includes add action

- **WHEN** the type editor is open for any type
- **THEN** the Templates section SHALL display a "+ Add Template" action row after the template list

### Requirement: Type editor navigates to template detail on Enter

When the cursor is on a template name in the Templates section and the user presses Enter, the TUI SHALL transition to `panelTemplate` mode showing the selected template's detail view. The type editor's state SHALL be preserved so the user can return via Esc.

#### Scenario: Enter on template opens detail view

- **WHEN** the cursor is on "📝 review" in the Templates section
- **AND** the user presses Enter
- **THEN** the right panel SHALL switch to `panelTemplate` mode
- **AND** the template "review" detail SHALL be displayed

#### Scenario: Esc from template detail returns to type editor

- **WHEN** the user is viewing a template in `panelTemplate` mode
- **AND** the user presses Esc
- **THEN** the right panel SHALL return to `panelTypeEditor` mode
- **AND** the type editor cursor SHALL be on the previously selected template

### Requirement: Template detail displays body and properties

The `panelTemplate` mode SHALL display a title panel, body panel, and properties panel mirroring the object detail layout. The title panel SHALL show "📝 type · template-name". The body panel SHALL show the template's markdown body content. The properties panel SHALL show the template's frontmatter property values.

#### Scenario: Template with body and properties

- **WHEN** template "review" for type "book" has body "## Notes\n" and properties `{status: draft}`
- **THEN** the title panel SHALL display "📝 book · review"
- **AND** the body panel SHALL display "## Notes"
- **AND** the properties panel SHALL display "status: draft"

#### Scenario: Template with body only

- **WHEN** template "simple" has body "## Template\n" and no properties
- **THEN** the body panel SHALL display "## Template"
- **AND** the properties panel SHALL be empty

#### Scenario: Template with properties only

- **WHEN** template "preset" has properties `{status: reading}` and empty body
- **THEN** the body panel SHALL be empty
- **AND** the properties panel SHALL display "status: reading"

### Requirement: Template properties panel shows schema-aware property list

The properties panel in `panelTemplate` mode SHALL show the union of template frontmatter properties and type schema-defined properties. Properties present in the template SHALL show their values. Schema properties not in the template SHALL show as empty placeholders. Immutable system properties (`created_at`, `updated_at`) SHALL be excluded. Mutable system properties (`name`, `description`, `tags`) SHALL be included if present in the template.

#### Scenario: Template property matches schema property

- **WHEN** the type schema defines property "status" with type "select"
- **AND** the template has `status: draft`
- **THEN** the properties panel SHALL display "status: draft"

#### Scenario: Schema property absent from template

- **WHEN** the type schema defines property "rating" with type "number"
- **AND** the template does not have "rating"
- **THEN** the properties panel SHALL display "rating:" with empty value

#### Scenario: Immutable system properties excluded

- **WHEN** the template has `created_at: 2020-01-01T00:00:00Z`
- **THEN** `created_at` SHALL NOT appear in the properties panel

### Requirement: Template body editing via textarea

In `panelTemplate` mode, pressing `e` SHALL enter body edit mode, activating a textarea pre-filled with the template's markdown body. Pressing Esc SHALL exit edit mode and save the changes via `Vault.SaveTemplate`. The textarea behavior SHALL mirror object body editing.

#### Scenario: Enter body edit mode

- **WHEN** the user is in template detail view mode
- **AND** presses `e`
- **THEN** the body panel SHALL switch to an editable textarea
- **AND** the textarea SHALL contain the current template body

#### Scenario: Save body on Esc

- **WHEN** the user is editing the template body
- **AND** presses Esc
- **THEN** the template SHALL be saved with the updated body via `Vault.SaveTemplate`
- **AND** the view SHALL return to read-only mode

#### Scenario: Cancel body edit discards changes

- **WHEN** the user is editing the template body
- **AND** presses Ctrl+C
- **THEN** the changes SHALL be discarded
- **AND** the view SHALL return to read-only mode with the original body

### Requirement: Template property editing via inline input

In the template properties panel, pressing Enter on a property SHALL open an inline text input for editing that property's value. Pressing Enter again SHALL confirm the edit. Pressing Esc SHALL cancel the edit. After confirming, the template SHALL be saved via `Vault.SaveTemplate`.

#### Scenario: Edit property value

- **WHEN** the cursor is on property "status: draft" in the properties panel
- **AND** the user presses Enter
- **THEN** an inline text input SHALL appear with current value "draft"

#### Scenario: Confirm property edit

- **WHEN** the user is editing property "status"
- **AND** changes the value to "reading" and presses Enter
- **THEN** the template SHALL be saved with `status: reading`
- **AND** the properties panel SHALL update to show "status: reading"

#### Scenario: Cancel property edit

- **WHEN** the user is editing property "status"
- **AND** presses Esc
- **THEN** the edit SHALL be cancelled
- **AND** the original value SHALL be preserved

#### Scenario: Clear property value

- **WHEN** the user is editing property "status" with value "draft"
- **AND** clears the input and presses Enter
- **THEN** the property SHALL be removed from the template frontmatter
- **AND** the template SHALL be saved

### Requirement: Template creation via add action

Pressing Enter on "+ Add Template" in the type editor SHALL start a creation flow: an inline text input for the template name. Pressing Enter confirms and creates an empty template file. Pressing Esc cancels. After creation, the template SHALL appear in the Templates section list.

#### Scenario: Create template with valid name

- **WHEN** the cursor is on "+ Add Template"
- **AND** the user presses Enter and types "meeting-notes"
- **AND** presses Enter to confirm
- **THEN** `templates/<type>/meeting-notes.md` SHALL be created as an empty file
- **AND** "📝 meeting-notes" SHALL appear in the Templates section

#### Scenario: Cancel template creation

- **WHEN** the user is entering a template name
- **AND** presses Esc
- **THEN** no template file SHALL be created
- **AND** the cursor SHALL return to "+ Add Template"

#### Scenario: Create template with duplicate name

- **WHEN** the user enters a name that matches an existing template
- **THEN** an error message SHALL be displayed
- **AND** no file SHALL be created

### Requirement: Template deletion with confirmation

In `panelTemplate` mode, pressing `d` SHALL show a delete confirmation prompt. Pressing `y` confirms deletion, removing the template file and returning to the type editor. Pressing `n` or Esc cancels and returns to the template detail view.

#### Scenario: Delete template with confirmation

- **WHEN** the user is viewing template "review" in detail mode
- **AND** presses `d`
- **THEN** a confirmation prompt SHALL appear: "Delete template review? (y/n)"

#### Scenario: Confirm deletion

- **WHEN** the delete confirmation is shown
- **AND** the user presses `y`
- **THEN** the template file SHALL be deleted via `Vault.DeleteTemplate`
- **AND** the right panel SHALL return to `panelTypeEditor` mode
- **AND** the template SHALL no longer appear in the Templates section

#### Scenario: Cancel deletion

- **WHEN** the delete confirmation is shown
- **AND** the user presses `n` or Esc
- **THEN** no file SHALL be deleted
- **AND** the view SHALL return to the template detail

### Requirement: Template detail title panel format

The title panel in `panelTemplate` mode SHALL follow the format: `📝 <type> · <template-name>`. If the type has an emoji defined, it SHALL NOT be shown (📝 replaces the type emoji to visually distinguish templates from objects).

#### Scenario: Title panel for template

- **WHEN** viewing template "review" of type "book" (emoji 📚)
- **THEN** the title panel SHALL display "📝 book · review"

### Requirement: Template detail help bar

The `panelTemplate` mode SHALL display a context-sensitive help bar showing available actions for the current mode.

#### Scenario: View mode help bar

- **WHEN** in template view mode
- **THEN** the help bar SHALL display "e:edit  d:delete  esc:back"

#### Scenario: Edit body mode help bar

- **WHEN** in template body edit mode
- **THEN** the help bar SHALL display "esc:save  ctrl+c:cancel"

#### Scenario: Property edit mode help bar

- **WHEN** in template property edit mode
- **THEN** the help bar SHALL display "enter:confirm  esc:cancel"

#### Scenario: Delete confirmation help bar

- **WHEN** in delete confirmation mode
- **THEN** the help bar SHALL display "y:confirm  n/esc:cancel"
