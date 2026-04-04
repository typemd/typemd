Feature: Use entry description override
  Shared property use entries can override the description field
  to provide type-specific documentation for reused properties.

  Background:
    Given a vault is ready

  Scenario: Use with description override is accepted
    Given a shared property file "due_date" with type "date" and emoji "📅"
    And a shared property file "priority" with type "select" and options "high,medium,low"
    And a type schema "project" with use "due_date" and description "Project deadline"
    When I validate all schemas
    Then schema "project" should have no errors

  Scenario: LoadType resolves use entry with description override
    Given a shared property file "due_date" with type "date" and emoji "📅"
    And a shared property file "priority" with type "select" and options "high,medium,low"
    And a type schema "project" with use "due_date" and description "Project deadline"
    When I load type "project"
    Then the loaded property "due_date" description should be "Project deadline"

  Scenario: LoadType preserves shared description when no override
    Given a shared properties file with described properties
    And a type schema "project" with use "due_date"
    When I load type "project"
    Then the loaded property "due_date" description should be "A date something is due"
