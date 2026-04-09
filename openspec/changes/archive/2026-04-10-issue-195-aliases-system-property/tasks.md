## 1. Core: BDD Scenarios for aliases system property

- [x] 1.1 Write BDD scenarios for `aliases` system property (stored, user-authored, reserved name, frontmatter order)
- [x] 1.2 Write BDD scenarios for alias-based name index lookup and wiki-link resolution by alias
- [x] 1.3 Implement step definitions for aliases scenarios (initially failing)

## 2. Core: Register aliases in system property registry

- [x] 2.1 Add `AliasesProperty = "aliases"` constant to `system_property.go`
- [x] 2.2 Register `aliases` in `systemProperties` slice as `list[text]` stored user-authored property, positioned after `tags`
- [x] 2.3 Add unit tests for `IsSystemProperty("aliases")` and `StoredPropertyNames()` order

## 3. Core: aliases frontmatter serialization

- [x] 3.1 Verify `writeFrontmatter` / `parseFrontmatter` correctly round-trips `aliases` as `[]string` via existing `keyOrder` mechanism (no struct change needed — uses `map[string]any`)
- [x] 3.2 Add unit tests for round-trip: aliases present, absent, and empty array treated as absent

## 4. Core: aliases indexed for name resolution

- [x] 4.1 Update `buildNameIndex` in `name_resolve.go` to index each alias from `obj.Properties["aliases"]` as additional lookup keys (slugified)
- [x] 4.2 Add unit tests for alias index building: single alias, multiple aliases, duplicate alias across objects (ambiguous)

## 5. Core: wiki-link resolution by alias

- [x] 5.1 Confirm `resolveWikiLinkTarget` uses `resolveByName` which already reads the name index — no change needed beyond step 4.1
- [x] 5.2 Add unit tests for wiki-link resolution: resolves by alias, alias lower priority than exact name match (covered by ambiguity in name index)
- [x] 5.3 Make BDD scenarios pass (alias resolution end-to-end via reconcile → name index → wikilink)

## 6. Core: schema validation rejects aliases as property name

- [x] 6.1 Confirm `type_schema_validate.go` already rejects `aliases` via `IsSystemProperty` check (since it's now registered) — verify with existing test
- [x] 6.2 Confirm `shared_properties.go` `ValidateSharedProperties` also rejects `aliases` via `IsSystemProperty` check — verify with existing test
- [x] 6.3 Add BDD or unit test scenarios for schema rejection of `aliases` property name

## 7. Verify and clean up

- [x] 7.1 Run `go test ./core/...` — all pass
- [x] 7.2 Run `go vet ./core/...` — clean
