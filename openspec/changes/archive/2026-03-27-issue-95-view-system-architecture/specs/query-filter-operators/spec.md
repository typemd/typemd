## ADDED Requirements

### Requirement: String property filter operators

String properties SHALL support the following filter operators: `is`, `is_not`, `contains`, `does_not_contain`, `starts_with`, `ends_with`, `is_empty`, `is_not_empty`.

#### Scenario: String "is" operator

- **WHEN** a filter rule has property "author" with operator "is" and value "Tolkien"
- **THEN** the query SHALL return only objects where the "author" property exactly equals "Tolkien"

#### Scenario: String "contains" operator

- **WHEN** a filter rule has property "author" with operator "contains" and value "Tolk"
- **THEN** the query SHALL return objects where the "author" property contains the substring "Tolk"

#### Scenario: String "starts_with" operator

- **WHEN** a filter rule has property "title" with operator "starts_with" and value "The"
- **THEN** the query SHALL return objects where the "title" property starts with "The"

#### Scenario: String "is_empty" operator

- **WHEN** a filter rule has property "author" with operator "is_empty"
- **THEN** the query SHALL return objects where the "author" property is missing, null, or an empty string

### Requirement: Number property filter operators

Number properties SHALL support the following filter operators: `eq`, `neq`, `gt`, `gte`, `lt`, `lte`, `is_empty`, `is_not_empty`.

#### Scenario: Number "gt" operator

- **WHEN** a filter rule has property "rating" with operator "gt" and value "4"
- **THEN** the query SHALL return only objects where the "rating" property is greater than 4

#### Scenario: Number "lte" operator

- **WHEN** a filter rule has property "pages" with operator "lte" and value "300"
- **THEN** the query SHALL return only objects where the "pages" property is less than or equal to 300

#### Scenario: Number "eq" operator

- **WHEN** a filter rule has property "rating" with operator "eq" and value "5"
- **THEN** the query SHALL return only objects where the "rating" property equals 5

### Requirement: Date property filter operators

Date and datetime properties SHALL support the following filter operators: `eq`, `before`, `after`, `on_or_before`, `on_or_after`, `is_empty`, `is_not_empty`.

#### Scenario: Date "after" operator

- **WHEN** a filter rule has property "published_date" with operator "after" and value "2025-01-01"
- **THEN** the query SHALL return only objects where "published_date" is after 2025-01-01

#### Scenario: Date "on_or_before" operator

- **WHEN** a filter rule has property "due_date" with operator "on_or_before" and value "2026-03-31"
- **THEN** the query SHALL return only objects where "due_date" is on or before 2026-03-31

#### Scenario: Date "is_empty" operator

- **WHEN** a filter rule has property "published_date" with operator "is_empty"
- **THEN** the query SHALL return objects where "published_date" is missing, null, or an empty string

### Requirement: Select property filter operators

Select properties SHALL support the following filter operators: `is`, `is_not`, `is_empty`, `is_not_empty`.

#### Scenario: Select "is" operator

- **WHEN** a filter rule has property "status" with operator "is" and value "reading"
- **THEN** the query SHALL return only objects where the "status" property exactly equals "reading"

#### Scenario: Select "is_not" operator

- **WHEN** a filter rule has property "status" with operator "is_not" and value "dropped"
- **THEN** the query SHALL return only objects where the "status" property does not equal "dropped"

### Requirement: Multi-select property filter operators

Multi-select properties SHALL support the following filter operators: `contains`, `does_not_contain`, `is_empty`, `is_not_empty`.

#### Scenario: Multi-select "contains" operator

- **WHEN** a filter rule has property "genres" with operator "contains" and value "fantasy"
- **THEN** the query SHALL return objects where the "genres" multi-select includes "fantasy"

#### Scenario: Multi-select "does_not_contain" operator

- **WHEN** a filter rule has property "genres" with operator "does_not_contain" and value "horror"
- **THEN** the query SHALL return objects where the "genres" multi-select does not include "horror"

### Requirement: Relation property filter operators

Relation properties SHALL support the following filter operators: `contains`, `does_not_contain`, `is_empty`, `is_not_empty`.

#### Scenario: Relation "contains" operator

- **WHEN** a filter rule has property "author" (relation type) with operator "contains" and value matching an object ID
- **THEN** the query SHALL return objects where the "author" relation includes the specified object

#### Scenario: Relation "is_empty" operator

- **WHEN** a filter rule has property "author" (relation type) with operator "is_empty"
- **THEN** the query SHALL return objects where the "author" relation has no linked objects

### Requirement: Checkbox property filter operators

Checkbox properties SHALL support the following filter operators: `is`, `is_not`.

#### Scenario: Checkbox "is" true

- **WHEN** a filter rule has property "completed" with operator "is" and value "true"
- **THEN** the query SHALL return only objects where the "completed" property is true

#### Scenario: Checkbox "is" false

- **WHEN** a filter rule has property "completed" with operator "is" and value "false"
- **THEN** the query SHALL return only objects where the "completed" property is false

### Requirement: Multiple filter rules use AND logic

When multiple filter rules are specified, they SHALL be combined with AND logic. An object must match ALL filter rules to be included in the result.

#### Scenario: Multiple filters combined

- **WHEN** filter rules are [{property: "status", operator: "is", value: "reading"}, {property: "rating", operator: "gt", value: "3"}]
- **THEN** the query SHALL return only objects where status is "reading" AND rating is greater than 3

### Requirement: Filter operator validation against property type

Filter validation SHALL reject operators that are not valid for the property's type. Validation requires resolving the property type from the type schema.

#### Scenario: Invalid operator for property type

- **WHEN** a filter rule uses operator "gt" on a string property "author"
- **THEN** validation SHALL return an error indicating "gt" is not a valid operator for string properties

#### Scenario: Valid operator for property type

- **WHEN** a filter rule uses operator "contains" on a string property "author"
- **THEN** validation SHALL accept without error

### Requirement: SQLiteObjectIndex translates filter operators to SQL

The `SQLiteObjectIndex` SHALL translate each filter operator to the appropriate SQL expression using `json_extract(properties, '$.<property>')`.

#### Scenario: "is" operator to SQL

- **WHEN** a filter has operator "is" with value "reading"
- **THEN** the SQL condition SHALL be `json_extract(properties, '$.status') = 'reading'`

#### Scenario: "contains" operator to SQL for string

- **WHEN** a filter has operator "contains" with value "Tolk"
- **THEN** the SQL condition SHALL use `LIKE '%Tolk%'`

#### Scenario: "gt" operator to SQL for number

- **WHEN** a filter has operator "gt" with value "4"
- **THEN** the SQL condition SHALL be `CAST(json_extract(properties, '$.rating') AS REAL) > 4`

#### Scenario: "is_empty" operator to SQL

- **WHEN** a filter has operator "is_empty"
- **THEN** the SQL condition SHALL check for NULL, empty string, or missing JSON key
