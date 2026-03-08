## Why

Users currently must use `go install` or download binaries manually from GitHub Releases to install typemd. Homebrew is the standard package manager on macOS, and supporting it provides a familiar, zero-friction installation path (`brew install typemd/tap/typemd-cli`). With GoReleaser and a release workflow already in place (#39), the infrastructure to automate this is ready.

## What Changes

- Create a `typemd/homebrew-tap` GitHub repository to host the Homebrew formula
- Add `brews` configuration to `.goreleaser.yml` so GoReleaser auto-publishes the formula on each release
- Configure a fine-grained PAT as a repository secret to allow cross-repo writes
- Update installation documentation (README, website) with Homebrew instructions

## Capabilities

### New Capabilities

- `homebrew-distribution`: Automated Homebrew formula generation and publishing via GoReleaser tap integration

### Modified Capabilities

None.

## Impact

- **New repository**: `typemd/homebrew-tap` on GitHub
- **Modified file**: `.goreleaser.yml` — new `brews` section
- **Modified file**: `.github/workflows/release.yml` — pass new token to GoReleaser
- **Secret**: New `HOMEBREW_TAP_TOKEN` repository secret
- **Documentation**: README and website install instructions updated
