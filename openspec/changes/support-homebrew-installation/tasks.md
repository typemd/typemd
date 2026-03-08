## 1. Repository Setup

- [x] 1.1 Create `typemd/homebrew-tap` public repository on GitHub with a README
- [x] 1.2 Create `Formula/` directory in the tap repo

## 2. GoReleaser Configuration

- [x] 2.1 Add `brews` section to `.goreleaser.yml` with formula name `typemd-cli`, homepage, description, and tap repository target
- [x] 2.2 Verify GoReleaser config is valid with `goreleaser check`

## 3. Token & Secret Setup

- [x] 3.1 Create a fine-grained PAT scoped to `typemd/homebrew-tap` with `contents: write` permission
- [x] 3.2 Add the PAT as `HOMEBREW_TAP_TOKEN` repository secret on `typemd/typemd`
- [x] 3.3 Update `.github/workflows/release.yml` to pass `HOMEBREW_TAP_TOKEN` to GoReleaser

## 4. Documentation

- [x] 4.1 Update README with Homebrew installation instructions
- [x] 4.2 Update website (`websites/site/src/pages/index.astro`) install block with Homebrew option
