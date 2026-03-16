## ADDED Requirements

### Requirement: Slugify converts natural-language names to valid slugs

The `Slugify(name string)` function SHALL convert a natural-language string to a valid slug suitable for filenames and ObjectIDs. The conversion SHALL:
1. Convert to lowercase
2. Replace spaces and underscores with hyphens
3. Remove characters that are not alphanumeric or hyphens
4. Collapse consecutive hyphens into a single hyphen
5. Trim leading and trailing hyphens

#### Scenario: Simple name with spaces

- **WHEN** `Slugify("Some Thought")` is called
- **THEN** it SHALL return `"some-thought"`

#### Scenario: Name with mixed case

- **WHEN** `Slugify("Clean Code")` is called
- **THEN** it SHALL return `"clean-code"`

#### Scenario: Name with underscores

- **WHEN** `Slugify("my_great_idea")` is called
- **THEN** it SHALL return `"my-great-idea"`

#### Scenario: Name with special characters

- **WHEN** `Slugify("What's the plan?")` is called
- **THEN** it SHALL return `"whats-the-plan"`

#### Scenario: Name with consecutive spaces

- **WHEN** `Slugify("too   many   spaces")` is called
- **THEN** it SHALL return `"too-many-spaces"`

#### Scenario: Name with leading/trailing whitespace

- **WHEN** `Slugify("  padded name  ")` is called
- **THEN** it SHALL return `"padded-name"`

#### Scenario: Already-slugified input is idempotent

- **WHEN** `Slugify("clean-code")` is called
- **THEN** it SHALL return `"clean-code"`

#### Scenario: Empty string

- **WHEN** `Slugify("")` is called
- **THEN** it SHALL return `""`

#### Scenario: Name with numbers

- **WHEN** `Slugify("Chapter 3 Notes")` is called
- **THEN** it SHALL return `"chapter-3-notes"`

#### Scenario: Name with non-ASCII characters

- **WHEN** `Slugify("café latte")` is called
- **THEN** it SHALL return `"caf-latte"`

### Requirement: ObjectService.Create applies slug conversion to filename

`ObjectService.Create()` SHALL apply `Slugify()` to the filename parameter when generating the ObjectID. The original (pre-slugified) input SHALL be used as the `name` property value.

#### Scenario: Natural-language name is slugified for filename

- **WHEN** `ObjectService.Create("idea", "Some Great Thought", "")` is called
- **THEN** the ObjectID filename SHALL contain `some-great-thought-<ULID>`
- **AND** the `name` property SHALL be `"Some Great Thought"`

#### Scenario: Pre-slugified name passes through unchanged

- **WHEN** `ObjectService.Create("idea", "already-slugified", "")` is called
- **THEN** the ObjectID filename SHALL contain `already-slugified-<ULID>`
- **AND** the `name` property SHALL be `"already-slugified"`
