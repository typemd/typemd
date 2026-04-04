## ADDED Requirements

### Requirement: Vault uses root-level types directory
The system SHALL resolve the types directory as `<vault_root>/types/` instead of `<vault_root>/.typemd/types/`.

#### Scenario: TypesDir returns root-level path
- **WHEN** a vault is initialized or opened
- **THEN** `Vault.TypesDir()` SHALL return `<vault_root>/types/`

### Requirement: Vault uses root-level properties directory
The system SHALL resolve the shared properties file as `<vault_root>/properties/properties.yaml` instead of `<vault_root>/.typemd/properties.yaml`.

#### Scenario: SharedPropertiesPath returns root-level path
- **WHEN** a vault is initialized or opened
- **THEN** `Vault.SharedPropertiesPath()` SHALL return `<vault_root>/properties/properties.yaml`

### Requirement: Init creates directories at vault root
The system SHALL create `types/` and `properties/` directories at vault root during initialization.

#### Scenario: tmd init creates root-level directories
- **WHEN** user runs `tmd init`
- **THEN** the system SHALL create `<vault_root>/types/` directory
- **AND** the system SHALL create `<vault_root>/properties/` directory
- **AND** the system SHALL NOT create `.typemd/types/` directory

### Requirement: Auto-migration of types directory
The system SHALL automatically move `.typemd/types/` to `types/` at vault root when opening a vault that has the old layout.

#### Scenario: Migrate types from old to new location
- **WHEN** `.typemd/types/` exists
- **AND** `types/` does NOT exist at vault root
- **THEN** the system SHALL move `.typemd/types/` to `<vault_root>/types/`
- **AND** the old `.typemd/types/` directory SHALL no longer exist

#### Scenario: No migration when already at new location
- **WHEN** `types/` exists at vault root
- **AND** `.typemd/types/` does NOT exist
- **THEN** the system SHALL NOT perform any migration for types

#### Scenario: Skip migration when old types directory does not exist
- **WHEN** `.typemd/types/` does NOT exist
- **AND** `types/` does NOT exist at vault root
- **THEN** the system SHALL NOT perform any migration for types
- **AND** the system SHALL NOT return an error

### Requirement: Auto-migration of properties file
The system SHALL automatically move `.typemd/properties.yaml` to `properties/properties.yaml` at vault root when opening a vault that has the old layout.

#### Scenario: Migrate properties from old to new location
- **WHEN** `.typemd/properties.yaml` exists
- **AND** `properties/properties.yaml` does NOT exist at vault root
- **THEN** the system SHALL create `<vault_root>/properties/` directory if needed
- **AND** the system SHALL move `.typemd/properties.yaml` to `<vault_root>/properties/properties.yaml`

#### Scenario: No migration when already at new location
- **WHEN** `properties/properties.yaml` exists at vault root
- **AND** `.typemd/properties.yaml` does NOT exist
- **THEN** the system SHALL NOT perform any migration for properties

#### Scenario: Skip migration when old properties file does not exist
- **WHEN** `.typemd/properties.yaml` does NOT exist
- **AND** `properties/properties.yaml` does NOT exist at vault root
- **THEN** the system SHALL NOT perform any migration for properties
- **AND** the system SHALL NOT return an error

### Requirement: Conflict detection for types
The system SHALL return an error if both old and new types directories exist simultaneously.

#### Scenario: Both types directories exist
- **WHEN** `.typemd/types/` exists
- **AND** `types/` exists at vault root
- **THEN** the system SHALL return an error indicating a conflict
- **AND** the error message SHALL instruct the user to resolve the conflict manually

### Requirement: Conflict detection for properties
The system SHALL return an error if both old and new properties paths exist simultaneously.

#### Scenario: Both properties paths exist
- **WHEN** `.typemd/properties.yaml` exists
- **AND** `properties/properties.yaml` exists at vault root
- **THEN** the system SHALL return an error indicating a conflict
- **AND** the error message SHALL instruct the user to resolve the conflict manually

### Requirement: Pre-check all paths before migration
The system SHALL check all paths for conflicts before moving anything, to avoid partial migration.

#### Scenario: Conflict in one path prevents all migration
- **WHEN** `.typemd/types/` exists and `types/` also exists (conflict)
- **AND** `.typemd/properties.yaml` exists and needs migration
- **THEN** the system SHALL return a conflict error
- **AND** the system SHALL NOT move the properties file

### Requirement: Migration removes empty old directories
After successful migration, the system SHALL clean up empty parent directories left behind.

#### Scenario: Old types directory removed after migration
- **WHEN** `.typemd/types/` is successfully moved to `types/`
- **THEN** `.typemd/types/` SHALL no longer exist
