## ADDED Requirements

### Requirement: Configurable global TUI keybindings
The TUI SHALL allow users to override any of the 18 global keybindings via a `tui.keybindings` map in `.typemd/config.yaml`. Action names SHALL be `snake_case` (`stats`, `schema_explore`, `toggle_props`, …). Unset actions SHALL keep their default bindings. An empty string value SHALL also be treated as unset.

#### Scenario: Override a single keybinding
- **WHEN** the user sets `tui.keybindings.stats: "ctrl+d"` in `.typemd/config.yaml`
- **THEN** pressing `ctrl+d` in the TUI SHALL activate the stats panel
- **AND** pressing `ctrl+s` SHALL no longer activate stats

#### Scenario: Unset action keeps its default
- **WHEN** the config file has no entry for `tui.keybindings.search`
- **THEN** pressing `/` SHALL still open the search input

#### Scenario: Empty string means default
- **WHEN** `tui.keybindings.search` is set to `""`
- **THEN** pressing `/` SHALL still open the search input

### Requirement: Keybinding validation with fallback
The TUI SHALL validate user-supplied keybindings at startup. Invalid action names, invalid key strings, and duplicate keys SHALL produce a warning toast. The affected action SHALL fall back to its default binding. The TUI SHALL NOT crash on invalid input.

#### Scenario: Unknown action name warns and is ignored
- **WHEN** `tui.keybindings.not_a_real_action: "ctrl+x"` is present
- **THEN** the TUI SHALL show a warning toast naming the unknown action
- **AND** every real action SHALL keep its default binding
- **AND** the TUI SHALL start normally

#### Scenario: Invalid key string falls back to default
- **WHEN** `tui.keybindings.stats: "crtl+s"` (typo) is set
- **THEN** the TUI SHALL show a warning toast naming the invalid key
- **AND** the stats action SHALL use its default key `ctrl+s`
- **AND** the TUI SHALL start normally

#### Scenario: Duplicate key across two actions warns
- **WHEN** `tui.keybindings.stats: "ctrl+d"` and `tui.keybindings.schema_explore: "ctrl+d"` are both set
- **THEN** the TUI SHALL show a warning toast naming the duplicated key and both actions

### Requirement: Help overlay reflects configured keys
The help popup and help bar SHALL display whichever key is currently bound to each action, not the compile-time default.

#### Scenario: Help popup shows overridden key
- **WHEN** `tui.keybindings.stats` is set to `ctrl+d`
- **AND** the user opens the help popup
- **THEN** the help popup SHALL list `ctrl+d` for the stats action, not `ctrl+s`

### Requirement: Keybindings integrate with config key registry
Every rebindable global keybinding SHALL be registered in the config key registry under the `tui.keybindings.<action>` prefix. `tmd config get tui.keybindings.<action>`, `tmd config set tui.keybindings.<action> <key>`, and the TUI Settings page SHALL all work for these keys.

#### Scenario: tmd config set updates a keybinding
- **WHEN** the user runs `tmd config set tui.keybindings.stats ctrl+d`
- **THEN** `.typemd/config.yaml` SHALL contain `tui.keybindings.stats: ctrl+d`
- **AND** `tmd config get tui.keybindings.stats` SHALL return `ctrl+d`

#### Scenario: ConfigKeysInfo exposes keybinding defaults
- **WHEN** `ConfigKeysInfo()` is called on a vault
- **THEN** the returned list SHALL include one entry per rebindable action under the `tui.keybindings.` prefix
- **AND** each entry's `Default` SHALL equal the compile-time default key for that action
- **AND** each entry's `Description` SHALL be non-empty
