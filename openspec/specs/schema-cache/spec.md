### Requirement: Vault caches type schemas in memory

The Vault SHALL maintain an in-memory cache of type schemas. `LoadType()` SHALL return the cached schema if available, otherwise load from the repository and cache the result.

#### Scenario: First load populates cache

- **WHEN** `LoadType("book")` is called and no cache entry exists for "book"
- **THEN** the Vault SHALL load from `ObjectRepository.GetSchema("book")` and cache the result

#### Scenario: Subsequent loads use cache

- **WHEN** `LoadType("book")` is called and a cache entry exists for "book"
- **THEN** the Vault SHALL return the cached schema without reading from disk

### Requirement: Schema cache invalidation on type mutation

The schema cache SHALL be invalidated when type schemas are modified through the Vault API.

#### Scenario: SaveType invalidates cache entry

- **WHEN** `SaveType()` is called for type "book"
- **THEN** the cache entry for "book" SHALL be invalidated
- **AND** the next `LoadType("book")` SHALL read from disk

#### Scenario: DeleteType invalidates cache entry

- **WHEN** `DeleteType()` is called for type "book"
- **THEN** the cache entry for "book" SHALL be invalidated

#### Scenario: MigrateSchemas invalidates entire cache

- **WHEN** `MigrateSchemas()` is called
- **THEN** all cache entries SHALL be invalidated

### Requirement: Schema cache invalidation on external file change

The schema cache SHALL be invalidated when `.typemd/types/` files are modified externally (outside the Vault API).

#### Scenario: External schema edit invalidates cache

- **WHEN** the watcher detects a file change in `.typemd/types/`
- **THEN** the entire schema cache SHALL be invalidated
