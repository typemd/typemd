## ADDED Requirements

### Requirement: ViewConfig struct defines view configuration

The `ViewConfig` struct SHALL have the following fields: `Name` (string, required), `Layout` (ViewLayout, required), `Filter` ([]FilterRule, optional), `Sort` ([]SortRule, optional), `GroupBy` (string, optional). `ViewLayout` SHALL be a string type with `ViewLayoutList` as the first defined constant with value `"list"`.

#### Scenario: ViewConfig with all fields

- **WHEN** a ViewConfig is created with name "by-rating", layout "list", filter on status=reading, sort by rating desc, group_by "genre"
- **THEN** all fields SHALL be accessible and correctly stored

#### Scenario: ViewConfig with minimal fields

- **WHEN** a ViewConfig is created with only name "default" and layout "list"
- **THEN** Filter SHALL be nil or empty, Sort SHALL be nil or empty, GroupBy SHALL be an empty string

### Requirement: FilterRule defines a single filter condition

The `FilterRule` struct SHALL have fields: `Property` (string, required), `Operator` (string, required), `Value` (string, optional — not required for is_empty/is_not_empty operators).

#### Scenario: FilterRule with value

- **WHEN** a FilterRule has property "status", operator "is", value "reading"
- **THEN** all fields SHALL be correctly stored

#### Scenario: FilterRule without value

- **WHEN** a FilterRule has property "status", operator "is_empty", and no value
- **THEN** the Value field SHALL be an empty string

### Requirement: SortRule defines a single sort criterion

The `SortRule` struct SHALL have fields: `Property` (string, required), `Direction` (string, required). Direction SHALL be either `"asc"` or `"desc"`.

#### Scenario: SortRule ascending

- **WHEN** a SortRule has property "name" and direction "asc"
- **THEN** both fields SHALL be correctly stored

#### Scenario: SortRule descending

- **WHEN** a SortRule has property "rating" and direction "desc"
- **THEN** both fields SHALL be correctly stored

### Requirement: ViewConfig serializes to and from YAML

ViewConfig SHALL serialize to YAML with the following field mapping: `name` → Name, `layout` → Layout, `filter` → Filter (omitted when empty), `sort` → Sort (omitted when empty), `group_by` → GroupBy (omitted when empty).

#### Scenario: Full ViewConfig to YAML

- **WHEN** a ViewConfig with name "by-rating", layout "list", one filter rule, one sort rule, and group_by "genre" is serialized
- **THEN** the YAML output SHALL contain all fields with correct values

#### Scenario: Minimal ViewConfig to YAML

- **WHEN** a ViewConfig with only name "default" and layout "list" is serialized
- **THEN** the YAML output SHALL contain only `name` and `layout` fields, without `filter`, `sort`, or `group_by`

#### Scenario: YAML to ViewConfig

- **WHEN** a YAML file containing `name: by-rating`, `layout: list`, `sort: [{property: rating, direction: desc}]` is parsed
- **THEN** the resulting ViewConfig SHALL have the correct Name, Layout, and Sort values

### Requirement: Vault ListViews returns all views for a type

`Vault.ListViews(typeName string)` SHALL return all saved ViewConfig objects for the given type by reading YAML files from `.typemd/types/<typeName>/views/`. If no views directory exists, it SHALL return an empty slice.

#### Scenario: Type with saved views

- **WHEN** `.typemd/types/book/views/` contains `default.yaml` and `by-rating.yaml`
- **THEN** `ListViews("book")` SHALL return two ViewConfig objects

#### Scenario: Type with no views directory

- **WHEN** `.typemd/types/note.yaml` exists as a single file (no directory)
- **THEN** `ListViews("note")` SHALL return an empty slice

### Requirement: Vault LoadView reads a specific view

`Vault.LoadView(typeName, viewName string)` SHALL read and parse the YAML file at `.typemd/types/<typeName>/views/<viewName>.yaml`. If the file does not exist, it SHALL return an error.

#### Scenario: Load existing view

- **WHEN** `.typemd/types/book/views/by-rating.yaml` exists with valid content
- **THEN** `LoadView("book", "by-rating")` SHALL return the parsed ViewConfig

#### Scenario: Load non-existent view

- **WHEN** `.typemd/types/book/views/missing.yaml` does not exist
- **THEN** `LoadView("book", "missing")` SHALL return an error

### Requirement: Vault SaveView writes a view to disk

`Vault.SaveView(typeName string, view *ViewConfig)` SHALL serialize the ViewConfig to YAML and write it to `.typemd/types/<typeName>/views/<view.Name>.yaml`. If the type is still in single-file format, SaveView SHALL trigger the directory migration first.

#### Scenario: Save view to existing directory

- **WHEN** `.typemd/types/book/` directory already exists
- **THEN** `SaveView("book", view)` SHALL write the view YAML to `.typemd/types/book/views/<name>.yaml`

#### Scenario: Save view triggers migration

- **WHEN** `.typemd/types/book.yaml` exists as a single file
- **THEN** `SaveView("book", view)` SHALL migrate `book.yaml` to `book/schema.yaml`, create `book/views/`, and write the view file

### Requirement: Vault DeleteView removes a view

`Vault.DeleteView(typeName, viewName string)` SHALL delete the YAML file at `.typemd/types/<typeName>/views/<viewName>.yaml`. Deleting the "default" view SHALL be allowed (the system generates an implicit default when no file exists).

#### Scenario: Delete existing view

- **WHEN** `.typemd/types/book/views/by-rating.yaml` exists
- **THEN** `DeleteView("book", "by-rating")` SHALL remove the file

#### Scenario: Delete non-existent view

- **WHEN** `.typemd/types/book/views/missing.yaml` does not exist
- **THEN** `DeleteView("book", "missing")` SHALL return an error

### Requirement: Vault DefaultView returns implicit default

`Vault.DefaultView(typeName string)` SHALL return a ViewConfig with name "default", layout "list", sort by "name" ascending, no filter, and no group_by. If a saved `views/default.yaml` exists, it SHALL return the saved version instead.

#### Scenario: No saved default view

- **WHEN** `.typemd/types/book/views/default.yaml` does not exist
- **THEN** `DefaultView("book")` SHALL return a ViewConfig with name "default", layout "list", sort [{property: "name", direction: "asc"}]

#### Scenario: Saved default view exists

- **WHEN** `.typemd/types/book/views/default.yaml` exists with custom sort by rating desc
- **THEN** `DefaultView("book")` SHALL return the saved ViewConfig with sort by rating desc
