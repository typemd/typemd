Feature: Type directory structure
  Type schemas use directory format (types/<name>/schema.yaml).

  # ── Load from directory format ──────────────────────────────

  Scenario: Load type schema from directory format
    Given a vault is ready
    And a type schema directory "movie" with schema content:
      """
      name: movie
      emoji: "\U0001F3AC"
      properties:
        - name: rating
          type: number
      """
    When I load type "movie"
    Then no error should occur
    And the loaded schema should have emoji "🎬"
    And the loaded schema should have 1 property

  # ── ListTypes with directory format ───────────────────────

  Scenario: ListTypes discovers directory format
    Given a vault is ready
    And a type schema directory "movie" with schema content:
      """
      name: movie
      properties: []
      """
    When I list all types
    Then the type list should contain "movie"

  # ── SaveType writes directory format ────────────────────────

  Scenario: SaveType creates directory format
    Given a vault is ready
    And a type schema "article" with no extra fields
    And the schema has a "title" string property
    When I save the type schema
    Then no error should occur
    And the type schema directory "article" should exist

  # ── DeleteType removes directory ────────────────────────────

  Scenario: Delete type removes entire directory
    Given a vault is ready
    And a type schema directory "scratch" with schema content:
      """
      name: scratch
      properties: []
      """
    When I delete type "scratch"
    Then no error should occur
    And the type schema directory "scratch" should not exist
