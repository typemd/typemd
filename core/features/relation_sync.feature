Feature: Relation sync and prefix resolution
  The Projector syncs schema-defined relation properties from frontmatter to
  the SQLite relations table. Relation values without ULID suffixes are treated
  as name references and auto-expanded to full IDs during sync.

  Background:
    Given a vault is ready with relation sync schemas

  # ── Name resolution ──────────────────────────────────────────────────

  Scenario: Relation name reference resolves to unique match
    Given a "person" object named "john-doe" exists
    And a "book" object named "clean-code" exists with author name reference "person/john-doe"
    When I run a full relation sync
    Then the book file should have author expanded to the person's full ID

  Scenario: Relation name reference with no match is left unchanged
    Given a "book" object named "clean-code" exists with author name reference "person/nobody"
    When I run a full relation sync
    Then the book file should still have author "person/nobody"
    And the sync result should have 1 unresolved reference

  Scenario: Relation name reference with ambiguous match is left unchanged
    Given a "person" object named "john" exists
    And another "person" object named "john" exists
    And a "book" object named "clean-code" exists with author name reference "person/john"
    When I run a full relation sync
    Then the book file should still have author "person/john"
    And the sync result should have 1 unresolved reference

  Scenario: Full ID relation value is kept as-is
    Given a "person" object named "john-doe" exists
    And a "book" object named "clean-code" exists linked to the person via "author"
    When I run a full relation sync
    Then the book file should have author referencing the person's full ID

  # ── Relation sync to index ───────────────────────────────────────────

  Scenario: Single-value relation synced from frontmatter
    Given a "person" object named "john-doe" exists
    And a "book" object named "clean-code" exists linked to the person via "author"
    When I run a full relation sync
    Then listing relations for "clean-code" should return 1 entries

  Scenario: Multi-value relation synced from frontmatter
    Given a "book" object named "book-a" exists
    And a "book" object named "book-b" exists
    And a "person" object named "john-doe" exists with books references to both books
    When I run a full relation sync
    Then listing relations for "john-doe" should return 2 entries

  Scenario: Non-relation property is not synced as relation
    Given a "book" object named "clean-code" exists
    When I run a full relation sync
    Then listing relations for "clean-code" should return 0 entries

  Scenario: Relation to non-existent object is skipped
    Given a "book" object named "clean-code" exists with author "person/gone-01jjjjjjjjjjjjjjjjjjjjjjjj"
    When I run a full relation sync
    Then listing relations for "clean-code" should return 0 entries

  # ── Auto-expand write-back ───────────────────────────────────────────

  Scenario: Multiple properties expanded in one file
    Given a "person" object named "john-doe" exists
    And a "person" object named "jane-smith" exists
    And a "book" object named "clean-code" exists with author "person/john-doe" and editor "person/jane-smith"
    When I run a full relation sync
    Then the book file should have author expanded to "john-doe"'s full ID
    And the book file should have editor expanded to "jane-smith"'s full ID
    And the sync result should have 2 expansions

  # ── SyncResult reporting ─────────────────────────────────────────────

  Scenario: SyncResult reports expansion count
    Given a "person" object named "john-doe" exists
    And a "book" object named "clean-code" exists with author name reference "person/john-doe"
    When I run a full relation sync
    Then the sync result should have 1 expansions

  Scenario: SyncResult with no prefix references
    Given a "person" object named "john-doe" exists
    And a "book" object named "clean-code" exists linked to the person via "author"
    When I run a full relation sync
    Then the sync result should have 0 expansions
    And the sync result should have 0 unresolved references
