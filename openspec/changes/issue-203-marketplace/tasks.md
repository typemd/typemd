## 1. Marketplace Structure

- [x] 1.1 Create `marketplace/.claude-plugin/marketplace.json` with marketplace metadata and empty plugins array
- [x] 1.2 Create `marketplace/README.md` with marketplace overview and installation instructions
- [x] 1.3 Create `marketplace/CONTRIBUTING.md` with submission process, naming rules, plugin structure, and quality requirements

## 2. Example Plugin: markdown-import

- [x] 2.1 Create plugin directory and `.claude-plugin/plugin.json` with metadata
- [x] 2.2 Create `skills/markdown-import/SKILL.md` with the skill content (vault context awareness, frontmatter generation, type detection, relation discovery, batch guidance)
- [x] 2.3 Create plugin `README.md` with usage examples and prerequisites
- [x] 2.4 Register the plugin in `marketplace.json` plugins array

## 3. CI Validation

- [x] 3.1 ~~Create validation shell script~~ → Use `claude plugin validate .` instead (removed custom script)
- [x] 3.2 Create `.github/workflows/validate-marketplace.yml` triggered on PRs modifying `marketplace/`

## 4. Subtree Publish

- [x] 4.1 Create `.github/workflows/publish-marketplace.yml` to subtree push `marketplace/` to `typemd/marketplace` on merge to main
