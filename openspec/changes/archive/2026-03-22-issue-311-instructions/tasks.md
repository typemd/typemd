## 1. Core: Embedded Skill Registry

- [x] 1.1 Write BDD scenarios for skill listing, loading, and frontmatter parsing (`core/features/`)
- [x] 1.2 Copy marketplace SKILL.md files to `core/skills/` and embed via `//go:embed`
- [x] 1.3 Implement `core/instructions.go`: Skill struct, SKILL.md parser, SkillRegistry with List/Get/GetRaw
- [x] 1.4 Implement BDD step definitions to make scenarios pass
- [x] 1.5 Add unit tests for edge cases (no frontmatter, unknown skill, malformed YAML)

## 2. Core: Vault Context Builder

- [x] 2.1 Write BDD scenarios for vault context building (type summaries with properties)
- [x] 2.2 Implement context builder: BuildSkillContext returning TypeSummary structs from Vault
- [x] 2.3 Implement BDD step definitions to make scenarios pass
- [x] 2.4 Add unit tests for edge cases (empty vault, built-in types only)

## 3. Core: Vault Override

- [x] 3.1 Write BDD scenarios for vault override (override present, override without frontmatter)
- [x] 3.2 Implement override loading in SkillRegistry: check `.typemd/instructions/<skill>.md`
- [x] 3.3 Implement BDD step definitions to make scenarios pass

## 4. CLI: Instructions Command

- [x] 4.1 Write BDD scenarios for `tmd instructions` CLI command (`cmd/features/`)
- [x] 4.2 Implement `cmd/instructions.go`: list mode, skill output mode, --skill flag, --json flag
- [x] 4.3 Implement BDD step definitions to make scenarios pass
- [x] 4.4 Add unit tests for error output (unknown skill lists available skills)
