## 1. SystemProperty: Add Immutable field

- [x] 1.1 Write BDD scenarios for system property immutability (registry contains immutable flags, IsImmutableSystemProperty behavior)
- [x] 1.2 Implement step definitions for immutability scenarios
- [x] 1.3 Add `Immutable bool` field to `SystemProperty` struct, set `created_at` and `updated_at` to `Immutable: true`
- [x] 1.4 Add `IsImmutableSystemProperty(name)` function
- [x] 1.5 Add unit tests for `IsImmutableSystemProperty` edge cases

## 2. Core: Template path helpers

- [x] 2.1 Add `TemplatesDir()`, `TypeTemplatesDir(typeName)`, and `TemplatePath(typeName, templateName)` methods to `Vault`
- [x] 2.2 Add unit tests for template path helpers

## 3. Core: ListTemplates

- [x] 3.1 Write BDD scenarios for template discovery (multiple templates, single template, no directory, empty directory)
- [x] 3.2 Implement step definitions for ListTemplates scenarios
- [x] 3.3 Implement `ListTemplates(typeName)` method on `Vault`
- [x] 3.4 Add unit tests for ListTemplates edge cases (non-.md files ignored, nested directories ignored)

## 4. Core: LoadTemplate

- [x] 4.1 Write BDD scenarios for template loading (frontmatter + body, body only, frontmatter only, not found)
- [x] 4.2 Implement step definitions for LoadTemplate scenarios
- [x] 4.3 Implement `LoadTemplate(typeName, templateName)` method on `Vault`

## 5. Core: NewObject template application

- [x] 5.1 Write BDD scenarios for template application (apply template, override schema default, immutable system props ignored, mutable system props applied, unknown props ignored, no template)
- [x] 5.2 Implement step definitions for NewObject template scenarios
- [x] 5.3 Extend `NewObject` signature to accept `templateName` parameter
- [x] 5.4 Implement template loading and property merge logic in `NewObject`
- [x] 5.5 Update all existing `NewObject` call sites to pass empty `templateName`

## 6. CLI: Template flag and selection

- [x] 6.1 Add `-t` / `--template` flag to `tmd object create` command
- [x] 6.2 Implement auto-apply logic for single template
- [x] 6.3 Implement interactive template selection for multiple templates
- [x] 6.4 Pass selected template name to `NewObject`
- [x] 6.5 Add unit tests for CLI template flag parsing
