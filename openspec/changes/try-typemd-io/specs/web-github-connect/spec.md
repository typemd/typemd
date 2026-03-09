## ADDED Requirements

### Requirement: Landing page with connection form
The landing page SHALL display a form with a repository URL input and an optional Personal Access Token field.

#### Scenario: Connect to public repo without token
- **WHEN** the user enters a public repo (e.g., `user/my-vault`) and clicks Connect without providing a token
- **THEN** the system SHALL fetch the vault data using unauthenticated GitHub API calls and display the vault browser

#### Scenario: Connect to private repo with token
- **WHEN** the user enters a private repo and a valid PAT, then clicks Connect
- **THEN** the system SHALL use the PAT for authentication and display the vault browser

#### Scenario: Invalid repo
- **WHEN** the user enters a repo that does not exist or is inaccessible
- **THEN** the system SHALL display an error message indicating the repo could not be found

#### Scenario: Invalid token
- **WHEN** the user enters a valid repo but an invalid PAT
- **THEN** the system SHALL display an error message indicating authentication failed

#### Scenario: Repo without typemd vault
- **WHEN** the user connects to a repo that has no `.typemd/` directory
- **THEN** the system SHALL display an error message indicating this is not a typemd vault

### Requirement: Token persistence with user consent
The system SHALL support optionally storing the PAT in localStorage when the user opts in.

#### Scenario: Remember token opted in
- **WHEN** the user checks "Remember token on this device" and connects successfully
- **THEN** the PAT SHALL be stored in localStorage and pre-filled on next visit

#### Scenario: Remember token opted out
- **WHEN** the user does not check "Remember token on this device"
- **THEN** the PAT SHALL be stored in memory only and lost when the page is closed

#### Scenario: Clear stored token
- **WHEN** the user clicks "Clear token" or "Disconnect"
- **THEN** the PAT SHALL be removed from both memory and localStorage

### Requirement: Connection status display
The system SHALL display the current connection status and provide a way to disconnect.

#### Scenario: Show connected repo
- **WHEN** the user is connected to a vault
- **THEN** the UI SHALL display the connected repo name and a disconnect button

#### Scenario: Disconnect
- **WHEN** the user clicks disconnect
- **THEN** the system SHALL clear the token from memory (and localStorage if stored), and return to the landing page

### Requirement: GitHub API rate limit awareness
The system SHALL monitor and display GitHub API rate limit status.

#### Scenario: Show remaining rate limit
- **WHEN** the user is connected
- **THEN** the UI SHALL display the remaining API calls (from `X-RateLimit-Remaining` response header)

#### Scenario: Rate limit exceeded
- **WHEN** the GitHub API returns a 403 rate limit error
- **THEN** the system SHALL display a message indicating the rate limit is exceeded and show when it resets
