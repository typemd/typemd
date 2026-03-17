Feature: Type Schema Version
  Type schemas support a version field for tracking schema evolution.

  # ── Serialization ──────────────────────────────────────────

  Scenario: Schema with version serializes to YAML
    Given a type schema "book" with no extra fields
    And the schema has version 1
    When I serialize the type schema
    Then the YAML output should contain "version: 1"

  Scenario: Schema without version omits version from YAML
    Given a type schema "note" with no extra fields
    When I serialize the type schema
    Then the YAML output should not contain "version:"

  Scenario: Schema version round-trips through marshal/unmarshal
    Given a type schema "book" with no extra fields
    And the schema has version 3
    When I serialize the type schema
    And I deserialize the YAML output back to a TypeSchema
    Then the round-trip schema version should be 3

  Scenario: Schema without version defaults to 0 on load
    Given a type schema "note" with no extra fields
    When I serialize the type schema
    And I deserialize the YAML output back to a TypeSchema
    Then the round-trip schema version should be 0

  # ── Validation ─────────────────────────────────────────────

  Scenario: Schema with positive version passes validation
    Given a type schema "book" with no extra fields
    And the schema has version 1
    When I validate the type schema
    Then no schema validation errors should occur

  Scenario: Schema with zero version passes validation
    Given a type schema "book" with no extra fields
    When I validate the type schema
    Then no schema validation errors should occur

  Scenario: Schema with negative version fails validation
    Given a type schema "book" with no extra fields
    And the schema has version -1
    When I validate the type schema
    Then a schema validation error should mention "non-negative"
