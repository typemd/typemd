## 1. Core: Built-in page type

- [x] 1.1 Write BDD scenarios for page type (load, deletion protection, listing)
- [x] 1.2 Implement BDD step definitions for page type scenarios
- [x] 1.3 Add `PageTypeName` constant and `page` entry in `defaultTypes`
- [x] 1.4 Add unit tests for page type edge cases (custom override, no custom properties)

## 2. CLI: Default type fallback

- [x] 2.1 Remove `resolveDefaultType()`, always use `page` as default type in `tmd init`
- [x] 2.2 Remove obsolete `TestResolveDefaultType` unit test

## 3. Documentation

- [x] 3.1 Update type schema docs and README to mention `page` as built-in type
