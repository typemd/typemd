## ADDED Requirements

### Requirement: tmd fix wikilinks expands shorthand links

The `tmd fix wikilinks` command SHALL walk all objects in the vault, resolve shorthand wiki-links to full IDs, and write back the expanded targets to source files.

#### Scenario: Shorthand links are expanded

- **WHEN** the vault contains objects with shorthand wiki-links and `tmd fix wikilinks` is run
- **THEN** all resolvable shorthand links are replaced with full IDs in the source files

#### Scenario: Summary is displayed after fix

- **WHEN** `tmd fix wikilinks` completes
- **THEN** the output shows the count of expanded links and the count of unresolved links

#### Scenario: Ambiguous links are reported

- **WHEN** a shorthand link matches multiple objects
- **THEN** the command reports the ambiguous link with the list of matching full IDs

#### Scenario: No changes needed

- **WHEN** all wiki-links in the vault are already full IDs
- **THEN** the command reports that no changes were needed
