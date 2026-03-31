## Requirements

### Requirement: Graph outputs valid DOT digraph format

The system SHALL output a valid DOT digraph with nodes and edges. The output SHALL begin with `digraph vault {` and end with `}`. Node identifiers SHALL be the full object ID enclosed in double quotes.

#### Scenario: Empty vault produces empty digraph

- **GIVEN** a vault with no objects
- **WHEN** `tmd graph` is executed
- **THEN** the output is a valid DOT digraph with no nodes or edges

#### Scenario: Single object with no relations produces one node

- **GIVEN** a vault with one object `book/clean-code-01abc` named "Clean Code" with emoji 📖
- **WHEN** `tmd graph` is executed
- **THEN** the output contains one node `"book/clean-code-01abc"` labeled `📖 Clean Code`

### Requirement: Relations are represented as labeled directed edges

The system SHALL create a directed edge for each relation. The edge label SHALL be the relation property name. Duplicate edges (from bidirectional relations producing two index rows) SHALL be deduplicated.

#### Scenario: Relation between two objects produces an edge

- **GIVEN** object `book/clean-code-01abc` has a relation `author` to `person/bob-01def`
- **WHEN** `tmd graph` is executed
- **THEN** the output contains an edge from `"book/clean-code-01abc"` to `"person/bob-01def"` labeled `author`

#### Scenario: Bidirectional relation produces one edge per direction stored

- **GIVEN** object `person/alice-01aaa` has a bidirectional relation `friend` to `person/bob-01bbb`
- **AND** the index stores both `(alice→bob, friend)` and `(bob→alice, friend)`
- **WHEN** `tmd graph` is executed
- **THEN** the output contains exactly one edge between the two objects labeled `friend`, not two

### Requirement: Wiki-links are represented as directed edges

The system SHALL create a directed edge for each resolved wiki-link. The edge label SHALL be `wikilink`. Unresolved wiki-links (empty `ToID`) SHALL be skipped.

#### Scenario: Wiki-link produces a directed edge

- **GIVEN** object `note/ideas-01abc` contains a wiki-link to `book/clean-code-01def`
- **WHEN** `tmd graph` is executed
- **THEN** the output contains an edge from `"note/ideas-01abc"` to `"book/clean-code-01def"` labeled `wikilink`

#### Scenario: Unresolved wiki-link is skipped

- **GIVEN** object `note/ideas-01abc` contains a wiki-link with empty ToID (unresolved target)
- **WHEN** `tmd graph` is executed
- **THEN** no edge is produced for that wiki-link

### Requirement: Type filter limits graph to specific types

The `--type` flag SHALL filter the graph to include only objects of the specified types. Edges are included only when both endpoints are in the filtered set. Multiple `--type` flags SHALL be combined (union).

#### Scenario: Type filter includes only matching objects

- **GIVEN** a vault with objects of type `book` and `person`
- **WHEN** `tmd graph --type book` is executed
- **THEN** only `book` objects appear as nodes
- **AND** only edges where both endpoints are `book` objects are included

#### Scenario: Multiple type filters are combined

- **GIVEN** a vault with objects of type `book`, `person`, and `tag`
- **WHEN** `tmd graph --type book --type person` is executed
- **THEN** `book` and `person` objects appear as nodes
- **AND** edges between `book` and `person` objects are included

### Requirement: Edge type flags control which edges appear

The `--no-relations` flag SHALL exclude relation edges. The `--no-wikilinks` flag SHALL exclude wiki-link edges. Both flags can be combined (resulting in nodes only).

#### Scenario: No-wikilinks flag excludes wiki-link edges

- **GIVEN** a vault with both relations and wiki-links
- **WHEN** `tmd graph --no-wikilinks` is executed
- **THEN** only relation edges appear in the output

#### Scenario: No-relations flag excludes relation edges

- **GIVEN** a vault with both relations and wiki-links
- **WHEN** `tmd graph --no-relations` is executed
- **THEN** only wiki-link edges appear in the output

#### Scenario: Both flags produce nodes only

- **GIVEN** a vault with objects, relations, and wiki-links
- **WHEN** `tmd graph --no-relations --no-wikilinks` is executed
- **THEN** nodes appear but no edges

### Requirement: CLI command outputs to stdout

The `tmd graph` command SHALL write DOT output to stdout, enabling piping to external tools (e.g., `dot -Tpng > graph.png`). Error messages SHALL be written to stderr.

#### Scenario: Output can be piped to file

- **WHEN** `tmd graph` is executed with stdout redirected
- **THEN** the redirected output is valid DOT format
