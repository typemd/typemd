## MODIFIED Requirements

### Requirement: CLI create command supports type flag

The `tmd object create` command SHALL accept an optional `--type` flag to specify the object type. When `--type` is provided, all positional arguments are treated as the name. The `--type` flag SHALL NOT have a `-t` short form (reserved by `--template`).

#### Scenario: Create with --type flag and name

- **WHEN** `tmd object create --type note "Meeting Notes"` is executed
- **THEN** the object SHALL be created with type `note` and name `"Meeting Notes"`

#### Scenario: Create with --type flag and no name

- **WHEN** `tmd object create --type idea` is executed
- **AND** the `idea` type has a name template
- **THEN** the object SHALL be created with the auto-generated name

#### Scenario: Create with --type flag overriding config default

- **WHEN** config has `cli.default_type: idea`
- **AND** `tmd object create --type note "Meeting Notes"` is executed
- **THEN** the object SHALL be created with type `note` (flag overrides config)

### Requirement: CLI create command type argument is optional

The `tmd object create` command SHALL accept 0 to 2 positional arguments. When the type is not provided as a positional argument, it SHALL be resolved from the `--type` flag or `cli.default_type` config.

#### Scenario: Zero args with config default type and name template

- **WHEN** `tmd object create` is executed with no arguments
- **AND** config has `cli.default_type: idea`
- **AND** the `idea` type has a name template
- **THEN** the object SHALL be created with type `idea` and auto-generated name

#### Scenario: Zero args without config or flag

- **WHEN** `tmd object create` is executed with no arguments
- **AND** no `--type` flag is provided
- **AND** no `cli.default_type` is configured
- **THEN** the command SHALL return an error indicating type is required

#### Scenario: One arg resolved as type (backward compatible)

- **WHEN** `tmd object create book` is executed
- **AND** `book` is a valid type in the vault
- **THEN** the object SHALL be created with type `book` (backward compatible behavior)

#### Scenario: One arg resolved as name with config default

- **WHEN** `tmd object create "Some Thought"` is executed
- **AND** `"Some Thought"` is NOT a valid type in the vault
- **AND** config has `cli.default_type: idea`
- **THEN** the object SHALL be created with type `idea` and name `"Some Thought"`

#### Scenario: One arg not a valid type and no default

- **WHEN** `tmd object create "Some Thought"` is executed
- **AND** `"Some Thought"` is NOT a valid type in the vault
- **AND** no `cli.default_type` is configured and no `--type` flag provided
- **THEN** the command SHALL return an error indicating unknown type

#### Scenario: Two args (backward compatible)

- **WHEN** `tmd object create book "Clean Code"` is executed
- **THEN** the object SHALL be created with type `book` and name `"Clean Code"` (unchanged behavior)
