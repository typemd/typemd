Feature: Shared Properties
  Users can define shared property definitions as individual YAML files
  in the properties/ directory. Type schemas can reference these via
  `use` to reuse property definitions across multiple types without duplication.

  Background:
    Given a vault is ready

  # ── Loading per-property files ───────────────────────────────────────────

  Scenario: Load single per-property file
    Given a shared property file "due_date" with type "date" and emoji "📅"
    When I load shared properties
    Then shared properties should contain 1 entry
    And shared property "due_date" should have type "date"

  Scenario: Load multiple per-property files
    Given a shared property file "due_date" with type "date" and emoji "📅"
    And a shared property file "priority" with type "select" and options "high,medium,low"
    When I load shared properties
    Then shared properties should contain 2 entries
    And shared property "due_date" should have type "date"
    And shared property "priority" should have type "select"

  Scenario: Shared properties directory does not exist
    When I load shared properties
    Then shared properties should contain 0 entries

  Scenario: Empty shared properties directory
    Given an empty shared properties directory
    When I load shared properties
    Then shared properties should contain 0 entries

  Scenario: Non-YAML files in properties directory are ignored
    Given a shared property file "rating" with type "number" and emoji "⭐"
    And a non-YAML file "README.md" in the properties directory
    When I load shared properties
    Then shared properties should contain 1 entry
    And shared property "rating" should have type "number"

  Scenario: Name field in per-property file is ignored
    Given a shared property file "due_date" with type "date" and name override "something_else"
    When I load shared properties
    Then shared property "due_date" should have type "date"

  # ── Shared properties validation ──────────────────────────────────────────

  Scenario: Valid shared properties pass validation
    Given a shared property file "due_date" with type "date" and emoji "📅"
    And a shared property file "priority" with type "select" and options "high,medium,low"
    When I validate all schemas
    Then shared properties should have no errors

  Scenario: Invalid property type in shared properties is rejected
    Given a shared property file "bad_prop" with type "invalid"
    When I validate all schemas
    Then shared properties should have errors

  Scenario: Reserved name "name" in shared properties is rejected
    Given a shared property file "name" with type "string"
    When I validate all schemas
    Then shared properties should have errors

  Scenario: Select without options in shared properties is rejected
    Given a shared property file "status" with type "select"
    When I validate all schemas
    Then shared properties should have errors

  # ── Use keyword in type schemas ───────────────────────────────────────────

  Scenario: Type schema with use entry references shared property
    Given a shared property file "due_date" with type "date" and emoji "📅"
    And a shared property file "priority" with type "select" and options "high,medium,low"
    And a type schema "project" with use "due_date"
    When I validate all schemas
    Then schema "project" should have no errors

  Scenario: Use with pin override is accepted
    Given a shared property file "due_date" with type "date" and emoji "📅"
    And a type schema "project" with use "due_date" and pin 1
    When I validate all schemas
    Then schema "project" should have no errors

  Scenario: Use with emoji override is accepted
    Given a shared property file "due_date" with type "date" and emoji "📅"
    And a type schema "project" with use "due_date" and emoji "🗓️"
    When I validate all schemas
    Then schema "project" should have no errors

  Scenario: Use with description override is accepted
    Given a shared property file "due_date" with type "date" and emoji "📅"
    And a type schema "project" with use "due_date" and description "Project deadline"
    When I validate all schemas
    Then schema "project" should have no errors

  Scenario: Use with disallowed type field is rejected
    Given a shared property file "due_date" with type "date" and emoji "📅"
    And a type schema "project" with use "due_date" and disallowed type override
    When I validate all schemas
    Then schema "project" should have errors

  Scenario: Use referencing non-existent shared property is rejected
    Given a shared property file "due_date" with type "date" and emoji "📅"
    And a type schema "project" with use "nonexistent"
    When I validate all schemas
    Then schema "project" should have errors

  Scenario: Local property name conflicting with shared property is rejected
    Given a shared property file "due_date" with type "date" and emoji "📅"
    And a type schema "project" with local property "due_date"
    When I validate all schemas
    Then schema "project" should have errors

  Scenario: Duplicate use entries are rejected
    Given a shared property file "due_date" with type "date" and emoji "📅"
    And a type schema "project" with duplicate use "due_date"
    When I validate all schemas
    Then schema "project" should have errors

  Scenario: Use and name on same entry are rejected
    Given a shared property file "due_date" with type "date" and emoji "📅"
    And a type schema "project" with both use and name on same entry
    When I validate all schemas
    Then schema "project" should have errors

  # ── LoadType resolution ───────────────────────────────────────────────────

  Scenario: LoadType resolves use entry with no overrides
    Given a shared property file "due_date" with type "date" and emoji "📅"
    And a type schema "project" with use "due_date"
    When I load type "project"
    Then the loaded type should have 1 property
    And the loaded property "due_date" should have type "date"
    And the loaded property "due_date" should have emoji "📅"

  Scenario: LoadType resolves use entry with pin override
    Given a shared property file "due_date" with type "date" and emoji "📅"
    And a type schema "project" with use "due_date" and pin 2
    When I load type "project"
    Then the loaded property "due_date" should have pin 2

  Scenario: LoadType resolves use entry with emoji override
    Given a shared property file "due_date" with type "date" and emoji "📅"
    And a type schema "project" with use "due_date" and emoji "🗓️"
    When I load type "project"
    Then the loaded property "due_date" should have emoji "🗓️"

  Scenario: LoadType resolves use entry with description override
    Given a shared property file with described properties
    And a type schema "project" with use "due_date" and description "Project deadline"
    When I load type "project"
    Then the loaded property "due_date" description should be "Project deadline"

  Scenario: LoadType resolves mixed use and name properties in order
    Given a shared property file "due_date" with type "date" and emoji "📅"
    And a type schema "project" with mixed use and name properties
    When I load type "project"
    Then the loaded type should have 3 properties
    And the loaded property at index 0 should be "title"
    And the loaded property at index 1 should be "due_date"
    And the loaded property at index 2 should be "budget"
