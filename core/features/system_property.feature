Feature: System property enforcement
  typemd manages system properties (name, description, created_at, updated_at,
  tags, locked, archived) and computed properties (object_type, links, backlinks,
  created_by, updated_by) automatically. Users cannot redefine them in schemas
  or set computed properties directly.

  # ── Schema validation ────────────────────────────────────────

  Scenario: Schema validation rejects created_at property
    Given a vault is ready
    And a type schema "bad" with a system property "created_at"
    When I validate all schemas
    Then schema "bad" should have errors

  Scenario: Schema validation rejects updated_at property
    Given a vault is ready
    And a type schema "bad" with a system property "updated_at"
    When I validate all schemas
    Then schema "bad" should have errors

  Scenario: Schema validation rejects description property
    Given a vault is ready
    And a type schema "bad" with a system property "description"
    When I validate all schemas
    Then schema "bad" should have errors

  Scenario: Schema validation rejects tags property
    Given a vault is ready
    And a type schema "bad" with a system property "tags"
    When I validate all schemas
    Then schema "bad" should have errors

  Scenario: Schema validation rejects computed property object_type
    Given a vault is ready
    And a type schema "bad" with a system property "object_type"
    When I validate all schemas
    Then schema "bad" should have errors

  Scenario: Schema validation rejects computed property links
    Given a vault is ready
    And a type schema "bad" with a system property "links"
    When I validate all schemas
    Then schema "bad" should have errors

  # ── Shared property validation ───────────────────────────────

  Scenario: Shared property validation rejects description
    Given a vault is ready
    And a shared properties file with a system property "description"
    When I validate all schemas
    Then shared properties should have errors

  Scenario: Shared property validation rejects created_at
    Given a vault is ready
    And a shared properties file with a system property "created_at"
    When I validate all schemas
    Then shared properties should have errors

  Scenario: Shared property validation rejects updated_at
    Given a vault is ready
    And a shared properties file with a system property "updated_at"
    When I validate all schemas
    Then shared properties should have errors

  Scenario: Shared property validation rejects tags
    Given a vault is ready
    And a shared properties file with a system property "tags"
    When I validate all schemas
    Then shared properties should have errors

  Scenario: Shared property validation rejects computed property object_type
    Given a vault is ready
    And a shared properties file with a system property "object_type"
    When I validate all schemas
    Then shared properties should have errors

  # ── Timestamp behavior ───────────────────────────────────────

  Scenario: New object has created_at and updated_at timestamps
    Given a vault is ready
    When I create a "book" object named "test-book"
    Then the object should have a "created_at" timestamp
    And the object should have an "updated_at" timestamp

  Scenario: created_at is not modified on save
    Given a vault is ready
    And a "book" object named "test-book" exists
    When I save the object
    Then the object "created_at" should not have changed

  Scenario: SaveObject updates updated_at
    Given a vault is ready
    And a "book" object named "test-book" exists
    When I save the object
    Then the object "updated_at" should be recent

  Scenario: SetProperty updates updated_at
    Given a vault is ready
    And a "book" object named "test-book" exists
    When I set property "title" to "Test" on the object
    Then the object "updated_at" should be recent

  # ── Description behavior ─────────────────────────────────────

  Scenario: New object does not have description
    Given a vault is ready
    When I create a "book" object named "no-desc-book"
    Then the object should not have property "description"

  Scenario: Object with description preserves it
    Given a vault is ready
    And a "book" object named "desc-book" exists
    When I set property "description" to "A great book" on the object
    Then the object property "description" should be "A great book"

  # ── Computed property enforcement ────────────────────────────

  Scenario: SetProperty rejects computed properties
    Given a vault is ready
    And a "book" object named "computed-test-book" exists
    When I set property "object_type" to "page" on the object
    Then the last error should mention "computed system property"

  Scenario: Frontmatter strips computed properties on save
    Given a vault is ready
    And a raw object file with a computed property exists
    When I save the raw object
    Then the raw object file should not contain "object_type"
