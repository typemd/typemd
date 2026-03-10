Feature: Shared Properties
  Users can define shared property definitions in .typemd/properties.yaml.
  Type schemas can reference these via `use` to reuse property definitions
  across multiple types without duplication.

  Background:
    Given a vault is ready

  # ── Loading shared properties ─────────────────────────────────────────────

  Scenario: Load shared properties file with multiple properties
    Given a shared properties file with "due_date" date and "priority" select properties
    When I load shared properties
    Then shared properties should contain 2 entries
    And shared property "due_date" should have type "date"
    And shared property "priority" should have type "select"

  Scenario: Shared properties file does not exist
    When I load shared properties
    Then shared properties should contain 0 entries

  Scenario: Empty shared properties file
    Given an empty shared properties file
    When I load shared properties
    Then shared properties should contain 0 entries

  # ── Shared properties validation ──────────────────────────────────────────

  Scenario: Valid shared properties pass validation
    Given a shared properties file with "due_date" date and "priority" select properties
    When I validate all schemas
    Then shared properties should have no errors

  Scenario: Duplicate shared property names are rejected
    Given a shared properties file with duplicate "due_date" properties
    When I validate all schemas
    Then shared properties should have errors

  Scenario: Invalid property type in shared properties is rejected
    Given a shared properties file with an invalid property type
    When I validate all schemas
    Then shared properties should have errors

  Scenario: Reserved name "name" in shared properties is rejected
    Given a shared properties file with a property named "name"
    When I validate all schemas
    Then shared properties should have errors

  Scenario: Select without options in shared properties is rejected
    Given a shared properties file with a select property missing options
    When I validate all schemas
    Then shared properties should have errors

  # ── Use keyword in type schemas ───────────────────────────────────────────

  Scenario: Type schema with use entry references shared property
    Given a shared properties file with "due_date" date and "priority" select properties
    And a type schema "project" with use "due_date"
    When I validate all schemas
    Then schema "project" should have no errors

  Scenario: Use with pin override is accepted
    Given a shared properties file with "due_date" date and "priority" select properties
    And a type schema "project" with use "due_date" and pin 1
    When I validate all schemas
    Then schema "project" should have no errors

  Scenario: Use with emoji override is accepted
    Given a shared properties file with "due_date" date and "priority" select properties
    And a type schema "project" with use "due_date" and emoji "🗓️"
    When I validate all schemas
    Then schema "project" should have no errors

  Scenario: Use with disallowed type field is rejected
    Given a shared properties file with "due_date" date and "priority" select properties
    And a type schema "project" with use "due_date" and disallowed type override
    When I validate all schemas
    Then schema "project" should have errors

  Scenario: Use referencing non-existent shared property is rejected
    Given a shared properties file with "due_date" date and "priority" select properties
    And a type schema "project" with use "nonexistent"
    When I validate all schemas
    Then schema "project" should have errors

  Scenario: Local property name conflicting with shared property is rejected
    Given a shared properties file with "due_date" date and "priority" select properties
    And a type schema "project" with local property "due_date"
    When I validate all schemas
    Then schema "project" should have errors

  Scenario: Duplicate use entries are rejected
    Given a shared properties file with "due_date" date and "priority" select properties
    And a type schema "project" with duplicate use "due_date"
    When I validate all schemas
    Then schema "project" should have errors

  Scenario: Use and name on same entry are rejected
    Given a shared properties file with "due_date" date and "priority" select properties
    And a type schema "project" with both use and name on same entry
    When I validate all schemas
    Then schema "project" should have errors

  # ── LoadType resolution ───────────────────────────────────────────────────

  Scenario: LoadType resolves use entry with no overrides
    Given a shared properties file with "due_date" date and "priority" select properties
    And a type schema "project" with use "due_date"
    When I load type "project"
    Then the loaded type should have 1 property
    And the loaded property "due_date" should have type "date"
    And the loaded property "due_date" should have emoji "📅"

  Scenario: LoadType resolves use entry with pin override
    Given a shared properties file with "due_date" date and "priority" select properties
    And a type schema "project" with use "due_date" and pin 2
    When I load type "project"
    Then the loaded property "due_date" should have pin 2

  Scenario: LoadType resolves use entry with emoji override
    Given a shared properties file with "due_date" date and "priority" select properties
    And a type schema "project" with use "due_date" and emoji "🗓️"
    When I load type "project"
    Then the loaded property "due_date" should have emoji "🗓️"

  Scenario: LoadType resolves mixed use and name properties in order
    Given a shared properties file with "due_date" date and "priority" select properties
    And a type schema "project" with mixed use and name properties
    When I load type "project"
    Then the loaded type should have 3 properties
    And the loaded property at index 0 should be "title"
    And the loaded property at index 1 should be "due_date"
    And the loaded property at index 2 should be "budget"
