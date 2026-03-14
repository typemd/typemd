## ADDED Requirements

### Requirement: Vault.SaveType persists a TypeSchema to YAML

`Vault.SaveType(schema *TypeSchema)` SHALL validate the schema, serialize it to YAML, and write it to `.typemd/types/<name>.yaml`. If validation fails, it SHALL return an error without writing. If the file already exists, it SHALL be overwritten.

#### Scenario: Save a valid type schema
- **WHEN** `SaveType` is called with a valid TypeSchema named "book" with emoji "📖" and properties [author(relation), genre(select)]
- **THEN** the file `.typemd/types/book.yaml` SHALL be created with valid YAML content
- **AND** the YAML SHALL contain `name: book`, `emoji: 📖`, and the properties array

#### Scenario: Save fails on invalid schema
- **WHEN** `SaveType` is called with a TypeSchema that has duplicate property names
- **THEN** it SHALL return a validation error and SHALL NOT write any file

#### Scenario: Save overwrites existing type
- **WHEN** `.typemd/types/book.yaml` already exists
- **AND** `SaveType` is called with an updated TypeSchema for "book"
- **THEN** the file SHALL be overwritten with the new content

### Requirement: TypeSchema YAML serialization handles NameTemplate

When serializing a TypeSchema that has a non-empty `NameTemplate`, the marshaler SHALL emit a `name` property entry containing only the `template` field. When `NameTemplate` is empty, no `name` property entry SHALL be emitted.

#### Scenario: Serialize with NameTemplate
- **WHEN** a TypeSchema has `NameTemplate: "{{ date:YYYY-MM-DD }}"`
- **THEN** the serialized YAML SHALL include a property entry `- name: name` with `template: "{{ date:YYYY-MM-DD }}"` and no other fields

#### Scenario: Serialize without NameTemplate
- **WHEN** a TypeSchema has an empty `NameTemplate`
- **THEN** the serialized YAML SHALL NOT include a `name` property entry

#### Scenario: Round-trip fidelity
- **WHEN** a TypeSchema is loaded from a YAML file, then serialized back to YAML
- **THEN** loading the serialized YAML SHALL produce an equivalent TypeSchema

### Requirement: TypeSchema YAML serialization preserves property order

The serialized YAML SHALL emit properties in their slice order from `TypeSchema.Properties`. System fields (name, plural, emoji, unique) SHALL appear before the properties array.

#### Scenario: Property order preserved
- **WHEN** a TypeSchema has properties [author, genre, rating] in that order
- **THEN** the serialized YAML SHALL list properties in the same order: author, genre, rating

### Requirement: TypeSchema YAML serialization omits zero-value optional fields

Optional fields with zero values SHALL be omitted from the serialized YAML using `omitempty` semantics: empty strings, zero integers, false booleans, nil slices.

#### Scenario: Minimal property serialization
- **WHEN** a property has only name "title" and type "string" with no emoji, pin, options, or relation fields
- **THEN** the serialized YAML for that property SHALL contain only `name: title` and `type: string`

#### Scenario: Full property serialization
- **WHEN** a property has name "author", type "relation", target "person", multiple true, bidirectional true, inverse "books", emoji "👤", pin 1
- **THEN** the serialized YAML SHALL include all non-zero fields

### Requirement: Vault.DeleteType removes a type schema file

`Vault.DeleteType(name string)` SHALL delete the file `.typemd/types/<name>.yaml`. It SHALL refuse to delete built-in types (currently only "tag").

#### Scenario: Delete user-defined type
- **WHEN** `DeleteType("note")` is called and `.typemd/types/note.yaml` exists
- **THEN** the file SHALL be deleted and no error SHALL be returned

#### Scenario: Delete built-in type rejected
- **WHEN** `DeleteType("tag")` is called
- **THEN** it SHALL return an error containing "cannot delete built-in type"
- **AND** no file SHALL be deleted

#### Scenario: Delete non-existent type
- **WHEN** `DeleteType("unknown")` is called and no such type file exists
- **THEN** it SHALL return an error indicating the type does not exist

### Requirement: ObjectRepository.DeleteSchema interface method

The `ObjectRepository` interface SHALL include a `DeleteSchema(name string) error` method. `LocalObjectRepository` SHALL implement it by deleting `.typemd/types/<name>.yaml`.

#### Scenario: DeleteSchema removes file
- **WHEN** `DeleteSchema("note")` is called
- **THEN** the file at `.typemd/types/note.yaml` SHALL be removed from disk

#### Scenario: DeleteSchema on missing file
- **WHEN** `DeleteSchema("unknown")` is called and no file exists
- **THEN** it SHALL return an error

### Requirement: Vault.CountObjectsByType returns object count

`Vault.CountObjectsByType(typeName string)` SHALL return the number of objects that belong to the specified type. It SHALL query the index for efficiency.

#### Scenario: Count objects for type with objects
- **WHEN** the vault contains 5 objects of type "note"
- **AND** `CountObjectsByType("note")` is called
- **THEN** it SHALL return 5

#### Scenario: Count objects for type with no objects
- **WHEN** the vault contains no objects of type "project"
- **AND** `CountObjectsByType("project")` is called
- **THEN** it SHALL return 0
