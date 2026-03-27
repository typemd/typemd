## 1. Core: Type Directory Structure Migration

- [x] 1.1 Write BDD scenarios for type directory format loading (single file, directory, auto-migration, precedence)
- [x] 1.2 Implement step definitions for type directory scenarios
- [x] 1.3 Update LocalObjectRepository to support directory format: LoadType checks `types/<name>/schema.yaml` then `types/<name>.yaml`
- [x] 1.4 Implement auto-migration: on load, migrate single file to directory format
- [x] 1.5 Update ListTypes to discover both single-file and directory formats
- [x] 1.6 Update SaveType to always write directory format and remove old single file
- [x] 1.7 Update DeleteType to remove entire directory including views
- [x] 1.8 Add unit tests for edge cases (both formats exist, empty directory, migration failure)

## 2. Core: ViewConfig Struct and YAML Serialization

- [x] 2.1 Write BDD scenarios for ViewConfig creation, YAML serialization/deserialization
- [x] 2.2 Implement step definitions for ViewConfig scenarios
- [x] 2.3 Create `core/view.go` with ViewConfig, ViewLayout, FilterRule, SortRule structs
- [x] 2.4 Implement YAML marshaling/unmarshaling with omitempty for optional fields
- [x] 2.5 Add unit tests for ViewConfig edge cases (empty fields, unknown layout, malformed YAML)

## 3. Core: Vault View CRUD

- [x] 3.1 Write BDD scenarios for ListViews, LoadView, SaveView, DeleteView, DefaultView
- [x] 3.2 Implement step definitions for Vault View CRUD scenarios
- [x] 3.3 Implement Vault.ListViews (read YAML files from views/ directory)
- [x] 3.4 Implement Vault.LoadView (read and parse single view YAML)
- [x] 3.5 Implement Vault.SaveView (write view YAML, trigger directory migration if needed)
- [x] 3.6 Implement Vault.DeleteView (remove view YAML file)
- [x] 3.7 Implement Vault.DefaultView (return saved default or generate implicit default)
- [x] 3.8 Add unit tests for edge cases (save to non-existent type, delete last view, concurrent access)

## 4. Core: QueryService Sort Support

- [x] 4.1 Write BDD scenarios for query with sort (single property, multiple properties, system properties, missing values)
- [x] 4.2 Implement step definitions for query sort scenarios
- [x] 4.3 Create QueryOptions and SortRule types, refactor QueryService.Query to accept QueryOptions
- [x] 4.4 Update SQLiteObjectIndex to generate ORDER BY clauses from SortRule
- [x] 4.5 Handle null/missing values in sort (NULLS LAST for asc, NULLS LAST for desc)
- [x] 4.6 Update all existing callers of QueryService.Query to use QueryOptions
- [x] 4.7 Add unit tests for SQL generation edge cases

## 5. Core: Filter Operator System

- [x] 5.1 Write BDD scenarios for filter operators by property type (string, number, date, select, multi_select, relation, checkbox)
- [x] 5.2 Implement step definitions for filter operator scenarios
- [x] 5.3 Define operator registry: mapping of property type → valid operators
- [x] 5.4 Implement filter operator to SQL translation in SQLiteObjectIndex
- [x] 5.5 Implement filter validation (reject invalid operator for property type)
- [x] 5.6 Implement AND logic for multiple filter rules
- [x] 5.7 Add unit tests for SQL generation, operator validation, and edge cases (is_empty, type casting)

## 6. TUI: View Mode

- [x] 6.1 Add `panelView` to right panel mode enum and ViewConfig state fields to app model
- [x] 6.2 Implement full-width View list rendering (object list with optional group headers)
- [x] 6.3 Implement View mode navigation stack (sidebar → view list → object detail → back)
- [x] 6.4 Implement group_by logic (group objects by property value, collapsible headers)
- [x] 6.5 Implement keyboard shortcut `v` to enter View mode for current type
- [x] 6.6 Implement view selection popup when type has multiple views
- [x] 6.7 Add Views section to type editor listing saved views + default

## 7. Integration and Documentation

- [x] 7.1 Run full test suite, fix any regressions from QueryService refactor
- [x] 7.2 Update CLAUDE.md with View and ViewConfig documentation
- [x] 7.3 Verify auto-migration works with existing test vaults
