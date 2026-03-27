## ADDED Requirements

### Requirement: Type supports directory structure

A type MAY be stored as either a single file (`.typemd/types/<name>.yaml`) or a directory (`.typemd/types/<name>/schema.yaml`). Both formats SHALL be recognized when loading types.

#### Scenario: Load type from directory

- **WHEN** `.typemd/types/book/schema.yaml` exists
- **THEN** `LoadType("book")` SHALL read and parse `book/schema.yaml`

#### Scenario: Load type from single file

- **WHEN** `.typemd/types/note.yaml` exists and `.typemd/types/note/` does not exist
- **THEN** `LoadType("note")` SHALL read and parse `note.yaml`

#### Scenario: Directory takes precedence over single file

- **WHEN** both `.typemd/types/book.yaml` and `.typemd/types/book/schema.yaml` exist
- **THEN** `LoadType("book")` SHALL use `book/schema.yaml`

### Requirement: Auto-migration from single file to directory

When a type is loaded and found in single-file format (`.typemd/types/<name>.yaml`), the system SHALL automatically migrate it to directory format: move `<name>.yaml` to `<name>/schema.yaml`.

#### Scenario: Single file auto-migrated

- **WHEN** `.typemd/types/book.yaml` exists and is loaded
- **THEN** the system SHALL create `.typemd/types/book/` directory, move `book.yaml` to `book/schema.yaml`, and remove the original `book.yaml`

#### Scenario: Directory format not migrated

- **WHEN** `.typemd/types/book/schema.yaml` already exists
- **THEN** no migration SHALL occur

### Requirement: ListTypes discovers both formats

`ListTypes()` SHALL discover types from both single-file (`.typemd/types/<name>.yaml`) and directory (`.typemd/types/<name>/schema.yaml`) formats. Each type SHALL appear only once in the result.

#### Scenario: Mixed format types

- **WHEN** `.typemd/types/` contains `note.yaml` (single file) and `book/schema.yaml` (directory)
- **THEN** `ListTypes()` SHALL return both "note" and "book"

#### Scenario: No duplicate when both formats exist

- **WHEN** both `.typemd/types/book.yaml` and `.typemd/types/book/schema.yaml` exist
- **THEN** `ListTypes()` SHALL return "book" only once

### Requirement: SaveType writes to directory format

`SaveType()` SHALL always write type schemas in directory format (`.typemd/types/<name>/schema.yaml`). If the type was previously in single-file format, the single file SHALL be removed after writing the directory version.

#### Scenario: Save type creates directory

- **WHEN** `SaveType()` is called for type "book"
- **THEN** the schema SHALL be written to `.typemd/types/book/schema.yaml`

#### Scenario: Save type removes old single file

- **WHEN** `SaveType()` is called for type "book" and `.typemd/types/book.yaml` exists
- **THEN** `.typemd/types/book.yaml` SHALL be removed after `.typemd/types/book/schema.yaml` is written

### Requirement: Views stored under type directory

View configuration files SHALL be stored at `.typemd/types/<typeName>/views/<viewName>.yaml`. The `views/` subdirectory SHALL be created automatically when the first view is saved.

#### Scenario: View file path

- **WHEN** a view named "by-rating" is saved for type "book"
- **THEN** the file SHALL be written to `.typemd/types/book/views/by-rating.yaml`

#### Scenario: Views directory created on first save

- **WHEN** `SaveView("book", view)` is called and `.typemd/types/book/views/` does not exist
- **THEN** the `views/` directory SHALL be created before writing the view file

### Requirement: DeleteType removes directory and views

`DeleteType()` SHALL remove the entire type directory including any views. If the type is in single-file format, it SHALL remove the single file.

#### Scenario: Delete type with views

- **WHEN** `.typemd/types/book/` contains `schema.yaml` and `views/by-rating.yaml`
- **THEN** `DeleteType("book")` SHALL remove the entire `book/` directory

#### Scenario: Delete type in single-file format

- **WHEN** `.typemd/types/note.yaml` exists
- **THEN** `DeleteType("note")` SHALL remove `note.yaml`
