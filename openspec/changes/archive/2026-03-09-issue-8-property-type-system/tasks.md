## 1. Property Struct and Schema Changes

- [x] 1.1 Add Option struct (Value, Label fields) and Options field to Property struct in `core/type_schema.go`
- [x] 1.2 Remove Values field from Property struct, update all references
- [x] 1.3 Expand property type allowlist to 9 types (string, number, date, datetime, url, checkbox, select, multi_select, relation)
- [x] 1.4 Add schema validation: select/multi_select require options, reject enum with guidance message
- [x] 1.5 Update default built-in types to use `select` + `options` instead of `enum` + `values`

## 2. Object Validation

- [x] 2.1 Add date validation (YYYY-MM-DD format, handle time.Time from YAML)
- [x] 2.2 Add datetime validation (ISO 8601 with time, handle time.Time from YAML)
- [x] 2.3 Add url validation (http:// or https:// prefix)
- [x] 2.4 Add checkbox validation (boolean only, reject string "true")
- [x] 2.5 Update select validation to use options[].value instead of values[]
- [x] 2.6 Add multi_select validation (list of values, each in options; coerce single string to list)

## 3. Migration (enum → select)

- [x] 3.1 Add enum-to-select migration logic in migrate.go: convert `type: enum` + `values` to `type: select` + `options`
- [x] 3.2 Support --dry-run for enum migration preview
- [x] 3.3 Update YAML serialization to output options format correctly

## 4. Display and Examples

- [x] 4.1 Update display.go for type-aware formatting (dates, checkboxes, URLs)
- [x] 4.2 Update example vault schemas (book.yaml, person.yaml) to use new types

## 5. BDD and Unit Tests

- [x] 5.1 Write BDD feature file for property type validation scenarios (core/features/property_type.feature)
- [x] 5.2 Implement BDD step definitions for property type scenarios
- [x] 5.3 Write unit tests for edge cases: YAML auto-parsing, type coercion, format validation
- [x] 5.4 Write unit tests for enum-to-select migration
