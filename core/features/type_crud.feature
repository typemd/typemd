Feature: Type CRUD
  Type schemas can be saved, deleted, loaded, and their objects counted.

  # ── DeleteSchema ───────────────────────────────────────────

  Scenario: Delete an existing type schema file
    Given a vault is ready
    And a type schema file "scratch" exists on disk
    When I delete schema "scratch"
    Then no error should occur
    And the type schema file "scratch" should not exist on disk

  Scenario: Delete a non-existent type schema file returns error
    Given a vault is ready
    When I delete schema "nonexistent"
    Then an error should occur

  # ── SaveType ───────────────────────────────────────────────

  Scenario: Save a valid type schema
    Given a vault is ready
    And a type schema "project" with no extra fields
    And the schema has a "status" select property with options "active,done"
    When I save the type schema
    Then no error should occur
    And the type schema file "project" should exist on disk
    And loading type "project" should return a schema with 1 property

  Scenario: Save fails on invalid schema
    Given a vault is ready
    And a type schema "" with no extra fields
    When I save the type schema
    Then an error should occur

  Scenario: Save overwrites existing type schema
    Given a vault is ready
    And a type schema "draft" with no extra fields
    When I save the type schema
    And I add a "priority" number property to the schema
    And I save the type schema
    Then no error should occur
    And loading type "draft" should return a schema with 1 property

  # ── DeleteType ─────────────────────────────────────────────

  Scenario: Delete a user-defined type
    Given a vault is ready
    And a type schema file "scratch" exists on disk
    When I delete type "scratch"
    Then no error should occur
    And the type schema file "scratch" should not exist on disk

  Scenario: Delete built-in type is rejected
    Given a vault is ready
    When I delete type "tag"
    Then an error should occur
    And the error message should contain "cannot delete built-in type"

  Scenario: Delete non-existent type returns error
    Given a vault is ready
    When I delete type "phantom"
    Then an error should occur

  # ── CountObjectsByType ─────────────────────────────────────

  Scenario: Count objects for type with objects
    Given a vault is ready
    And a "book" object named "go-book" exists
    And a "book" object named "rust-book" exists
    When I count objects of type "book"
    Then the count should be 2

  Scenario: Count objects for type with no objects
    Given a vault is ready
    When I count objects of type "project"
    Then the count should be 0

  # ── LoadType ──────────────────────────────────────────────

  Scenario: Custom type emoji overrides built-in default
    Given a vault is ready
    And a custom "tag" type schema with emoji "🔖"
    When I load type "tag"
    Then no error should occur
    And the loaded schema should have emoji "🔖"

  Scenario: Custom tag schema without unique field defaults to false
    Given a vault is ready
    And a custom tag type schema without unique field
    When I load type "tag"
    Then no error should occur
    And the loaded schema should have unique false

  Scenario: Loading undefined type returns error
    Given a vault is ready
    When I load type "nonexistent"
    Then an error should occur

  Scenario: Tag type loads with default emoji and unique
    Given a vault is ready
    When I load type "tag"
    Then no error should occur
    And the loaded schema should have emoji "🏷️"
    And the loaded schema should have unique true
