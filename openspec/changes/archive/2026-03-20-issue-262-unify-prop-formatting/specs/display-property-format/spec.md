## ADDED Requirements

### Requirement: DisplayProperty provides value-only formatting via FormatValue

DisplayProperty SHALL provide a `FormatValue() string` method that returns the formatted value without key prefix. This method SHALL handle all property types: checkbox, date, datetime, multi_select, relation, backlink, reverse relation, and default.

#### Scenario: FormatValue returns value without key for string property
- **WHEN** a DisplayProperty has Key "author" and Value "Robert Martin" and Type "string"
- **THEN** FormatValue() SHALL return "Robert Martin"

#### Scenario: FormatValue returns value without key for date property
- **WHEN** a DisplayProperty has Key "published" and Value time 2024-01-15 and Type "date"
- **THEN** FormatValue() SHALL return "2024-01-15"

#### Scenario: FormatValue returns value without key for datetime property
- **WHEN** a DisplayProperty has Key "created" and Value time 2024-01-15T10:30:00 and Type "datetime"
- **THEN** FormatValue() SHALL return "2024-01-15T10:30:00"

#### Scenario: FormatValue returns value without key for multi_select property
- **WHEN** a DisplayProperty has Key "tags" and Value ["go", "cli"] and Type "multi_select"
- **THEN** FormatValue() SHALL return "[go, cli]"

#### Scenario: FormatValue returns value without key for relation property
- **WHEN** a DisplayProperty has IsRelation true and Value "person/robert-martin-01abc"
- **THEN** FormatValue() SHALL return "→ person/robert-martin"

#### Scenario: FormatValue returns value without key for backlink property
- **WHEN** a DisplayProperty has IsBacklink true and FromID "note/my-note-01abc"
- **THEN** FormatValue() SHALL return "⟵ note/my-note"

#### Scenario: FormatValue returns value without key for reverse relation property
- **WHEN** a DisplayProperty has IsReverse true and FromID "book/clean-code-01abc"
- **THEN** FormatValue() SHALL return "← book/clean-code"

#### Scenario: FormatValue returns empty string for nil value
- **WHEN** a DisplayProperty has Value nil
- **THEN** FormatValue() SHALL return ""

### Requirement: Checkbox properties display as checkmark or empty

DisplayProperty SHALL format checkbox (bool) values as "✓" for true and empty string for false. This applies to both `FormatValue()` and `Format()`.

#### Scenario: Checkbox true displays as checkmark
- **WHEN** a DisplayProperty has Type "checkbox" and Value true
- **THEN** FormatValue() SHALL return "✓"

#### Scenario: Checkbox false displays as empty string
- **WHEN** a DisplayProperty has Type "checkbox" and Value false
- **THEN** FormatValue() SHALL return ""

#### Scenario: Format includes checkmark with key
- **WHEN** a DisplayProperty has Key "active" and Type "checkbox" and Value true
- **THEN** Format() SHALL return "active: ✓"

#### Scenario: Format shows key with empty value for false checkbox
- **WHEN** a DisplayProperty has Key "active" and Type "checkbox" and Value false
- **THEN** Format() SHALL return "active: "

### Requirement: Format delegates to FormatValue for value portion

`Format()` SHALL compose its output as `key + ": " + FormatValue()` for standard properties. Backlink, reverse relation, and nil value formatting SHALL also delegate to FormatValue.

#### Scenario: Format composes key and FormatValue for string
- **WHEN** a DisplayProperty has Key "title" and Value "Hello" and Type "string"
- **THEN** Format() SHALL return "title: Hello"

#### Scenario: Format composes key and FormatValue for date
- **WHEN** a DisplayProperty has Key "published" and Value time 2024-01-15 and Type "date"
- **THEN** Format() SHALL return "published: 2024-01-15"

### Requirement: View mode table uses FormatValue for property columns

View mode table rows SHALL use `DisplayProperty.FormatValue()` instead of the local `formatPropValue()` function for formatting property values in table columns and preview panels.

#### Scenario: View mode table formats date property correctly
- **WHEN** a view mode table column displays a date property with value 2024-06-15
- **THEN** the cell SHALL display "2024-06-15" (not Go default time format)

#### Scenario: View mode table formats multi_select property correctly
- **WHEN** a view mode table column displays a multi_select property with values ["a", "b"]
- **THEN** the cell SHALL display "[a, b]"

#### Scenario: View mode table formats checkbox property as checkmark
- **WHEN** a view mode table column displays a checkbox property with value true
- **THEN** the cell SHALL display "✓"

#### Scenario: View mode preview formats properties consistently
- **WHEN** the view mode preview panel displays property values
- **THEN** it SHALL use FormatValue() for consistent formatting with table columns
