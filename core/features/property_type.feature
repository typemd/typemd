Feature: Property Type System
  Objects have typed properties that are validated against their type schema.
  The system supports 9 property types: string, number, date, datetime, url,
  checkbox, select, multi_select, and relation.

  Background:
    Given a vault is ready

  # ── Schema validation ────────────────────────────────────────────────────

  Scenario: All 9 property types are accepted in a schema
    Given a type schema with all 9 property types
    When I validate all schemas
    Then schema "complete" should have no errors

  Scenario: Enum type is rejected with migration guidance
    Given a type schema "legacy" with an enum property
    When I validate all schemas
    Then schema "legacy" should have errors

  Scenario: Select type requires options
    Given a type schema "bad" with a select property missing options
    When I validate all schemas
    Then schema "bad" should have errors

  # ── Object validation: date ──────────────────────────────────────────────

  Scenario: Valid date property passes validation
    Given a type schema "event" with a date property
    And an "event" object named "birthday" exists with raw property "date" set to "2026-01-15"
    When I validate the object against its schema
    Then the object should have no validation errors

  Scenario: Invalid date format is rejected
    Given a type schema "event" with a date property
    And an "event" object named "birthday" exists with raw property "date" set to "01/15/2026"
    When I validate the object against its schema
    Then the object should have validation errors

  # ── Object validation: url ───────────────────────────────────────────────

  Scenario: Valid URL passes validation
    Given a type schema "bookmark" with a url property
    And a "bookmark" object named "example" exists with raw property "link" set to "https://example.com"
    When I validate the object against its schema
    Then the object should have no validation errors

  Scenario: URL without http scheme is rejected
    Given a type schema "bookmark" with a url property
    And a "bookmark" object named "example" exists with raw property "link" set to "ftp://example.com"
    When I validate the object against its schema
    Then the object should have validation errors

  # ── Object validation: select ────────────────────────────────────────────

  Scenario: Valid select value passes validation
    Given a type schema "book" with a select status property
    And a "book" object named "test" exists with raw property "status" set to "reading"
    When I validate the object against its schema
    Then the object should have no validation errors

  Scenario: Invalid select value is rejected
    Given a type schema "book" with a select status property
    And a "book" object named "test" exists with raw property "status" set to "unknown"
    When I validate the object against its schema
    Then the object should have validation errors

  # ── Schema migration: enum → select ─────────────────────────────────────

  Scenario: Enum schemas are migrated to select
    Given a type schema "book" with an enum property
    When I migrate schemas
    Then the "book" schema should use select instead of enum
