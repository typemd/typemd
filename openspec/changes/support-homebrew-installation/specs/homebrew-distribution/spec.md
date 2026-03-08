## ADDED Requirements

### Requirement: GoReleaser publishes Homebrew formula on release

When a new version is released via `v*` tag, GoReleaser SHALL automatically generate and push a Homebrew formula to the `typemd/homebrew-tap` repository.

#### Scenario: New release triggers formula update

- **WHEN** a `v*` tag is pushed and the release workflow completes
- **THEN** GoReleaser pushes an updated `Formula/typemd-cli.rb` to `typemd/homebrew-tap` with the correct version, URLs, and SHA256 checksums

### Requirement: Users can install typemd CLI via Homebrew

Users SHALL be able to install the `tmd` binary using Homebrew with the custom tap.

#### Scenario: First-time installation

- **WHEN** a user runs `brew install typemd/tap/typemd-cli`
- **THEN** Homebrew installs the `tmd` binary and it is available in the user's PATH

#### Scenario: Upgrade to new version

- **WHEN** a new release is published and user runs `brew upgrade typemd-cli`
- **THEN** the `tmd` binary is updated to the latest version

### Requirement: Formula metadata is accurate

The generated formula SHALL include correct metadata for discoverability.

#### Scenario: Formula info displays correct details

- **WHEN** a user runs `brew info typemd-cli`
- **THEN** the output shows the homepage as `https://typemd.io`, a description of "A local-first CLI knowledge management tool", and the current version

### Requirement: Cross-repo token is configured

The release workflow SHALL use a dedicated token with minimal permissions to push to the tap repository.

#### Scenario: Token has least-privilege access

- **WHEN** the release workflow runs GoReleaser
- **THEN** GoReleaser uses the `HOMEBREW_TAP_TOKEN` secret which has `contents: write` permission scoped only to `typemd/homebrew-tap`
