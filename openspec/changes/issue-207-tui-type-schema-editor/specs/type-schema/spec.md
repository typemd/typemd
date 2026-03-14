## ADDED Requirements

### Requirement: TypeSchema supports YAML serialization

A `MarshalTypeSchema(schema *TypeSchema) ([]byte, error)` function SHALL serialize a TypeSchema struct to YAML bytes suitable for writing to `.typemd/types/<name>.yaml`.

```
TypeSchema struct                     YAML output
┌──────────────────────┐              ┌──────────────────────────┐
│ Name: "book"         │              │ name: book               │
│ Plural: "books"      │──serialize──▶│ plural: books            │
│ Emoji: "📖"          │              │ emoji: "📖"              │
│ Unique: false        │              │ properties:              │
│ NameTemplate: ""     │              │   - name: author         │
│ Properties:          │              │     type: relation       │
│   [{Name: "author",  │              │     target: person       │
│     Type: "relation",│              │   - name: genre          │
│     Target: "person"}│              │     type: select         │
│    {Name: "genre",   │              │     options:             │
│     Type: "select",  │              │       - value: fiction   │
│     Options: [...]}] │              │       - value: non-fic   │
└──────────────────────┘              └──────────────────────────┘
```

#### Scenario: Serialize complete TypeSchema
- **WHEN** `MarshalTypeSchema` is called with a TypeSchema with name, plural, emoji, and properties
- **THEN** it SHALL return valid YAML containing all non-zero fields

#### Scenario: Serialize TypeSchema with NameTemplate
- **WHEN** `MarshalTypeSchema` is called with a TypeSchema where `NameTemplate` is `"{{ date:YYYY-MM-DD }}"`
- **THEN** the YAML SHALL include a property entry with `name: name` and `template: "{{ date:YYYY-MM-DD }}"`
- **AND** the `name` entry SHALL appear before other properties

#### Scenario: Omit zero-value optional fields
- **WHEN** `MarshalTypeSchema` is called with a TypeSchema where plural is empty and unique is false
- **THEN** the YAML SHALL NOT contain `plural` or `unique` keys

### Requirement: ObjectRepository interface includes DeleteSchema

The `ObjectRepository` interface SHALL include `DeleteSchema(name string) error` for removing type schema YAML files.

#### Scenario: Interface compliance
- **WHEN** a type implements `ObjectRepository`
- **THEN** it SHALL provide a `DeleteSchema(name string) error` method

### Requirement: LocalObjectRepository implements DeleteSchema

`LocalObjectRepository.DeleteSchema(name string)` SHALL delete the file at `.typemd/types/<name>.yaml`. If the file does not exist, it SHALL return an error.

#### Scenario: Delete existing schema file
- **WHEN** `.typemd/types/note.yaml` exists
- **AND** `DeleteSchema("note")` is called
- **THEN** the file SHALL be removed

#### Scenario: Delete non-existent schema file
- **WHEN** no `.typemd/types/unknown.yaml` exists
- **AND** `DeleteSchema("unknown")` is called
- **THEN** it SHALL return an error
