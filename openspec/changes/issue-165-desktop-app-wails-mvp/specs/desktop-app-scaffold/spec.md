## ADDED Requirements

### Requirement: Wails project structure under app/

The `app/` directory SHALL contain a Wails v3 project with a Go backend entry point and a web frontend directory.

#### Scenario: Project directory structure
- **WHEN** the `app/` directory is inspected
- **THEN** it contains `main.go`, `frontend/`, and Wails configuration files

### Requirement: AppService binds core library to frontend

The desktop app SHALL expose an `AppService` struct that wraps `core.Vault` and provides methods callable from the web frontend.

#### Scenario: AppService provides object listing
- **WHEN** the frontend calls the bound `ListObjects` method
- **THEN** it returns a list of objects from the vault with type, ID, and display name

#### Scenario: AppService initializes vault
- **WHEN** the desktop app starts with a vault path
- **THEN** the `AppService` initializes a `core.Vault` and makes it available for queries

### Requirement: macOS binary builds successfully

The project SHALL produce a runnable macOS binary via the Wails build system.

#### Scenario: Build produces executable
- **WHEN** `go build ./app/` is run
- **THEN** it produces an executable binary without errors
