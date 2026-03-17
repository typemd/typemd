Feature: Type Schema Version
  Type schemas support a version field for tracking schema evolution.

  # ── Serialization ──────────────────────────────────────────

  Scenario: Schema with version serializes to YAML
    Given a versioned type schema "book" with version 1
    When I serialize the versioned schema
    Then the versioned YAML output should contain "version: 1"

  Scenario: Schema without version omits version from YAML
    Given a versioned type schema "note" with version 0
    When I serialize the versioned schema
    Then the versioned YAML output should not contain "version:"

  Scenario: Schema version round-trips through marshal/unmarshal
    Given a versioned type schema "book" with version 3
    When I serialize the versioned schema
    And I deserialize the versioned YAML output
    Then the round-trip schema version should be 3

  Scenario: Schema without version defaults to 0 on load
    Given a versioned type schema "note" with version 0
    When I serialize the versioned schema
    And I deserialize the versioned YAML output
    Then the round-trip schema version should be 0

  # ── Validation ─────────────────────────────────────────────

  Scenario: Schema with positive version passes validation
    Given a versioned type schema "book" with version 1
    When I validate the versioned schema
    Then no schema validation errors should occur

  Scenario: Schema with zero version passes validation
    Given a versioned type schema "book" with version 0
    When I validate the versioned schema
    Then no schema validation errors should occur

  Scenario: Schema with negative version fails validation
    Given a versioned type schema "book" with version -1
    When I validate the versioned schema
    Then a schema validation error should mention "non-negative"
