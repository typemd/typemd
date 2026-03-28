## 1. Core: Config

- [x] 1.1 Write BDD scenarios for date_format/datetime_format config keys (get, set, defaults)
- [x] 1.2 Implement BDD step definitions for config scenarios
- [x] 1.3 Add `DateFormat`/`DatetimeFormat` fields to `VaultConfig` and register config keys in `config.go`

## 2. Core: Display formatting

- [x] 2.1 Write BDD scenarios for FormatValue with configurable date/datetime formats
- [x] 2.2 Implement BDD step definitions for display format scenarios
- [x] 2.3 Export `ConvertDateFormat()` in `name_template.go`
- [x] 2.4 Add `DateFormat`/`DatetimeFormat` fields to `DisplayProperty` and update `FormatValue()` to use them
- [x] 2.5 Update `BuildDisplayProperties()` in `QueryService` to populate format fields from vault config
- [x] 2.6 Add unit tests for FormatValue edge cases (empty format, unrecognized tokens, nil value)

## 3. Core: Default datetime format change

- [x] 3.1 Update existing tests that assert datetime format `2006-01-02T15:04:05` to expect `2006-01-02 15:04:05` (space separator default)
