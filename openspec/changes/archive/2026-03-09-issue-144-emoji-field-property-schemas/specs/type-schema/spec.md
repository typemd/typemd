## MODIFIED Requirements

### Requirement: Property struct supports options field

The Property struct SHALL support an `options` field containing an array of Option objects. Each Option SHALL have a required `value` string field and an optional `label` string field. When `label` is omitted, it SHALL default to the `value`. The Property struct SHALL also support an optional `emoji` string field for compact visual identification.

#### Scenario: Property with options defined
- **WHEN** a type schema property has `options: [{value: reading, label: Reading}]`
- **THEN** the loaded Property SHALL have Options with one entry where Value is "reading" and Label is "Reading"

#### Scenario: Property with options without label
- **WHEN** a type schema property has `options: [{value: reading}]`
- **THEN** the loaded Property SHALL have Options with one entry where Value is "reading" and Label defaults to "reading"

#### Scenario: Property with emoji and options
- **WHEN** a type schema property has `emoji: 📊` and `options: [{value: active}]`
- **THEN** the loaded Property SHALL have Emoji "📊" and Options with one entry where Value is "active"
