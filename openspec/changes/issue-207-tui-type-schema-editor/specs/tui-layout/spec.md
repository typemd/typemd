## ADDED Requirements

### Requirement: Right panel supports multiple view modes

The TUI right panel SHALL support three view modes controlled by a `rightPanelMode` enum: `panelEmpty` (no content selected), `panelObject` (object detail view — existing behavior), and `panelTypeEditor` (type editor view).

```
┌─ rightPanelMode state transitions ──────────────────────────┐
│                                                             │
│  panelEmpty ──Enter on object──▶ panelObject                │
│  panelEmpty ──Enter on header──▶ panelTypeEditor            │
│  panelObject ──Enter on header──▶ panelTypeEditor           │
│  panelObject ──Esc (deselect)──▶ panelEmpty                 │
│  panelTypeEditor ──Esc──▶ panelEmpty                        │
│  panelTypeEditor ──Enter on object──▶ panelObject           │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

#### Scenario: Panel mode transitions
- **WHEN** the right panel is in `panelEmpty` mode
- **AND** the user presses `Enter` on a type group header
- **THEN** the right panel SHALL switch to `panelTypeEditor` mode

#### Scenario: Object selection from type editor
- **WHEN** the right panel is in `panelTypeEditor` mode
- **AND** the user navigates to an object in the sidebar and presses `Enter`
- **THEN** the right panel SHALL switch to `panelObject` mode

### Requirement: Sidebar "+ New Type" row

A non-collapsible "+ New Type" row SHALL appear at the bottom of the sidebar, below all type groups. It SHALL be selectable via cursor navigation.

```
┌─ Sidebar ──────────────┐
│ ▶ 📖 Books (3)         │
│ ▶ 👤 People (2)        │
│ ▶ 📝 Notes (5)         │
│ ▶ 🏷️ Tags (8)          │
│                        │
│ + New Type             │  ← always at bottom
└────────────────────────┘
```

#### Scenario: New Type row appears in sidebar
- **WHEN** the sidebar is rendered
- **THEN** a "+ New Type" row SHALL appear below the last type group

#### Scenario: New Type row is cursor-navigable
- **WHEN** the cursor is on the last type group header
- **AND** the user presses `↓`
- **THEN** the cursor SHALL move to the "+ New Type" row

## MODIFIED Requirements

### Requirement: TUI startup initializes from restored state
The TUI `Start()` function SHALL attempt to load session state from `.typemd/tui-state.yaml` before applying default values. Restored state values take precedence over hardcoded defaults. If no state file exists or loading fails, the TUI SHALL use the current default behavior (first group expanded, first object selected).

The right panel mode SHALL always initialize to `panelEmpty` on startup, regardless of saved state. The panel mode is determined by user interaction after startup, not persisted.

#### Scenario: Startup with saved state
- **WHEN** the TUI starts and `.typemd/tui-state.yaml` contains valid state
- **THEN** the TUI SHALL initialize with the restored state instead of hardcoded defaults
- **AND** the right panel mode SHALL be `panelEmpty`

#### Scenario: Startup without saved state
- **WHEN** the TUI starts and no `.typemd/tui-state.yaml` exists
- **THEN** the TUI SHALL initialize with current defaults (first group expanded, cursor at top, focus on left panel)
- **AND** the right panel mode SHALL be `panelEmpty`

### Requirement: Space toggles expand/collapse on type headers

The `Space` key on a type group header SHALL toggle the group's expand/collapse state. This separates the structural action (expand/collapse) from the primary action (`Enter` = open editor).

```
┌─ Key behavior on sidebar items ─────────────────────────────┐
│                                                             │
│  Item type        Enter              Space                  │
│  ─────────────    ─────────────────  ──────────────────     │
│  Type header      Open type editor   Toggle expand/collapse │
│  Object           Select object      Select object          │
│  + New Type       Start creation     Start creation         │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

#### Scenario: Space toggles group
- **WHEN** the cursor is on a collapsed type group header
- **AND** the user presses `Space`
- **THEN** the group SHALL expand to show its objects

#### Scenario: Enter opens type editor
- **WHEN** the cursor is on a type group header
- **AND** the user presses `Enter`
- **THEN** the right panel SHALL switch to type editor mode for that type
