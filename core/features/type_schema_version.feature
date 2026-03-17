Feature: Type Schema Version
  Type schemas support a semver-style "major.minor" version field for tracking schema evolution.

  # ── Serialization ──────────────────────────────────────────

  Scenario: Schema with version serializes to YAML
    Given a type schema "book" with no extra fields
    And the schema has version "1.0"
    When I serialize the type schema
    Then the YAML output should contain "version: \"1.0\""

  Scenario: Schema without version omits version from YAML
    Given a type schema "note" with no extra fields
    When I serialize the type schema
    Then the YAML output should not contain "version:"

  Scenario: Schema version round-trips through marshal/unmarshal
    Given a type schema "book" with no extra fields
    And the schema has version "2.3"
    When I serialize the type schema
    And I deserialize the YAML output back to a TypeSchema
    Then the round-trip schema version should be "2.3"

  Scenario: Schema without version defaults to 0.0 on load
    Given a type schema "note" with no extra fields
    When I serialize the type schema
    And I deserialize the YAML output back to a TypeSchema
    Then the round-trip schema version should be "0.0"

  # ── Validation ─────────────────────────────────────────────

  Scenario: Schema with valid version passes validation
    Given a type schema "book" with no extra fields
    And the schema has version "2.3"
    When I validate the type schema
    Then no schema validation errors should occur

  Scenario: Schema with zero version passes validation
    Given a type schema "book" with no extra fields
    When I validate the type schema
    Then no schema validation errors should occur

  Scenario: Single number version fails validation
    Given a type schema "book" with no extra fields
    And the schema has version "1"
    When I validate the type schema
    Then a schema validation error should mention "major.minor"

  Scenario: Three segments version fails validation
    Given a type schema "book" with no extra fields
    And the schema has version "1.0.0"
    When I validate the type schema
    Then a schema validation error should mention "major.minor"

  Scenario: Leading zeros version fails validation
    Given a type schema "book" with no extra fields
    And the schema has version "01.0"
    When I validate the type schema
    Then a schema validation error should mention "major.minor"

  Scenario: Negative number version fails validation
    Given a type schema "book" with no extra fields
    And the schema has version "-1.0"
    When I validate the type schema
    Then a schema validation error should mention "major.minor"

  Scenario: Non-numeric version fails validation
    Given a type schema "book" with no extra fields
    And the schema has version "abc"
    When I validate the type schema
    Then a schema validation error should mention "major.minor"

  # ── Comparison ─────────────────────────────────────────────

  Scenario: Higher major version is greater
    Then comparing version "2.0" with "1.3" should return 1

  Scenario: Higher minor version is greater
    Then comparing version "1.2" with "1.1" should return 1

  Scenario: Equal versions
    Then comparing version "1.1" with "1.1" should return 0

  Scenario: Lower version is less
    Then comparing version "0.1" with "1.0" should return -1
