## 1. Core: Color field on TypeSchema

- [x] 1.1 Write BDD scenarios for type schema color (preset names, hex format, validation, serialization)
- [x] 1.2 Implement BDD step definitions for color scenarios
- [x] 1.3 Add Color field to TypeSchema struct, marshalSchema struct, ValidateSchema(), and MarshalTypeSchema()
- [x] 1.4 Add ValidColorPresets() function and color validation logic (preset list + hex regex)
- [x] 1.5 Add unit tests for color edge cases (case sensitivity, invalid hex lengths, mixed case hex, empty string)

## 2. Core: Description field on TypeSchema and Property

- [x] 2.1 Write BDD scenarios for TypeSchema description and Property description (parsing, serialization, empty default)
- [x] 2.2 Implement BDD step definitions for description scenarios
- [x] 2.3 Add Description field to TypeSchema struct, Property struct, and marshalSchema struct
- [x] 2.4 Update MarshalTypeSchema() to include description in output
- [x] 2.5 Add unit tests for description edge cases (empty string omitted, multiline, special characters)

## 3. Core: Use entry description override

- [x] 3.1 Write BDD scenarios for use entry with description override (accepted, resolved with/without override)
- [x] 3.2 Implement BDD step definitions for use entry description scenarios
- [x] 3.3 Update validateUseOverrides() to allow description field
- [x] 3.4 Update use entry resolution in LoadType() to apply description override
- [x] 3.5 Add unit tests for use entry description override edge cases

## 4. TUI: Type editor meta fields

- [x] 4.1 Update metaFieldCount from 4 to 6 and add Color (index 3) and Description (index 5) to type editor view
- [x] 4.2 Add inline editing support for Color and Description meta fields
- [x] 4.3 Verify TUI type editor cursor navigation works correctly with new fields

## 5. Integration verification

- [x] 5.1 Run full test suite (go test ./...) and fix any regressions
- [x] 5.2 Manual verification: create type with color and description, verify YAML output
