## ADDED Requirements

### Requirement: Config key descriptions
Each config key in the registry SHALL have a human-readable `Description` field. The `ConfigKeysInfo()` function SHALL return metadata (key, description, default, current value) for all registered config keys.

#### Scenario: All config keys have descriptions
- **WHEN** `ConfigKeysInfo()` is called on a vault
- **THEN** every returned entry SHALL have a non-empty `Description`

#### Scenario: Config key info includes current value
- **WHEN** a config key has been set to a non-default value
- **THEN** `ConfigKeysInfo()` SHALL return both the current value and the default value for that key

### Requirement: Config page activation
The TUI SHALL provide a `,` keybinding to enter the config settings page. The config page SHALL use a full-width `panelConfig` right panel mode.

#### Scenario: Enter config page
- **WHEN** the user presses `,` in normal mode
- **THEN** the TUI SHALL switch to `panelConfig` mode and display the config settings page

#### Scenario: Exit config page
- **WHEN** the user presses `Esc` while on the config page (not in edit popup)
- **THEN** the TUI SHALL return to the previous panel mode and restore sidebar focus

### Requirement: Category navigation
The config page SHALL display a left column with config categories (General, CLI, TUI, AI, Web) and a right column showing settings for the selected category.

#### Scenario: Browse categories
- **WHEN** the config page is displayed
- **THEN** the left column SHALL list categories derived from config key prefixes
- **AND** the first category SHALL be selected by default

#### Scenario: Switch categories
- **WHEN** the user navigates categories with `j`/`k` in the left column
- **THEN** the right column SHALL update to show settings for the newly selected category

#### Scenario: Switch columns
- **WHEN** the user presses `Tab`
- **THEN** focus SHALL alternate between the category column and the settings column

### Requirement: Settings browsing
Within a category, the settings column SHALL list all config keys belonging to that category with their current values.

#### Scenario: Display settings list
- **WHEN** a category is selected
- **THEN** each setting in that category SHALL display its key name, description, and current value (or default if unset)

#### Scenario: Navigate settings
- **WHEN** the user presses `j`/`k` in the settings column
- **THEN** the cursor SHALL move between settings within the current category

### Requirement: Setting editing via popup
Pressing `Enter` on a setting SHALL open a `CenteredPopup` for editing. The popup SHALL show the key name, description, default value, and an input field pre-filled with the current value.

#### Scenario: Edit string setting
- **WHEN** the user presses `Enter` on a string-type setting
- **THEN** a popup SHALL appear with a text input pre-filled with the current value
- **AND** pressing `Enter` in the popup SHALL save the new value immediately via `SetConfigValue()`
- **AND** the popup SHALL close and the settings list SHALL reflect the updated value

#### Scenario: Edit boolean setting
- **WHEN** the user presses `Enter` on a boolean-type setting (e.g., `tui.toast.show_warnings`)
- **THEN** the value SHALL cycle through: "true" → "false" → "unset" (or "unset" → "true" → "false" if currently unset)
- **AND** "unset" SHALL remove the key from config.yaml, restoring the default

#### Scenario: Cancel editing
- **WHEN** the user presses `Esc` while the edit popup is open
- **THEN** the popup SHALL close without saving any changes

#### Scenario: Clear to default
- **WHEN** the user saves an empty string for a setting
- **THEN** the key SHALL be removed from config.yaml (restoring default behavior)

### Requirement: Help integration
The `,` keybinding SHALL appear in the TUI help popup under the appropriate section.

#### Scenario: Help shows settings key
- **WHEN** the user opens the help popup (press `?`)
- **THEN** the `,` key SHALL be listed with description "settings"

### Requirement: Config page help bar
While on the config page, a contextual help bar SHALL display navigation hints.

#### Scenario: Config page help bar content
- **WHEN** the config page is active
- **THEN** the help bar SHALL show navigation keys: `j/k` navigate, `Tab` switch column, `Enter` edit, `Esc` back
