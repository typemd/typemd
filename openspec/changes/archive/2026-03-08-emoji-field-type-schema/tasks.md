## 1. Core: TypeSchema Struct

- [x] 1.1 Add `Emoji string` field with `yaml:"emoji"` tag to `TypeSchema` struct in `core/type_schema.go`
- [x] 1.2 Update built-in default types (book → 📚, person → 👤, note → 📝) in `core/type_schema.go`
- [x] 1.3 Add unit tests for emoji field parsing and default values in `core/type_schema_test.go`

## 2. CLI: Display Emoji

- [x] 2.1 Update `tmd type show` to display emoji alongside type name in `cmd/type_show.go`
- [x] 2.2 Update `tmd type list` to display emoji alongside type names in `cmd/type_list.go`
- [x] 2.3 Add tests for CLI emoji display

## 3. Examples and Documentation

- [x] 3.1 Update example type YAML files in `examples/book-vault/.typemd/types/` to include emoji field
